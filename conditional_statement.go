package pipeline

type Statement[T any] struct {
	label string
	fn    func(T) bool
}

func (s *Statement[T]) Name() string {
	return s.label
}

func (s *Statement[T]) Evaluate(v T) bool {
	return s.fn != nil && s.fn(v)
}

// NewSimpleStatement News a statement represented by the given name, that will evaluate to the given evaluation
func NewSimpleStatement[T any](name string, evaluation func(T) bool) *Statement[T] {
	return &Statement[T]{
		label: name,
		fn:    evaluation,
	}
}

// NewAnonymousStatement News an anonymous statement with no representation, that will evaluate to the given evaluation
func NewAnonymousStatement[T any](evaluation func(T) bool) *Statement[T] {
	return &Statement[T]{
		fn: evaluation,
	}
}
