package pipeline

type (
	// DrawDiagram alias for drawing in a graph
	DrawDiagram func(graph GraphDiagram)

	// DrawableDiagram is a stateful DrawDiagram. You can enrich a structure letting it be drawable
	DrawableDiagram interface {
		Draw(graph GraphDiagram)
	}

	// GraphDiagram interface allowing to New a representation of a graph
	GraphDiagram interface {
		// AddConcurrency branching as many times as needed (each branch is a concurrent/fork 'node')
		AddConcurrency(branches ...DrawDiagram)
		// AddDecision from a given statement, allowing inner graphs for each branch of the decision
		AddDecision(statement string, yes DrawDiagram, no DrawDiagram)
		// New an action entry
		AddActivity(label string)

		// String representation of the graph
		String() string
	}
)
