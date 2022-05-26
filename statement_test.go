package pipeline_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/saantiaguilera/go-pipeline"
)

// This examples shows a simple statement that lets
// us evaluate it with a given input to yield
// a boolean result
//
// This example uses dummy data to showcase as simple as possible this scenario.
func ExampleStatement() {
	stmt := pipeline.NewStatement(
		"is_number_odd",
		func(ctx context.Context, i int) bool {
			return i%2 != 0
		},
	)

	out := stmt.Evaluate(context.Background(), 25)

	fmt.Println(out)
	// output: true
}

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
