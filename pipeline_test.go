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

var (
	stepMux      = sync.Mutex{}
	containerMux = sync.Mutex{}
)

type (
	mockContainer[T any] struct {
		mock.Mock
	}

	SimpleExecutor[T any] struct{}
)

func (m *mockContainer[T]) Draw(graph pipeline.GraphDiagram) {
	_ = m.Called(graph)
}

func (m *mockContainer[T]) Visit(ex pipeline.Executor[T], ctx T) error {
	args := m.Called(ex, ctx)

	return args.Error(0)
}

func (s SimpleExecutor[T]) Run(runnable pipeline.Step[T], in T) error {
	return runnable.Run(in)
}

func NewStep(data int, arr **[]int) pipeline.Step[int] {
	return pipeline.NewStep("", func(t int) error {
		stepMux.Lock()
		tmp := append(**arr, data)
		*arr = &tmp
		stepMux.Unlock()
		time.Sleep(time.Duration(100/(data+1)) * time.Millisecond) // Force a trap / yield
		return nil
	})
}

func NewContainer(data int, arr **[]int) pipeline.Container[int] {
	container := new(mockContainer[int])
	container.On("Visit", SimpleExecutor[int]{}, 1).Run(func(args mock.Arguments) {
		containerMux.Lock()
		tmp := append(**arr, data)
		*arr = &tmp
		containerMux.Unlock()
		time.Sleep(5 * time.Millisecond) // Force a possible trap / yield
	}).Return(nil).Once()

	return container
}

func TestPipeline_GivenAPipeline_WhenRunning_TheContainerIsRan(t *testing.T) {
	expectedErr := errors.New("error")
	pipe := pipeline.NewClient[interface{}](SimpleExecutor[interface{}]{})

	container := new(mockContainer[interface{}])
	container.On("Visit", SimpleExecutor[interface{}]{}, mock.Anything).Return(expectedErr).Once()

	err := pipe.Run(container, 1)

	assert.Equal(t, expectedErr, err)
	container.AssertExpectations(t)
}

func TestPipeline_GivenAPipeline_WhenRunning_TheContainerIsRanWithTheSameContext(t *testing.T) {
	expectedErr := errors.New("error")
	pipe := pipeline.NewClient[interface{}](SimpleExecutor[interface{}]{})
	v := 1
	container := new(mockContainer[interface{}])
	container.On("Visit", SimpleExecutor[interface{}]{}, mock.Anything).Run(func(args mock.Arguments) {
		*args.Get(1).(*int) = 2
	}).Return(expectedErr).Once()

	err := pipe.Run(container, &v)

	assert.Equal(t, expectedErr, err)
	assert.Equal(t, 2, v)
	container.AssertExpectations(t)
}
