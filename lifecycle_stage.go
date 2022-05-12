package pipeline

// BeforeStage is an alias for before hooks of a stage about to be executed with a given context.
// If the hook fails, the stage won't be executed
type BeforeStage[T any] func(stage Stage[T], in T) error

// AfterStage is an alias for after hooks of a stage. If the stage fails, one can recover from here or fallback to a new error.
// Also, this stage can fail, thus failing the execution (note that this is a blob of a stage, so if a hook fails, the stage fails too).
// The provided context was the resulting one after the stage was executed
type AfterStage[T any] func(stage Stage[T], in T, err error) error

// Blob structure that allows us to decorate a stage with pre/post hooks
// Note that we can compose many lifecycle stages if we want to have multiple hooks. Such as:
// lifecycleStage := NewLifecycleStage(realStage, aBeforeHook, anAfterHook)
// lifecycleStage = NewBeforeStageLifecycle(lifecycleStage, anotherBeforeHook)
// lifecycleStage = NewAfterStageLifecycle(lifecycleStage, anotherAfterHook)
type lifecycleStage[T any] struct {
	Before BeforeStage[T]
	After  AfterStage[T]
	Stage  Stage[T]
}

func (l *lifecycleStage[T]) Draw(graph GraphDiagram) {
	l.Stage.Draw(graph)
}

// Run the hooks and the stage, validating errors along the way and mutating the stage error in case it failed.
func (l *lifecycleStage[T]) Run(executor Executor[T], in T) error {
	if l.Before != nil {
		err := l.Before(l.Stage, in)

		if err != nil {
			return err
		}
	}

	err := l.Stage.Run(executor, in)

	if l.After != nil {
		err = l.After(l.Stage, in, err)
	}
	return err
}

// NewBeforeStageLifecycle News a lifecycle stage with a before hook
func NewBeforeStageLifecycle[T any](stage Stage[T], before BeforeStage[T]) Stage[T] {
	return &lifecycleStage[T]{
		Before: before,
		Stage:  stage,
	}
}

// NewAfterStageLifecycle News a lifecycle stage with an after hook
func NewAfterStageLifecycle[T any](stage Stage[T], after AfterStage[T]) Stage[T] {
	return &lifecycleStage[T]{
		After: after,
		Stage: stage,
	}
}

// NewStageLifecycle News a lifecycle stage with a before and an after hook
func NewStageLifecycle[T any](stage Stage[T], before BeforeStage[T], after AfterStage[T]) Stage[T] {
	return &lifecycleStage[T]{
		Before: before,
		After:  after,
		Stage:  stage,
	}
}
