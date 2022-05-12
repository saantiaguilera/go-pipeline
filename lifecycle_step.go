package pipeline

// BeforeStep is an alias for before hooks of a step about to be executed with a given context.
// If the hook fails, the step won't be run.
type BeforeStep[T any] func(step Step[T], in T) error

// AfterStep is an alias for after hooks of a step. If the step fails, one can recover from here or fallback to a new error.
// Also, this step can fail, thus failing the unit (note that this is a blob of a step, so if a hook fails, the step fails too).
// The context provided is the resulting one after the step was executed
type AfterStep[T any] func(step Step[T], in T, err error) error

// Blob structure that allows us to decorate a step with pre/post hooks
// Note that we can compose many lifecycle steps if we want to have multiple hooks. Such as:
// lifecycleStep := NewLifecycleStep(realStep, aBeforeHook, anAfterHook)
// lifecycleStep = NewBeforeStepLifecycle(lifecycleStep, anotherBeforeHook)
// lifecycleStep = NewAfterStepLifecycle(lifecycleStep, anotherAfterHook)
type lifecycleStep[T any] struct {
	Before BeforeStep[T]
	After  AfterStep[T]
	Step   Step[T]
}

// Run the hooks and the step, validating errors along the way and mutating the step error in case it failed.
func (l *lifecycleStep[T]) Run(in T) error {
	if l.Before != nil {
		err := l.Before(l.Step, in)

		if err != nil {
			return err
		}
	}

	err := l.Step.Run(in)

	if l.After != nil {
		err = l.After(l.Step, in, err)
	}
	return err
}

// Name is delegated to the real step
func (l *lifecycleStep[T]) Name() string {
	return l.Step.Name()
}

// NewBeforeStepLifecycle News a lifecycle step with a before hook
func NewBeforeStepLifecycle[T any](step Step[T], before BeforeStep[T]) Step[T] {
	return &lifecycleStep[T]{
		Before: before,
		Step:   step,
	}
}

// NewAfterStepLifecycle News a lifecycle step with an after hook
func NewAfterStepLifecycle[T any](step Step[T], after AfterStep[T]) Step[T] {
	return &lifecycleStep[T]{
		After: after,
		Step:  step,
	}
}

// NewStepLifecycle News a lifecycle step with a before and an after hook
func NewStepLifecycle[T any](step Step[T], before BeforeStep[T], after AfterStep[T]) Step[T] {
	return &lifecycleStep[T]{
		Before: before,
		After:  after,
		Step:   step,
	}
}
