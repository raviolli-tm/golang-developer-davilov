package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")

	})

	t.Run("tasks without errors via require.Eventually", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * 35
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				for i := 0; i < 100_000_000; i++ {
					j := i
					i = j
				}
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		runtime.GOMAXPROCS(2)
		workersCount := 5
		maxErrorsCount := 1
		fmt.Println(sumTime)

		start := time.Now()
		for i := 0; i < 100_000_000; i++ {
			j := i
			i = j
		}
		elapsedTime := time.Since(start)
		fmt.Println("------", elapsedTime)

		start = time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)

		elapsedTime = time.Since(start)
		fmt.Println(elapsedTime)

		require.NoError(t, err)
		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")

	})
}
