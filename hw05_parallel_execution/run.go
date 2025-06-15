package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	var wg sync.WaitGroup
	wg.Add(n)
	lenTask := int32(len(tasks))
	var taskCount int32
	var errCount int32

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for {
				if int32(m) <= atomic.LoadInt32(&errCount) {
					return
				}
				idx := atomic.AddInt32(&taskCount, 1)
				if idx <= lenTask {
					err := tasks[idx-1]()
					if err != nil {
						atomic.AddInt32(&errCount, 1)
					}
				} else {
					return
				}
			}
		}()
	}

	wg.Wait()
	if errCount > int32(m) {
		return ErrErrorsLimitExceeded
	}
	return nil
}
