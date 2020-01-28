package paint_example_test

import "fmt"

type CalculateVolumeStep struct {
	Depth  int
	Width  int
	Height int

	Volume int
}

func (s *CalculateVolumeStep) Run() error {
	s.Volume = s.Depth * s.Width * s.Height
	fmt.Printf("Getting volume: %d\n", s.Volume)
	return nil
}

func (s *CalculateVolumeStep) Name() string {
	return "get volume"
}
