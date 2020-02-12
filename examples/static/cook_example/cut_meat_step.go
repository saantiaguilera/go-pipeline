package main

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
	meatSize, _ := ctx.GetInt(tagMeatSize)
	ovenSize, _ := ctx.GetInt(tagOvenSize)
	fmt.Printf("Cutting meat of size %d into %d\n", meatSize, ovenSize)
	time.Sleep(1 * time.Second)

	ctx.Set(tagMeatSize, ovenSize)
	return nil
}

func createCutMeatStep() pipeline.Step {
	return &cutMeatStep{}
}

func createMeatTooBigStatement() func(ctx pipeline.Context) bool {
	return func(ctx pipeline.Context) bool {
		meatSize, _ := ctx.GetInt(tagMeatSize)
		ovenSize, _ := ctx.GetInt(tagOvenSize)
		return meatSize > ovenSize
	}
}
