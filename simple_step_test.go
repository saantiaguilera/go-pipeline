package pipeline_test

import (
	"errors"
	"testing"

	"github.com/saantiaguilera/go-pipeline"
	"github.com/stretchr/testify/assert"
)

func TestSimpleStep_GivenAName_WhenGettingItsName_ThenItsTheExpected(t *testing.T) {
	expectedName := "test_name"
	step := pipeline.CreateSimpleStep(expectedName, nil)

	name := step.Name()

	assert.Equal(t, expectedName, name)
}

func TestSimpleStep_GivenARunFunc_WhenRunning_ThenItsCalled(t *testing.T) {
	called := false
	run := func(ctx pipeline.Context) error {
		called = true
		return nil
	}
	step := pipeline.CreateSimpleStep("", run)

	_ = step.Run(nil)

	assert.True(t, called)
}

func TestSimpleStep_GivenARunFuncThatErrors_WhenRunning_ThenErrorIsReturned(t *testing.T) {
	expectedErr := errors.New("some error")
	run := func(ctx pipeline.Context) error {
		return expectedErr
	}
	step := pipeline.CreateSimpleStep("", run)

	err := step.Run(nil)

	assert.Equal(t, expectedErr, err)
}
