package pipeline

import (
	"reflect"
	"runtime"
	"strings"
)

// Statement is an alias for a function that returns a boolean
type Statement func() bool

type conditionalGroup struct {
	Statement Statement
	True      Stage
	False     Stage
}

func (c *conditionalGroup) Draw(graph GraphDiagram) {
	name := runtime.FuncForPC(reflect.ValueOf(c.Statement).Pointer()).Name()
	name = name[strings.LastIndexByte(name, '.')+1:]
	if i := strings.LastIndexByte(name, '-'); i > 0 {
		name = name[:i]
	}

	graph.AddDecision(
		name,
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
