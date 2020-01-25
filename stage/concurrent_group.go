package pipeline_stage

import (
	"github.com/saantiaguilera/go-pipeline"
	"sync"
)

type concurrentGroup []pipeline.Stage

func (s concurrentGroup) Run(executor pipeline.Executor) error {
	var wg sync.WaitGroup
	var finalErr error

	run := func(stage pipeline.Stage) {
		err := stage.Run(executor)

		if err != nil {
			finalErr = err
		}

		wg.Done()
	}

	wg.Add(len(s))
	for _, c := range s {
		go run(c)
	}

	wg.Wait()
	return finalErr
}

func CreateConcurrentGroup(stages ...pipeline.Stage) pipeline.Stage {
	var stage concurrentGroup = stages
	return &stage
}
