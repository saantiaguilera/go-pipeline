/*
Package pipeline is a pure Go client library for building pipelines in a declarative way.
It includes a high level API to easily build, execute and draw a graph/template of a desired structure.

If the defined implementations of the API are insufficient, one can create their own implementation adding new behaviours,
such as circuit-breaker executors, panic recover executors, new-relic step decorators, among any idea or feature you would
like.

Supported structure

Below you can find the atomic types this API exposes. This elements are mandatory for creating and executing a graph in a pipeline.

Step

A step is a single unit of work. It's an alias for Runnable

	type createUserStep struct {
		Service      user.Service
		InputStream  chan *model.AuthenticatedUser
		OutputStream chan *model.User
	}

	func (s *createUserStep) Name() string {
		return "create_user_step"
	}

	func (s *createUserStep) Run() error {
		defer close(s.OutputStream)
		user, err := s.Service.Create(<-s.InputStream)
		if err == nil {
			s.OutputStream <- user
		}
		return err
	}

	func CreateUserStep(service user.Service, input chan *model.AuthenticatedUser) pipeline.Step {
		return &createUserStep{
			Service:      service,
			InputStream:  input,
			OutputStream: make(chan *model.User, 1),
		}
	}

Stage

A stage contains a collection of steps. The collection will be executed according to the stage implementation (eg. concurrently, sequentially, condition-based, etc).

To create one of the already defined stages, we can simply invoke its constructor function. For example, for a sequential stage:

	stage := pipeline.CreateSequentialStage(
		user.CreateGetAuthenticatedUserStep(request, userChan),
		user.CreateUserStep(userService, userChan),
	)

Since a stage is nothing more than an interface, you can create your own implementations abiding that contract.

Stage group

A stage group is a collection of stages. It's also a Stage itself, thus a group can be used in the same scope as a stage.

	stage := pipeline.CreateSequentialGroup(
		pipeline.CreateConcurrentGroup(
			pipeline.CreateSequentialStage(
				CreateBoilEggsStep(numberOfEggs, eggsChan),
				CreateCutEggsStep(eggsChan),
			),
			pipeline.CreateSequentialStage(
				CreateWashCarrotsStep(numberOfCarrots, carrotsChan),
				CreateCutCarrotsStep(carrotsChan),
			),
		),
		pipeline.CreateSequentialStage(
			CreateMakeSaladStep(carrotsChan, eggsChan, saladChan),
		),
	)

Again, as long as you abide the Stage contract, you can create your own ones.

Defined structures

We already implement out of the box some structures that are pretty much mandatory. You can make your own custom
implementation to create behaviours we are not currently defining. The provided ones are:

- Concurrent: Run stages or steps concurrently

- Sequential: Run stages or steps sequentially

- Conditional: Run either a flow or another one depending on the evaluation of a statement

- Tracer: Add tracers to a stage or step

- Lifecycle: Add lifecycle methods (before / after) to a stage or step

With the use of them a graph-like / template structure can be achieved, that will be executed by a Pipeline through an Executor.

Pipeline

A pipeline is a contract for executing a root stage (that will internally execute nested stages, thus evaluating the whole graph).

Executor

An executor is a contract capable of running Runnable (the single unit of work.)

This is useful if we want to add global step hooks / circuit-breakers / tracers / etc to the step's graph.

Drawing

The package comes with a handy drawing API for representing the created graphs.

	diagram := pipeline.CreateUMLActivityGraphDiagram()
	renderer := pipeline.CreateUMLActivityRenderer(pipeline.UMLOptions{
		Type: pipeline.UMLFormatSVG,
	})
	file, _ := os.Create("outputFile.svg")

	yourStage.Draw(diagram)

	err := renderer.Render(diagram, file)
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
