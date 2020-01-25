package pipeline_stage

import (
	"github.com/saantiaguilera/go-pipeline"
)

type concurrentGroup []pipeline.Stage

func (s concurrentGroup) Run(executor pipeline.Executor) error {
	return spawnAsync(len(s), func(index int) error {
		return s[index].Run(executor)
	})
}

// Create a stage that will run each of the stages concurrently.
// The stage will wait for all of the stages to finish before returning.
//
// If one of them fails, the stage will wait until everyone finishes and after that return the error.
// If more than one fails, then the error will be the one delivered by the last failure.
func CreateConcurrentGroup(stages ...pipeline.Stage) pipeline.Stage {
	var stage concurrentGroup = stages
	return &stage
}
