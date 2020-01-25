package pipeline_stage

import "github.com/saantiaguilera/go-pipeline"

type sequentialStage []pipeline.Step

func (s sequentialStage) Run(executor pipeline.Executor) error {
	return runSync(len(s), func(index int) error {
		return executor.Run(s[index])
	})
}

// Create a stage that will run each of the steps sequentially. If one of them fails, the operation will abort immediately
func CreateSequentialStage(steps ...pipeline.Step) pipeline.Stage {
	var stage sequentialStage = steps
	return &stage
}
