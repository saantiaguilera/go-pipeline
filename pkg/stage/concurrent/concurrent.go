package concurrent

import "sync"

// Spawn a number of workers asynchronously, waiting for all of them to finish.
// After they're all done, if one of them failed the error is returned.
// If more than one fails, only the last one is returned
func spawnAsync(workers int, run func(index int) error) error {
	var wg sync.WaitGroup
	var finalErr error

	spawn := func(index int) {
		err := run(index)

		if err != nil {
			finalErr = err
		}

		wg.Done()
	}

	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go spawn(i)
	}

	wg.Wait()
	return finalErr
}
