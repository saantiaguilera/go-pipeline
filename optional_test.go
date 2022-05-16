package pipeline_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/saantiaguilera/go-pipeline"
)

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
