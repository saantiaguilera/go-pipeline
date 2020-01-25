package pipeline_core

import "github.com/saantiaguilera/go-pipeline"

type pipe struct {
	Before   []func(stage pipeline.Stage) error
	After    []func(stage pipeline.Stage, err error) error
	Executor pipeline.Executor
}

func (p *pipe) Run(stage pipeline.Stage) error {
	for _, before := range p.Before {
		err := before(stage)

		if err != nil {
			return err
		}
	}

	err := stage.Run(p.Executor)

	for _, after := range p.After {
		err = after(stage, err)

		if err != nil {
			return err
		}
	}

	return err
}

func (p *pipe) AddOnBeforePipelineRun(beforePipeline func(stage pipeline.Stage) error) {
	p.Before = append(p.Before, beforePipeline)
}

func (p *pipe) AddAfterPipelineRun(afterPipeline func(stage pipeline.Stage, err error) error) {
	p.After = append(p.After, afterPipeline)
}

func CreatePipeline(executor pipeline.Executor) pipeline.Pipeline {
	return &pipe{
		Executor: executor,
		Before:   []func(stage pipeline.Stage) error{},
		After:    []func(stage pipeline.Stage, err error) error{},
	}
}
