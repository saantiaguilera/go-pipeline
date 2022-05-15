package pipeline

import (
	"context"
	"fmt"
)

type (
	ConditionalStep[I, O any] struct {
		statement conditionalStatement[I]
		trueCn    Step[I, O]
		falseCn   Step[I, O]
	}

	conditionalStatement[T any] interface {
		Name() string
		Evaluate(context.Context, T) bool
	}
)

// NewConditionalStep creates a conditional step that will run a statement. If it holds true, then the "true" step will be run.
// Else, the "false" step will be called.
// If a statement is nil, then it will be considered to hold false (thus, the "false" step is called)
// If one of the steps is nil and the statement is such, then nothing will happen.
func NewConditionalStep[I, O any](statement conditionalStatement[I], t, f Step[I, O]) ConditionalStep[I, O] {
	return ConditionalStep[I, O]{
		statement: statement,
		trueCn:    t,
		falseCn:   f,
	}
}

func (c ConditionalStep[I, O]) Draw(graph GraphDiagram) {
	graph.AddDecision(
		c.statement.Name(),
		func(graph GraphDiagram) {
			if c.trueCn != nil {
				c.trueCn.Draw(graph)
			}
		},
		func(graph GraphDiagram) {
			if c.falseCn != nil {
				c.falseCn.Draw(graph)
			}
		},
	)
}

func (c ConditionalStep[I, O]) Run(ctx context.Context, in I) (O, error) {
	ok := c.statement.Evaluate(ctx, in)
	if ok {
		if c.trueCn != nil {
			return c.trueCn.Run(ctx, in)
		}
	} else {
		if c.falseCn != nil {
			return c.falseCn.Run(ctx, in)
		}
	}
	return *new(O), fmt.Errorf("conditional step '%s' cannot run since the evaluated condition (%v) has a nil branch", c.statement.Name(), ok)
}
