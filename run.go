package pipeline

import "sync"

type (
	// Atomic error that allows set/get concurrent operations
	atomicErr struct {
		sync.RWMutex
		value error
	}
)

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
