package cook_example_test

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

func (s *cutMeatStep) Run() error {
	fmt.Printf("Cutting meat of size %d into %d\n", s.MeatSize, s.OvenSize)
	time.Sleep(1 * time.Second)

	s.Stream <- s.OvenSize
	return nil
}

func CreateCutMeatStep(meatSize, ovenSize int, meatChan chan int) pipeline.Step {
	return &cutMeatStep{
		MeatSize: meatSize,
		OvenSize: ovenSize,
		Stream:   meatChan,
	}
}

type addMeatStep struct {
	MeatSize int
	OvenSize int
	Stream   chan int
}

func (s *addMeatStep) Name() string {
	return "cut_meat_step"
}

func (s *addMeatStep) Run() error {
	fmt.Printf("Adding %d meat\n", s.OvenSize-s.MeatSize)
	time.Sleep(1 * time.Second) // Simulate time it takes to do this action

	s.Stream <- s.OvenSize
	return nil
}

func CreateAddMeatStep(meatSize, ovenSize int, meatChan chan int) pipeline.Step {
	return &addMeatStep{
		MeatSize: meatSize,
		OvenSize: ovenSize,
		Stream:   meatChan,
	}
}

func IsMeatTooBigForTheOven(meatSize, ovenSize int) bool {
	return meatSize > ovenSize
}
