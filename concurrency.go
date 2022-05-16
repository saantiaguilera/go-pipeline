package pipeline

import "sync"

type (
	// Atomic error that allows set/get concurrent operations
	atomicErr struct {
		sync.RWMutex
		value error
	}

	// atomicSlice that allows storing of values of type T and appending them concurrently
	atomicSlice[T any] struct {
		sync.RWMutex
		v []T
	}
)

// Read slice T stored inside
func (s *atomicSlice[T]) Get() []T {
	s.RLock()
	defer s.RUnlock()
	return s.v
}

// Append safely inside T slice
func (s *atomicSlice[T]) Append(t T) {
	s.Lock()
	defer s.Unlock()
	s.v = append(s.v, t)
}

// Get the error stored
func (e *atomicErr) Get() error {
	e.RLock()
	defer e.RUnlock()
	return e.value
}

// Set the error stored if there's not one already stored
func (e *atomicErr) SetIfNil(err error) {
	e.Lock()
	defer e.Unlock()
	if e.value == nil {
		e.value = err
	}
}

// Run a number of workers concurrently, waiting for all of them to finish.
// After they're all done, if one of them failed the error is returned.
// If more than one fails, the first error is returned
func runConcurrently[R, O any](workers []R, run func(R) (O, error)) ([]O, error) {
	runIndexed := func(wg *sync.WaitGroup, errResult *atomicErr, mergedRes *atomicSlice[O], i int) {
		res, err := run(workers[i])

		if err != nil {
			errResult.SetIfNil(err)
		}

		mergedRes.Append(res)
		wg.Done()
	}

	var wg sync.WaitGroup
	var errResult atomicErr
	var mergedRes atomicSlice[O]

	wg.Add(len(workers))
	if len(workers) > 1 {
		for i := 0; i < len(workers); i++ {
			go runIndexed(&wg, &errResult, &mergedRes, i)
		}
		wg.Wait()
	} else { // avoid concurrency, no need to spawn and wait just use current
		runIndexed(&wg, &errResult, &mergedRes, 0)
	}

	return mergedRes.Get(), errResult.Get()
}
