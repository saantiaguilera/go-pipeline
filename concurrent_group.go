package pipeline

type concurrentGroup []Stage

func (s concurrentGroup) createStageGraphActivity(drawable DrawableDiagram) DrawDiagram {
	return func(graph GraphDiagram) {
		drawable.Draw(graph)
	}
}

func (s concurrentGroup) Draw(graph GraphDiagram) {
	if len(s) > 0 {
		var forkStages []DrawDiagram
		for _, stage := range s {
			forkStages = append(forkStages, s.createStageGraphActivity(stage))
		}

		graph.AddConcurrency(forkStages...)
	}
}

func (s concurrentGroup) Run(executor Executor, ctx Context) error {
	return spawnAsync(len(s), func(index int) error {
		return s[index].Run(executor, ctx)
	})
}

// CreateConcurrentGroup creates a stage that will run each of the stages concurrently.
// The stage will wait for all of the stages to finish before returning.
//
// If one of them fails, the stage will wait until everyone finishes and after that return the error.
// If more than one fails, then the error will be the one delivered by the last failure.
func CreateConcurrentGroup(stages ...Stage) Stage {
	var stage concurrentGroup = stages
	return &stage
}
