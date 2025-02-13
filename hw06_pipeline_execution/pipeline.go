package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := make(Bi)
	defer close(out)

	interOut := in
	for _, stage := range stages {
		interOut = stageProcess(stage, interOut, done)
	}
	return interOut
}

func stageProcess(stage Stage, in In, done In) Out {
	interIn := make(Bi)
	go func(interIn Bi, in In, done In) {
		defer close(interIn)
		for item := range in {
			select {
			case <-done:
				go func(in In) {
					//nolint:revive
					for range in {
					}
				}(in)
				return
			default:
				interIn <- item
			}
		}
	}(interIn, in, done)

	out := stage(interIn)
	return out
}
