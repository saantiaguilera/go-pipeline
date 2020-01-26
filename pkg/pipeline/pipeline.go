package pipeline

import (
	"github.com/saantiaguilera/go-pipeline/pkg"
)

type pipe struct {
	Before   []func(stage pkg.Stage) error
	After    []func(stage pkg.Stage, err error) error
	Executor pkg.Executor
}

func (p *pipe) Run(stage pkg.Stage) error {
	for _, before := range p.Before {
		err := before(stage)

		if err != nil {
			return err
		}
	}

	err := stage.Run(p.Executor)

	for _, after := range p.After {
		err = after(stage, err)
	}

	return err
}

func (p *pipe) AddBeforeRunHook(beforePipeline func(stage pkg.Stage) error) {
	p.Before = append(p.Before, beforePipeline)
}

func (p *pipe) AddAfterRunHook(afterPipeline func(stage pkg.Stage, err error) error) {
	p.After = append(p.After, afterPipeline)
}

// Create a pipeline with a given executor
func CreatePipeline(executor pkg.Executor) pkg.Pipeline {
	return &pipe{
		Executor: executor,
		Before:   []func(stage pkg.Stage) error{},
		After:    []func(stage pkg.Stage, err error) error{},
	}
}
