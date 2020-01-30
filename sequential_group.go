package pipeline

type sequentialGroup []Stage

func (s sequentialGroup) Run(executor Executor) error {
	return runSync(len(s), func(index int) error {
		return s[index].Run(executor)
	})
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
