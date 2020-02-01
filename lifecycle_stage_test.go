package pipeline_test

import (
	"errors"
	"testing"

	"github.com/saantiaguilera/go-pipeline"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLifecycleStage_GivenAfterFunc_WhenRun_ThenFuncAreCalled(t *testing.T) {
	called := false
	afterFunc := func(stage pipeline.Stage, err error) error {
		called = true
		return nil
	}
	stage := new(mockStage)
	stage.On("Run", SimpleExecutor{}).Return(nil).Once()

	lifecycleStage := pipeline.CreateAfterStageLifecycle(stage, afterFunc)
	err := lifecycleStage.Run(SimpleExecutor{})

	assert.Nil(t, err)
	assert.True(t, called)
	stage.AssertExpectations(t)
}

func TestLifecycleStage_GivenAfterFuncErroring_WhenRun_ThenErrorIsReturned(t *testing.T) {
	expectedErr := errors.New("some error")
	afterFunc := func(stage pipeline.Stage, err error) error {
		return expectedErr
	}
	stage := new(mockStage)
	stage.On("Run", SimpleExecutor{}).Return(nil).Once()

	lifecycleStage := pipeline.CreateAfterStageLifecycle(stage, afterFunc)
	err := lifecycleStage.Run(SimpleExecutor{})

	assert.Equal(t, expectedErr, err)
	stage.AssertExpectations(t)
}

func TestLifecycleStage_GivenAfterFuncRecoveringError_WhenRun_ThenFuncCanRecover(t *testing.T) {
	expectedErr := errors.New("some error")
	var retrievedErr error
	afterFunc := func(stage pipeline.Stage, err error) error {
		retrievedErr = err
		return nil
	}
	stage := new(mockStage)
	stage.On("Run", SimpleExecutor{}).Return(expectedErr).Once()

	lifecycleStage := pipeline.CreateAfterStageLifecycle(stage, afterFunc)
	err := lifecycleStage.Run(SimpleExecutor{})

	assert.Nil(t, err)
	assert.Equal(t, expectedErr, retrievedErr)
	stage.AssertExpectations(t)
}

func TestLifecycleStage_GivenBeforeFunc_WhenRun_ThenFuncAreCalled(t *testing.T) {
	called := false
	beforeFunc := func(stage pipeline.Stage) error {
		called = true
		return nil
	}
	stage := new(mockStage)
	stage.On("Run", SimpleExecutor{}).Return(nil).Once()

	lifecycleStage := pipeline.CreateBeforeStageLifecycle(stage, beforeFunc)
	err := lifecycleStage.Run(SimpleExecutor{})

	assert.Nil(t, err)
	assert.True(t, called)
	stage.AssertExpectations(t)
}

func TestLifecycleStage_GivenBeforeFuncReturningError_WhenRun_ThenErrorIsReturned(t *testing.T) {
	expectedErr := errors.New("some error")
	beforeFunc := func(stage pipeline.Stage) error {
		return expectedErr
	}
	stage := new(mockStage)

	lifecycleStage := pipeline.CreateBeforeStageLifecycle(stage, beforeFunc)
	err := lifecycleStage.Run(SimpleExecutor{})

	assert.Equal(t, expectedErr, err)
	stage.AssertExpectations(t)
}

func TestLifecycleStage_GivenBeforeFuncReturningError_WhenRun_ThenStageAndAfterAreNotRan(t *testing.T) {
	expectedErr := errors.New("some error")
	called := false
	beforeFunc := func(stage pipeline.Stage) error {
		return expectedErr
	}
	afterFunc := func(stage pipeline.Stage, err error) error {
		called = true
		return nil
	}
	stage := new(mockStage)

	lifecycleStage := pipeline.CreateStageLifecycle(stage, beforeFunc, afterFunc)
	err := lifecycleStage.Run(SimpleExecutor{})

	assert.Equal(t, expectedErr, err)
	assert.False(t, called)
	stage.AssertExpectations(t)
}

func TestLifecycleStage_GivenAStage_WhenRun_ThenStageIsRun(t *testing.T) {
	expectedErr := errors.New("some error")
	beforeFunc := func(stage pipeline.Stage) error {
		return nil
	}
	afterFunc := func(stage pipeline.Stage, err error) error {
		return err
	}
	stage := new(mockStage)
	stage.On("Run", SimpleExecutor{}).Return(expectedErr).Once()

	lifecycleStage := pipeline.CreateStageLifecycle(stage, beforeFunc, afterFunc)
	err := lifecycleStage.Run(SimpleExecutor{})

	assert.Equal(t, expectedErr, err)
	stage.AssertExpectations(t)
}

func TestLifecycleStage_GivenComposition_WhenRun_ThenCompositionBehavesAsAnArray(t *testing.T) {
	var callings []string
	before := func(stage pipeline.Stage) error {
		callings = append(callings, "before")
		return nil
	}
	after := func(stage pipeline.Stage, err error) error {
		callings = append(callings, "after")
		return err
	}

	stage := new(mockStage)
	stage.On("Run", SimpleExecutor{}).Run(func(args mock.Arguments) {
		callings = append(callings, "stage")
	}).Return(nil).Once()

	lifecycleStage := pipeline.CreateStageLifecycle(stage, before, after)
	lifecycleStage = pipeline.CreateBeforeStageLifecycle(lifecycleStage, before)
	lifecycleStage = pipeline.CreateAfterStageLifecycle(lifecycleStage, after)
	lifecycleStage = pipeline.CreateBeforeStageLifecycle(lifecycleStage, before)

	err := lifecycleStage.Run(SimpleExecutor{})

	assert.Nil(t, err)
	assert.Len(t, callings, 6)
	assert.Equal(t, []string{"before", "before", "before", "stage", "after", "after"}, callings)
	stage.AssertExpectations(t)
}

func TestLifecycleStage_GivenAGraphToDraw_WhenDrawn_ThenDelegatesToInnerStage(t *testing.T) {
	before := func(stage pipeline.Stage) error {
		return nil
	}

	mockGraphDiagram := new(mockGraphDiagram)

	stage := new(mockStage)
	stage.On("Draw", mockGraphDiagram).Once()

	lifecycleStage := pipeline.CreateBeforeStageLifecycle(stage, before)
	lifecycleStage.Draw(mockGraphDiagram)

	stage.AssertExpectations(t)
	mockGraphDiagram.AssertExpectations(t)
}
