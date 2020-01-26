package examples_test

import (
	"fmt"
	"github.com/saantiaguilera/go-pipeline/examples/steps"
	"github.com/saantiaguilera/go-pipeline/pkg"
	"github.com/saantiaguilera/go-pipeline/pkg/pipeline"
	"github.com/saantiaguilera/go-pipeline/pkg/stage/concurrent"
	"github.com/saantiaguilera/go-pipeline/pkg/stage/conditional"
	"github.com/saantiaguilera/go-pipeline/pkg/stage/sequential"
	"github.com/saantiaguilera/go-pipeline/pkg/stage/trace"
	"github.com/saantiaguilera/go-pipeline/pkg/step/lifecycle"
	trace2 "github.com/saantiaguilera/go-pipeline/pkg/step/trace"
	"github.com/stretchr/testify/assert"
	"testing"
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

func Graph() pkg.Stage {
	// We use steps with before hooks to bind data (thus making a flow), but you can adopt any method of communication
	// between steps such as:
	// - channels (at creation you bind, a step produces and another consumes)
	// - pointers (at creation you set a same pointer to both steps, one sets data into, the other consumes)
	// - interfaces (each step knows a X interface which implements the consumer step, when called the data is applied
	// to the implementing step, which will use it later.

	// As long as you dont couple steps between each other, everything works

	widthStep := &steps.GetWidthStep{}
	heightStep := &steps.GetHeightStep{}
	depthStep := &steps.GetDepthStep{}

	calculateVolumeStep := &steps.CalculateVolumeStep{}
	calculateSurfaceStep := &steps.CalculateSurfaceStep{}

	calculatePriceToPaintSurfaceStep := &steps.GetPriceToPaintSurfaceStep{}
	calculatePriceToPaintVolumeStep := &steps.GetPriceToPaintVolumeStep{}

	recordPriceSurfaceStep := &steps.RecordPriceStep{}
	recordPriceVolumeStep := &steps.RecordPriceStep{}

	evaluateStep := &steps.EvaluateStep{}

	acceptSurfacePaintingStep := &steps.AcceptSurfacePaintingStep{}
	acceptVolumePaintingStep := &steps.AcceptVolumePaintingStep{}

	paintSurfaceStep := &steps.PaintSurfaceStep{}
	paintVolumeStep := &steps.PaintVolumeStep{}

	return sequential.CreateSequentialGroup(
		trace.CreateTracedStage("measurement_stage", concurrent.CreateConcurrentStage(
			widthStep,
			heightStep,
			depthStep,
		)),
		concurrent.CreateConcurrentStage(
			lifecycle.CreateBeforeStepLifecycle(
				calculateVolumeStep,
				func(step pkg.Step) error {
					calculateVolumeStep.Width = widthStep.Width
					calculateVolumeStep.Height = heightStep.Height
					calculateVolumeStep.Depth = depthStep.Depth
					return nil
				},
			),
			lifecycle.CreateBeforeStepLifecycle(
				calculateSurfaceStep,
				func(step pkg.Step) error {
					calculateSurfaceStep.Width = widthStep.Width
					calculateSurfaceStep.Height = heightStep.Height
					return nil
				},
			),
		),
		concurrent.CreateConcurrentGroup(
			sequential.CreateSequentialStage(
				lifecycle.CreateBeforeStepLifecycle(
					calculatePriceToPaintSurfaceStep,
					func(step pkg.Step) error {
						calculatePriceToPaintSurfaceStep.Surface = calculateSurfaceStep.Surface
						return nil
					},
				),
				lifecycle.CreateBeforeStepLifecycle(
					recordPriceSurfaceStep,
					func(step pkg.Step) error {
						recordPriceSurfaceStep.Price = calculatePriceToPaintSurfaceStep.Price
						return nil
					},
				),
			),
			sequential.CreateSequentialStage(
				lifecycle.CreateBeforeStepLifecycle(
					calculatePriceToPaintVolumeStep,
					func(step pkg.Step) error {
						calculatePriceToPaintVolumeStep.Volume = calculateVolumeStep.Volume
						return nil
					},
				),
				lifecycle.CreateBeforeStepLifecycle(
					recordPriceVolumeStep,
					func(step pkg.Step) error {
						recordPriceVolumeStep.Price = calculatePriceToPaintVolumeStep.Price
						return nil
					},
				),
			),
		),
		sequential.CreateSequentialStage(
			lifecycle.CreateBeforeStepLifecycle(
				evaluateStep,
				func(step pkg.Step) error {
					evaluateStep.SurfacePrice = calculatePriceToPaintSurfaceStep.Price
					evaluateStep.VolumePrice = calculatePriceToPaintVolumeStep.Price
					return nil
				},
			),
		),
		conditional.CreateConditionalGroup(
			func() bool {
				return evaluateStep.ShouldPaint
			},
			sequential.CreateSequentialGroup(
				concurrent.CreateConcurrentStage(
					acceptSurfacePaintingStep,
					acceptVolumePaintingStep,
				),
				concurrent.CreateConcurrentStage(
					lifecycle.CreateBeforeStepLifecycle(
						paintSurfaceStep,
						func(step pkg.Step) error {
							paintSurfaceStep.Surface = calculateSurfaceStep.Surface
							return nil
						},
					),
					trace2.CreateTracedStep(lifecycle.CreateBeforeStepLifecycle(
						paintVolumeStep,
						func(step pkg.Step) error {
							paintVolumeStep.Volume = calculateVolumeStep.Volume
							return nil
						},
					)),
				),
			),
			sequential.CreateSequentialStage(),
		),
	)
}

// You could have your own executor using hystrix or whatever.
// Decorate it with tracers / circuit-breakers / loggers / new-relic / etc.
type SampleExecutor struct{}

func (t *SampleExecutor) Run(cmd pkg.Runnable) error {
	fmt.Printf("Running task %s\n", cmd.Name())
	return cmd.Run()
}

func Test_Pipeline(t *testing.T) {
	p := pipeline.CreatePipeline(&SampleExecutor{})

	err := p.Run(Graph())

	assert.Nil(t, err)
}
