package pipeline_stage

import "github.com/saantiaguilera/go-pipeline"

type sequentialStage []pipeline.Step

func (s sequentialStage) Run(executor pipeline.Executor) error {
	return runSync(len(s), func(index int) error {
		return executor.Run(s[index])
	})
}

func CreateSequentialStage(steps ...pipeline.Step) pipeline.Stage {
	var stage sequentialStage = steps
	return &stage
}
