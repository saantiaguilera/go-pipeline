package paint_example_test

import (
	"fmt"
	"time"

	"github.com/saantiaguilera/go-pipeline"
)

type RecordPriceStep struct {
	PriceType pipeline.Tag
}

func (s *RecordPriceStep) Run(ctx pipeline.Context) error {
	price, _ := ctx.GetFloat64(s.PriceType)
	time.Sleep(100 * time.Millisecond) // Simulate time it takes to use a service
	fmt.Printf("Recording price: %f\n", price)
	return nil
}

func (s *RecordPriceStep) Name() string {
	return "record a price"
}
