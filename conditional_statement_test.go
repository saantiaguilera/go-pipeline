package pipeline_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/saantiaguilera/go-pipeline"
)

func TestStatement_GivenAnAnonymousStatement_WhenNamed_ThenReturnsEmpty(t *testing.T) {
	statement := pipeline.NewAnonymousStatement(func(in int) bool {
		return true
	})

	assert.Empty(t, statement.Name())
}

func TestStatement_GivenAnAnonymousStatement_WhenEvaluated_ThenEvaluatesPassed(t *testing.T) {
	statement := pipeline.NewAnonymousStatement(func(in int) bool {
		return true
	})

	assert.True(t, statement.Evaluate(1))
}

func TestStatement_GivenAStatement_WhenNamed_ThenReturnsName(t *testing.T) {
	statement := pipeline.NewSimpleStatement("some name", func(in int) bool {
		return true
	})

	assert.Equal(t, "some name", statement.Name())
}

func TestStatement_GivenAStatement_WhenEvaluated_ThenEvaluatesPassed(t *testing.T) {
	statement := pipeline.NewSimpleStatement("some name", func(in int) bool {
		return true
	})

	assert.True(t, statement.Evaluate(1))
}

func TestStatement_GivenAnAnonymousStatementWithNoFunc_WhenEvaluated_ThenReturnsFalse(t *testing.T) {
	statement := pipeline.NewAnonymousStatement[int](nil)

	assert.False(t, statement.Evaluate(1))
}

func TestStatement_GivenAStatementWithNoFunc_WhenEvaluated_ThenReturnsFalse(t *testing.T) {
	statement := pipeline.NewSimpleStatement[int]("some name", nil)

	assert.False(t, statement.Evaluate(1))
}
