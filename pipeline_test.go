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

func (m mockPipeline) Run(stage pipeline.Stage, ctx pipeline.Context) error {
	args := m.Called(stage, ctx)
	return args.Error(0)
}

type mockStep struct {
	mock.Mock
}

func (m *mockStep) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *mockStep) Run(ctx pipeline.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

type mockStage struct {
	mock.Mock
}

func (m *mockStage) Draw(graph pipeline.GraphDiagram) {
	_ = m.Called(graph)
}

func (m *mockStage) Run(executor pipeline.Executor, ctx pipeline.Context) error {
	args := m.Called(executor, ctx)

	return args.Error(0)
}

type SimpleExecutor struct{}

func (s SimpleExecutor) Run(runnable pipeline.Runnable, ctx pipeline.Context) error {
	return runnable.Run(ctx)
}

var stepMux = sync.Mutex{}

func createStep(data int, arr **[]int) pipeline.Step {
	step := new(mockStep)
	step.On("Run", &mockContext{}).Run(func(args mock.Arguments) {
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
	stage.On("Run", SimpleExecutor{}, &mockContext{}).Run(func(args mock.Arguments) {
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
	stage.On("Run", SimpleExecutor{}, mock.Anything).Return(expectedErr).Once()

	err := pipe.Run(stage, &mockContext{})

	assert.Equal(t, expectedErr, err)
	stage.AssertExpectations(t)
}

func TestPipeline_GivenAPipeline_WhenRunning_TheStageIsRanWithTheSameContext(t *testing.T) {
	var tag pipeline.Tag = "tag"
	expectedErr := errors.New("error")
	pipe := pipeline.CreatePipeline(SimpleExecutor{})
	ctx := new(mockContext)
	ctx.On("Set", tag, "value").Once()

	stage := new(mockStage)
	stage.On("Run", SimpleExecutor{}, ctx).Run(func(args mock.Arguments) {
		args.Get(1).(pipeline.Context).Set(tag, "value")
	}).Return(expectedErr).Once()

	err := pipe.Run(stage, ctx)

	assert.Equal(t, expectedErr, err)
	ctx.AssertExpectations(t)
	stage.AssertExpectations(t)
}
