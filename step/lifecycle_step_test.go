package pipeline_step_test

import (
	"errors"
	"github.com/saantiaguilera/go-pipeline"
	pipeline_step "github.com/saantiaguilera/go-pipeline/step"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type mockStep struct {
	mock.Mock
}

func (m *mockStep) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *mockStep) Run() error {
	args := m.Called()
	return args.Error(0)
}

func TestLifecycleStep_GivenAfterFunc_WhenRun_ThenFuncAreCalled(t *testing.T) {
	called := false
	afterFunc := func(step pipeline.Step, err error) error {
		called = true
		return nil
	}
	step := new(mockStep)
	step.On("Run").Return(nil).Once()

	lifecycleStep := pipeline_step.CreateAfterStepLifecycle(step, afterFunc)
	err := lifecycleStep.Run()

	assert.Nil(t, err)
	assert.True(t, called)
	step.AssertExpectations(t)
}

func TestLifecycleStep_GivenAfterFuncErroring_WhenRun_ThenErrorIsReturned(t *testing.T) {
	expectedErr := errors.New("some error")
	afterFunc := func(step pipeline.Step, err error) error {
		return expectedErr
	}
	step := new(mockStep)
	step.On("Run").Return(nil).Once()

	lifecycleStep := pipeline_step.CreateAfterStepLifecycle(step, afterFunc)
	err := lifecycleStep.Run()

	assert.Equal(t, expectedErr, err)
	step.AssertExpectations(t)
}

func TestLifecycleStep_GivenAfterFuncRecoveringError_WhenRun_ThenFuncCanRecover(t *testing.T) {
	expectedErr := errors.New("some error")
	var retrievedErr error
	afterFunc := func(step pipeline.Step, err error) error {
		retrievedErr = err
		return nil
	}
	step := new(mockStep)
	step.On("Run").Return(expectedErr).Once()

	lifecycleStep := pipeline_step.CreateAfterStepLifecycle(step, afterFunc)
	err := lifecycleStep.Run()

	assert.Nil(t, err)
	assert.Equal(t, expectedErr, retrievedErr)
	step.AssertExpectations(t)
}

func TestLifecycleStep_GivenBeforeFunc_WhenRun_ThenFuncAreCalled(t *testing.T) {
	called := false
	beforeFunc := func(step pipeline.Step) error {
		called = true
		return nil
	}
	step := new(mockStep)
	step.On("Run").Return(nil).Once()

	lifecycleStep := pipeline_step.CreateBeforeStepLifecycle(step, beforeFunc)
	err := lifecycleStep.Run()

	assert.Nil(t, err)
	assert.True(t, called)
	step.AssertExpectations(t)
}

func TestLifecycleStep_GivenBeforeFuncReturningError_WhenRun_ThenErrorIsReturned(t *testing.T) {
	expectedErr := errors.New("some error")
	beforeFunc := func(step pipeline.Step) error {
		return expectedErr
	}
	step := new(mockStep)

	lifecycleStep := pipeline_step.CreateBeforeStepLifecycle(step, beforeFunc)
	err := lifecycleStep.Run()

	assert.Equal(t, expectedErr, err)
	step.AssertExpectations(t)
}

func TestLifecycleStep_GivenBeforeFuncReturningError_WhenRun_ThenStepAndAfterAreNotRan(t *testing.T) {
	expectedErr := errors.New("some error")
	called := false
	beforeFunc := func(step pipeline.Step) error {
		return expectedErr
	}
	afterFunc := func(step pipeline.Step, err error) error {
		called = true
		return nil
	}
	step := new(mockStep)

	lifecycleStep := pipeline_step.CreateStepLifecycle(step, beforeFunc, afterFunc)
	err := lifecycleStep.Run()

	assert.Equal(t, expectedErr, err)
	assert.False(t, called)
	step.AssertExpectations(t)
}

func TestLifecycleStep_GivenAStep_WhenRun_ThenStepIsRun(t *testing.T) {
	expectedErr := errors.New("some error")
	beforeFunc := func(step pipeline.Step) error {
		return nil
	}
	afterFunc := func(step pipeline.Step, err error) error {
		return err
	}
	step := new(mockStep)
	step.On("Run").Return(expectedErr).Once()

	lifecycleStep := pipeline_step.CreateStepLifecycle(step, beforeFunc, afterFunc)
	err := lifecycleStep.Run()

	assert.Equal(t, expectedErr, err)
	step.AssertExpectations(t)
}

func TestLifecycleStep_GivenAStep_WhenNamed_ThenStepIsDelegated(t *testing.T) {
	expectedName := "mock step"
	step := new(mockStep)
	step.On("Name").Return(expectedName).Once()

	lifecycleStep := pipeline_step.CreateStepLifecycle(step, func(step pipeline.Step) error {
		return nil
	}, func(step pipeline.Step, err error) error {
		return err
	})

	assert.Equal(t, expectedName, lifecycleStep.Name())
	step.AssertExpectations(t)
}

func TestLifecycleStep_GivenComposition_WhenRun_ThenCompositionBehavesAsAnArray(t *testing.T) {
	var callings []string
	before := func(step pipeline.Step) error {
		callings = append(callings, "before")
		return nil
	}
	after := func(step pipeline.Step, err error) error {
		callings = append(callings, "after")
		return err
	}

	step := new(mockStep)
	step.On("Run").Run(func(args mock.Arguments) {
		callings = append(callings, "step")
	}).Return(nil).Once()

	lifecycleStep := pipeline_step.CreateStepLifecycle(step, before, after)
	lifecycleStep = pipeline_step.CreateBeforeStepLifecycle(lifecycleStep, before)
	lifecycleStep = pipeline_step.CreateAfterStepLifecycle(lifecycleStep, after)
	lifecycleStep = pipeline_step.CreateBeforeStepLifecycle(lifecycleStep, before)

	err := lifecycleStep.Run()

	assert.Nil(t, err)
	assert.Len(t, callings, 6)
	assert.Equal(t, []string{"before", "before", "before", "step", "after", "after"}, callings)
	step.AssertExpectations(t)
}

