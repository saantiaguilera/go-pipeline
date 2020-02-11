package cook_example_test

import (
	"fmt"
	"time"

	"github.com/saantiaguilera/go-pipeline"
)

type putMeatInTheOvenStep struct {
	Meat chan int
}

func (s *putMeatInTheOvenStep) Name() string {
	return "put_meat_in_the_oven_step"
}

func (s *putMeatInTheOvenStep) Run(ctx pipeline.Context) error {
	meat := <-s.Meat
	fmt.Printf("Putting in the oven %d meat\n", meat)
	time.Sleep(1 * time.Second) // Simulate time it takes to do this action

	s.Meat <- meat
	return nil
}

func CreatePutMeatInOvenStep(meatChan chan int) pipeline.Step {
	return &putMeatInTheOvenStep{
		Meat: meatChan,
	}
}
