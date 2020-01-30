package pipeline

import (
	"reflect"
	"runtime"
	"strings"
)

type conditionalStage struct {
	Statement Statement
	True      Step
	False     Step
}

func (c *conditionalStage) Draw(graph GraphDiagram) {
	name := runtime.FuncForPC(reflect.ValueOf(c.Statement).Pointer()).Name()
	name = name[strings.LastIndexByte(name, '.')+1:]
	if i := strings.LastIndexByte(name, '-'); i > 0 {
		name = name[:i]
	}

	graph.AddDecision(
		name,
		func(graph GraphDiagram) {
			if c.True != nil {
				graph.AddActivity(c.True.Name())
			}
		},
		func(graph GraphDiagram) {
			if c.False != nil {
				graph.AddActivity(c.False.Name())
			}
		},
	)
}

func (c *conditionalStage) Run(executor Executor) error {
	if c.Statement != nil && c.Statement() {
		if c.True != nil {
			return executor.Run(c.True)
		}
	} else {
		if c.False != nil {
			return executor.Run(c.False)
		}
	}
	return nil
}

// CreateConditionalStage creates a conditional stage that will run a statement. If it holds true, then the "true" step will be run.
// Else, the "false" step will be called.
// If a statement is nil, then it will be considered to hold false (thus, the "false" step is called)
// If one of the steps is nil and the statement is such, then nothing will happen.
func CreateConditionalStage(statement Statement, true Step, false Step) Stage {
	return &conditionalStage{
		Statement: statement,
		True:      true,
		False:     false,
	}
}
