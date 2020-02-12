package main

import (
	"fmt"
	"time"

	"github.com/saantiaguilera/go-pipeline"
)

type getPriceToPaintVolumeStep struct{}

func (s *getPriceToPaintVolumeStep) Run(ctx pipeline.Context) error {
	volume, _ := ctx.GetInt(tagVolume)
	price := float64(volume) * 348.1923

	time.Sleep(100 * time.Millisecond) // Simulate time it takes to get this from a service
	fmt.Printf("Getting price for painting volume: %f\n", price)

	ctx.Set(tagVolumePrice, price)
	return nil
}

func (s *getPriceToPaintVolumeStep) Name() string {
	return "get price for painting volume"
}
