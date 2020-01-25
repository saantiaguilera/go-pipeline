package steps

import "fmt"

type PaintSurfaceStep struct {
	Surface int
}

func (s *PaintSurfaceStep) Run() error {
	fmt.Printf("Painting %d surface\n", s.Surface)
	return nil
}

func (s *PaintSurfaceStep) Name() string {
	return "paint surface painting"
}
