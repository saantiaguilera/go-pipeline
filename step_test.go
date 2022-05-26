package pipeline_test

import (
	"context"
	"errors"
	"fmt"
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

	noopStep[T any] struct{}
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

func (s noopStep[T]) Draw(pipeline.Graph) {
	// nothing
}

func (s noopStep[T]) Run(_ context.Context, in T) (T, error) {
	return in, nil
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

// The following example shows a simple unit step that runs a unit of work
// with a given input and yields a result of a different type or an error
// depending on the execution
//
// This example uses dummy data to showcase as simple as possible this scenario.
//
// Note: we use several UnitStep to showcase as it allows us to
// easily run dummy code, but it could use any type of step you want
// as long as it implements pipeline.Step[I, O]
func ExampleUnitStep() {
	type (
		InData  any
		OutData any
	)
	step := pipeline.NewUnitStep(
		"do_something",
		func(ctx context.Context, in InData) (OutData, error) {
			// do something with input
			return OutData(in), nil
		},
	)

	out, err := step.Run(context.Background(), InData("example"))

	fmt.Println(out, err)
	// output: example <nil>
}

// Benchmark for traversing a unit step. This is simply used so that future changes can
// easily reflect how they affected the performance
//
// goos: darwin
// goarch: amd64
// pkg: github.com/saantiaguilera/go-pipeline
// cpu: Intel(R) Core(TM) i7-1068NG7 CPU @ 2.30GHz
// BenchmarkUnitStep-8   	 5496482	       168.5 ns/op	       0 B/op	       0 allocs/op
func BenchmarkUnitStep(b *testing.B) {
	var err error
	s := pipeline.NewUnitStep(
		"name",
		func(ctx context.Context, a any) (any, error) {
			return a, nil
		},
	)
	ctx := context.Background()
	in := 0

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		_, err = s.Run(ctx, in)
		b.StopTimer()

		if err != nil {
			b.Fail()
		}
	}
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
