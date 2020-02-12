package main

import (
	"fmt"

	"github.com/saantiaguilera/go-pipeline"
)

type getDepthStep struct{}

func (s *getDepthStep) Run(ctx pipeline.Context) error {
	fmt.Print("Getting depth: 1234\n")
	ctx.Set(tagDepth, 1234)
	return nil
}

func (s *getDepthStep) Name() string {
	return "get depth"
}
