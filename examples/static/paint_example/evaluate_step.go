package main

import (
	"fmt"
	"time"

	"github.com/saantiaguilera/go-pipeline"
)

type evaluateStep struct{}

func (s *evaluateStep) Run(ctx pipeline.Context) error {
	volumePrice, _ := ctx.GetFloat64(tagVolumePrice)
	surfacePrice, _ := ctx.GetFloat64(tagSurfacePrice)
	time.Sleep(100 * time.Millisecond) // Simulate time it takes to use a service
	fmt.Printf("Evaluating prices: %f and %f\n", volumePrice, surfacePrice)
	return nil
}

func (s *evaluateStep) Name() string {
	return "evaluate prices"
}
