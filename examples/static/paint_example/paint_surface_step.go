package main

import (
	"fmt"

	"github.com/saantiaguilera/go-pipeline"
)

type paintSurfaceStep struct{}

func (s *paintSurfaceStep) Run(ctx pipeline.Context) error {
	surface, _ := ctx.Get(tagSurface)
	fmt.Printf("Painting %d surface\n", surface)
	return nil
}

func (s *paintSurfaceStep) Name() string {
	return "paint surface painting"
}
