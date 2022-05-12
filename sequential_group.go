package pipeline

type sequentialGroup[T any] []Stage[T]

func (s sequentialGroup[T]) Run(executor Executor[T], in T) error {
	for _, stage := range s {
		err := stage.Run(executor, in)

		if err != nil {
			return err
		}
	}
	return nil
}

func (s sequentialGroup[T]) Draw(graph GraphDiagram) {
	for _, stage := range s {
		stage.Draw(graph)
	}
}

// NewSequentialGroup News a stage that will run each of stages sequentially. If one of them fails, the operation will abort immediately
func NewSequentialGroup[T any](stages ...Stage[T]) Stage[T] {
	var stage sequentialGroup[T] = stages
	return &stage
}
