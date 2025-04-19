package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var wg sync.WaitGroup
	intChan := make(chan int)
	errorChan := make(chan error)
	stopChan := make(chan struct{})

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(stopChan <-chan struct{}) {
			defer wg.Done()
			for {
				select {
				case <-stopChan:
					return
				case idx := <-intChan:
					errorChan <- tasks[idx]()
				}
			}

		}(stopChan)
	}

	go func() {
		for i := 0; i < len(tasks); i++ {
			select {
			case <-stopChan:
				return
			default:
				intChan <- i
			}

		}
		close(stopChan)
	}()

	go func() {
		for err := range errorChan {
			fmt.Println("err", m, err)
			if err != nil {
				m--
				if m == 0 {
					close(stopChan)
				}
			}
		}
	}()

	wg.Wait()
	close(errorChan)
	if m > 0 {
		return nil
	}
	return ErrErrorsLimitExceeded

}
