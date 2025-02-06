package hw06pipelineexecution

import (
	"sync"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	wg := sync.WaitGroup{}
	out := make(Bi)

	go func() {
		defer close(out)

		for item := range in {
			wg.Add(1)

			go func(item interface{}) {
				defer wg.Done()

				current := make(Bi)
				go func() {
					current <- item
					close(current)
				}()

				var currentIn In = current

				for _, stage := range stages {
					select {
					case <-done:
						return
					default:
						currentIn = stage(currentIn)
					}
				}
				for resultItem := range currentIn {
					select {
					case <-done:
						return
					default:
						out <- resultItem
					}
				}
			}(item)
		}
		wg.Wait()
	}()

	return out
}
