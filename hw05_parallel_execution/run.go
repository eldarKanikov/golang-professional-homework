package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	wg := sync.WaitGroup{}
	stopChan := make(chan struct{})
	taskChan := make(chan Task)
	var errorsCount atomic.Int32
	var once sync.Once
	for i := 0; i < n; i++ {
		wg.Add(1)
		initWorker(&wg, stopChan, taskChan, m, &errorsCount, &once)
	}

	putTasksToChannel(tasks, taskChan, stopChan)

	wg.Wait()

	if m > 0 && int(errorsCount.Load()) >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func initWorker(wg *sync.WaitGroup, stopChan chan struct{}, taskChan chan Task, m int, errorsCount *atomic.Int32,
	once *sync.Once,
) {
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stopChan:
				return
			default:
				task, ok := <-taskChan
				if !ok {
					return
				}
				err := task()
				handleError(m, err, errorsCount, once, stopChan)
			}
		}
	}()
}

func putTasksToChannel(tasks []Task, taskChan chan Task, stopChan chan struct{}) {
	for _, task := range tasks {
		select {
		case <-stopChan:
			break
		case taskChan <- task:
		}
	}
	close(taskChan)
}

func handleError(m int, err error, errorsCount *atomic.Int32, once *sync.Once, stopChan chan struct{}) {
	if m > 0 {
		if err != nil {
			errorsCount.Add(1)
		}
		if int(errorsCount.Load()) >= m {
			once.Do(func() {
				close(stopChan)
			})
		}
	}
}
