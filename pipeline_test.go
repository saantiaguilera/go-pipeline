package pipeline_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/saantiaguilera/go-pipeline"
)

type mockStep[T any] struct {
	mock.Mock
}

func (m *mockStep[T]) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *mockStep[T]) Run(in T) error {
	args := m.Called(in)
	return args.Error(0)
}

type mockStage[T any] struct {
	mock.Mock
}

func (m *mockStage[T]) Draw(graph pipeline.GraphDiagram) {
	_ = m.Called(graph)
}

func (m *mockStage[T]) Run(executor pipeline.Executor[T], ctx T) error {
	args := m.Called(executor, ctx)

	return args.Error(0)
}

type SimpleExecutor[T any] struct{}

func (s SimpleExecutor[T]) Run(runnable pipeline.Step[T], in T) error {
	return runnable.Run(in)
}

var stepMux = sync.Mutex{}

func NewStep(data int, arr **[]int) pipeline.Step[int] {
	step := new(mockStep[int])
	step.On("Run", 1).Run(func(args mock.Arguments) {
		stepMux.Lock()
		tmp := append(**arr, data)
		*arr = &tmp
		stepMux.Unlock()
		time.Sleep(time.Duration(100/(data+1)) * time.Millisecond) // Force a trap / yield
	}).Return(nil).Once()

	return step
}

var stageMux = sync.Mutex{}

func NewStage(data int, arr **[]int) pipeline.Stage[int] {
	stage := new(mockStage[int])
	stage.On("Run", SimpleExecutor[int]{}, 1).Run(func(args mock.Arguments) {
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
	pipe := pipeline.NewClient[interface{}](SimpleExecutor[interface{}]{})

	stage := new(mockStage[interface{}])
	stage.On("Run", SimpleExecutor[interface{}]{}, mock.Anything).Return(expectedErr).Once()

	err := pipe.Run(stage, 1)

	assert.Equal(t, expectedErr, err)
	stage.AssertExpectations(t)
}

func TestPipeline_GivenAPipeline_WhenRunning_TheStageIsRanWithTheSameContext(t *testing.T) {
	expectedErr := errors.New("error")
	pipe := pipeline.NewClient[interface{}](SimpleExecutor[interface{}]{})
	v := 1
	stage := new(mockStage[interface{}])
	stage.On("Run", SimpleExecutor[interface{}]{}, mock.Anything).Run(func(args mock.Arguments) {
		*args.Get(1).(*int) = 2
	}).Return(expectedErr).Once()

	err := pipe.Run(stage, &v)

	assert.Equal(t, expectedErr, err)
	assert.Equal(t, 2, v)
	stage.AssertExpectations(t)
}
