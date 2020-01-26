package lifecycle_test

import (
	"errors"
	"github.com/saantiaguilera/go-pipeline/pkg"
	"github.com/saantiaguilera/go-pipeline/pkg/step/lifecycle"
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
	afterFunc := func(step pkg.Step, err error) error {
		called = true
		return nil
	}
	step := new(mockStep)
	step.On("Run").Return(nil).Once()

	lifecycleStep := lifecycle.CreateAfterStepLifecycle(step, afterFunc)
	err := lifecycleStep.Run()

	assert.Nil(t, err)
	assert.True(t, called)
	step.AssertExpectations(t)
}

func TestLifecycleStep_GivenAfterFuncErroring_WhenRun_ThenErrorIsReturned(t *testing.T) {
	expectedErr := errors.New("some error")
	afterFunc := func(step pkg.Step, err error) error {
		return expectedErr
	}
	step := new(mockStep)
	step.On("Run").Return(nil).Once()

	lifecycleStep := lifecycle.CreateAfterStepLifecycle(step, afterFunc)
	err := lifecycleStep.Run()

	assert.Equal(t, expectedErr, err)
	step.AssertExpectations(t)
}

func TestLifecycleStep_GivenAfterFuncRecoveringError_WhenRun_ThenFuncCanRecover(t *testing.T) {
	expectedErr := errors.New("some error")
	var retrievedErr error
	afterFunc := func(step pkg.Step, err error) error {
		retrievedErr = err
		return nil
	}
	step := new(mockStep)
	step.On("Run").Return(expectedErr).Once()

	lifecycleStep := lifecycle.CreateAfterStepLifecycle(step, afterFunc)
	err := lifecycleStep.Run()

	assert.Nil(t, err)
	assert.Equal(t, expectedErr, retrievedErr)
	step.AssertExpectations(t)
}

func TestLifecycleStep_GivenBeforeFunc_WhenRun_ThenFuncAreCalled(t *testing.T) {
	called := false
	beforeFunc := func(step pkg.Step) error {
		called = true
		return nil
	}
	step := new(mockStep)
	step.On("Run").Return(nil).Once()

	lifecycleStep := lifecycle.CreateBeforeStepLifecycle(step, beforeFunc)
	err := lifecycleStep.Run()

	assert.Nil(t, err)
	assert.True(t, called)
	step.AssertExpectations(t)
}

func TestLifecycleStep_GivenBeforeFuncReturningError_WhenRun_ThenErrorIsReturned(t *testing.T) {
	expectedErr := errors.New("some error")
	beforeFunc := func(step pkg.Step) error {
		return expectedErr
	}
	step := new(mockStep)

	lifecycleStep := lifecycle.CreateBeforeStepLifecycle(step, beforeFunc)
	err := lifecycleStep.Run()

	assert.Equal(t, expectedErr, err)
	step.AssertExpectations(t)
}

func TestLifecycleStep_GivenBeforeFuncReturningError_WhenRun_ThenStepAndAfterAreNotRan(t *testing.T) {
	expectedErr := errors.New("some error")
	called := false
	beforeFunc := func(step pkg.Step) error {
		return expectedErr
	}
	afterFunc := func(step pkg.Step, err error) error {
		called = true
		return nil
	}
	step := new(mockStep)

	lifecycleStep := lifecycle.CreateStepLifecycle(step, beforeFunc, afterFunc)
	err := lifecycleStep.Run()

	assert.Equal(t, expectedErr, err)
	assert.False(t, called)
	step.AssertExpectations(t)
}

func TestLifecycleStep_GivenAStep_WhenRun_ThenStepIsRun(t *testing.T) {
	expectedErr := errors.New("some error")
	beforeFunc := func(step pkg.Step) error {
		return nil
	}
	afterFunc := func(step pkg.Step, err error) error {
		return err
	}
	step := new(mockStep)
	step.On("Run").Return(expectedErr).Once()

	lifecycleStep := lifecycle.CreateStepLifecycle(step, beforeFunc, afterFunc)
	err := lifecycleStep.Run()

	assert.Equal(t, expectedErr, err)
	step.AssertExpectations(t)
}

func TestLifecycleStep_GivenAStep_WhenNamed_ThenStepIsDelegated(t *testing.T) {
	expectedName := "mock step"
	step := new(mockStep)
	step.On("Name").Return(expectedName).Once()

	lifecycleStep := lifecycle.CreateStepLifecycle(step, func(step pkg.Step) error {
		return nil
	}, func(step pkg.Step, err error) error {
		return err
	})

	assert.Equal(t, expectedName, lifecycleStep.Name())
	step.AssertExpectations(t)
}

func TestLifecycleStep_GivenComposition_WhenRun_ThenCompositionBehavesAsAnArray(t *testing.T) {
	var callings []string
	before := func(step pkg.Step) error {
		callings = append(callings, "before")
		return nil
	}
	after := func(step pkg.Step, err error) error {
		callings = append(callings, "after")
		return err
	}

	step := new(mockStep)
	step.On("Run").Run(func(args mock.Arguments) {
		callings = append(callings, "step")
	}).Return(nil).Once()

	lifecycleStep := lifecycle.CreateStepLifecycle(step, before, after)
	lifecycleStep = lifecycle.CreateBeforeStepLifecycle(lifecycleStep, before)
	lifecycleStep = lifecycle.CreateAfterStepLifecycle(lifecycleStep, after)
	lifecycleStep = lifecycle.CreateBeforeStepLifecycle(lifecycleStep, before)

	err := lifecycleStep.Run()

	assert.Nil(t, err)
	assert.Len(t, callings, 6)
	assert.Equal(t, []string{"before", "before", "before", "step", "after", "after"}, callings)
	step.AssertExpectations(t)
}

