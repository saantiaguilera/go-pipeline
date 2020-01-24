package pipeline_stage

import (
	"github.com/saantiaguilera/go-pipeline"
	"sync"
)

type Concurrent []pipeline.Stage

func (s *Concurrent) Run(executor pipeline.Executor) error {
	var wg sync.WaitGroup
	var finalErr error

	run := func(stage pipeline.Stage) {
		err := stage.Run(executor)

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