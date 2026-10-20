package hw05parallelexecution

import (
	"errors"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunEdgeCases(t *testing.T) {
	t.Run("zero error limit without errors", func(t *testing.T) {
		var ran int32
		tasks := make([]Task, 5)
		for i := range tasks {
			tasks[i] = func() error {
				atomic.AddInt32(&ran, 1)
				return nil
			}
		}

		require.NoError(t, Run(tasks, 2, 0))
		require.Equal(t, int32(5), ran)
	})

	t.Run("zero error limit with an error", func(t *testing.T) {
		tasks := []Task{
			func() error { return errors.New("boom") },
			func() error { return nil },
		}

		require.ErrorIs(t, Run(tasks, 1, 0), ErrErrorsLimitExceeded)
	})

	t.Run("exactly m errors", func(t *testing.T) {
		tasks := []Task{
			func() error { return errors.New("first") },
			func() error { return errors.New("second") },
			func() error { return nil },
		}

		require.ErrorIs(t, Run(tasks, 1, 2), ErrErrorsLimitExceeded)
	})

	t.Run("fewer than m errors", func(t *testing.T) {
		tasks := []Task{
			func() error { return errors.New("first") },
			func() error { return nil },
		}

		require.NoError(t, Run(tasks, 1, 2))
	})

	t.Run("empty task list", func(t *testing.T) {
		require.NoError(t, Run(nil, 3, 1))
	})
}
