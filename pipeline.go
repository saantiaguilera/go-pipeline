/*
Package go-pipeline is a pure Go client library for building pipelines in a declarative way (similar to Keras or other
deep learning frameworks). It includes a high level API for easily building a graph/template of a structure to represent.

If the defined implementations of the API are insufficient, one can create their own implementation adding new behaviours,
such as circuit-breaker executors, panic recover executors, new-relic step decorators, among any idea or feature you would
like.

Supported structure

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

With them one can build a graph-like / template structure, that will be executed by a Pipeline through an Executor

Pipeline

A pipeline is a contract for executing a root stage (that will internally execute nested stages, thus evaluating the whole graph).

Executor

An executor is a contract capable of running Runnable (the single unit of work.)

This is useful if we want to add global step hooks / circuit-breakers / tracers / etc to the step's graph.
*/
package pipeline

// Pipeline contract for running a graph/template.
type Pipeline interface {
	// Run a stage graph. This method is blocking until the stage finishes.
	// Returns an error denoting that the stage couldn't complete (and its reason)
	Run(stage Stage) error
}

type pipeline struct {
	Executor Executor
}

func (p *pipeline) Run(stage Stage) error {
	return stage.Run(p.Executor)
}

// CreatePipeline creates a pipeline with a given executor
func CreatePipeline(executor Executor) Pipeline {
	return &pipeline{
		Executor: executor,
	}
}
