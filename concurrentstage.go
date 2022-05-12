package pipeline

type (
	ConcurrentStage[T any] []Step[T]
)

// NewConcurrentStage News a stage that will run each of the steps concurrently.
// The stage will wait for all of the steps to finish before returning.
//
// If one of them fails, the stage will wait until everyone finishes and after that return the error.
// If more than one fails, then the error will be the one delivered by the last failure.
func NewConcurrentStage[T any](steps ...Step[T]) ConcurrentStage[T] {
	return steps
}

func (s ConcurrentStage[T]) NewStepGraphActivity(name string) DrawDiagram {
	return func(graph GraphDiagram) {
		graph.AddActivity(name)
	}
}

func (s ConcurrentStage[T]) Draw(graph GraphDiagram) {
	if len(s) > 0 {
		var forkSteps []DrawDiagram
		for _, step := range s {
			forkSteps = append(forkSteps, s.NewStepGraphActivity(step.Name()))
		}

		graph.AddConcurrency(forkSteps...)
	}
}

func (s ConcurrentStage[T]) Run(executor Executor[T], in T) error {
	return spawnAsync(len(s), func(index int) error {
		return executor.Run(s[index], in)
	})
}
