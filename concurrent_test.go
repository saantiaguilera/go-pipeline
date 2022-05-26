package pipeline_test

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/saantiaguilera/go-pipeline"
)

type (
	count32 int32
)

func (c *count32) increment() int32 {
	return atomic.AddInt32((*int32)(c), 1)
}

// This example shows a same step that is run many times concurrently.
//
// The example uses dummy data to better showcase the immutability of the graph and step
// since we can produce different values and later reduce them without needing to
// take into account goroutines, mutexes, waitgroups to synchronize data at the end
//
// Note: we use several UnitStep to showcase as it allows us to
// easily run dummy code, but it could use any type of step you want
// as long as it implements pipeline.Step[I, O]
func ExampleConcurrentStep_same() {
	step := pipeline.NewUnitStep(
		"half_increase",
		func(ctx context.Context, i int) (float32, error) {
			return float32(i) * 1.5, nil
		},
	)
	reduce := func(ctx context.Context, a, b float32) (float32, error) {
		return a + b, nil
	}
	ctx := context.Background()
	in := 1

	pipe := pipeline.NewConcurrentStep(
		[]pipeline.Step[int, float32]{
			step, step, step, step, step, step,
			step, step, step, step, step, step,
		},
		reduce,
	)

	out, err := pipe.Run(ctx, in)

	fmt.Println(out, err)
	// output:
	// 18 <nil>
}

// The example uses dummy data and simulates a specific (random) scenario
// were we need to get a resource that is created from two different data sources (and we can
// obtain them concurrently to improve the performance) and finally reduce it to a single
// output
//
// Note: we use several UnitStep to showcase as it allows us to
// easily run dummy code, but it could use any type of step you want
// as long as it implements pipeline.Step[I, O]
func ExampleConcurrentStep_different() {
	type DriverID int
	type Driver struct {
		Person  any
		Vehicle any
	}

	gp := pipeline.NewUnitStep(
		"get_person",
		func(ctx context.Context, i DriverID) (Driver, error) {
			// do something with input
			return Driver{
				Person: true,
			}, nil
		},
	)
	gv := pipeline.NewUnitStep(
		"get_vehicle",
		func(ctx context.Context, i DriverID) (Driver, error) {
			// do something with input
			return Driver{
				Vehicle: true,
			}, nil
		},
	)
	reduce := func(ctx context.Context, a, b Driver) (Driver, error) {
		if b.Person != nil {
			a.Person = b.Person
		}
		if b.Vehicle != nil {
			a.Vehicle = b.Vehicle
		}
		return a, nil
	}
	ctx := context.Background()
	in := DriverID(1)

	pipe := pipeline.NewConcurrentStep(
		[]pipeline.Step[DriverID, Driver]{gp, gv},
		reduce,
	)

	out, err := pipe.Run(ctx, in)

	fmt.Println(out, err)
	// output:
	// {true true} <nil>
}

