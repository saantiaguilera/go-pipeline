package pipeline_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/saantiaguilera/go-pipeline"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

type mockStage struct {
	mock.Mock
}

func (m *mockStage) Run(executor pipeline.Executor) error {
	args := m.Called(executor)

	return args.Error(0)
}

type SimpleExecutor struct{}

func (s SimpleExecutor) Run(runnable pipeline.Runnable) error {
	return runnable.Run()
}

var stepMux = sync.Mutex{}

func createStep(data int, arr **[]int) pipeline.Step {
	step := new(mockStep)
	step.On("Run").Run(func(args mock.Arguments) {
		stepMux.Lock()
		tmp := append(**arr, data)
		*arr = &tmp
		stepMux.Unlock()
		time.Sleep(time.Duration(100/(data+1)) * time.Millisecond) // Force a trap / yield
	}).Return(nil).Once()

	return step
}

var stageMux = sync.Mutex{}

func createStage(data int, arr **[]int) pipeline.Stage {
	stage := new(mockStage)
	stage.On("Run", SimpleExecutor{}).Run(func(args mock.Arguments) {
		stageMux.Lock()
		tmp := append(**arr, data)
		*arr = &tmp
		stageMux.Unlock()
		time.Sleep(5 * time.Millisecond) // Force a possible trap / yield
	}).Return(nil).Once()

	return stage
}

func TestPipeline_GivenAPipelineAddingABeforeHook_WhenRunning_ThenTheBeforeHookIsCalled(t *testing.T) {
	called := false
	pipe := pipeline.CreatePipeline(SimpleExecutor{})
	pipe.AddBeforeRunHook(func(stage pipeline.Stage) error {
		called = true
		return nil
	})

	stage := new(mockStage)
	stage.On("Run", SimpleExecutor{}).Return(nil).Once()

	err := pipe.Run(stage)

	assert.Nil(t, err)
	assert.True(t, called)
	stage.AssertExpectations(t)
}

func TestPipeline_GivenAPipelineAddingAFailingBeforeHook_WhenRunning_ThenItFailsWithoutCallingTheStage(t *testing.T) {
	expectedError := errors.New("error")

	pipe := pipeline.CreatePipeline(SimpleExecutor{})
	pipe.AddBeforeRunHook(func(stage pipeline.Stage) error {
		return expectedError
	})

	stage := new(mockStage)

	err := pipe.Run(stage)

	assert.Equal(t, expectedError, err)
	stage.AssertExpectations(t)
}

func TestPipeline_GivenAPipelineAddingABeforeHookAndAFailingBeforeHook_WhenRunning_ThenItFailsWithoutCallingTheStageButCallingTheFirstHook(t *testing.T) {
	called := false
	expectedError := errors.New("error")

	pipe := pipeline.CreatePipeline(SimpleExecutor{})
	pipe.AddBeforeRunHook(func(stage pipeline.Stage) error {
		called = true
		return nil
	})
	pipe.AddBeforeRunHook(func(stage pipeline.Stage) error {
		return expectedError
	})

	stage := new(mockStage)

	err := pipe.Run(stage)

	assert.True(t, called)
	assert.Equal(t, expectedError, err)
	stage.AssertExpectations(t)
}

func TestPipeline_GivenAPipeline_WhenRunning_TheStageIsRan(t *testing.T) {
	expectedErr := errors.New("error")
	pipe := pipeline.CreatePipeline(SimpleExecutor{})

	stage := new(mockStage)
	stage.On("Run", SimpleExecutor{}).Return(expectedErr).Once()

	err := pipe.Run(stage)

	assert.Equal(t, expectedErr, err)
	stage.AssertExpectations(t)
}

