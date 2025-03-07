package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	var errorCount int64
	taskChan := make(chan Task)
	wg := sync.WaitGroup{}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskChan {
				if atomic.LoadInt64(&errorCount) >= int64(m) {
					break
				}

				err := task()
				if err != nil {
					atomic.AddInt64(&errorCount, 1)
				}
			}
		}()
	}

	for _, task := range tasks {
		if atomic.LoadInt64(&errorCount) >= int64(m) {
			break
		}
		taskChan <- task
	}
	close(taskChan)
	wg.Wait()

	if atomic.LoadInt64(&errorCount) >= int64(m) {
		return ErrErrorsLimitExceeded
	}
	return nil
}
