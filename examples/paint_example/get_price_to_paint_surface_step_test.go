package paint_example_test

import (
	"fmt"
	"time"

	"github.com/saantiaguilera/go-pipeline"
)

type GetPriceToPaintSurfaceStep struct{}

func (s *GetPriceToPaintSurfaceStep) Run(ctx pipeline.Context) error {
	surface, _ := ctx.GetInt(TagSurface)
	price := float64(surface) * 348.1923

	time.Sleep(100 * time.Millisecond) // Simulate time it takes to get this from a service
	fmt.Printf("Getting price for painting surface: %f\n", price)

	ctx.Set(TagSurfacePrice, price)
	return nil
}

func (s *GetPriceToPaintSurfaceStep) Name() string {
	return "get price for painting surface"
}
