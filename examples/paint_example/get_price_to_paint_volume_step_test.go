package paint_example_test

import (
	"fmt"
	"time"

	"github.com/saantiaguilera/go-pipeline"
)

type GetPriceToPaintVolumeStep struct{}

func (s *GetPriceToPaintVolumeStep) Run(ctx pipeline.Context) error {
	volume, _ := ctx.GetInt(TagVolume)
	price := float64(volume) * 348.1923

	time.Sleep(100 * time.Millisecond) // Simulate time it takes to get this from a service
	fmt.Printf("Getting price for painting volume: %f\n", price)

	ctx.Set(TagVolumePrice, price)
	return nil
}

func (s *GetPriceToPaintVolumeStep) Name() string {
	return "get price for painting volume"
}