// Benchmark for traversing a concurrent step. This is simply used so that future changes can
// easily reflect how they affected the performance
//
// goos: darwin
// goarch: amd64
// pkg: github.com/saantiaguilera/go-pipeline
// cpu: Intel(R) Core(TM) i7-1068NG7 CPU @ 2.30GHz
// BenchmarkConcurrentStep-8   	  559716	      1990 ns/op	     384 B/op	       4 allocs/op
func BenchmarkConcurrentStep(b *testing.B) {
	var err error
	s := pipeline.NewConcurrentStep(
		[]pipeline.Step[any, any]{noopStep[any]{}, noopStep[any]{}},
		func(ctx context.Context, a1, a2 any) (any, error) {
			return a1, nil
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

func TestConcurrentStep_GivenStepsWithoutErrors_WhenRun_ThenAllStepsAreRunConcurrently(t *testing.T) {
	arr := &[]int{}
	var expectedArr []int
	var steps []pipeline.Step[int, int]
	for i := 0; i < 100; i++ {
		steps = append(steps, NewStep[int, int](i, &arr))
		expectedArr = append(expectedArr, i)
	}

	step := pipeline.NewConcurrentStep(steps, func(ctx context.Context, a, b int) (int, error) {
		return 0, nil
	})

	v, err := step.Run(context.Background(), 1)

	assert.Nil(t, err)
	assert.NotEqual(t, expectedArr, *arr)
	assert.Equal(t, len(expectedArr), len(*arr))
	assert.Equal(t, 0, v)
}

func TestConcurrentStep_GivenStepsWithErrors_WhenRun_ThenAllStepsAreRun(t *testing.T) {
	expectedErr := errors.New("error")
	var times count32
	step := pipeline.NewUnitStep("", func(ctx context.Context, t int) (int, error) {
		times.increment()
		return t, expectedErr
	})
	cstep := pipeline.NewConcurrentStep(
		[]pipeline.Step[int, int]{
			step, step, step, step, step,
			step, step, step, step, step,
		},
		func(ctx context.Context, a, b int) (int, error) {
			return a + b, nil
		},
	)

	v, err := cstep.Run(context.Background(), 1)

	assert.Equal(t, expectedErr, err)
	assert.Equal(t, count32(10), times)
	assert.Equal(t, 0, v)
}

func TestConcurrentStep_GivenStepsWithValues_WhenRun_ThenReducesValues(t *testing.T) {
	step := pipeline.NewUnitStep("", func(ctx context.Context, t int) (int, error) {
		return t, nil
	})
	cstep := pipeline.NewConcurrentStep(
		[]pipeline.Step[int, int]{
			step, step, step, step, step,
			step, step, step, step, step,
		},
		func(ctx context.Context, a, b int) (int, error) {
			return a + b, nil
		},
	)

	v, err := cstep.Run(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, 10, v)
}

func TestConcurrentStep_GivenAnErrorReducing_WhenRun_ThenErrors(t *testing.T) {
	step := pipeline.NewUnitStep("", func(ctx context.Context, t int) (int, error) {
		return t, nil
	})
	cstep := pipeline.NewConcurrentStep(
		[]pipeline.Step[int, int]{
			step, step, step, step, step,
			step, step, step, step, step,
		},
		func(ctx context.Context, a, b int) (int, error) {
			return 0, errors.New("oops")
		},
	)

	v, err := cstep.Run(context.Background(), 1)

	assert.Error(t, err)
	assert.Zero(t, v)
}

func TestConcurrentStep_GivenAnEmptyConcurrentStep_WhenRun_ThenErrors(t *testing.T) {
	cstep := pipeline.NewConcurrentStep(
		[]pipeline.Step[int, int]{},
		func(ctx context.Context, a, b int) (int, error) {
			return 0, nil
		},
	)

	v, err := cstep.Run(context.Background(), 1)

	assert.Error(t, err)
	assert.Zero(t, v)
}

func TestConcurrentStep_GivenASingleConcurrentStep_WhenRun_ThenDoesntReduce(t *testing.T) {
	step := pipeline.NewUnitStep("", func(ctx context.Context, t int) (int, error) {
		return 1000, nil
	})
	cstep := pipeline.NewConcurrentStep(
		[]pipeline.Step[int, int]{step},
		func(ctx context.Context, a, b int) (int, error) {
			return 0, nil
		},
	)

	v, err := cstep.Run(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, 1000, v)
}

func TestConcurrentStep_GivenAGraphToDraw_WhenDrawn_ThenConcurrentActionsAreApplied(t *testing.T) {
	mockGraph := new(mockGraph)
	var diagrams []pipeline.GraphDrawer
	innerStep := pipeline.NewUnitStep[int, int]("testname", nil)
	mockGraph.On("AddConcurrency", mock.MatchedBy(func(obj []pipeline.GraphDrawer) bool {
		diagrams = obj
		return true
	})).Once()

	step := pipeline.NewConcurrentStep(
		[]pipeline.Step[int, int]{
			innerStep, innerStep, innerStep, innerStep, innerStep, innerStep,
		},
		func(ctx context.Context, a, b int) (int, error) {
			return 1, nil
		},
	)

	step.Draw(mockGraph)

	assert.Len(t, diagrams, 6)
	mockGraph.AssertExpectations(t)
}

func TestConcurrentStep_GivenAGraphToDraw_WhenDrawn_ThenConcurrentStepsAreAddedAsActionsByTheirNames(t *testing.T) {
	mockGraph := new(mockGraph)
	innerStep := pipeline.NewUnitStep[int, int]("testname", nil)
	mockGraph.On("AddActivity", "testname").Times(5)
	mockGraph.On("AddConcurrency", mock.MatchedBy(func(obj interface{}) bool {
		return true
	})).Run(func(args mock.Arguments) {
		diagrams := args.Get(0).([]pipeline.GraphDrawer)
		for _, d := range diagrams {
			d(mockGraph)
		}
	})
	step := pipeline.NewConcurrentStep(
		[]pipeline.Step[int, int]{
			innerStep, innerStep, innerStep, innerStep, innerStep,
		},
		func(ctx context.Context, a, b int) (int, error) {
			return 1, nil
		},
	)

	step.Draw(mockGraph)

	mockGraph.AssertExpectations(t)
}
