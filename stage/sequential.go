package pipeline_stage

import "github.com/saantiaguilera/go-pipeline"

type Sequential []pipeline.Stage

func (s *Sequential) Run(executor pipeline.Executor) error {
	for _, stage := range *s {
		err := stage.Run(executor)

		if err != nil {
			return err
		}
	}
	return nil
}