package pipeline

type (
	SequentialContainer[T any] []Container[T]
)

// NewSequentialContainer creates container that will run each of the steps sequentially. If one of them fails, the operation will abort immediately
func NewSequentialContainer[T any](units ...Container[T]) SequentialContainer[T] {
	return units
}

func (s SequentialContainer[T]) Visit(ex Executor[T], in T) error {
	for _, c := range s {
		if err := c.Visit(ex, in); err != nil {
			return err
		}
	}
	return nil
}

func (s SequentialContainer[T]) Draw(graph GraphDiagram) {
	for _, c := range s {
		c.Draw(graph)
	}
}
