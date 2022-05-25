package pipeline

type (
	// DrawableGraph is a contract for drawing in graphs
	DrawableGraph interface {
		Draw(graph Graph)
	}

	// Graph interface allowing to create a representation/drawing of a graph
	Graph interface {
		// AddConcurrency branching as many times as needed (each branch is a concurrent/fork 'node')
		AddConcurrency(branches ...GraphDrawer)
		// AddDecision from a given statement, allowing inner graphs for each branch of the decision
		AddDecision(statement string, yes GraphDrawer, no GraphDrawer)
		// Create an action entry
		AddActivity(label string)
	}

	// GraphDrawer alias for Draw(Graph) functions
	GraphDrawer = func(Graph)
)
