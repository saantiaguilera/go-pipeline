package pipeline_step

import "github.com/saantiaguilera/go-pipeline"

// Alias for before hooks of a step. If the hook fails, the step won't be run
type BeforeStep func(step pipeline.Step) error

// Alias for after hooks of a step. If the step fails, one can recover from here or fallback to a new error.
// Also, this step can fail, thus failing the unit (note that this is a blob of a step, so if a hook fails, the step fails too).
type AfterStep func(step pipeline.Step, err error) error

// Blob structure that allows us to decorate a step with pre/post hooks
// Note that we can compose many lifecycle steps if we want to have multiple hooks. Such as:
// lifecycleStep := CreateLifecycleStep(realStep, aBeforeHook, anAfterHook)
// lifecycleStep = CreateBeforeStepLifecycle(lifecycleStep, anotherBeforeHook)
// lifecycleStep = CreateAfterStepLifecycle(lifecycleStep, anotherAfterHook)
type lifecycleStep struct {
	Before  BeforeStep
	After   AfterStep
	Step    pipeline.Step
}

// Run the hooks and the step, validating errors along the way and mutating the step error in case it failed.
func (l *lifecycleStep) Run() error {
	if l.Before != nil {
		err := l.Before(l.Step)

		if err != nil {
			return err
		}
	}

	err := l.Step.Run()

	if l.After != nil {
		err = l.After(l.Step, err)
	}
	return err
}

// Name is delegated to the real step
func (l *lifecycleStep) Name() string {
	return l.Step.Name()
}

// Create a lifecycle step with a before hook
func CreateBeforeStepLifecycle(step pipeline.Step, before BeforeStep) pipeline.Step {
	return &lifecycleStep{
		Before: before,
		Step:   step,
	}
}

// Create a lifecycle step with an after hook
func CreateAfterStepLifecycle(step pipeline.Step, after AfterStep) pipeline.Step {
	return &lifecycleStep{
		After:  after,
		Step:   step,
	}
}

// Create a lifecycle step with a before and an after hook
func CreateStepLifecycle(step pipeline.Step, before BeforeStep, after AfterStep) pipeline.Step {
	return &lifecycleStep{
		Before: before,
		After:  after,
		Step:   step,
	}
}