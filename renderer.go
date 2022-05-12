package pipeline

import (
	"io"
)

type (
	// DiagramRenderer for creating renderings of graphs
	DiagramRenderer interface {
		// Render of the given stage, in the given output
		Render(graphDiagram GraphDiagram, output io.WriteCloser) error
	}
)
