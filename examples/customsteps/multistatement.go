package main

import (
	"context"
	"errors"

	"github.com/saantiaguilera/go-pipeline"
)

const (
	exampleChoiceA exampleChoice = iota
	exampleChoiceB
	exampleChoiceC
)

type (
	// multiStatementStep is a step that allows multiple branchings from a statement
	multiStatementStep[I, O any] struct {
		stmt                      func(context.Context, I) exampleChoice
		choiceA, choiceB, choiceC pipeline.Step[I, O]
	}

	exampleChoice int
)

func (c multiStatementStep[I, O]) Draw(graph pipeline.Graph) {
	graph.AddDecision(
		"choice A",
		func(graph pipeline.Graph) {
			if c.choiceA != nil {
				c.choiceA.Draw(graph)
			}
		},
		func(graph pipeline.Graph) {
			graph.AddDecision(
				"choice B else C",
				func(graph pipeline.Graph) {
					if c.choiceB != nil {
						c.choiceB.Draw(graph)
					}
				},
				func(g pipeline.Graph) {
					if c.choiceC != nil {
						c.choiceC.Draw(graph)
					}
				},
			)
		},
	)
}

func (c multiStatementStep[I, O]) Run(ctx context.Context, in I) (O, error) {
	switch c.stmt(ctx, in) {
	case exampleChoiceA:
		return c.choiceA.Run(ctx, in)
	case exampleChoiceB:
		return c.choiceB.Run(ctx, in)
	case exampleChoiceC:
		return c.choiceC.Run(ctx, in)
	}
	return *new(O), errors.New("unknown choice returned")
}
