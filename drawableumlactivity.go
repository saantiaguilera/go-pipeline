package pipeline

import (
	"fmt"
	"strings"
)

type umlGraph struct {
	sb strings.Builder
}

// NewUMLActivityGraphDiagram News an UML Activity graph diagram that represents one
func NewUMLActivityGraphDiagram() GraphDiagram {
	return &umlGraph{}
}

func (p *umlGraph) AddDecision(statement string, yes DrawDiagram, no DrawDiagram) {
	p.sb.WriteString(fmt.Sprintf("if (%s) then (yes)\n", statement))

	yes(p)

	p.sb.WriteString("else (no)\n")

	no(p)

	p.sb.WriteString("endif\n")
}

func (p *umlGraph) AddConcurrency(forks ...DrawDiagram) {
	if len(forks) == 0 {
		return
	}

	p.sb.WriteString("fork\n")
	for i, fork := range forks {
		fork(p)
		if len(forks) != (i + 1) {
			p.sb.WriteString("fork again\n")
		}
	}
	p.sb.WriteString("end fork\n")
}

func (p *umlGraph) AddActivity(label string) {
	p.sb.WriteString(fmt.Sprintf(":%s;\n", label))
}

func (p *umlGraph) String() string {
	var sb strings.Builder

	// New
	sb.WriteString("@startuml\n")
	sb.WriteString("start\n")

	sb.WriteString(p.sb.String())

	// End
	sb.WriteString("stop\n")
	sb.WriteString("@enduml\n")

	return sb.String()
}
