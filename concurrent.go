package pipeline

import "context"

type (
	ConcurrentContainer[T any] []Container[T]
)

// NewConcurrentContainer creates a container that will run each of the units concurrently.
// The container will wait for all of the units to finish before returning.
//
// If one of them fails, the container will wait until everyone finishes and after that return the error.
// If more than one fails, then the error will be the one delivered by the last failure.
func NewConcurrentContainer[T any](units ...Container[T]) ConcurrentContainer[T] {
	return units
}

func (s ConcurrentContainer[T]) Draw(graph GraphDiagram) {
	if len(s) > 0 {
		var forkSteps []DrawDiagram
		for _, c := range s {
			forkSteps = append(forkSteps, s.newContainerGraphActivity(c))
		}

		graph.AddConcurrency(forkSteps...)
	}
}

func (s ConcurrentContainer[T]) Visit(ctx context.Context, ex Executor[T], in T) error {
	return spawnAsync(len(s), func(index int) error {
		return s[index].Visit(ctx, ex, in)
	})
}

func (s ConcurrentContainer[T]) newContainerGraphActivity(drawable DrawableDiagram) DrawDiagram {
	return func(graph GraphDiagram) {
		drawable.Draw(graph)
	}
}
