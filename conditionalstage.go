package pipeline

type (
	ConditionalStage[T any] struct {
		statement *Statement[T]
		trueCn    Step[T]
		falseCn   Step[T]
	}
)

// NewConditionalStage News a conditional stage that will run a statement. If it holds true, then the "true" step will be run.
// Else, the "false" step will be called.
// If a statement is nil, then it will be considered to hold false (thus, the "false" step is called)
// If one of the steps is nil and the statement is such, then nothing will happen.
func NewConditionalStage[T any](statement *Statement[T], t, f Step[T]) *ConditionalStage[T] {
	return &ConditionalStage[T]{
		statement: statement,
		trueCn:    t,
		falseCn:   f,
	}
}

func (c *ConditionalStage[T]) Draw(graph GraphDiagram) {
	graph.AddDecision(
		c.statement.Name(),
		func(graph GraphDiagram) {
			if c.trueCn != nil {
				graph.AddActivity(c.trueCn.Name())
			}
		},
		func(graph GraphDiagram) {
			if c.falseCn != nil {
				graph.AddActivity(c.falseCn.Name())
			}
		},
	)
}

func (c *ConditionalStage[T]) Run(executor Executor[T], in T) error {
	if c.statement != nil && c.statement.Evaluate(in) {
		if c.trueCn != nil {
			return executor.Run(c.trueCn, in)
		}
	} else {
		if c.falseCn != nil {
			return executor.Run(c.falseCn, in)
		}
	}
	return nil
}
