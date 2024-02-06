package worker

import (
	"fmt"
	"sync"

	"gitlab.mvalley.com/wind/rime-utils/internal/pkg/storage"
)

type Job interface {
	Run()
	Stop()
	GetID() string
	SaveSyncTask() error
}

type Worker struct {
	jobQueue chan Job
	Mu       sync.Mutex
	jobMap   map[string]Job
	RunCount chan int
	storage  *storage.Storage
}

func NewWorker(storage *storage.Storage, workers int) *Worker {
	return &Worker{
		jobQueue: make(chan Job, 100000), // 消息队列
		jobMap:   make(map[string]Job),
		RunCount: make(chan int, workers), // 最大运行数
		storage:  storage,
	}
}

func (w *Worker) Run() {
	go func() {
		for job := range w.jobQueue {
			w.RunCount <- 1
			go func(job Job) {
				defer func() {
					<-w.RunCount
				}()
				ex := w.Exist(job.GetID())
				if !ex {
					return
				}
				job.Run()
				err := job.SaveSyncTask()
				if err != nil {
					fmt.Println(err)
				}
				w.Done(job.GetID())
			}(job)
		}
	}()
}

// 添加任务
func (w *Worker) Append(job Job) {
	w.Mu.Lock()
	w.jobMap[job.GetID()] = job
	w.Mu.Unlock()
	w.jobQueue <- job

}

// 判断同步任务是否存在
func (w *Worker) Exist(id string) bool {
	w.Mu.Lock()
	defer w.Mu.Unlock()
	_, ex := w.jobMap[id]
	return ex
}

// 移除工作任务
func (w *Worker) Done(id string) {
	// 再移除
	w.Mu.Lock()
	defer w.Mu.Unlock()
	delete(w.jobMap, id)
}

// 停止工作任务
func (w *Worker) Stop(id string) {
	w.Mu.Lock()
	defer w.Mu.Unlock()
	job, ex := w.jobMap[id]
	if !ex {
		return
	}
	job.Stop()
}
