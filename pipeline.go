package pipeline

type Named interface {
	Name() string
}

type Pipeline interface {
	Run(stage Stage) error
	AddOnBeforePipelineRun(beforePipeline func(stage Stage) error)
	AddAfterPipelineRun(afterPipeline func(stage Stage, err error) error)
}

type pipeline struct {
	Before   []func(stage Stage) error
	After    []func(stage Stage, err error) error
	Executor Executor
}

func (p *pipeline) Run(stage Stage) error {
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

func (p *pipeline) AddOnBeforePipelineRun(beforePipeline func(stage Stage) error) {
	p.Before = append(p.Before, beforePipeline)
}

func (p *pipeline) AddAfterPipelineRun(afterPipeline func(stage Stage, err error) error) {
	p.After = append(p.After, afterPipeline)
}

func CreatePipeline(executor Executor) Pipeline {
	return &pipeline{
		Executor: executor,
		Before:   []func(stage Stage) error{},
		After:    []func(stage Stage, err error) error{},
	}
}