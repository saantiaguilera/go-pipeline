package main

import (
	"fmt"

	"github.com/saantiaguilera/go-pipeline"
)

type paintVolumeStep struct{}

func (s *paintVolumeStep) Run(ctx pipeline.Context) error {
	volume, _ := ctx.Get(tagVolume)
	fmt.Printf("Painting %d volume\n", volume)
	return nil
}

func (s *paintVolumeStep) Name() string {
	return "paint volume painting"
}
