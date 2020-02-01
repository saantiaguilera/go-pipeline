package pipeline

type concurrentStage []Step

func (s concurrentStage) createStepGraphActivity(name string) DrawDiagram {
	return func(graph GraphDiagram) {
		graph.AddActivity(name)
	}
}

func (s concurrentStage) Draw(graph GraphDiagram) {
	if len(s) > 0 {
		var forkSteps []DrawDiagram
		for _, step := range s {
			forkSteps = append(forkSteps, s.createStepGraphActivity(step.Name()))
		}

		graph.AddConcurrency(forkSteps...)
	}
}

func (s concurrentStage) Run(executor Executor) error {
	return spawnAsync(len(s), func(index int) error {
		return executor.Run(s[index])
	})
}

// CreateConcurrentStage creates a stage that will run each of the steps concurrently.
// The stage will wait for all of the steps to finish before returning.
//
// If one of them fails, the stage will wait until everyone finishes and after that return the error.
// If more than one fails, then the error will be the one delivered by the last failure.
func CreateConcurrentStage(steps ...Step) Stage {
	var stage concurrentStage = steps
	return &stage
}
