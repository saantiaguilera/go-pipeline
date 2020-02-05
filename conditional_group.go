package pipeline

type conditionalGroup struct {
	Statement Statement
	True      Stage
	False     Stage
}

func (c *conditionalGroup) Draw(graph GraphDiagram) {
	graph.AddDecision(
		c.Statement.Name(),
		func(graph GraphDiagram) {
			if c.True != nil {
				c.True.Draw(graph)
			}
		},
		func(graph GraphDiagram) {
			if c.False != nil {
				c.False.Draw(graph)
			}
		},
	)
}

func (c *conditionalGroup) Run(executor Executor, ctx Context) error {
	if c.Statement != nil && c.Statement.Evaluate(ctx) {
		if c.True != nil {
			return c.True.Run(executor, ctx)
		}
	} else {
		if c.False != nil {
			return c.False.Run(executor, ctx)
		}
	}
	return nil
}

// CreateConditionalGroup creates a conditional stage that will run a statement. If it holds true, then the "true" stage will be run.
// Else, the "false" one.
// If a statement is nil, then it will be considered to hold false (thus, the "false" stage is called)
// If one of the stages is nil and the statement is such, then nothing will happen.
func CreateConditionalGroup(statement Statement, true Stage, false Stage) Stage {
	return &conditionalGroup{
		Statement: statement,
		True:      true,
		False:     false,
	}
}
