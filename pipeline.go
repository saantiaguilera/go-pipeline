package pipeline

type Named interface {
	Name() string
}

type Pipeline interface {
	Run(stage Stage) error
	AddOnBeforePipelineRun(beforePipeline func(stage Stage) error)
	AddAfterPipelineRun(afterPipeline func(stage Stage, err error) error)
}
