package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/saantiaguilera/go-pipeline"
)

var render = flag.Bool("pipeline.render", false, "render pipeline")

// Graph creates dynamic workflow for this sample. It's all in a single func completely coupled for showing purposes
// you should probably decouple this into more atomic ones (eg. a func for making the salad that returns a Step)
func Graph(numberOfCarrots, numberOfEggs, meatSize, ovenSize int) pipeline.Step {
	// Channels as a mean for communicating input / output.
	// Use whatever you prefer here. It doesn't even have to be channels, but this is an example.
	eggsChan := make(chan int, 1)
	carrotsChan := make(chan int, 1)
	saladChan := make(chan int, 1)
	meatChan := make(chan int, 1)

	// Complete step. Its sequential because we can't serve
	// before all the others are done.
	graph := pipeline.NewSequentialGroup(
		// Concurrent step, given we are 3, we can do the salad / meat separately
		pipeline.NewConcurrentGroup(
			// This will be the salad flow. It can be done concurrently with the meat
			pipeline.NewSequentialGroup(
				// Eggs and carrots can be operated concurrently too
				pipeline.NewConcurrentGroup(
					// Sequential step for the eggs flow
					pipeline.NewSequentialStep(
						NewBoilEggsStep(numberOfEggs, eggsChan),
						NewCutEggsStep(eggsChan),
					),
					// Another sequential step for the carrots (eggs and carrots will be concurrent though!)
					pipeline.NewSequentialStep(
						NewWashCarrotsStep(numberOfCarrots, carrotsChan),
						NewCutCarrotsStep(carrotsChan),
					),
				),
				// This is sequential. When carrots and eggs are done, this will run
				pipeline.NewSequentialStep(
					NewMakeSaladStep(eggsChan, carrotsChan, saladChan),
				),
			),
			// Another sequential step for the meat (concurrently with salad)
			pipeline.NewSequentialGroup(
				// If we end up cutting the meat, we can optimize it with the oven operation
				pipeline.NewConcurrentGroup(
					// Conditional step, the meat might be too big
					pipeline.NewConditionalStep(
						pipeline.NewStatement("is_meat_too_big", NewMeatTooBigStatement(meatSize, ovenSize)),
						// True:
						NewCutMeatStep(meatSize, ovenSize, meatChan),
						// False:
						pipeline.NewStep("leave_meat_as_it_is", func(ctx pipeline.Context) error {
							meatChan <- meatSize // Simply pass the meat size
							return nil
						}),
					),
					pipeline.NewSequentialStep(
						NewTurnOnOvenStep(),
					),
				),
				pipeline.NewSequentialStep(
					NewPutMeatInOvenStep(meatChan),
				),
			),
		),
		// When everything is done. Serve
		pipeline.NewSequentialStep(
			NewServeStep(meatChan, saladChan),
		),
	)

	return graph
}

// You could have your own executor using hystrix or whatever.
// Decorate it with tracers / circuit-breakers / loggers / new-relic / etc.
type sampleExecutor struct{}

func (t *sampleExecutor) Run(cmd pipeline.Runnable, ctx pipeline.Context) error {
	fmt.Printf("Running task %s\n", cmd.Name()) // Log when the task starts running
	err := cmd.Run(ctx)
	fmt.Printf("Finished task %s\n", cmd.Name()) // Log when the task ends running
	return err
}

// RunGraphRendering represents the graph in UML Activity and renders it as an SVG file (template.svg)
func RunGraphRendering() {
	if *render {
		diagram := pipeline.NewUMLActivityGraphDiagram()
		renderer := pipeline.NewUMLActivityRenderer(pipeline.UMLOptions{
			Type: pipeline.UMLFormatSVG,
		})
		file, _ := os.New("template.svg")

		Graph(0, 0, 0, 0).Draw(diagram)

		err := renderer.Render(diagram, file)

		if err != nil {
			panic(err)
		}
	}
}

// RunPipeline runs the provided pipeline.
// Output: (one of many)
//
// Turning oven on
// Washing 8 carrots
// Adding 400 meat
// Boiling 5 eggs
// Cutting 5 eggs into 25 pieces
// Putting in the oven 500 meat
// Cutting 8 carrots into 40 pieces
// Making salad with 25 eggs and 40 carrots
// Serving 65 of salad and 500 of meat
func RunPipeline() {
	// New a pipeline with our own traced executor (no circuit breakers, nothing). Also a graph, its stateless
	// so it can be re-run as many times as we like
	pipe := pipeline.NewPipeline(&sampleExecutor{})
	graph := Graph(8, 5, 600, 500)

	// Initial context input data.
	// Since this is dynamic, we can place stuff that is completely agnostic to the graph here.
	ctx := pipeline.NewContext()

	// Run and assert.
	err := pipe.Run(graph, ctx)

	if err != nil {
		panic(err)
	}
}

func main() {
	RunGraphRendering()
	RunPipeline()
}
