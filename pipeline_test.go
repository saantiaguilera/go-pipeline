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

type mockPipeline struct {
	mock.Mock
}

func (m mockPipeline) Run(stage pipeline.Stage) error {
	args := m.Called(stage)
	return args.Error(0)
}

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

func (m *mockStage) Draw(graph pipeline.GraphDiagram) {
	_ = m.Called(graph)
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

func TestPipeline_GivenAPipeline_WhenRunning_TheStageIsRan(t *testing.T) {
	expectedErr := errors.New("error")
	pipe := pipeline.CreatePipeline(SimpleExecutor{})

	stage := new(mockStage)
	stage.On("Run", SimpleExecutor{}).Return(expectedErr).Once()

	err := pipe.Run(stage)

	assert.Equal(t, expectedErr, err)
	stage.AssertExpectations(t)
}
