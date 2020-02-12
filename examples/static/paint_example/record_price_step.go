package main

import (
	"fmt"
	"time"

	"github.com/saantiaguilera/go-pipeline"
)

type recordPriceStep struct {
	PriceType pipeline.Tag
}

func (s *recordPriceStep) Run(ctx pipeline.Context) error {
	price, _ := ctx.GetFloat64(s.PriceType)
	time.Sleep(100 * time.Millisecond) // Simulate time it takes to use a service
	fmt.Printf("Recording price: %f\n", price)
	return nil
}

func (s *recordPriceStep) Name() string {
	return "record a price"
}
