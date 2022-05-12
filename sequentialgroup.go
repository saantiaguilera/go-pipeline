package pipeline

type (
	SequentialGroup[T any] []Stage[T]
)

// NewSequentialGroup News a stage that will run each of stages sequentially. If one of them fails, the operation will abort immediately
func NewSequentialGroup[T any](stages ...Stage[T]) SequentialGroup[T] {
	return stages
}

func (s SequentialGroup[T]) Run(executor Executor[T], in T) error {
	for _, stage := range s {
		err := stage.Run(executor, in)

		if err != nil {
			return err
		}
	}
	return nil
}

func (s SequentialGroup[T]) Draw(graph GraphDiagram) {
	for _, stage := range s {
		stage.Draw(graph)
	}
}