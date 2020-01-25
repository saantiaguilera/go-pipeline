package pipeline_stage

import "github.com/saantiaguilera/go-pipeline"

type sequentialGroup []pipeline.Stage

func (s sequentialGroup) Run(executor pipeline.Executor) error {
	return runSync(len(s), func(index int) error {
		return s[index].Run(executor)
	})
}

// Create a stage that will run each of stages sequentially. If one of them fails, the operation will abort immediately
func CreateSequentialGroup(stages ...pipeline.Stage) pipeline.Stage {
	var stage sequentialGroup = stages
	return &stage
}
