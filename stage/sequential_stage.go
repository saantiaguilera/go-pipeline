package pipeline_stage

import "github.com/saantiaguilera/go-pipeline"

type sequentialStage []pipeline.Step

func (s sequentialStage) Run(executor pipeline.Executor) error {
	for _, step := range s {
		err := executor.Run(step)

		if err != nil {
			return err
		}
	}
	return nil
}

func CreateSequentialStage(steps ...pipeline.Step) pipeline.Stage {
	var stage sequentialStage = steps
	return &stage
}
