package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, stopCh chan struct{}) {
	for {
		select {
		case <-stopCh:
			fmt.Printf("Worker %d: Stopping\n", id)
			return
		default:
			// 模拟工作
			fmt.Printf("Worker %d: Working...\n", id)
			time.Sleep(time.Second * 100)
		}
	}
}

func main() {
	numWorkers := 3
	stopCh := make(chan struct{})
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			worker(id, stopCh)
		}(i)
	}

	// 主协程等待一段时间后发送停止信号
	time.Sleep(3 * time.Second)
	close(stopCh)

	// 等待所有协程完成
	wg.Wait()

	fmt.Println("Main goroutine exiting.")
}
