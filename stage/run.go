package pipeline_stage

import (
	"sync"
)

func runAsync(workers int, run func(index int) error) error {
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

func runSync(workers int, run func(index int) error) error {
	for i := 0; i < workers; i++ {
		err := run(i)

		if err != nil {
			return err
		}
	}
	return nil
}