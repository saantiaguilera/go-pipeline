package steps

import "fmt"

type PaintVolumeStep struct {
	Volume int
}

func (s *PaintVolumeStep) Run() error {
	fmt.Printf("Painting %d volume\n", s.Volume)
	return nil
}

func (s *PaintVolumeStep) Name() string {
	return "paint volume painting"
}
