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

// Graph creates static workflow for this sample. It's all in a single func completely coupled for showing purposes
// you should probably decouple this into more atomic ones (eg. a func for calculating sizes that returns a Container)
func Graph() pipeline.Container {
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

	return pipeline.NewSequentialGroup(
		pipeline.NewTracedContainer("measurement_container", pipeline.NewConcurrentContainer(
			widthStep,
			heightStep,
			depthStep,
		)),
		pipeline.NewConcurrentContainer(
			calculateVolumeStep,
			calculateSurfaceStep,
		),
		pipeline.NewConcurrentGroup(
			pipeline.NewSequentialContainer(
				calculatePriceToPaintSurfaceStep,
				recordPriceSurfaceStep,
			),
			pipeline.NewSequentialContainer(
				calculatePriceToPaintVolumeStep,
				recordPriceVolumeStep,
			),
		),
		pipeline.NewSequentialContainer(
			evaluateStep,
		),
		pipeline.NewConditionalGroup(
			pipeline.NewStatement("should_paint", func(ctx pipeline.Context) bool {
				volumePrice, _ := ctx.GetFloat64(tagVolumePrice)
				surfacePrice, _ := ctx.GetFloat64(tagSurfacePrice)
				return volumePrice+surfacePrice < 100000
			}),
			pipeline.NewSequentialGroup(
				pipeline.NewConcurrentContainer(
					acceptSurfacePaintingStep,
					acceptVolumePaintingStep,
				),
				pipeline.NewConcurrentContainer(
					pipeline.NewTracedStep(paintSurfaceStep),
					pipeline.NewTracedStep(paintVolumeStep),
				),
			),
			pipeline.NewSequentialContainer(),
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
		diagram := pipeline.NewUMLActivityGraphDiagram()
		renderer := pipeline.NewUMLActivityRenderer(pipeline.UMLOptions{
			Type: pipeline.UMLFormatSVG,
		})
		file, _ := os.New("template.svg")

		Graph().Draw(diagram)

		err := renderer.Render(diagram, file)

		if err != nil {
			panic(err)
		}
	}
}

// RunPipeline runs the provided pipeline.
func RunPipeline() {
	// New initial data, this can be Newd once and reused multiple times
	pipe := pipeline.NewPipeline(&sampleExecutor{})
	graph := Graph()

	// New context to be used on a single graph evaluation. Feed inputs for the starting steps here.
	ctx := pipeline.NewContext()

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
