package pipeline

type Named interface {
	Name() string
}

type Pipeline struct {
	Before   func() error
	After    func(stage Stage, err error) error
	Executor Executor
}

func (p *Pipeline) Run(stage Stage) error {
	if p.Before != nil {
		err := p.Before()

		if err != nil {
			return err
		}
	}

	err := stage.Run(p.Executor)

	if p.After != nil {
		err = p.After(stage, err)
	}

	return err
}