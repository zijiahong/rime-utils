package utils

import "sync"

type Job interface {
	Run()
	Stop(id string)
}

type MysqlJob struct {
	id     string
	fn     func(stopCh chan struct{})
	stopCh chan struct{}
}

func NewMysqlWork(id string, fn func(stopCh chan struct{})) *MysqlJob {

	return &MysqlJob{
		id:     id,
		fn:     fn,
		stopCh: make(chan struct{}),
		// wg:       ,
	}
}

type Worker struct {
	jobQueue chan Job
	RunCount int
	wg       *sync.WaitGroup
}

func NewWorker(workers int) *Worker {
	var wg sync.WaitGroup
	return &Worker{
		jobQueue: make(chan Job, workers),
		wg:       &wg,
	}
}

func (w *Worker) Start() {
	defer w.wg.Done()
	for job := range w.jobQueue {
		job.Run()

	}
}

func (w *Worker) Move(id string) {}
