package pipeline

import "context"

type (
	Statement[T any] struct {
		label string
		fn    func(context.Context, T) bool
	}
)

// NewStatement creates a statement represented by the given name, that will evaluate to the given evaluation
func NewStatement[T any](name string, eval func(context.Context, T) bool) Statement[T] {
	return Statement[T]{
		label: name,
		fn:    eval,
	}
}

// NewAnonymousStatement creates an anonymous statement with no representation, that will evaluate to the given evaluation
func NewAnonymousStatement[T any](eval func(context.Context, T) bool) Statement[T] {
	return NewStatement("", eval)
}

func (s Statement[T]) Name() string {
	return s.label
}

func (s Statement[T]) Evaluate(ctx context.Context, v T) bool {
	return s.fn != nil && s.fn(ctx, v)
}
