package paint_example_test

import (
	"fmt"

	"github.com/saantiaguilera/go-pipeline"
)

type CalculateVolumeStep struct{}

func (s *CalculateVolumeStep) Run(ctx pipeline.Context) error {
	width, _ := ctx.GetInt(TagWidth)
	height, _ := ctx.GetInt(TagHeight)
	depth, _ := ctx.GetInt(TagDepth)
	volume := width * height * depth
	fmt.Printf("Getting volume: %d\n", volume)

	ctx.Set(TagVolume, volume)
	return nil
}

func (s *CalculateVolumeStep) Name() string {
	return "get volume"
}
