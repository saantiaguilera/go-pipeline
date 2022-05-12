package pipeline

import "context"

type (
	ConditionalContainer[T any] struct {
		statement Statement[T]
		trueCn    Container[T]
		falseCn   Container[T]
	}
)

// NewConditionalContainer creates a conditional container that will run a statement. If it holds true, then the "true" step will be run.
// Else, the "false" step will be called.
// If a statement is nil, then it will be considered to hold false (thus, the "false" step is called)
// If one of the steps is nil and the statement is such, then nothing will happen.
func NewConditionalContainer[T any](statement Statement[T], t, f Container[T]) ConditionalContainer[T] {
	return ConditionalContainer[T]{
		statement: statement,
		trueCn:    t,
		falseCn:   f,
	}
}

func (c ConditionalContainer[T]) Draw(graph GraphDiagram) {
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

func (c ConditionalContainer[T]) Visit(ctx context.Context, ex Executor[T], in T) error {
	if c.statement.Evaluate(in) {
		if c.trueCn != nil {
			return c.trueCn.Visit(ctx, ex, in)
		}
	} else {
		if c.falseCn != nil {
			return c.falseCn.Visit(ctx, ex, in)
		}
	}
	return nil
}
