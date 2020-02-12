package main

import (
	"fmt"

	"github.com/saantiaguilera/go-pipeline"
)

type getWidthStep struct{}

func (s *getWidthStep) Run(ctx pipeline.Context) error {
	fmt.Print("Getting width: 1234\n")
	ctx.Set(tagWidth, 1234)
	return nil
}

func (s *getWidthStep) Name() string {
	return "get width"
}
