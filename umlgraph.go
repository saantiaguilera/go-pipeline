package pipeline

import (
	"fmt"
	"strings"
)

// UMLGraph represents a graph that can render itself into UML
type UMLGraph struct {
	sb strings.Builder
}

// NewUMLGraph createsn UML Activity graph diagram that represents one
func NewUMLGraph() *UMLGraph {
	return &UMLGraph{}
}

func (p *UMLGraph) AddDecision(statement string, yes GraphDrawer, no GraphDrawer) {
	p.sb.WriteString(fmt.Sprintf("if (%s) then (yes)\n", statement))

	yes(p)

	p.sb.WriteString("else (no)\n")

	no(p)

	p.sb.WriteString("endif\n")
}

func (p *UMLGraph) AddConcurrency(forks ...GraphDrawer) {
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

func (p *UMLGraph) AddActivity(label string) {
	p.sb.WriteString(fmt.Sprintf(":%s;\n", label))
}

func (p *UMLGraph) String() string {
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
