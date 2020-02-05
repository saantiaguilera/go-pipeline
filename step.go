package pipeline

// Runnable interface for making a unit work.
type Runnable interface {

	// Named because a runnable is named
	Named

	// Run the unit, returns error if it fails to complete successfully
	// The provided context is the same across all units, so it's useful to store and retrieve data as a mean of
	// communication between different units
	Run(ctx Context) error
}

// Step is an Alias for runnable.
type Step Runnable
