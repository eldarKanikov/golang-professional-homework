package hw06pipelineexecution

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	sleepPerStage = time.Millisecond * 100
	fault         = sleepPerStage / 2
)

var isFullTesting = true

func generateStages() []Stage {
	// Stage generator
	g := func(_ string, f func(v interface{}) interface{}) Stage {
		return func(in In) Out {
			out := make(Bi)
			go func() {
				defer close(out)
				for v := range in {
					time.Sleep(sleepPerStage)
					out <- f(v)
				}
			}()
			return out
		}
	}

	return []Stage{
		g("Dummy", func(v interface{}) interface{} { return v }),
		g("Multiplier (* 2)", func(v interface{}) interface{} { return v.(int) * 2 }),
		g("Adder (+ 100)", func(v interface{}) interface{} { return v.(int) + 100 }),
		g("Stringifier", func(v interface{}) interface{} { return strconv.Itoa(v.(int)) }),
	}
}

func TestPipeline(t *testing.T) {

	stages := generateStages()

	t.Run("simple case", func(t *testing.T) {
		in := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 10)
		start := time.Now()
		for s := range ExecutePipeline(in, nil, stages...) {
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)

		require.Equal(t, []string{"102", "104", "106", "108", "110"}, result)
		require.Less(t,
			int64(elapsed),
			// ~0.8s for processing 5 values in 4 stages (100ms every) concurrently
			int64(sleepPerStage)*int64(len(stages)+len(data)-1)+int64(fault))
	})

	t.Run("done case", func(t *testing.T) {
		in := make(Bi)
		done := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		// Abort after 200ms
		abortDur := sleepPerStage * 2
		go func() {
			<-time.After(abortDur)
			close(done)
		}()

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 10)
		start := time.Now()
		for s := range ExecutePipeline(in, done, stages...) {
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)

		require.Len(t, result, 0)
		require.Less(t, int64(elapsed), int64(abortDur)+int64(fault))
	})
}

func TestPanic(t *testing.T) {
	t.Run("stage with panic in the middle", func(t *testing.T) {
		in := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		go func() {
			defer close(in)
			for _, v := range data {
				in <- v
			}
		}()

		panicChan := make(chan interface{}, 1)

		safeG := func(name string, f func(v interface{}) interface{}) Stage {
			return func(in In) Out {
				out := make(Bi)
				go func() {
					defer func() {
						if r := recover(); r != nil {
							panicChan <- r
							close(panicChan)
						}
						close(out)
					}()

					for v := range in {
						time.Sleep(sleepPerStage)
						out <- f(v)
					}
				}()
				return out
			}
		}

		stages2 := []Stage{
			safeG("Dummy", func(v interface{}) interface{} { return v }),
			safeG("With Panic", func(v interface{}) interface{} {
				if v.(int) == 3 {
					panic("surprise!  panic: value is 3")
				}
				return v
			}),
			safeG("Multiplier (* 2)", func(v interface{}) interface{} { return v.(int) * 2 }),
		}

		var results []interface{}
		for v := range ExecutePipeline(in, nil, stages2...) {
			results = append(results, v)
		}

		require.Len(t, results, 2)
		require.Equal(t, 2, results[0].(int))
		require.Equal(t, 4, results[1].(int))

		select {
		case panicVal := <-panicChan:
			require.Equal(t, panicVal.(string), "surprise!  panic: value is 3")
		}
	})
}

func TestEmptyCases(t *testing.T) {
	stages := generateStages()

	t.Run("empty stages, same result", func(t *testing.T) {
		in := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		var result []interface{}
		for s := range ExecutePipeline(in, nil, []Stage{}...) {
			result = append(result, s)
		}

		expected := make([]interface{}, len(data))
		for i, v := range data {
			expected[i] = v
		}

		require.Len(t, result, 5)
		require.Equal(t, expected, result)
	})

	t.Run("empty in", func(t *testing.T) {
		in := make(Bi)
		done := make(Bi)

		go func() {
			close(in)
		}()

		result := make([]interface{}, 0, 10)
		for s := range ExecutePipeline(in, done, stages...) {
			result = append(result, s.(string))
		}
		require.Len(t, result, 0)
		require.Len(t, stages, 4)
	})
}

func TestAllStageStop(t *testing.T) {
	if !isFullTesting {
		return
	}
	wg := sync.WaitGroup{}
	// Stage generator
	g := func(_ string, f func(v interface{}) interface{}) Stage {
		return func(in In) Out {
			out := make(Bi)
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer close(out)
				for v := range in {
					time.Sleep(sleepPerStage)
					out <- f(v)
				}
			}()
			return out
		}
	}

	stages := []Stage{
		g("Dummy", func(v interface{}) interface{} { return v }),
		g("Multiplier (* 2)", func(v interface{}) interface{} { return v.(int) * 2 }),
		g("Adder (+ 100)", func(v interface{}) interface{} { return v.(int) + 100 }),
		g("Stringifier", func(v interface{}) interface{} { return strconv.Itoa(v.(int)) }),
	}

	t.Run("done case", func(t *testing.T) {
		in := make(Bi)
		done := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		// Abort after 200ms
		abortDur := sleepPerStage * 2
		go func() {
			<-time.After(abortDur)
			close(done)
		}()

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 10)
		for s := range ExecutePipeline(in, done, stages...) {
			result = append(result, s.(string))
		}
		wg.Wait()

		require.Len(t, result, 0)
	})
}
