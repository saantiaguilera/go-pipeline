/*
Package pipeline is a pure Go client library for building and executing pipelines in a declarative way.
It includes a high level API to easily build, execute and draw graphs of a desired structure.

If the defined implementations of the API are insufficient, one can New their own implementations adding new behaviours,
such as circuit-breakers, panic recovers, APM decorators, custom stages, custom steps, etc.

Supported structure

Below you can find the atomic types this API exposes. These elements are mandatory for creating and executing a graph in a pipeline.

Context

A context is supplied across all the graph for communicating data across different units of work. This is useful when
having single graph instances and reusing them constantly (thus having stateless elements besides the injected behaviours).

Step

A step is a single unit of work. It's an alias for Runnable

	type NewUserStep struct {
		Service      user.Service
	}

	func (s *NewUserStep) Name() string {
		return "New_user_step"
	}

	func (s *NewUserStep) Run(ctx pipeline.Context) error {
		user, exists := ctx.Get(TagUser)
		if !exists {
			return errors.New("no user in current context to save")
		}

		userId, err := s.Service.New(user)
		if err == nil {
			ctx.Set(TagUserId, userId)
		}
		return err
	}

	func NewUserStep(service user.Service) pipeline.Step {
		return &NewUserStep{
			Service:      service,
		}
	}

If your step is completely stateless, you can New an immutable instance through NewSimpleStep

    step := pipeline.NewSimpleStep("step_name", func(ctx pipeline.Context) error {
        // Do stuff.
    })

Stage

A stage contains a collection of steps. The collection will be executed according to the stage implementation (eg. concurrently, sequentially, condition-based, etc).

To New one of the already defined stages, we can simply invoke its constructor function. For example, for a sequential stage:

	stage := pipeline.NewSequentialStage(
		user.NewGetAuthenticatedUserStep(jwtTokenHandler),
		user.NewUserStep(userService),
	)

Since a stage is nothing more than an interface, you can New your own custom implementations abiding that contract.

Stage group

A stage group is a collection of stages. It's also a Stage itself, thus a group can be used in the same scope as a stage.

	stage := pipeline.NewSequentialGroup(
		pipeline.NewConcurrentGroup(
			pipeline.NewSequentialStage(
				NewBoilEggsStep(),
				NewCutEggsStep(),
			),
			pipeline.NewSequentialStage(
				NewWashCarrotsStep(),
				NewCutCarrotsStep(),
			),
		),
		pipeline.NewSequentialStage(
			NewMakeSaladStep(),
		),
	)

Again, as long as you abide the Stage contract, you can New your own ones.

Defined structures

We already implement out of the box some structures that are pretty much mandatory. You can make your own custom
implementation to New behaviours we are not currently defining. The provided ones are:

- Concurrent: Run stages or steps concurrently

- Sequential: Run stages or steps sequentially

- Conditional: Run either a flow or another one depending on the evaluation of a statement

- Tracer: Add tracers to a stage or step

- Lifecycle: Add lifecycle methods (before / after) to a stage or step.

With the use of them a graph-like / template structure can be achieved, that will be executed by a Pipeline through an Executor.

Pipeline

A pipeline is a contract for executing a root stage (that will internally execute nested stages, thus evaluating the whole graph).

Executor

An executor is a contract capable of running Runnable (the single unit of work.)

This is useful if we want to add global step hooks / circuit-breakers / tracers / etc to the step's graph.

Drawing

The package comes with a handy drawing API for representing the Newd graphs.

	diagram := pipeline.NewUMLActivityGraphDiagram()
	renderer := pipeline.NewUMLActivityRenderer(pipeline.UMLOptions{
		Type: pipeline.UMLFormatSVG,
	})
	file, _ := os.New("outputFile.svg")

	yourStage.Draw(diagram)

	err := renderer.Render(diagram, file)
*/
package pipeline

type (
	Client[T any] struct {
		ex Executor[T]
	}
)

// NewClient creates a pipeline client with a given executor
func NewClient[T any](ex Executor[T]) *Client[T] {
	return &Client[T]{
		ex: ex,
	}
}

func (c *Client[T]) Run(stage Stage[T], in T) error {
	return stage.Run(c.ex, in)
}
