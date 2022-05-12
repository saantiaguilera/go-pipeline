package pipeline

type ConditionalGroup[T any] struct {
	statement *Statement[T]
	trueCn    Stage[T]
	falseCn   Stage[T]
}

func (c *ConditionalGroup[T]) Draw(graph GraphDiagram) {
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

func (c *ConditionalGroup[T]) Run(executor Executor[T], in T) error {
	if c.statement != nil && c.statement.Evaluate(in) {
		if c.trueCn != nil {
			return c.trueCn.Run(executor, in)
		}
	} else {
		if c.falseCn != nil {
			return c.falseCn.Run(executor, in)
		}
	}
	return nil
}

// NewConditionalGroup News a conditional stage that will run a statement. If it holds true, then the "true" stage will be run.
// Else, the "false" one.
// If a statement is nil, then it will be considered to hold false (thus, the "false" stage is called)
// If one of the stages is nil and the statement is such, then nothing will happen.
func NewConditionalGroup[T any](statement *Statement[T], t, f Stage[T]) *ConditionalGroup[T] {
	return &ConditionalGroup[T]{
		statement: statement,
		trueCn:    t,
		falseCn:   f,
	}
}
