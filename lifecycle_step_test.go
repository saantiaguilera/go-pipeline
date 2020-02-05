package pipeline_test

import (
	"errors"
	"testing"

	"github.com/saantiaguilera/go-pipeline"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLifecycleStep_GivenAfterFunc_WhenRun_ThenFuncAreCalled(t *testing.T) {
	called := false
	afterFunc := func(step pipeline.Step, ctx pipeline.Context, err error) error {
		called = true
		return nil
	}
	step := new(mockStep)
	step.On("Run", &mockContext{}).Return(nil).Once()

	lifecycleStep := pipeline.CreateAfterStepLifecycle(step, afterFunc)
	err := lifecycleStep.Run(&mockContext{})

	assert.Nil(t, err)
	assert.True(t, called)
	step.AssertExpectations(t)
}

func TestLifecycleStep_GivenAfterFuncErroring_WhenRun_ThenErrorIsReturned(t *testing.T) {
	expectedErr := errors.New("some error")
	afterFunc := func(step pipeline.Step, ctx pipeline.Context, err error) error {
		return expectedErr
	}
	step := new(mockStep)
	step.On("Run", &mockContext{}).Return(nil).Once()

	lifecycleStep := pipeline.CreateAfterStepLifecycle(step, afterFunc)
	err := lifecycleStep.Run(&mockContext{})

	assert.Equal(t, expectedErr, err)
	step.AssertExpectations(t)
}

func TestLifecycleStep_GivenAfterFuncRecoveringError_WhenRun_ThenFuncCanRecover(t *testing.T) {
	expectedErr := errors.New("some error")
	var retrievedErr error
	afterFunc := func(step pipeline.Step, ctx pipeline.Context, err error) error {
		retrievedErr = err
		return nil
	}
	step := new(mockStep)
	step.On("Run", &mockContext{}).Return(expectedErr).Once()

	lifecycleStep := pipeline.CreateAfterStepLifecycle(step, afterFunc)
	err := lifecycleStep.Run(&mockContext{})

	assert.Nil(t, err)
	assert.Equal(t, expectedErr, retrievedErr)
	step.AssertExpectations(t)
}

func TestLifecycleStep_GivenBeforeFunc_WhenRun_ThenFuncAreCalled(t *testing.T) {
	called := false
	beforeFunc := func(step pipeline.Step, ctx pipeline.Context) error {
		called = true
		return nil
	}
	step := new(mockStep)
	step.On("Run", &mockContext{}).Return(nil).Once()

	lifecycleStep := pipeline.CreateBeforeStepLifecycle(step, beforeFunc)
	err := lifecycleStep.Run(&mockContext{})

	assert.Nil(t, err)
	assert.True(t, called)
	step.AssertExpectations(t)
}

func TestLifecycleStep_GivenBeforeFuncReturningError_WhenRun_ThenErrorIsReturned(t *testing.T) {
	expectedErr := errors.New("some error")
	beforeFunc := func(step pipeline.Step, ctx pipeline.Context) error {
		return expectedErr
	}
	step := new(mockStep)

	lifecycleStep := pipeline.CreateBeforeStepLifecycle(step, beforeFunc)
	err := lifecycleStep.Run(&mockContext{})

	assert.Equal(t, expectedErr, err)
	step.AssertExpectations(t)
}

func TestLifecycleStep_GivenBeforeFuncReturningError_WhenRun_ThenStepAndAfterAreNotRan(t *testing.T) {
	expectedErr := errors.New("some error")
	called := false
	beforeFunc := func(step pipeline.Step, ctx pipeline.Context) error {
		return expectedErr
	}
	afterFunc := func(step pipeline.Step, ctx pipeline.Context, err error) error {
		called = true
		return nil
	}
	step := new(mockStep)

	lifecycleStep := pipeline.CreateStepLifecycle(step, beforeFunc, afterFunc)
	err := lifecycleStep.Run(&mockContext{})

	assert.Equal(t, expectedErr, err)
	assert.False(t, called)
	step.AssertExpectations(t)
}

func TestLifecycleStep_GivenAStep_WhenRun_ThenStepIsRun(t *testing.T) {
	expectedErr := errors.New("some error")
	beforeFunc := func(step pipeline.Step, ctx pipeline.Context) error {
		return nil
	}
	afterFunc := func(step pipeline.Step, ctx pipeline.Context, err error) error {
		return err
	}
	step := new(mockStep)
	step.On("Run", &mockContext{}).Return(expectedErr).Once()

	lifecycleStep := pipeline.CreateStepLifecycle(step, beforeFunc, afterFunc)
	err := lifecycleStep.Run(&mockContext{})

	assert.Equal(t, expectedErr, err)
	step.AssertExpectations(t)
}

func TestLifecycleStep_GivenAStep_WhenNamed_ThenStepIsDelegated(t *testing.T) {
	expectedName := "mock step"
	step := new(mockStep)
	step.On("Name").Return(expectedName).Once()

	lifecycleStep := pipeline.CreateStepLifecycle(step, func(step pipeline.Step, ctx pipeline.Context) error {
		return nil
	}, func(step pipeline.Step, ctx pipeline.Context, err error) error {
		return err
	})

	assert.Equal(t, expectedName, lifecycleStep.Name())
	step.AssertExpectations(t)
}

func TestLifecycleStep_GivenComposition_WhenRun_ThenCompositionBehavesAsAnArray(t *testing.T) {
	var callings []string
	before := func(step pipeline.Step, ctx pipeline.Context) error {
		callings = append(callings, "before")
		return nil
	}
	after := func(step pipeline.Step, ctx pipeline.Context, err error) error {
		callings = append(callings, "after")
		return err
	}

	step := new(mockStep)
	step.On("Run", &mockContext{}).Run(func(args mock.Arguments) {
		callings = append(callings, "step")
	}).Return(nil).Once()

	lifecycleStep := pipeline.CreateStepLifecycle(step, before, after)
	lifecycleStep = pipeline.CreateBeforeStepLifecycle(lifecycleStep, before)
	lifecycleStep = pipeline.CreateAfterStepLifecycle(lifecycleStep, after)
	lifecycleStep = pipeline.CreateBeforeStepLifecycle(lifecycleStep, before)

	err := lifecycleStep.Run(&mockContext{})

	assert.Nil(t, err)
	assert.Len(t, callings, 6)
	assert.Equal(t, []string{"before", "before", "before", "step", "after", "after"}, callings)
	step.AssertExpectations(t)
}
