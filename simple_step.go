package pipeline

// Simple step structure. A simple step is a stateless unit of work (just a function to run).
type simpleStep struct {
	name string
	run  func(ctx Context) error
}

func (s *simpleStep) Name() string {
	return s.name
}

func (s *simpleStep) Run(ctx Context) error {
	return s.run(ctx)
}

// CreateSimpleStep creates an immutable stateless unit of work based on a function that matches the Runnable contract.
// You can use this implementation when your use-cases will be completely stateless (they don't rely on a service
// or anything that can be injected at the start and stay immutable for the lifetime of the process)
func CreateSimpleStep(name string, run func(ctx Context) error) Step {
	return &simpleStep{
		name: name,
		run:  run,
	}
}
