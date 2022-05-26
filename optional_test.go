package pipeline_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/saantiaguilera/go-pipeline"
)

// The following example shows an optional step were it will be run if the evaluated
// statement yields true, otherwise it will return the same input it was provided
// this type of constructed step doesn't allow output mutation as we don't know
// how to default the mutated output if the step is skipped
//
// This example uses dummy data to showcase as simple as possible this scenario.
//
// Note: we use several UnitStep to showcase as it allows us to
// easily run dummy code, but it could use any type of step you want
// as long as it implements pipeline.Step[I, O]
func ExampleOptionalStep() {
	type User any
	stmt := pipeline.NewStatement(
		"check_something",
		func(ctx context.Context, in User) bool {
			// check and return if we should run optional or not
			return in == User(1)
		},
	)
	of := pipeline.NewUnitStep(
		"optional_case",
		func(ctx context.Context, in User) (User, error) {
			// do something with input
			return User(true), nil
		},
	)
	ctx := context.Background()

	pipe := pipeline.NewOptionalStep[User](stmt, of)

	// Skips optional step (returns same input)
	out, err := pipe.Run(ctx, User(12))
	fmt.Println(out, err)

	// Runs optional step (returns step output)
	out, err = pipe.Run(ctx, User(1))
	fmt.Println(out, err)

	// output:
	// 12 <nil>
	// true <nil>
}

// The following example allows us to define an optional step with different
// output from its input, and in case the optional step is skipped it will
// call a provided default function for it to return the default output you want for that case
//
// This example uses dummy data to showcase as simple as possible this scenario.
//
// Note: we use several UnitStep to showcase as it allows us to
// easily run dummy code, but it could use any type of step you want
// as long as it implements pipeline.Step[I, O]
func ExampleOptionalStep_default() {
	type User any
	type Data any
	stmt := pipeline.NewStatement(
		"check_something",
		func(ctx context.Context, in User) bool {
			// check and return if we should run optional or not
			return in == User(1)
		},
	)
	of := pipeline.NewUnitStep(
		"optional_case",
		func(ctx context.Context, in User) (Data, error) {
			// do something with input
			return Data(true), nil
		},
	)
	def := func(ctx context.Context, in User) (Data, error) {
		// create default value of type Data (gets run if we
		// skip the step as we don't know how to default Data type)
		return Data(false), nil
	}
	ctx := context.Background()

	pipe := pipeline.NewOptionalStepWithDefault[User, Data](stmt, of, def)

	// Skips optional step (returns default)
	out, err := pipe.Run(ctx, User(12))
	fmt.Println(out, err)

	// Runs optional step (returns step output)
	out, err = pipe.Run(ctx, User(1))
	fmt.Println(out, err)

	// output:
	// false <nil>
	// true <nil>
}

func TestOptionalStep_GivenNilStatement_WhenRun_ThenDefaults(t *testing.T) {
	run := false
	step := pipeline.NewOptionalStep[any](pipeline.NewAnonymousStatement[any](nil), pipeline.NewUnitStep("", func(_ context.Context, _ any) (any, error) {
		return nil, nil
	}))

	res, err := step.Run(context.Background(), 1)

	assert.Nil(t, err)
	assert.Equal(t, 1, res)
	assert.False(t, run)
}

func TestOptionalStep_GivenStatementTrue_WhenRun_ThenEvaluatesStep(t *testing.T) {
	run := false
	ev := pipeline.NewUnitStep("", func(ctx context.Context, t any) (any, error) {
		run = true
		return 25, nil
	})
	step := pipeline.NewOptionalStep[any](pipeline.NewAnonymousStatement(func(ctx context.Context, in any) bool {
		return true
	}), ev)

	v, err := step.Run(context.Background(), 1)

	assert.Nil(t, err)
	assert.Equal(t, 25, v)
	assert.True(t, run)
}

func TestOptionalStep_GivenNilStatementWithDefault_WhenRun_ThenDefaults(t *testing.T) {
	run := false
	step := pipeline.NewOptionalStepWithDefault[any, any](pipeline.NewAnonymousStatement[any](nil), pipeline.NewUnitStep("", func(_ context.Context, _ any) (any, error) {
		return nil, nil
	}), func(_ context.Context, _ any) (any, error) {
		return 25, nil
	})

	res, err := step.Run(context.Background(), 1)

	assert.Nil(t, err)
	assert.Equal(t, 25, res)
	assert.False(t, run)
}

func TestOptionalStep_GivenStatementTrueWithDefault_WhenRun_ThenEvaluatesStep(t *testing.T) {
	run := false
	ev := pipeline.NewUnitStep("", func(ctx context.Context, t any) (any, error) {
		run = true
		return 25, nil
	})
	step := pipeline.NewOptionalStepWithDefault[any, any](pipeline.NewAnonymousStatement(func(ctx context.Context, in any) bool {
		return true
	}), ev, func(_ context.Context, _ any) (any, error) {
		return 50, nil
	})

	v, err := step.Run(context.Background(), 1)

	assert.Nil(t, err)
	assert.Equal(t, 25, v)
	assert.True(t, run)
}

func TestOptionalStep_GivenAGraphToDrawWithAnonymouseStatement_WhenDrawn_ThenConditionGetsEmptyName(t *testing.T) {
	statement := pipeline.NewAnonymousStatement(func(ctx context.Context, in any) bool {
		return true
	})
	mockGraph := new(mockGraph)
	mockGraph.On(
		"AddDecision",
		"",
		mock.MatchedBy(func(obj any) bool {
			return true
		}), mock.MatchedBy(func(obj any) bool {
			return true
		}),
	)
	trueStep := pipeline.NewUnitStep[any, any]("", nil)
	step := pipeline.NewOptionalStepWithDefault[any, any](statement, trueStep, func(_ context.Context, _ any) (any, error) {
		return nil, nil
	})

	step.Draw(mockGraph)

	mockGraph.AssertExpectations(t)
}

func TestOptionalStep_GivenAGraphToDraw_WhenDrawn_ThenConditionGetsNameOfStatement(t *testing.T) {
	mockGraph := new(mockGraph)
	mockGraph.On(
		"AddDecision",
		"SomeFuncName",
		mock.MatchedBy(func(obj any) bool {
			return true
		}), mock.MatchedBy(func(obj any) bool {
			return true
		}),
	)
	trueStep := pipeline.NewUnitStep[any, any]("", nil)
	step := pipeline.NewOptionalStep[any](pipeline.NewStatement[any]("SomeFuncName", nil), trueStep)

	step.Draw(mockGraph)

	mockGraph.AssertExpectations(t)
}

func TestOptionalStep_GivenAGraphToDraw_WhenDrawn_ThenConditionIsAppliedWithBothBranches(t *testing.T) {
	mockGraph := new(mockGraph)
	mockGraph.On("AddActivity", "truestep").Once()
	mockGraph.On(
		"AddDecision",
		mock.Anything,
		mock.MatchedBy(func(obj any) bool {
			return true
		}), mock.MatchedBy(func(obj any) bool {
			return true
		}),
	).Run(func(args mock.Arguments) {
		args.Get(1).(pipeline.GraphDrawer)(mockGraph)
		args.Get(2).(pipeline.GraphDrawer)(mockGraph)
	})
	trueStep := pipeline.NewUnitStep[any, any]("truestep", nil)
	step := pipeline.NewOptionalStep[any](pipeline.NewAnonymousStatement(func(ctx context.Context, in any) bool {
		return true
	}), trueStep)

	step.Draw(mockGraph)

	mockGraph.AssertExpectations(t)
}
