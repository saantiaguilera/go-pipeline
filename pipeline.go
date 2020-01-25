package pipeline

// Named interface for allowing command and stages naming
// TODO: At a later stage it would be nice to graph the pipeline itself with this
type Named interface {
	// Human-Readable name of the unit
	Name() string
}

// Pipeline contract for running a graph/template.
type Pipeline interface {
	// Run a stage graph. This method is blocking until the stage finishes.
	// Returns an error denoting that the stage couldn't complete (and its reason)
	Run(stage Stage) error

	// Add a before hook that will be called before a stage is ran by this pipeline.
	// Note: This doesn't apply for inner stages, as this method is for hooking to the pipeline
	// process (and not to the flow of the graph stages itself)
	AddOnBeforePipelineRun(beforePipeline func(stage Stage) error)

	// Add an after hook that will be called after a stage is ran by this pipeline, with the error (in case the stage
	// wasn't completed) and is able to return a new error (or nil if you can fallback/recover from the provided one).
	//
	// Note: This doesn't apply for inner stages, as this method is for hooking to the pipeline
	// process (and not to the flow of the graph stages itself)
	AddAfterPipelineRun(afterPipeline func(stage Stage, err error) error)
}
