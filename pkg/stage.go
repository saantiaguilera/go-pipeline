package pkg

// Stage is a grouping of units (steps / stages / etc) allowing one to create a workflow/template/graph of a given
// problem.
// A stage can be run with a given executor
type Stage interface {

	// Run a stage with a given executor. Returns an error if this stage fails to complete.
	Run(executor Executor) error
}
