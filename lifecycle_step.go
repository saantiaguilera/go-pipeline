package pipeline

// BeforeStep is an alias for before hooks of a step about to be executed with a given context.
// If the hook fails, the step won't be run.
type BeforeStep func(step Step, ctx Context) error

// AfterStep is an alias for after hooks of a step. If the step fails, one can recover from here or fallback to a new error.
// Also, this step can fail, thus failing the unit (note that this is a blob of a step, so if a hook fails, the step fails too).
// The context provided is the resulting one after the step was executed
type AfterStep func(step Step, ctx Context, err error) error

// Blob structure that allows us to decorate a step with pre/post hooks
// Note that we can compose many lifecycle steps if we want to have multiple hooks. Such as:
// lifecycleStep := CreateLifecycleStep(realStep, aBeforeHook, anAfterHook)
// lifecycleStep = CreateBeforeStepLifecycle(lifecycleStep, anotherBeforeHook)
// lifecycleStep = CreateAfterStepLifecycle(lifecycleStep, anotherAfterHook)
type lifecycleStep struct {
	Before BeforeStep
	After  AfterStep
	Step   Step
}

// Run the hooks and the step, validating errors along the way and mutating the step error in case it failed.
func (l *lifecycleStep) Run(ctx Context) error {
	if l.Before != nil {
		err := l.Before(l.Step, ctx)

		if err != nil {
			return err
		}
	}

	err := l.Step.Run(ctx)

	if l.After != nil {
		err = l.After(l.Step, ctx, err)
	}
	return err
}

// Name is delegated to the real step
func (l *lifecycleStep) Name() string {
	return l.Step.Name()
}

// CreateBeforeStepLifecycle creates a lifecycle step with a before hook
func CreateBeforeStepLifecycle(step Step, before BeforeStep) Step {
	return &lifecycleStep{
		Before: before,
		Step:   step,
	}
}

// CreateAfterStepLifecycle creates a lifecycle step with an after hook
func CreateAfterStepLifecycle(step Step, after AfterStep) Step {
	return &lifecycleStep{
		After: after,
		Step:  step,
	}
}

// CreateStepLifecycle creates a lifecycle step with a before and an after hook
func CreateStepLifecycle(step Step, before BeforeStep, after AfterStep) Step {
	return &lifecycleStep{
		Before: before,
		After:  after,
		Step:   step,
	}
}
