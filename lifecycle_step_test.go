package pipeline_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/saantiaguilera/go-pipeline"
)

func TestLifecycleStep_GivenAfterFunc_WhenRun_ThenFuncAreCalled(t *testing.T) {
	called := false
	afterFunc := func(step pipeline.Step[interface{}], in interface{}, err error) error {
		called = true
		return nil
	}
	step := new(mockStep[interface{}])
	step.On("Run", 1).Return(nil).Once()

	lifecycleStep := pipeline.NewAfterStepLifecycle[interface{}](step, afterFunc)
	err := lifecycleStep.Run(1)

	assert.Nil(t, err)
	assert.True(t, called)
	step.AssertExpectations(t)
}

func TestLifecycleStep_GivenAfterFuncErroring_WhenRun_ThenErrorIsReturned(t *testing.T) {
	expectedErr := errors.New("some error")
	afterFunc := func(step pipeline.Step[interface{}], in interface{}, err error) error {
		return expectedErr
	}
	step := new(mockStep[interface{}])
	step.On("Run", 1).Return(nil).Once()

	lifecycleStep := pipeline.NewAfterStepLifecycle[interface{}](step, afterFunc)
	err := lifecycleStep.Run(1)

	assert.Equal(t, expectedErr, err)
	step.AssertExpectations(t)
}

func TestLifecycleStep_GivenAfterFuncRecoveringError_WhenRun_ThenFuncCanRecover(t *testing.T) {
	expectedErr := errors.New("some error")
	var retrievedErr error
	afterFunc := func(step pipeline.Step[interface{}], in interface{}, err error) error {
		retrievedErr = err
		return nil
	}
	step := new(mockStep[interface{}])
	step.On("Run", 1).Return(expectedErr).Once()

	lifecycleStep := pipeline.NewAfterStepLifecycle[interface{}](step, afterFunc)
	err := lifecycleStep.Run(1)

	assert.Nil(t, err)
	assert.Equal(t, expectedErr, retrievedErr)
	step.AssertExpectations(t)
}

func TestLifecycleStep_GivenBeforeFunc_WhenRun_ThenFuncAreCalled(t *testing.T) {
	called := false
	beforeFunc := func(step pipeline.Step[interface{}], in interface{}) error {
		called = true
		return nil
	}
	step := new(mockStep[interface{}])
	step.On("Run", 1).Return(nil).Once()

	lifecycleStep := pipeline.NewBeforeStepLifecycle[interface{}](step, beforeFunc)
	err := lifecycleStep.Run(1)

	assert.Nil(t, err)
	assert.True(t, called)
	step.AssertExpectations(t)
}

func TestLifecycleStep_GivenBeforeFuncReturningError_WhenRun_ThenErrorIsReturned(t *testing.T) {
	expectedErr := errors.New("some error")
	beforeFunc := func(step pipeline.Step[interface{}], in interface{}) error {
		return expectedErr
	}
	step := new(mockStep[interface{}])

	lifecycleStep := pipeline.NewBeforeStepLifecycle[interface{}](step, beforeFunc)
	err := lifecycleStep.Run(1)

	assert.Equal(t, expectedErr, err)
	step.AssertExpectations(t)
}

func TestLifecycleStep_GivenBeforeFuncReturningError_WhenRun_ThenStepAndAfterAreNotRan(t *testing.T) {
	expectedErr := errors.New("some error")
	called := false
	beforeFunc := func(step pipeline.Step[interface{}], in interface{}) error {
		return expectedErr
	}
	afterFunc := func(step pipeline.Step[interface{}], in interface{}, err error) error {
		called = true
		return nil
	}
	step := new(mockStep[interface{}])

	lifecycleStep := pipeline.NewStepLifecycle[interface{}](step, beforeFunc, afterFunc)
	err := lifecycleStep.Run(1)

	assert.Equal(t, expectedErr, err)
	assert.False(t, called)
	step.AssertExpectations(t)
}

func TestLifecycleStep_GivenAStep_WhenRun_ThenStepIsRun(t *testing.T) {
	expectedErr := errors.New("some error")
	beforeFunc := func(step pipeline.Step[interface{}], in interface{}) error {
		return nil
	}
	afterFunc := func(step pipeline.Step[interface{}], in interface{}, err error) error {
		return err
	}
	step := new(mockStep[interface{}])
	step.On("Run", 1).Return(expectedErr).Once()

	lifecycleStep := pipeline.NewStepLifecycle[interface{}](step, beforeFunc, afterFunc)
	err := lifecycleStep.Run(1)

	assert.Equal(t, expectedErr, err)
	step.AssertExpectations(t)
}

func TestLifecycleStep_GivenAStep_WhenNamed_ThenStepIsDelegated(t *testing.T) {
	expectedName := "mock step"
	step := new(mockStep[interface{}])
	step.On("Name").Return(expectedName).Once()

	lifecycleStep := pipeline.NewStepLifecycle[interface{}](step, func(step pipeline.Step[interface{}], in interface{}) error {
		return nil
	}, func(step pipeline.Step[interface{}], in interface{}, err error) error {
		return err
	})

	assert.Equal(t, expectedName, lifecycleStep.Name())
	step.AssertExpectations(t)
}

func TestLifecycleStep_GivenComposition_WhenRun_ThenCompositionBehavesAsAnArray(t *testing.T) {
	var callings []string
	before := func(step pipeline.Step[interface{}], in interface{}) error {
		callings = append(callings, "before")
		return nil
	}
	after := func(step pipeline.Step[interface{}], in interface{}, err error) error {
		callings = append(callings, "after")
		return err
	}

	step := new(mockStep[interface{}])
	step.On("Run", 1).Run(func(args mock.Arguments) {
		callings = append(callings, "step")
	}).Return(nil).Once()

	lifecycleStep := pipeline.NewStepLifecycle[interface{}](step, before, after)
	lifecycleStep = pipeline.NewBeforeStepLifecycle(lifecycleStep, before)
	lifecycleStep = pipeline.NewAfterStepLifecycle(lifecycleStep, after)
	lifecycleStep = pipeline.NewBeforeStepLifecycle(lifecycleStep, before)

	err := lifecycleStep.Run(1)

	assert.Nil(t, err)
	assert.Len(t, callings, 6)
	assert.Equal(t, []string{"before", "before", "before", "step", "after", "after"}, callings)
	step.AssertExpectations(t)
}
