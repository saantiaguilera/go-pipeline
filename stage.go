package pipeline

// Stage is a grouping of units (steps / stages / etc) allowing one to New a workflow/template/graph of a given
// problem.
// A stage can be run with a given executor
type Stage[T any] interface {
	DrawableDiagram

	// Run a stage with a given executor. Returns an error if this stage fails to complete.
	// A context is provided as a mean of communication between different stages and units of work
	Run(Executor[T], T) error
}
