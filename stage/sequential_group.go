package pipeline_stage

import "github.com/saantiaguilera/go-pipeline"

type sequentialGroup []pipeline.Stage

func (s sequentialGroup) Run(executor pipeline.Executor) error {
	for _, stage := range s {
		err := stage.Run(executor)

		if err != nil {
			return err
		}
	}
	return nil
}

func CreateSequentialGroup(stages ...pipeline.Stage) pipeline.Stage {
	var stage sequentialGroup = stages
	return &stage
}
