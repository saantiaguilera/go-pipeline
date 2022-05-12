package pipeline

// Executor interface for running a Step.
// The executor should be used as a mean for decorating or changing the behaviour of a command execution
// Such as adding tracing capabilities, logging, circuit-breakers, timeouts, fallbacks, etc.
type Executor[T any] interface {
	// Run a runnable with a given context. Returns an error in case the runnable (or the execution itself) fails.
	Run(Step[T], T) error
}
