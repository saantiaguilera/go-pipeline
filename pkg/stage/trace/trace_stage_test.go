package trace_test

import (
	"bytes"
	"errors"
	"github.com/saantiaguilera/go-pipeline/pkg/api"
	"github.com/saantiaguilera/go-pipeline/pkg/stage/trace"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

type TestStage struct{}

func (t TestStage) Run(executor api.Executor) error {
	return errors.New("some error")
}

type TestExecutor struct{}

func (t TestExecutor) Run(runnable api.Runnable) error {
	return runnable.Run()
}

func TestTrace_GivenAStageToTrace_WhenRun_ThenOutputsInnerStageErr(t *testing.T) {
	stage := trace.CreateTracedStage("test name", &TestStage{})

	err := stage.Run(&TestExecutor{})

	assert.NotNil(t, err)
	assert.Equal(t, "some error", err.Error())
}

func TestTrace_GivenAStageToTrace_WhenRun_ThenSpecificFormatIsUsed(t *testing.T) {
	writer := bytes.NewBufferString("")
	stage := trace.CreateTracedStageWithWriter("test name", &TestStage{}, writer)
	validator := regexp.MustCompile(`^\[STAGE] \d{4}-\d{2}-\d{2} - \d{2}:\d{2}:\d{2} \| test name \| [.\d]+[µnm]s \| Failure: some error\n$`)

	_ = stage.Run(&TestExecutor{})

	output := writer.Bytes()

	assert.True(t, validator.Match(output))
}
