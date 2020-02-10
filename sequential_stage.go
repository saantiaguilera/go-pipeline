package pipeline

type sequentialStage []Step

func (s sequentialStage) Run(executor Executor, ctx Context) error {
	for _, step := range s {
		err := executor.Run(step, ctx)

		if err != nil {
			return err
		}
	}
	return nil
}

func (s sequentialStage) Draw(graph GraphDiagram) {
	for _, step := range s {
		graph.AddActivity(step.Name())
	}
}

// CreateSequentialStage creates a stage that will run each of the steps sequentially. If one of them fails, the operation will abort immediately
func CreateSequentialStage(steps ...Step) Stage {
	var stage sequentialStage = steps
	return &stage
}
