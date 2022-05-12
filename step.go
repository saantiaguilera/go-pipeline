package pipeline

import "github.com/google/uuid"

type (
	// Step interface for making a unit work.
	Step[T any] struct {
		id   string
		name string
		run  func(T) error
	}

	// Container is a grouping of units (steps / other containers / etc) allowing one to create a workflow/template/graph of a given
	// problem.
	// A container can be run with a given executor
	Container[T any] interface {
		DrawableDiagram

		// Visit a container with a given executor. Returns an error if this container fails to complete.
		// A T is provided as a mean of communication between different containers and units of work
		Visit(Executor[T], T) error
	}
)

// NewStep createsn immutable stateless unit of work based on a function that matches the Runnable contract.
// You can use this implementation when your use-cases will be completely stateless (they don't rely on a service
// or anything that can be injected at the start and stay immutable for the lifetime of the process)
func NewStep[T any](name string, run func(T) error) Step[T] {
	return Step[T]{
		id:   uuid.New().String(),
		name: name,
		run:  run,
	}
}

func (s Step[T]) ID() string {
	return s.id
}

func (s Step[T]) Name() string {
	return s.name
}

func (s Step[T]) Run(in T) error {
	return s.run(in)
}

func (s Step[T]) Visit(ex Executor[T], in T) error {
	return ex.Run(s, in)
}

func (s Step[T]) Draw(graph GraphDiagram) {
	graph.AddActivity(s.Name())
}
