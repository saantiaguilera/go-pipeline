package pipeline_test

import (
	"testing"

	"github.com/saantiaguilera/go-pipeline"
	"github.com/stretchr/testify/assert"
)

func TestStatement_GivenAnAnonymousStatement_WhenNamed_ThenReturnsEmpty(t *testing.T) {
	statement := pipeline.CreateAnonymousStatement(func(ctx pipeline.Context) bool {
		return true
	})

	assert.Empty(t, statement.Name())
}

func TestStatement_GivenAnAnonymousStatement_WhenEvaluated_ThenEvaluatesPassed(t *testing.T) {
	statement := pipeline.CreateAnonymousStatement(func(ctx pipeline.Context) bool {
		return true
	})

	assert.True(t, statement.Evaluate(&mockContext{}))
}

func TestStatement_GivenAStatement_WhenNamed_ThenReturnsName(t *testing.T) {
	statement := pipeline.CreateSimpleStatement("some name", func(ctx pipeline.Context) bool {
		return true
	})

	assert.Equal(t, "some name", statement.Name())
}

func TestStatement_GivenAStatement_WhenEvaluated_ThenEvaluatesPassed(t *testing.T) {
	statement := pipeline.CreateSimpleStatement("some name", func(ctx pipeline.Context) bool {
		return true
	})

	assert.True(t, statement.Evaluate(&mockContext{}))
}

func TestStatement_GivenAnAnonymousStatementWithNoFunc_WhenEvaluated_ThenReturnsFalse(t *testing.T) {
	statement := pipeline.CreateAnonymousStatement(nil)

	assert.False(t, statement.Evaluate(&mockContext{}))
}

func TestStatement_GivenAStatementWithNoFunc_WhenEvaluated_ThenReturnsFalse(t *testing.T) {
	statement := pipeline.CreateSimpleStatement("some name", nil)

	assert.False(t, statement.Evaluate(&mockContext{}))
}
