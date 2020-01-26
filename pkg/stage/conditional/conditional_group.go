package conditional

import (
	"github.com/saantiaguilera/go-pipeline/pkg"
)

// Alias for a function that returns a boolean
type Statement func() bool

type conditionalGroup struct {
	Statement Statement
	True      pkg.Stage
	False     pkg.Stage
}

func (c *conditionalGroup) Run(executor pkg.Executor) error {
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

// Create a conditional stage that will run a statement. If it holds true, then the "true" stage will be run.
// Else, the "false" one.
// If a statement is nil, then it will be considered to hold false (thus, the "false" stage is called)
// If one of the stages is nil and the statement is such, then nothing will happen.
func CreateConditionalGroup(statement Statement, true pkg.Stage, false pkg.Stage) pkg.Stage {
	return &conditionalGroup{
		Statement: statement,
		True:      true,
		False:     false,
	}
}