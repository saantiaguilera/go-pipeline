package steps

import "fmt"

type CalculateSurfaceStep struct {
	Width  int
	Height int

	Surface int
}

func (s *CalculateSurfaceStep) Run() error {
	s.Surface = s.Width * s.Height
	fmt.Printf("Getting surface: %d\n", s.Surface)
	return nil
}

func (s *CalculateSurfaceStep) Name() string {
	return "get surface"
}
