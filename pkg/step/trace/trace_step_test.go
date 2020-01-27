package trace_test

import (
	"bytes"
	"errors"
	"github.com/saantiaguilera/go-pipeline/pkg/api"
	"github.com/saantiaguilera/go-pipeline/pkg/step/trace"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

type TestStep struct{}

func (t TestStep) Name() string {
	return "test name"
}

func (t TestStep) Run() error {
	return errors.New("some error")
}

type TestExecutor struct{}

func (t TestExecutor) Run(runnable api.Runnable) error {
	return runnable.Run()
}

func TestTrace_GivenAStepToTrace_WhenRun_ThenOutputsInnerStepErr(t *testing.T) {
	step := trace.CreateTracedStep(&TestStep{})

	err := TestExecutor{}.Run(step)

	assert.NotNil(t, err)
	assert.Equal(t, "some error", err.Error())
}

func TestTrace_GivenAStepToTrace_WhenRun_ThenSpecificFormatIsUsed(t *testing.T) {
	writer := bytes.NewBufferString("")
	step := trace.CreateTracedStepWithWriter(&TestStep{}, writer)
	validator := regexp.MustCompile(`^\[STEP] \d{4}-\d{2}-\d{2} - \d{2}:\d{2}:\d{2} \| test name \| [.\d]+[µnm]s \| Failure: some error\n$`)

	_ = TestExecutor{}.Run(step)

	output := writer.Bytes()

	assert.True(t, validator.Match(output))
}
