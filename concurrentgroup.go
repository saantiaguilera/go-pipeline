package pipeline

type (
	ConcurrentGroup[T any] []Stage[T]
)

// NewConcurrentGroup News a stage that will run each of the stages concurrently.
// The stage will wait for all of the stages to finish before returning.
//
// If one of them fails, the stage will wait until everyone finishes and after that return the error.
// If more than one fails, then the error will be the one delivered by the last failure.
func NewConcurrentGroup[T any](stages ...Stage[T]) ConcurrentGroup[T] {
	return stages
}

func (s ConcurrentGroup[T]) NewStageGraphActivity(drawable DrawableDiagram) DrawDiagram {
	return func(graph GraphDiagram) {
		drawable.Draw(graph)
	}
}

func (s ConcurrentGroup[T]) Draw(graph GraphDiagram) {
	if len(s) > 0 {
		var forkStages []DrawDiagram
		for _, stage := range s {
			forkStages = append(forkStages, s.NewStageGraphActivity(stage))
		}

		graph.AddConcurrency(forkStages...)
	}
}

func (s ConcurrentGroup[T]) Run(executor Executor[T], in T) error {
	return spawnAsync(len(s), func(index int) error {
		return s[index].Run(executor, in)
	})
}
