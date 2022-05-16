package pipeline_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/saantiaguilera/go-pipeline"
)

type (
	mockStep[I, O any] struct {
		mock.Mock
	}
)

var (
	stepMux = sync.Mutex{}
)

func (m *mockStep[I, O]) ID() string {
	args := m.Called()
	return args.String(0)
}

func (m *mockStep[I, O]) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *mockStep[I, O]) Run(ctx context.Context, in I) (O, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return *new(O), args.Error(1)
	}
	return args.Get(0).(O), args.Error(1)
}

func (m *mockStep[I, O]) Draw(g pipeline.Graph) {
	_ = m.Called(g)
}

func NewStep[I, O any](data int, arr **[]int) pipeline.Step[I, O] {
	step := new(mockStep[I, O])
	step.On("Run", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		stepMux.Lock()
		tmp := append(**arr, data)
		*arr = &tmp
		stepMux.Unlock()
		time.Sleep(time.Duration(100/(data+1)) * time.Millisecond) // Force a trap / yield
	}).Return(*new(O), nil).Once()
	return step
}

func TestUnitStep_GivenAName_WhenGettingItsName_ThenItsTheExpected(t *testing.T) {
	expectedName := "test_name"
	step := pipeline.NewUnitStep[any, any](expectedName, nil)

	name := step.Name()

	assert.Equal(t, expectedName, name)
}

func TestUnitStep_GivenARunFunc_WhenRunning_ThenItsCalled(t *testing.T) {
	called := false
	run := func(context.Context, any) (any, error) {
		called = true
		return nil, nil
	}
	step := pipeline.NewUnitStep("", run)

	_, _ = step.Run(context.Background(), 1)

	assert.True(t, called)
}

func TestUnitStep_GivenARunFuncThatErrors_WhenRunning_ThenErrorIsReturned(t *testing.T) {
	expectedErr := errors.New("some error")
	run := func(context.Context, any) (any, error) {
		return nil, expectedErr
	}
	step := pipeline.NewUnitStep("", run)

	_, err := step.Run(context.Background(), 1)

	assert.Equal(t, expectedErr, err)
}

func TestUnitStep_GivenOne_ThenHasID(t *testing.T) {
	step := pipeline.NewUnitStep[any, any]("", nil)

	id := step.ID()

	assert.NotEmpty(t, id)
}

func TestUnitStep_GivenACancelledContext_WhenRunning_ThenFailsWithoutRunning(t *testing.T) {
	called := false
	step := pipeline.NewUnitStep("", func(ctx context.Context, t any) (any, error) {
		called = true
		return nil, nil
	})
	ctx, canc := context.WithCancel(context.Background())
	canc()

	v, _ := step.Run(ctx, 1)

	assert.Nil(t, v)
	assert.False(t, called)
}
