package main

import (
	"fmt"

	"github.com/saantiaguilera/go-pipeline"
)

type calculateVolumeStep struct{}

func (s *calculateVolumeStep) Run(ctx pipeline.Context) error {
	width, _ := ctx.GetInt(tagWidth)
	height, _ := ctx.GetInt(tagHeight)
	depth, _ := ctx.GetInt(tagDepth)
	volume := width * height * depth
	fmt.Printf("Getting volume: %d\n", volume)

	ctx.Set(tagVolume, volume)
	return nil
}

func (s *calculateVolumeStep) Name() string {
	return "get volume"
}
