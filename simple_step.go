package pipeline

// Simple step structure. A simple step is a stateless unit of work (just a function to run).
type simpleStep[T any] struct {
	name string
	run  func(T) error
}

func (s *simpleStep[T]) Name() string {
	return s.name
}

func (s *simpleStep[T]) Run(in T) error {
	return s.run(in)
}

// NewSimpleStep News an immutable stateless unit of work based on a function that matches the Runnable contract.
// You can use this implementation when your use-cases will be completely stateless (they don't rely on a service
// or anything that can be injected at the start and stay immutable for the lifetime of the process)
func NewSimpleStep[T any](name string, run func(T) error) Step[T] {
	return &simpleStep[T]{
		name: name,
		run:  run,
	}
}
