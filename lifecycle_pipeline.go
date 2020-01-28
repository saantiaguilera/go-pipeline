package pipeline

// BeforePipeline is an alias for before hooks of a stage about to be run by a pipeline. If the hook fails, the stage
// won't be executed
type BeforePipeline func(stage Stage) error

// AfterPipeline is an alias for after hooks of a stage about to be run by a pipeline. If the pipeline fails,
// one can recover from here or fallback to a new error. Also, this pipeline can fail, thus failing the execution (note
// that this is a blob of a pipeline, so if a hook fails, the pipeline execution fails too).
type AfterPipeline func(stage Stage, err error) error

// Blob structure that allows us to decorate a pipeline with pre/post hooks
// Note that we can compose many lifecycle pipelines if we want to have multiple hooks. Such as:
// lifecyclePipeline := CreateLifecyclePipeline(realPipeline, aBeforeHook, anAfterHook)
// lifecyclePipeline = CreateBeforePipelineLifecycle(lifecyclePipeline, anotherBeforeHook)
// lifecyclePipeline = CreateAfterPipelineLifecycle(lifecyclePipeline, anotherAfterHook)
// lifecyclePipeline.Run(stage)
type lifecyclePipeline struct {
	Before   BeforePipeline
	After    AfterPipeline
	Pipeline Pipeline
}

// Run the hooks and the stage, validating errors along the way and mutating the pipeline error in case it failed.
func (l *lifecyclePipeline) Run(stage Stage) error {
	if l.Before != nil {
		err := l.Before(stage)

		if err != nil {
			return err
		}
	}

	err := l.Pipeline.Run(stage)

	if l.After != nil {
		err = l.After(stage, err)
	}
	return err
}

// CreateBeforePipelineLifecycle creates a lifecycle pipeline with a before hook
func CreateBeforePipelineLifecycle(pipeline Pipeline, before BeforePipeline) Pipeline {
	return &lifecyclePipeline{
		Before:   before,
		Pipeline: pipeline,
	}
}

// CreateAfterPipelineLifecycle creates a lifecycle pipeline with an after hook
func CreateAfterPipelineLifecycle(pipeline Pipeline, after AfterPipeline) Pipeline {
	return &lifecyclePipeline{
		After:    after,
		Pipeline: pipeline,
	}
}

// CreatePipelineLifecycle creates a lifecycle pipeline with a before and an after hook
func CreatePipelineLifecycle(pipeline Pipeline, before BeforePipeline, after AfterPipeline) Pipeline {
	return &lifecyclePipeline{
		Before:   before,
		After:    after,
		Pipeline: pipeline,
	}
}
