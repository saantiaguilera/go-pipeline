package pipeline

type Runnable interface {
	Named
	Run() error
}

type Step Runnable