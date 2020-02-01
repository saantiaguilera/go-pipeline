package pipeline

import (
	"io"
)

// DiagramRenderer for creating renderings of graphs
type DiagramRenderer interface {
	// Render of the given stage, in the given output
	Render(graphDiagram GraphDiagram, output io.WriteCloser) error
}
