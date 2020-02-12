package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/saantiaguilera/go-pipeline"
)

/*
GetHeight

GetWidth

GetDepth
--------
CalculateVolume

CalculateSurface
--------
GetPriceToPaintSurface
RecordPrice

GetPriceToPaintVolume
RecordPrice
---------
SendToEvaluation
---------
if paint
  AcceptSurfacePainting

  AcceptVolumePainting
  --------
  PaintSurface

  PaintVolume
else
  nothing
*/

var render = flag.Bool("pipeline.render", false, "render pipeline")

// Graph creates a static workflow for this sample. It's all in a single func completely coupled for showing purposes
// you should probably decouple this into more atomic ones (eg. a func for calculating sizes that returns a Stage)
func Graph() pipeline.Stage {
	widthStep := &getWidthStep{}
	heightStep := &getHeightStep{}
	depthStep := &getDepthStep{}

	calculateVolumeStep := &calculateVolumeStep{}
	calculateSurfaceStep := &calculateSurfaceStep{}

	calculatePriceToPaintSurfaceStep := &getPriceToPaintSurfaceStep{}
	calculatePriceToPaintVolumeStep := &getPriceToPaintVolumeStep{}

	recordPriceSurfaceStep := &recordPriceStep{
		PriceType: tagSurfacePrice,
	}
	recordPriceVolumeStep := &recordPriceStep{
		PriceType: tagVolumePrice,
	}

	evaluateStep := &evaluateStep{}

	acceptSurfacePaintingStep := &acceptSurfacePaintingStep{}
	acceptVolumePaintingStep := &acceptVolumePaintingStep{}

	paintSurfaceStep := &paintSurfaceStep{}
	paintVolumeStep := &paintVolumeStep{}

	return pipeline.CreateSequentialGroup(
		pipeline.CreateTracedStage("measurement_stage", pipeline.CreateConcurrentStage(
			widthStep,
			heightStep,
			depthStep,
		)),
		pipeline.CreateConcurrentStage(
			calculateVolumeStep,
			calculateSurfaceStep,
		),
		pipeline.CreateConcurrentGroup(
			pipeline.CreateSequentialStage(
				calculatePriceToPaintSurfaceStep,
				recordPriceSurfaceStep,
			),
			pipeline.CreateSequentialStage(
				calculatePriceToPaintVolumeStep,
				recordPriceVolumeStep,
			),
		),
		pipeline.CreateSequentialStage(
			evaluateStep,
		),
		pipeline.CreateConditionalGroup(
			pipeline.CreateSimpleStatement("should_paint", func(ctx pipeline.Context) bool {
				volumePrice, _ := ctx.GetFloat64(tagVolumePrice)
				surfacePrice, _ := ctx.GetFloat64(tagSurfacePrice)
				return volumePrice+surfacePrice < 100000
			}),
			pipeline.CreateSequentialGroup(
				pipeline.CreateConcurrentStage(
					acceptSurfacePaintingStep,
					acceptVolumePaintingStep,
				),
				pipeline.CreateConcurrentStage(
					pipeline.CreateTracedStep(paintSurfaceStep),
					pipeline.CreateTracedStep(paintVolumeStep),
				),
			),
			pipeline.CreateSequentialStage(),
		),
	)
}

// You could have your own executor using hystrix or whatever.
// Decorate it with tracers / circuit-breakers / loggers / new-relic / etc.
type sampleExecutor struct{}

func (t *sampleExecutor) Run(cmd pipeline.Runnable, ctx pipeline.Context) error {
	fmt.Printf("Running task %s\n", cmd.Name())
	return cmd.Run(ctx)
}

// RunGraphRendering represents the graph in UML Activity and renders it as an SVG file (template.svg)
func RunGraphRendering() {
	if *render {
		diagram := pipeline.CreateUMLActivityGraphDiagram()
		renderer := pipeline.CreateUMLActivityRenderer(pipeline.UMLOptions{
			Type: pipeline.UMLFormatSVG,
		})
		file, _ := os.Create("template.svg")

		Graph().Draw(diagram)

		err := renderer.Render(diagram, file)

		if err != nil {
			panic(err)
		}
	}
}

// RunPipeline runs the provided pipeline.
func RunPipeline() {
	// Create initial data, this can be created once and reused multiple times
	pipe := pipeline.CreatePipeline(&sampleExecutor{})
	graph := Graph()

	// Create context to be used on a single graph evaluation. Feed inputs for the starting steps here.
	ctx := pipeline.CreateContext()

	// Evaluate graph and assert errors.
	err := pipe.Run(graph, ctx)

	if err != nil {
		panic(err)
	}
}

func main() {
	RunGraphRendering()
	RunPipeline()
}
