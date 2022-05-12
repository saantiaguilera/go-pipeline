package pipeline

// Step interface for making a unit work.
type Step[T any] interface {

	// Named because a Step is named
	Named

	// Run the unit, returns error if it fails to complete successfully
	// The provided input is the same across all units, so it's useful to store and retrieve data as a mean of
	// communication between different units
	Run(in T) error
}
