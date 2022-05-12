package pipeline

type (
	SequentialStage[T any] []Step[T]
)

// NewSequentialStage News a stage that will run each of the steps sequentially. If one of them fails, the operation will abort immediately
func NewSequentialStage[T any](steps ...Step[T]) SequentialStage[T] {
	return steps
}

func (s SequentialStage[T]) Run(executor Executor[T], in T) error {
	for _, step := range s {
		err := executor.Run(step, in)

		if err != nil {
			return err
		}
	}
	return nil
}

func (s SequentialStage[T]) Draw(graph GraphDiagram) {
	for _, step := range s {
		graph.AddActivity(step.Name())
	}
}
