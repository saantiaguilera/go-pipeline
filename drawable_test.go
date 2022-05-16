package pipeline_test

import (
	"github.com/stretchr/testify/mock"

	"github.com/saantiaguilera/go-pipeline"
)

type mockGraph struct {
	mock.Mock
}

func (m *mockGraph) AddConcurrency(branches ...pipeline.GraphDrawer) {
	_ = m.Called(branches)
}

func (m *mockGraph) AddDecision(statement string, yes pipeline.GraphDrawer, no pipeline.GraphDrawer) {
	_ = m.Called(statement, yes, no)
}

func (m *mockGraph) AddActivity(label string) {
	_ = m.Called(label)
}

func (m *mockGraph) String() string {
	args := m.Called()
	return args.String(0)
}
