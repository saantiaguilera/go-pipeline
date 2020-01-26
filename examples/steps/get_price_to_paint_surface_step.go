package steps

import (
	"fmt"
	"time"
)

type GetPriceToPaintSurfaceStep struct {
	Surface int

	Price float64
}

func (s *GetPriceToPaintSurfaceStep) Run() error {
	s.Price = float64(s.Surface) * 348.1923
	time.Sleep(100 * time.Millisecond) // Simulate time it takes to get this from a service
	fmt.Printf("Getting price for painting surface: %f\n", s.Price)
	return nil
}

func (s *GetPriceToPaintSurfaceStep) Name() string {
	return "get price for painting surface"
}
