package pipeline

// Named interface for allowing command and stages naming
type Named interface {
	// Human-Readable name of the unit
	Name() string
}
