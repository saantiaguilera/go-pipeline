package pipeline

type Lifecycle interface {

	Before(stage Stage) error

	After(stage Stage, err error) error

}

type Stage interface {
	Run(executor Executor) error
}
