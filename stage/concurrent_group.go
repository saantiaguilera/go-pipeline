package pipeline_stage

import (
	"github.com/saantiaguilera/go-pipeline"
)

type concurrentGroup []pipeline.Stage

func (s concurrentGroup) Run(executor pipeline.Executor) error {
	return runAsync(len(s), func(index int) error {
		return s[index].Run(executor)
	})
}

func CreateConcurrentGroup(stages ...pipeline.Stage) pipeline.Stage {
	var stage concurrentGroup = stages
	return &stage
}
