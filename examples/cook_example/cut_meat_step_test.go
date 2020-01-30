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

func CreateMeatTooBigStatement(meatSize, ovenSize int) func() bool {
	s := &MeatTooBig{
		MeatSize: meatSize,
		OvenSize: ovenSize,
	}
	return s.IsMeatTooBigForTheOven
}

type MeatTooBig struct {
	MeatSize int
	OvenSize int
}

func (m *MeatTooBig) IsMeatTooBigForTheOven() bool {
	return m.MeatSize > m.OvenSize
}
