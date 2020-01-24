package pipeline_step

import "github.com/saantiaguilera/go-pipeline"

type LifecycleStep struct {
	Before  func() error
	After   func(step pipeline.Step, err error) error
	Step    pipeline.Step
}

func (l *LifecycleStep) Run() error {
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

func (l *LifecycleStep) Name() string {
	return l.Step.Name()
}
