package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline2(in In, done In, stages ...Stage) Out {
	if len(stages) == 1 {
		return stages[0](in)
	}

	return ExecutePipeline2(stages[0](in), done, stages[1:]...)
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	preparedOut := make(Bi)
	outChan := stages[0](in)
	channel := ExecutePipeline2(outChan, done, stages[1:]...)

	go func(out Out, preparedOut Bi) {
		defer close(preparedOut)
		for {
			select {
			case v, ok := <-out:
				if !ok {
					return
				}
				preparedOut <- v
			case <-done:
				return
			}
		}
	}(channel, preparedOut)

	return preparedOut
}
