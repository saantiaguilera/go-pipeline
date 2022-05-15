package pipeline

import "sync"

type (
	// Atomic error that allows set/get concurrent operations
	atomicErr struct {
		sync.RWMutex
		value error
	}

	mergedResult[T any] struct {
		mut   *sync.RWMutex
		v     []T
	}
)

// newMergedResult creates a new value that can be safely mutated by different peers
func newMergedResult[T any](cap int) *mergedResult[T] {
	return &mergedResult[T]{
		mut:   new(sync.RWMutex),
		v:     make([]T, 0, cap),
	}
}

// Read slice T stored inside
func (c *mergedResult[T]) Read() []T {
	c.mut.RLock()
	defer c.mut.RUnlock()
	return c.v
}

// Append safely inside T slice
func (c *mergedResult[T]) Append(t T) {
	c.mut.Lock()
	defer c.mut.Unlock()
	c.v = append(c.v, t)
}

// Get the error stored
func (e *atomicErr) Get() error {
	e.RLock()
	defer e.RUnlock()
	return e.value
}

// Set the error stored
func (e *atomicErr) SetIfNil(err error) {
	e.Lock()
	defer e.Unlock()
	if e.value == nil {
		e.value = err
	}
}

// Spawn a number of workers asynchronously, waiting for all of them to finish.
// After they're all done, if one of them failed the error is returned.
// If more than one fails, the first error is returned
func spawnAsync[R, O any](workers []R, run func(R) (O, error)) ([]O, error) {
	var wg sync.WaitGroup
	var errResult atomicErr
	mergedRes := newMergedResult[O](len(workers))

	spawn := func(i int) {
		res, err := run(workers[i])

		if err != nil {
			errResult.SetIfNil(err)
		}

		mergedRes.Append(res)
		wg.Done()
	}

	wg.Add(len(workers))
	if len(workers) > 1 {
		for i := 0; i < len(workers); i++ {
			go spawn(i)
		}
		wg.Wait()
	} else { // avoid concurrency, no need to spawn and wait just use current
		spawn(0)
	}

	return mergedRes.Read(), errResult.Get()
}
