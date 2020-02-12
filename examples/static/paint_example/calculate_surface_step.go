package main

import (
	"fmt"

	"github.com/saantiaguilera/go-pipeline"
)

type calculateSurfaceStep struct{}

func (s *calculateSurfaceStep) Run(ctx pipeline.Context) error {
	width, _ := ctx.GetInt(tagWidth)
	height, _ := ctx.GetInt(tagHeight)
	surface := width * height
	fmt.Printf("Getting surface: %d\n", surface)

	ctx.Set(tagSurface, surface)
	return nil
}

func (s *calculateSurfaceStep) Name() string {
	return "get surface"
}