func TestPipeline_GivenAPipelineAddingAnAfterHook_WhenRunning_TheHookIsCalled(t *testing.T) {
	called := false
	pipe := pipeline.CreatePipeline(SimpleExecutor{})
	pipe.AddAfterRunHook(func(stage pipeline.Stage, err error) error {
		called = true
		return nil
	})

	stage := new(mockStage)
	stage.On("Run", SimpleExecutor{}).Return(nil).Once()

	err := pipe.Run(stage)

	assert.Nil(t, err)
	assert.True(t, called)
	stage.AssertExpectations(t)
}

func TestPipeline_GivenAPipelineAddingAFailureAfterHook_WhenRunning_ThenItFails(t *testing.T) {
	expectedErr := errors.New("error")
	pipe := pipeline.CreatePipeline(SimpleExecutor{})
	pipe.AddAfterRunHook(func(stage pipeline.Stage, err error) error {
		return expectedErr
	})

	stage := new(mockStage)
	stage.On("Run", SimpleExecutor{}).Return(nil).Once()

	err := pipe.Run(stage)

	assert.Equal(t, expectedErr, err)
	stage.AssertExpectations(t)
}

func TestPipeline_GivenAPipelineAddingAnAfterHookAndAFailingAfterHook_WhenRunning_ThenItFailsCallingAllHooks(t *testing.T) {
	called := false
	expectedError := errors.New("error")

	pipe := pipeline.CreatePipeline(SimpleExecutor{})
	pipe.AddAfterRunHook(func(stage pipeline.Stage, err error) error {
		return err
	})
	pipe.AddAfterRunHook(func(stage pipeline.Stage, err error) error {
		called = true
		return expectedError
	})

	stage := new(mockStage)
	stage.On("Run", SimpleExecutor{}).Return(nil).Once()

	err := pipe.Run(stage)

	assert.True(t, called)
	assert.Equal(t, expectedError, err)
	stage.AssertExpectations(t)
}

func TestPipeline_GivenAPipelineAddingAFailureAfterHookAndARecoveringAfterHook_WhenRunning_ThenItItRecovers(t *testing.T) {
	expectedError := errors.New("error")

	pipe := pipeline.CreatePipeline(SimpleExecutor{})
	pipe.AddAfterRunHook(func(stage pipeline.Stage, err error) error {
		return expectedError // First fails
	})
	pipe.AddAfterRunHook(func(stage pipeline.Stage, err error) error {
		return nil // Second recovers
	})

	stage := new(mockStage)
	stage.On("Run", SimpleExecutor{}).Return(nil).Once()

	err := pipe.Run(stage)

	assert.Nil(t, err)
	stage.AssertExpectations(t)
}

func TestPipeline_GivenAPipelineAddingARecoverAfterHook_WhenRunningAndFailing_ThenItRecovers(t *testing.T) {
	expectedErr := errors.New("error")
	var recoveredErr error
	pipe := pipeline.CreatePipeline(SimpleExecutor{})
	pipe.AddAfterRunHook(func(stage pipeline.Stage, err error) error {
		recoveredErr = err
		return nil
	})

	stage := new(mockStage)
	stage.On("Run", SimpleExecutor{}).Return(expectedErr).Once()

	err := pipe.Run(stage)

	assert.Nil(t, err)
	assert.Equal(t, recoveredErr, expectedErr)
	stage.AssertExpectations(t)
}

func TestPipeline_GivenAPipelineWithBothHooks_WhenRunning_TheStageAndHooksAreCalled(t *testing.T) {
	beforeCalled := false
	afterCalled := false

	pipe := pipeline.CreatePipeline(SimpleExecutor{})
	pipe.AddBeforeRunHook(func(stage pipeline.Stage) error {
		beforeCalled = true
		return nil
	})
	pipe.AddAfterRunHook(func(stage pipeline.Stage, err error) error {
		afterCalled = true
		return nil
	})

	stage := new(mockStage)
	stage.On("Run", SimpleExecutor{}).Return(nil).Once()

	err := pipe.Run(stage)

	assert.Nil(t, err)
	assert.True(t, beforeCalled)
	assert.True(t, afterCalled)
	stage.AssertExpectations(t)
}
