package pipeline

type sequentialStage []Step

func (s sequentialStage) Run(executor Executor) error {
	return runSync(len(s), func(index int) error {
		return executor.Run(s[index])
	})
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
