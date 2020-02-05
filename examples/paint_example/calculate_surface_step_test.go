package paint_example_test

import (
	"fmt"

	"github.com/saantiaguilera/go-pipeline"
)

type CalculateSurfaceStep struct{}

func (s *CalculateSurfaceStep) Run(ctx pipeline.Context) error {
	width, _ := ctx.GetInt(TagWidth)
	height, _ := ctx.GetInt(TagHeight)
	surface := width * height
	fmt.Printf("Getting surface: %d\n", surface)

	ctx.Set(TagSurface, surface)
	return nil
}

func (s *CalculateSurfaceStep) Name() string {
	return "get surface"
}
