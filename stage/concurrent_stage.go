package pipeline_stage

import (
	"github.com/saantiaguilera/go-pipeline"
)

type concurrentStage []pipeline.Step

func (s concurrentStage) Run(executor pipeline.Executor) error {
	return runAsync(len(s), func(index int) error {
		return executor.Run(s[index])
	})
}

func CreateConcurrentStage(steps ...pipeline.Step) pipeline.Stage {
	var stage concurrentStage = steps
	return &stage
}
