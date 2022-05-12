package pipeline

import "sync"

type (
	// Atomic error that allows set/get concurrent operations
	atomicErr struct {
		sync.RWMutex
		value error
	}

	ConcurrentValue[T any] struct {
		mut   *sync.RWMutex
		v     T
		dirty bool
	}
)

// NewConcurrentValue creates a new value that can be safely mutated by different peers
func NewConcurrentValue[T any](v T) *ConcurrentValue[T] {
	return &ConcurrentValue[T]{
		mut:   new(sync.RWMutex),
		v:     v,
		dirty: false,
	}
}

// Read a value T stored inside
func (c *ConcurrentValue[T]) Read() T {
	c.mut.RLock()
	defer c.mut.RUnlock()
	return c.v
}

// Write safely input inside T value
func (c *ConcurrentValue[T]) Write(f func(*T)) {
	c.mut.Lock()
	defer c.mut.Unlock()
	f(&c.v)
	c.dirty = true
}

// Dirty returns true if the value has been changed at least once
func (c *ConcurrentValue[T]) Dirty() bool {
	c.mut.RLock()
	defer c.mut.RUnlock()
	return c.dirty
}

// Get the error stored
func (e *atomicErr) Get() error {
	e.RLock()
	defer e.RUnlock()
	return e.value
}

// Set the error stored
func (e *atomicErr) Set(err error) {
	e.Lock()
	e.value = err
	e.Unlock()
}

// Spawn a number of workers asynchronously, waiting for all of them to finish.
// After they're all done, if one of them failed the error is returned.
// If more than one fails, only the last one is returned
func spawnAsync(workers int, run func(index int) error) error {
	var wg sync.WaitGroup
	var errResult atomicErr

	spawn := func(index int) {
		err := run(index)

		if err != nil && errResult.Get() == nil {
			errResult.Set(err)
		}

		wg.Done()
	}

	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go spawn(i)
	}

	wg.Wait()
	return errResult.Get()
}
