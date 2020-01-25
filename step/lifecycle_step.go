package pipeline_step

import "github.com/saantiaguilera/go-pipeline"

type BeforeStep func() error
type AfterStep func(step pipeline.Step, err error) error

type lifecycleStep struct {
	Before  BeforeStep
	After   AfterStep
	Step    pipeline.Step
}

func (l *lifecycleStep) Run() error {
	if l.Before != nil {
		err := l.Before()

		if err != nil {
			return err
		}
	}

	err := l.Step.Run()

	if l.After != nil {
		err = l.After(l.Step, err)
	}
	return err
}

func (l *lifecycleStep) Name() string {
	return l.Step.Name()
}

func CreateBeforeStepLifecycle(step pipeline.Step, before BeforeStep) pipeline.Step {
	return &lifecycleStep{
		Before: before,
		Step:   step,
	}
}

func CreateAfterStepLifecycle(step pipeline.Step, after AfterStep) pipeline.Step {
	return &lifecycleStep{
		After:  after,
		Step:   step,
	}
}

func CreateStepLifecycle(step pipeline.Step, before BeforeStep, after AfterStep) pipeline.Step {
	return &lifecycleStep{
		Before: before,
		After:  after,
		Step:   step,
	}
}