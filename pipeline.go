/*
Package go-pipeline is a pure Go client library for building pipelines in a declarative way (similar to Keras or other
deep learning frameworks). It includes a high level API for easily building a graph/template of a structure to represent.

If the defined implementations of the API are insufficient, one can create their own implementation adding new behaviours,
such as circuit-breaker executors, panic recover executors, new-relic step decorators, among any idea or feature you would
like.

### Supported structure

The high level structures defined by the API are:
- Step: Single unit of work. Alias for Runnable
- Stage: Collection of steps
- Stage group: Collection of stages (also a Stage itself)

The defined structures can be implemented to create several behaviours, such as:
- Concurrent: Run stages or steps concurrently
- Sequential: Run stages or steps sequentially
- Conditional: Run either a flow or another one depending on the evaluation of a statement
- Tracer: Add tracers to a stage or step
- Lifecycle: Add lifecycle methods (before / after) to a stage or step

With them one can build a graph-like / template structure, that will be executed by a `Pipeline` through an `Executor`

### Pipeline

A pipeline is a contract for executing a root stage (that will internally execute nested stages, thus evaluating the whole graph).

If needed, one can add before / after execution hooks to enrich it (eg. Decorating requests pre execution, recovering errors, tracing, etc)

### Executor

An executor is a contract capable of running `Runnable` (the single unit of work.)

This is useful if we want to add global step hooks / circuit-breakers / tracers / etc to the step's graph.
*/
package pipeline

// Lifecycle is a contract for attaching hooks to the lifecycle of a pipeline execution
type Lifecycle interface {
	// AddBeforeRunHook adds a before hook that will be called before a stage is ran by this pipeline.
	// Note: This doesn't apply for inner stages, as this method is for hooking to the pipeline
	// process (and not to the flow of the graph stages itself)
	AddBeforeRunHook(beforePipeline func(stage Stage) error)

	// AddAfterRunHook adds an after hook that will be called after a stage is ran by this pipeline, with the error (in case the stage
	// wasn't completed) and is able to return a new error (or nil if you can fallback/recover from the provided one).
	//
	// Note: This doesn't apply for inner stages, as this method is for hooking to the pipeline
	// process (and not to the flow of the graph stages itself)
	AddAfterRunHook(afterPipeline func(stage Stage, err error) error)
}

// Pipeline contract for running a graph/template.
type Pipeline interface {
	Lifecycle

	// Run a stage graph. This method is blocking until the stage finishes.
	// Returns an error denoting that the stage couldn't complete (and its reason)
	Run(stage Stage) error
}

type pipeline struct {
	Before   []func(stage Stage) error
	After    []func(stage Stage, err error) error
	Executor Executor
}

func (p *pipeline) Run(stage Stage) error {
	for _, before := range p.Before {
		err := before(stage)

		if err != nil {
			return err
		}
	}

	err := stage.Run(p.Executor)

	for _, after := range p.After {
		err = after(stage, err)
	}

	return err
}

func (p *pipeline) AddBeforeRunHook(beforePipeline func(stage Stage) error) {
	p.Before = append(p.Before, beforePipeline)
}

func (p *pipeline) AddAfterRunHook(afterPipeline func(stage Stage, err error) error) {
	p.After = append(p.After, afterPipeline)
}

// CreatePipeline creates a pipeline with a given executor
func CreatePipeline(executor Executor) Pipeline {
	return &pipeline{
		Executor: executor,
		Before:   []func(stage Stage) error{},
		After:    []func(stage Stage, err error) error{},
	}
}
