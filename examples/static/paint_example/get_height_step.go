package main

import (
	"fmt"

	"github.com/saantiaguilera/go-pipeline"
)

type getHeightStep struct{}

func (s *getHeightStep) Run(ctx pipeline.Context) error {
	fmt.Print("Getting height: 1234\n")
	ctx.Set(tagHeight, 1234)
	return nil
}

func (s *getHeightStep) Name() string {
	return "get height"
}
