package pipeline_test

import (
	"errors"
	"testing"

	"github.com/saantiaguilera/go-pipeline"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLifecyclePipeline_GivenAfterFunc_WhenRun_ThenFuncAreCalled(t *testing.T) {
	called := false
	afterFunc := func(stage pipeline.Stage, ctx pipeline.Context, err error) error {
		called = true
		return nil
	}
	stage := &mockStage{}

	mockPipeline := new(mockPipeline)
	mockPipeline.On("Run", stage, &mockContext{}).Return(nil).Once()

	lifecyclePipeline := pipeline.CreateAfterPipelineLifecycle(mockPipeline, afterFunc)
	err := lifecyclePipeline.Run(stage, &mockContext{})

	assert.Nil(t, err)
	assert.True(t, called)
	mockPipeline.AssertExpectations(t)
}

func TestLifecyclePipeline_GivenAfterFuncErroring_WhenRun_ThenErrorIsReturned(t *testing.T) {
	expectedErr := errors.New("some error")
	afterFunc := func(stage pipeline.Stage, ctx pipeline.Context, err error) error {
		return expectedErr
	}
	stage := &mockStage{}

	mockPipeline := new(mockPipeline)
	mockPipeline.On("Run", stage, &mockContext{}).Return(nil).Once()

	lifecyclePipeline := pipeline.CreateAfterPipelineLifecycle(mockPipeline, afterFunc)
	err := lifecyclePipeline.Run(stage, &mockContext{})

	assert.Equal(t, expectedErr, err)
	mockPipeline.AssertExpectations(t)
}

func TestLifecyclePipeline_GivenAfterFuncRecoveringError_WhenRun_ThenFuncCanRecover(t *testing.T) {
	expectedErr := errors.New("some error")
	var retrievedErr error
	afterFunc := func(stage pipeline.Stage, ctx pipeline.Context, err error) error {
		retrievedErr = err
		return nil
	}
	stage := &mockStage{}

	mockPipeline := new(mockPipeline)
	mockPipeline.On("Run", stage, &mockContext{}).Return(expectedErr).Once()

	lifecyclePipeline := pipeline.CreateAfterPipelineLifecycle(mockPipeline, afterFunc)
	err := lifecyclePipeline.Run(stage, &mockContext{})

	assert.Nil(t, err)
	assert.Equal(t, expectedErr, retrievedErr)
	mockPipeline.AssertExpectations(t)
}

func TestLifecyclePipeline_GivenBeforeFunc_WhenRun_ThenFuncAreCalled(t *testing.T) {
	called := false
	beforeFunc := func(stage pipeline.Stage, ctx pipeline.Context) error {
		called = true
		return nil
	}
	stage := &mockStage{}

	mockPipeline := new(mockPipeline)
	mockPipeline.On("Run", stage, &mockContext{}).Return(nil).Once()

	lifecyclePipeline := pipeline.CreateBeforePipelineLifecycle(mockPipeline, beforeFunc)
	err := lifecyclePipeline.Run(stage, &mockContext{})

	assert.Nil(t, err)
	assert.True(t, called)
	mockPipeline.AssertExpectations(t)
}

func TestLifecyclePipeline_GivenBeforeFuncReturningError_WhenRun_ThenErrorIsReturned(t *testing.T) {
	expectedErr := errors.New("some error")
	beforeFunc := func(stage pipeline.Stage, ctx pipeline.Context) error {
		return expectedErr
	}
	stage := &mockStage{}

	mockPipeline := new(mockPipeline)

	lifecyclePipeline := pipeline.CreateBeforePipelineLifecycle(mockPipeline, beforeFunc)
	err := lifecyclePipeline.Run(stage, &mockContext{})

	assert.Equal(t, expectedErr, err)
	mockPipeline.AssertExpectations(t)
}

func TestLifecyclePipeline_GivenBeforeFuncReturningError_WhenRun_ThenPipelineAndAfterAreNotRan(t *testing.T) {
	expectedErr := errors.New("some error")
	called := false
	beforeFunc := func(stage pipeline.Stage, ctx pipeline.Context) error {
		return expectedErr
	}
	afterFunc := func(stage pipeline.Stage, ctx pipeline.Context, err error) error {
		called = true
		return nil
	}
	stage := &mockStage{}

	mockPipeline := new(mockPipeline)

	lifecyclePipeline := pipeline.CreatePipelineLifecycle(mockPipeline, beforeFunc, afterFunc)
	err := lifecyclePipeline.Run(stage, &mockContext{})

	assert.Equal(t, expectedErr, err)
	assert.False(t, called)
	mockPipeline.AssertExpectations(t)
}

func TestLifecyclePipeline_GivenAPipeline_WhenRun_ThenPipelineIsRun(t *testing.T) {
	expectedErr := errors.New("some error")
	beforeFunc := func(stage pipeline.Stage, ctx pipeline.Context) error {
		return nil
	}
	afterFunc := func(stage pipeline.Stage, ctx pipeline.Context, err error) error {
		return err
	}
	stage := &mockStage{}

	mockPipeline := new(mockPipeline)
	mockPipeline.On("Run", stage, &mockContext{}).Return(expectedErr).Once()

	lifecyclePipeline := pipeline.CreatePipelineLifecycle(mockPipeline, beforeFunc, afterFunc)
	err := lifecyclePipeline.Run(stage, &mockContext{})

	assert.Equal(t, expectedErr, err)
	mockPipeline.AssertExpectations(t)
}

func TestLifecyclePipeline_GivenComposition_WhenRun_ThenCompositionBehavesAsAnArray(t *testing.T) {
	var callings []string
	before := func(stage pipeline.Stage, ctx pipeline.Context) error {
		callings = append(callings, "before")
		return nil
	}
	after := func(stage pipeline.Stage, ctx pipeline.Context, err error) error {
		callings = append(callings, "after")
		return err
	}
	stage := &mockStage{}

	mockPipeline := new(mockPipeline)
	mockPipeline.On("Run", stage, &mockContext{}).Run(func(args mock.Arguments) {
		callings = append(callings, "pipeline")
	}).Return(nil).Once()

	lifecyclePipeline := pipeline.CreatePipelineLifecycle(mockPipeline, before, after)
	lifecyclePipeline = pipeline.CreateBeforePipelineLifecycle(lifecyclePipeline, before)
	lifecyclePipeline = pipeline.CreateAfterPipelineLifecycle(lifecyclePipeline, after)
	lifecyclePipeline = pipeline.CreateBeforePipelineLifecycle(lifecyclePipeline, before)

	err := lifecyclePipeline.Run(stage, &mockContext{})

	assert.Nil(t, err)
	assert.Len(t, callings, 6)
	assert.Equal(t, []string{"before", "before", "before", "pipeline", "after", "after"}, callings)
	mockPipeline.AssertExpectations(t)
}
