package paint_example_test

import (
	"fmt"

	"github.com/saantiaguilera/go-pipeline"
)

type PaintSurfaceStep struct{}

func (s *PaintSurfaceStep) Run(ctx pipeline.Context) error {
	surface, _ := ctx.Get(TagSurface)
	fmt.Printf("Painting %d surface\n", surface)
	return nil
}

func (s *PaintSurfaceStep) Name() string {
	return "paint surface painting"
}
