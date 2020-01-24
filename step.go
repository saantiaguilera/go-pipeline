package pipeline

type Runnable interface {
	Run() error
}

type Step interface {
	Runnable
	Named
}