package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	interOut := in
	for _, stage := range stages {
		interOut = stageProcess(stage, interOut, done)
	}
	return interOut
}

func stageProcess(stage Stage, in In, done In) Out {
	interIn := make(Bi)
	go func() {
		for {
			select {
			case <-done:
				go func() {
					for item := range in {
						_ = item
					}
				}()
				close(interIn)
				return
			case item, ok := <-in:
				if !ok {
					close(interIn)
					return
				}
				interIn <- item
			}
		}
	}()

	return stage(interIn)
}
