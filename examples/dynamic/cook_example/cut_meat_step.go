package main

import (
	"fmt"
	"time"

	"github.com/saantiaguilera/go-pipeline"
)

type cutMeatStep struct {
	MeatSize int
	OvenSize int
	Stream   chan int
}

func (s *cutMeatStep) Name() string {
	return "cut_meat_step"
}

func (s *cutMeatStep) Run(ctx pipeline.Context) error {
	fmt.Printf("Cutting meat of size %d into %d\n", s.MeatSize, s.OvenSize)
	time.Sleep(1 * time.Second)

	s.Stream <- s.OvenSize
	return nil
}

func NewCutMeatStep(meatSize, ovenSize int, meatChan chan int) pipeline.Step {
	return &cutMeatStep{
		MeatSize: meatSize,
		OvenSize: ovenSize,
		Stream:   meatChan,
	}
}

func NewMeatTooBigStatement(meatSize, ovenSize int) func(ctx pipeline.Context) bool {
	s := &meatTooBig{
		MeatSize: meatSize,
		OvenSize: ovenSize,
	}
	return s.isMeatTooBigForTheOven
}

type meatTooBig struct {
	MeatSize int
	OvenSize int
}

func (m *meatTooBig) isMeatTooBigForTheOven(ctx pipeline.Context) bool {
	return m.MeatSize > m.OvenSize
}
