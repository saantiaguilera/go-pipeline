package pipeline

// DrawDiagram alias for drawing in a graph
type DrawDiagram func(graph GraphDiagram)

// DrawableDiagram is a stateful DrawDiagram. You can enrich a structure letting it be drawable
type DrawableDiagram interface {
	Draw(graph GraphDiagram)
}

// GraphDiagram interface allowing to create a representation of a graph
type GraphDiagram interface {
	// AddConcurrency branching as many times as needed (each branch is a concurrent/fork 'node')
	AddConcurrency(branches ...DrawDiagram)
	// AddDecision from a given statement, allowing inner graphs for each branch of the decision
	AddDecision(statement string, yes DrawDiagram, no DrawDiagram)
	// Create an action entry
	AddActivity(label string)

	// String representation of the graph
	String() string
}
