package pipeline_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/saantiaguilera/go-pipeline"
)

func TestStatement_GivenAnAnonymousStatement_WhenNamed_ThenReturnsEmpty(t *testing.T) {
	statement := pipeline.NewAnonymousStatement(func(ctx context.Context, in int) bool {
		return true
	})

	assert.Empty(t, statement.Name())
}

func TestStatement_GivenAnAnonymousStatement_WhenEvaluated_ThenEvaluatesPassed(t *testing.T) {
	statement := pipeline.NewAnonymousStatement(func(ctx context.Context, in int) bool {
		return true
	})

	assert.True(t, statement.Evaluate(context.Background(), 1))
}

func TestStatement_GivenAStatement_WhenNamed_ThenReturnsName(t *testing.T) {
	statement := pipeline.NewStatement("some name", func(ctx context.Context, in int) bool {
		return true
	})

	assert.Equal(t, "some name", statement.Name())
}

func TestStatement_GivenAStatement_WhenEvaluated_ThenEvaluatesPassed(t *testing.T) {
	statement := pipeline.NewStatement("some name", func(ctx context.Context, in int) bool {
		return true
	})

	assert.True(t, statement.Evaluate(context.Background(), 1))
}

func TestStatement_GivenAnAnonymousStatementWithNoFunc_WhenEvaluated_ThenReturnsFalse(t *testing.T) {
	statement := pipeline.NewAnonymousStatement[int](nil)

	assert.False(t, statement.Evaluate(context.Background(), 1))
}

func TestStatement_GivenAStatementWithNoFunc_WhenEvaluated_ThenReturnsFalse(t *testing.T) {
	statement := pipeline.NewStatement[int]("some name", nil)

	assert.False(t, statement.Evaluate(context.Background(), 1))
}
