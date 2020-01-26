package pkg

// Executor interface for running units (usually in the form of Step).
// The executor should be used as a mean for decorating or changing the behaviour of a command execution
// Such as adding tracing capabilities, logging, circuit-breakers, timeouts, fallbacks, etc.
type Executor interface {

	// Run a runnable. Returns an error in case the runnable (or the execution itself) fails.
	Run(runnable Runnable) error
}
