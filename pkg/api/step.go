package api

// Runnable interface for making a unit work.
type Runnable interface {

	// A runnable is named
	Named

	// Run the unit, returns error if it fails to complete successfully
	Run() error

}

// Step is an Alias for runnable.
type Step Runnable