package paint_example_test

import (
	"fmt"
	"time"
)

type GetPriceToPaintVolumeStep struct {
	Volume int

	Price float64
}

func (s *GetPriceToPaintVolumeStep) Run() error {
	s.Price = float64(s.Volume) * 348.1923
	time.Sleep(100 * time.Millisecond) // Simulate time it takes to get this from a service
	fmt.Printf("Getting price for painting volume: %f\n", s.Price)
	return nil
}

func (s *GetPriceToPaintVolumeStep) Name() string {
	return "get price for painting volume"
}
