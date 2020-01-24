package pipeline_stage

import "github.com/saantiaguilera/go-pipeline"

type SequentialStage []pipeline.Step

func (s *SequentialStage) Run(executor pipeline.Executor) error {
	for _, step := range *s {
		err := executor.Run(step)

		if err != nil {
			return err
		}
	}
	return nil
}