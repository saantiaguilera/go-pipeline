package pipeline

type Stage interface {
	Run(executor Executor) error
}
