package pipeline

// Statement is an alias for a function that returns a boolean
type Statement func() bool

type conditionalGroup struct {
	Statement Statement
	True      Stage
	False     Stage
}

func (c *conditionalGroup) Run(executor Executor) error {
	if c.Statement != nil && c.Statement() {
		if c.True != nil {
			return c.True.Run(executor)
		}
	} else {
		if c.False != nil {
			return c.False.Run(executor)
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