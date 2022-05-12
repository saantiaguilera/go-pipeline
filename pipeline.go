/*
Package pipeline is a pure Go client library for building and executing pipelines in a declarative way.
It includes a high level API to easily build, execute and draw graphs of a desired structure.

If the defined implementations of the API are insufficient, one can New their own implementations adding new behaviours,
such as circuit-breakers, panic recovers, APM decorators, custom containers, custom steps, etc.

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

If your step is completely stateless, you can New an immutable instance through NewStep

    step := pipeline.NewStep("step_name", func(ctx pipeline.Context) error {
        // Do stuff.
    })

Container

A container contains a collection of steps. The collection will be executed according to the container implementation (eg. concurrently, sequentially, condition-based, etc).

To New one of the already defined containers, we can simply invoke its constructor function. For example, for a sequential container:

	container := pipeline.NewSequentialContainer(
		user.NewGetAuthenticatedUserStep(jwtTokenHandler),
		user.NewUserStep(userService),
	)

Since a container is nothing more than an interface, you can New your own custom implementations abiding that contract.

Container group

A container group is a collection of containers. It's also a Container itself, thus a group can be used in the same scope as a container.

	container := pipeline.NewSequentialGroup(
		pipeline.NewConcurrentGroup(
			pipeline.NewSequentialContainer(
				NewBoilEggsStep(),
				NewCutEggsStep(),
			),
			pipeline.NewSequentialContainer(
				NewWashCarrotsStep(),
				NewCutCarrotsStep(),
			),
		),
		pipeline.NewSequentialContainer(
			NewMakeSaladStep(),
		),
	)

Again, as long as you abide the Container contract, you can New your own ones.

Defined structures

We already implement out of the box some structures that are pretty much mandatory. You can make your own custom
implementation to New behaviours we are not currently defining. The provided ones are:

- Concurrent: Run containers or steps concurrently

- Sequential: Run containers or steps sequentially

- Conditional: Run either a flow or another one depending on the evaluation of a statement

- Tracer: Add tracers to a container or step

- Lifecycle: Add lifecycle methods (before / after) to a container or step.

With the use of them a graph-like / template structure can be achieved, that will be executed by a Pipeline through an Executor.

Pipeline

A pipeline is a contract for executing a root container (that will internally execute nested containers, thus evaluating the whole graph).

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

	yourContainer.Draw(diagram)

	err := renderer.Render(diagram, file)
*/
package pipeline

import "context"

type (
	Client[T any] struct {
		ex Executor[T]
	}
)

// NewClient creates a pipeline client with a given executor
func NewClient[T any](ex Executor[T]) Client[T] {
	return Client[T]{
		ex: ex,
	}
}

func (c Client[T]) Run(ctx context.Context, cn Container[T], in T) error {
	return cn.Visit(ctx, c.ex, in)
}
