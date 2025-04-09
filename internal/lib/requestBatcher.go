package lib

import (
	"log"
	"sync"
	"sync/atomic"
)

var activeWorkers int32

func RunBatch(tasks []func(), batchSize int) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, batchSize)

	taskLen := len(tasks)

	for i, task := range tasks {
		wg.Add(1)
		go func(task func()) {
			defer wg.Done()

			// Acquire semaphore
			sem <- struct{}{}
			//LOG

			curr := atomic.AddInt32(&activeWorkers, 1)
			log.Printf("START: task %v/%v (active workers: %d)", i, taskLen, curr)

			defer func() {
				<-sem // Release after work

				//LOG
				curr := atomic.AddInt32(&activeWorkers, -1)
				log.Printf("DONE: task %v/%v (active workers: %d)", i, taskLen, curr)
			}()

			task()
		}(task)
	}

	wg.Wait()
}

func RunBatchSizeClosure(batchSize int) func(tasks []func()) {
	return func(tasks []func()) {
		RunBatch(tasks, batchSize)
	}
}
