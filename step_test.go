package pipeline_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/saantiaguilera/go-pipeline"
)

func TestSimpleStep_GivenAName_WhenGettingItsName_ThenItsTheExpected(t *testing.T) {
	expectedName := "test_name"
	step := pipeline.NewStep[interface{}](expectedName, nil)

	name := step.Name()

	assert.Equal(t, expectedName, name)
}

func TestSimpleStep_GivenARunFunc_WhenRunning_ThenItsCalled(t *testing.T) {
	called := false
	run := func(int) error {
		called = true
		return nil
	}
	step := pipeline.NewStep("", run)

	_ = step.Run(1)

	assert.True(t, called)
}

func TestSimpleStep_GivenARunFuncThatErrors_WhenRunning_ThenErrorIsReturned(t *testing.T) {
	expectedErr := errors.New("some error")
	run := func(int) error {
		return expectedErr
	}
	step := pipeline.NewStep("", run)

	err := step.Run(1)

	assert.Equal(t, expectedErr, err)
}
