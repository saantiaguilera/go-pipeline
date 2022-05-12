package pipeline

type (
	Statement[T any] struct {
		label string
		fn    func(T) bool
	}
)

// NewStatement creates statement represented by the given name, that will evaluate to the given evaluation
func NewStatement[T any](name string, eval func(T) bool) Statement[T] {
	return Statement[T]{
		label: name,
		fn:    eval,
	}
}

// NewAnonymousStatement createsn anonymous statement with no representation, that will evaluate to the given evaluation
func NewAnonymousStatement[T any](eval func(T) bool) Statement[T] {
	return Statement[T]{
		fn: eval,
	}
}

func (s Statement[T]) Name() string {
	return s.label
}

func (s Statement[T]) Evaluate(v T) bool {
	return s.fn != nil && s.fn(v)
}
