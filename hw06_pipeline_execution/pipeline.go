package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

// ExecutePipeline connects the stages and returns the output of the last one.
// Closing done stops the pipeline: every stage stops reading its input and
// closes its output, so no goroutines are left behind.
func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in

	for _, stage := range stages {
		out = executeStage(stage, out, done)
	}

	return out
}

func executeStage(stage Stage, in In, done In) Out {
	stageIn := make(Bi)
	go proxyInput(in, stageIn, done)

	stageOut := stage(stageIn)
	out := make(Bi)
	go proxyOutput(stageOut, out, done)

	return out
}

func proxyInput(in In, out Bi, done In) {
	defer close(out)

	for {
		if cancelled(done) {
			drain(in)
			return
		}

		select {
		case <-done:
			drain(in)
			return
		case value, ok := <-in:
			if !ok {
				return
			}

			select {
			case out <- value:
			case <-done:
				drain(in)
				return
			}
		}
	}
}

func proxyOutput(in In, out Bi, done In) {
	defer close(out)

	for {
		if cancelled(done) {
			discard(in)
			return
		}

		select {
		case <-done:
			// The stage input is already closed by proxyInput, so the stage
			// will close its output soon; drain it to let the stage finish.
			discard(in)
			return
		case value, ok := <-in:
			if !ok {
				return
			}

			select {
			case out <- value:
			case <-done:
				discard(in)
				return
			}
		}
	}
}

func cancelled(done In) bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

// drain consumes the rest of the input in the background so a producer
// blocked on sending does not leak after cancellation.
func drain(in In) {
	if in == nil {
		return
	}

	go func() {
		for range in {
		}
	}()
}

func discard(in In) {
	for range in {
	}
}
