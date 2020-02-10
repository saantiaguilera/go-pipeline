package pipeline

type sequentialGroup []Stage

func (s sequentialGroup) Run(executor Executor, ctx Context) error {
	for _, stage := range s {
		err := stage.Run(executor, ctx)

		if err != nil {
			return err
		}
	}
	return nil
}

func (s sequentialGroup) Draw(graph GraphDiagram) {
	for _, stage := range s {
		stage.Draw(graph)
	}
}

// CreateSequentialGroup creates a stage that will run each of stages sequentially. If one of them fails, the operation will abort immediately
func CreateSequentialGroup(stages ...Stage) Stage {
	var stage sequentialGroup = stages
	return &stage
}
