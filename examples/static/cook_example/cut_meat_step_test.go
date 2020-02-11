package cook_example_test

import (
	"fmt"
	"time"

	"github.com/saantiaguilera/go-pipeline"
)

type cutMeatStep struct{}

func (s *cutMeatStep) Name() string {
	return "cut_meat_step"
}

func (s *cutMeatStep) Run(ctx pipeline.Context) error {
	meatSize, _ := ctx.GetInt(TagMeatSize)
	ovenSize, _ := ctx.GetInt(TagOvenSize)
	fmt.Printf("Cutting meat of size %d into %d\n", meatSize, ovenSize)
	time.Sleep(1 * time.Second)

	ctx.Set(TagMeatSize, ovenSize)
	return nil
}

func CreateCutMeatStep() pipeline.Step {
	return &cutMeatStep{}
}

func CreateMeatTooBigStatement() func(ctx pipeline.Context) bool {
	return func(ctx pipeline.Context) bool {
		meatSize, _ := ctx.GetInt(TagMeatSize)
		ovenSize, _ := ctx.GetInt(TagOvenSize)
		return meatSize > ovenSize
	}
}
