package pipeline_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/saantiaguilera/go-pipeline"
)

// The following example shows a sequence between two steps were each of them is
// run sequentially with the output of the previous one as input
//
// This example uses dummy data to showcase as simple as possible this scenario.
//
// Note: we use several UnitStep to showcase as it allows us to
// easily run dummy code, but it could use any type of step you want
// as long as it implements pipeline.Step[I, O]
func ExampleSequentialStep() {
	type DriverID int
	type Driver any
	type NotificationReceipt any
	gd := pipeline.NewUnitStep(
		"get_driver",
		func(ctx context.Context, id DriverID) (Driver, error) {
			// do something with input
			return Driver(id), nil
		},
	)
	sn := pipeline.NewUnitStep(
		"send_notification",
		func(ctx context.Context, d Driver) (NotificationReceipt, error) {
			// do something with input
			return NotificationReceipt(25), nil
		},
	)

	pipe := pipeline.NewSequentialStep[DriverID, Driver, NotificationReceipt](gd, sn)

	out, err := pipe.Run(context.Background(), DriverID(1234))

	fmt.Println(out, err)
	// output:
	// 25 <nil>
}

// Benchmark for traversing a sequential step. This is simply used so that future changes can
// easily reflect how they affected the performance
//
// goos: darwin
// goarch: amd64
// pkg: github.com/saantiaguilera/go-pipeline
// cpu: Intel(R) Core(TM) i7-1068NG7 CPU @ 2.30GHz
// BenchmarkSequentialStep-8   	 7136264	       172.5 ns/op	       0 B/op	       0 allocs/op
func BenchmarkSequentialStep(b *testing.B) {
	var err error
	s := pipeline.NewSequentialStep[any, any, any](
		noopStep[any]{},
		noopStep[any]{},
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

func TestSequentialStep_GivenTwoSteps_WhenRun_ThenBehavesSequentially(t *testing.T) {
	start := new(mockStep[int, string])
	start.On("Run", mock.Anything, 1).Return("test", nil)
	end := new(mockStep[string, bool])
	end.On("Run", mock.Anything, "test").Return(true, nil)
	step := pipeline.NewSequentialStep[int, string, bool](start, end)

	v, err := step.Run(context.Background(), 1)

	assert.Nil(t, err)
	assert.True(t, v)
	mock.AssertExpectationsForObjects(t, start, end)
}

func TestSequentialStep_GivenStepsWithErrors_WhenRun_ThenStepsAreHaltedAfterError(t *testing.T) {
	expectedErr := errors.New("oops")
	start := new(mockStep[int, string])
	start.On("Run", mock.Anything, 1).Return("", expectedErr)
	end := new(mockStep[string, bool])
	step := pipeline.NewSequentialStep[int, string, bool](start, end)

	v, err := step.Run(context.Background(), 1)

	assert.Equal(t, expectedErr, err)
	assert.Empty(t, v)
	mock.AssertExpectationsForObjects(t, start, end)
}

func TestSequentialStep_GivenAGraphToDraw_WhenDrawn_ThenStepsAreAddedAsActivitiesByTheirNames(t *testing.T) {
	mockGraph := new(mockGraph)
	mockGraph.On("AddActivity", "1").Once()
	mockGraph.On("AddActivity", "2").Once()
	first := pipeline.NewUnitStep[any, any]("1", nil)
	second := pipeline.NewUnitStep[any, any]("2", nil)
	initStep := pipeline.NewSequentialStep[any, any, any](first, second)

	initStep.Draw(mockGraph)

	mockGraph.AssertExpectations(t)
}
