package hw06pipelineexecution

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestExecutePipelineEdgeCases(t *testing.T) {
	t.Run("without stages returns input", func(t *testing.T) {
		in := make(Bi)
		go func() {
			defer close(in)
			in <- 1
			in <- 2
		}()

		result := make([]interface{}, 0)
		for value := range ExecutePipeline(in, nil) {
			result = append(result, value)
		}

		require.Equal(t, []interface{}{1, 2}, result)
	})

	t.Run("one stage", func(t *testing.T) {
		in := make(Bi)
		go func() {
			defer close(in)
			for i := 0; i < 3; i++ {
				in <- i
			}
		}()

		stage := func(in In) Out {
			out := make(Bi)
			go func() {
				defer close(out)
				for value := range in {
					out <- value.(int) * 2
				}
			}()
			return out
		}

		result := make([]int, 0)
		for value := range ExecutePipeline(in, nil, stage) {
			result = append(result, value.(int))
		}

		require.Equal(t, []int{0, 2, 4}, result)
	})

	t.Run("stops on done without leaking", func(t *testing.T) {
		in := make(Bi)
		done := make(Bi)

		stage := func(in In) Out {
			out := make(Bi)
			go func() {
				defer close(out)
				for value := range in {
					time.Sleep(5 * time.Millisecond)
					out <- value
				}
			}()
			return out
		}

		go func() {
			defer close(in)
			for i := 0; i < 1000; i++ {
				in <- i
			}
		}()

		out := ExecutePipeline(in, done, stage, stage)
		go func() {
			time.Sleep(30 * time.Millisecond)
			close(done)
		}()

		start := time.Now()
		result := make([]interface{}, 0)
		for value := range out {
			result = append(result, value)
		}

		require.Less(t, len(result), 1000)
		require.Less(t, time.Since(start), 500*time.Millisecond)
	})
}
