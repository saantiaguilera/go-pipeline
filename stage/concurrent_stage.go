package pipeline_stage

import (
	"github.com/saantiaguilera/go-pipeline"
	"sync"
)

type ConcurrentStage []pipeline.Step

func (s *ConcurrentStage) Run(executor pipeline.Executor) error {
	var wg sync.WaitGroup
	var finalErr error

	run := func(step pipeline.Step) {
		err := executor.Run(step)

		if err != nil {
			finalErr = err
		}

		wg.Done()
	}

	wg.Add(len(*s))
	for _, c := range *s {
		go run(c)
	}

	wg.Wait()
	return finalErr
}