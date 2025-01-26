package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
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

		require.Eventually(t, func() bool {
			return runTasksCount == int32(tasksCount)
		}, 1*time.Second, 250*time.Millisecond, "not all tasks were completed")

		require.Eventually(t, func() bool {
			return int64(elapsedTime) <= int64(sumTime/2)
		}, 1*time.Second, 250*time.Millisecond, "tasks were run sequentially?")
	})

	t.Run("test m values", func(t *testing.T) {
		taskCount := 10
		var runTasksCount atomic.Int32
		taskSleep := 100 * time.Millisecond
		tasks := make([]Task, 0, taskCount)

		for i := 0; i < taskCount; i++ {
			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				runTasksCount.Add(1)
				if i%2 == 0 {
					return nil
				}
				return errors.New("error")
			})
		}

		workersCount := 5
		maxErrorsCount := 0
		err := Run(tasks, workersCount, maxErrorsCount)
		require.NoError(t, err)

		maxErrorsCount = -1
		err = Run(tasks, workersCount, maxErrorsCount)
		require.NoError(t, err)

		maxErrorsCount = 1
		err = Run(tasks, workersCount, maxErrorsCount)
		require.Error(t, err)

		maxErrorsCount = taskCount / 2
		err = Run(tasks, workersCount, maxErrorsCount)
		require.Error(t, err)

		maxErrorsCount = taskCount/2 + 1
		err = Run(tasks, workersCount, maxErrorsCount)
		require.NoError(t, err)
	})
}
