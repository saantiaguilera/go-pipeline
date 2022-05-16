package pipeline_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/saantiaguilera/go-pipeline"
)

func TestConditionalStep_GivenNilStatement_WhenRun_FalseIsRun(t *testing.T) {
	run := false
	falseStep := pipeline.NewUnitStep("", func(ctx context.Context, t any) (any, error) {
		run = true
		return nil, nil
	})
	trueStep := pipeline.NewUnitStep[any, any]("", nil)

	step := pipeline.NewConditionalStep[any, any](pipeline.NewAnonymousStatement[any](nil), trueStep, falseStep)

	_, err := step.Run(context.Background(), 1)

	assert.Nil(t, err)
	assert.True(t, run)
}

func TestConditionalStep_GivenStatementTrue_WhenRun_TrueIsRun(t *testing.T) {
	run := false
	falseStep := pipeline.NewUnitStep[any, any]("", nil)
	trueStep := pipeline.NewUnitStep("", func(ctx context.Context, t any) (any, error) {
		run = true
		return nil, nil
	})
	step := pipeline.NewConditionalStep[any, any](pipeline.NewAnonymousStatement(func(ctx context.Context, in any) bool {
		return true
	}), trueStep, falseStep)

	_, err := step.Run(context.Background(), 1)

	assert.Nil(t, err)
	assert.True(t, run)
}

func TestConditionalStep_GivenStatementFalse_WhenRun_FalseIsRun(t *testing.T) {
	run := false
	falseStep := pipeline.NewUnitStep("", func(ctx context.Context, t any) (any, error) {
		run = true
		return nil, nil
	})
	trueStep := pipeline.NewUnitStep[any, any]("", nil)
	step := pipeline.NewConditionalStep[any, any](pipeline.NewAnonymousStatement(func(ctx context.Context, in any) bool {
		return false
	}), trueStep, falseStep)

	_, err := step.Run(context.Background(), 1)

	assert.Nil(t, err)
	assert.True(t, run)
}

func TestConditionalStep_GivenStatementTrueAndNilTrue_WhenRun_ThenErrors(t *testing.T) {
	falseStep := pipeline.NewUnitStep[any, any]("", nil)
	step := pipeline.NewConditionalStep[any, any](pipeline.NewAnonymousStatement(func(ctx context.Context, in any) bool {
		return true
	}), nil, falseStep)

	_, err := step.Run(context.Background(), 1)

	assert.Error(t, err)
}

func TestConditionalStep_GivenStatementFalseNilFalse_WhenRun_ThenErrors(t *testing.T) {
	trueStep := pipeline.NewUnitStep[any, any]("", nil)
	step := pipeline.NewConditionalStep[any, any](pipeline.NewAnonymousStatement(func(ctx context.Context, in any) bool {
		return false
	}), trueStep, nil)

	_, err := step.Run(context.Background(), 1)

	assert.Error(t, err)
}

func TestConditionalStep_GivenStatementTrueWithTrueError_WhenRun_TrueErrorReturned(t *testing.T) {
	trueErr := errors.New("error")
	falseStep := pipeline.NewUnitStep[any, any]("", nil)
	trueStep := pipeline.NewUnitStep("", func(ctx context.Context, t any) (any, error) {
		return nil, trueErr
	})
	step := pipeline.NewConditionalStep[any, any](pipeline.NewAnonymousStatement(func(ctx context.Context, in any) bool {
		return true
	}), trueStep, falseStep)

	_, err := step.Run(context.Background(), 1)

	assert.Equal(t, trueErr, err)
}

func TestConditionalStep_GivenStatementFalseWithFalseError_WhenRun_FalseErrorReturned(t *testing.T) {
	falseErr := errors.New("error")
	trueStep := pipeline.NewUnitStep[any, any]("", nil)
	falseStep := pipeline.NewUnitStep("", func(ctx context.Context, t any) (any, error) {
		return nil, falseErr
	})
	step := pipeline.NewConditionalStep[any, any](pipeline.NewAnonymousStatement(func(ctx context.Context, in any) bool {
		return false
	}), trueStep, falseStep)

	_, err := step.Run(context.Background(), 1)

	assert.Equal(t, falseErr, err)
}

func TestConditionalStep_GivenAGraphToDrawWithAnonymouseStatement_WhenDrawn_ThenConditionGetsEmptyName(t *testing.T) {
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
	falseStep := pipeline.NewUnitStep[any, any]("", nil)
	trueStep := pipeline.NewUnitStep[any, any]("", nil)
	step := pipeline.NewConditionalStep[any, any](statement, trueStep, falseStep)

	step.Draw(mockGraph)

	mockGraph.AssertExpectations(t)
}

func TestConditionalStep_GivenAGraphToDraw_WhenDrawn_ThenConditionGetsNameOfStatement(t *testing.T) {
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
	falseStep := pipeline.NewUnitStep[any, any]("", nil)
	trueStep := pipeline.NewUnitStep[any, any]("", nil)
	step := pipeline.NewConditionalStep[any, any](pipeline.NewStatement[any]("SomeFuncName", nil), trueStep, falseStep)

	step.Draw(mockGraph)

	mockGraph.AssertExpectations(t)
}

func TestConditionalStep_GivenAGraphToDraw_WhenDrawn_ThenConditionIsAppliedWithBothBranches(t *testing.T) {
	mockGraph := new(mockGraph)
	mockGraph.On("AddActivity", "truestep").Once()
	mockGraph.On("AddActivity", "falsestep").Once()
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
	falseStep := pipeline.NewUnitStep[any, any]("falsestep", nil)
	trueStep := pipeline.NewUnitStep[any, any]("truestep", nil)
	step := pipeline.NewConditionalStep[any, any](pipeline.NewAnonymousStatement(func(ctx context.Context, in any) bool {
		return true
	}), trueStep, falseStep)

	step.Draw(mockGraph)

	mockGraph.AssertExpectations(t)
}

func TestConditionalStep_GivenAGraphToDraw_WhenDrawnAndTrueExecuted_ThenTrueBranchIsNilValidated(t *testing.T) {
	mockGraph := new(mockGraph)
	mockGraph.On("AddActivity", "falsestep").Once()
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
	falseStep := pipeline.NewUnitStep[any, any]("falsestep", nil)
	step := pipeline.NewConditionalStep[any, any](pipeline.NewAnonymousStatement(func(ctx context.Context, in any) bool {
		return true
	}), nil, falseStep)

	step.Draw(mockGraph)

	mockGraph.AssertExpectations(t)
}

func TestConditionalStep_GivenAGraphToDraw_WhenDrawnAndFalseExecuted_ThenFalseBranchIsNilValidated(t *testing.T) {
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
	step := pipeline.NewConditionalStep[any, any](pipeline.NewAnonymousStatement(func(ctx context.Context, in any) bool {
		return true
	}), trueStep, nil)

	step.Draw(mockGraph)

	mockGraph.AssertExpectations(t)
}
