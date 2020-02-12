package main

import (
	"fmt"
	"time"

	"github.com/saantiaguilera/go-pipeline"
)

type getPriceToPaintSurfaceStep struct{}

func (s *getPriceToPaintSurfaceStep) Run(ctx pipeline.Context) error {
	surface, _ := ctx.GetInt(tagSurface)
	price := float64(surface) * 348.1923

	time.Sleep(100 * time.Millisecond) // Simulate time it takes to get this from a service
	fmt.Printf("Getting price for painting surface: %f\n", price)

	ctx.Set(tagSurfacePrice, price)
	return nil
}

func (s *getPriceToPaintSurfaceStep) Name() string {
	return "get price for painting surface"
}
