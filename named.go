package pipeline

// Named interface for allowing command and stages naming
// TODO: At a later stage it would be nice to graph the pipeline itself with this
type Named interface {
	// Human-Readable name of the unit
	Name() string
}
