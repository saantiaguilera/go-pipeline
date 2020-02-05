package paint_example_test

import (
	"fmt"

	"github.com/saantiaguilera/go-pipeline"
)

type PaintVolumeStep struct{}

func (s *PaintVolumeStep) Run(ctx pipeline.Context) error {
	volume, _ := ctx.Get(TagVolume)
	fmt.Printf("Painting %d volume\n", volume)
	return nil
}

func (s *PaintVolumeStep) Name() string {
	return "paint volume painting"
}
