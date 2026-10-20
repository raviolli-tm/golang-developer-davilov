package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors
// from tasks. Tasks that have already been started are allowed to finish.
func Run(tasks []Task, n, m int) error {
	if n <= 0 {
		if len(tasks) == 0 {
			return nil
		}
		n = 1
	}

	if m < 0 {
		m = 0
	}

	var (
		wg       sync.WaitGroup
		mu       sync.Mutex
		next     int
		errCount int
		stopped  bool
	)

	limitReached := func() bool {
		if m == 0 {
			return errCount > 0
		}
		return errCount >= m
	}

	worker := func() {
		defer wg.Done()

		for {
			mu.Lock()
			if stopped || next >= len(tasks) {
				mu.Unlock()
				return
			}
			task := tasks[next]
			next++
			mu.Unlock()

			if err := task(); err != nil {
				mu.Lock()
				errCount++
				if limitReached() {
					stopped = true
				}
				mu.Unlock()
			}
		}
	}

	wg.Add(n)
	for i := 0; i < n; i++ {
		go worker()
	}
	wg.Wait()

	mu.Lock()
	defer mu.Unlock()

	if m == 0 {
		if errCount > 0 {
			return ErrErrorsLimitExceeded
		}
		return nil
	}

	if errCount >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}
