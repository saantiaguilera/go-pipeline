package pipeline

type Executor interface {

	Run(runnable Runnable) error

}

