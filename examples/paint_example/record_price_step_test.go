package paint_example_test

import (
	"fmt"
	"time"
)

type RecordPriceStep struct {
	Price float64
}

func (s *RecordPriceStep) Run() error {
	time.Sleep(100 * time.Millisecond) // Simulate time it takes to use a service
	fmt.Printf("Recording price: %f\n", s.Price)
	return nil
}

func (s *RecordPriceStep) Name() string {
	return "record a price"
}
