package pipeline_test

import (
	"github.com/saantiaguilera/go-pipeline"
	"github.com/stretchr/testify/mock"
)

type mockGraphDiagram struct {
	mock.Mock
}

func (m *mockGraphDiagram) AddConcurrency(branches ...pipeline.DrawDiagram) {
	_ = m.Called(branches)
}

func (m *mockGraphDiagram) AddDecision(statement string, yes pipeline.DrawDiagram, no pipeline.DrawDiagram) {
	_ = m.Called(statement, yes, no)
}

func (m *mockGraphDiagram) AddActivity(label string) {
	_ = m.Called(label)
}

func (m *mockGraphDiagram) String() string {
	args := m.Called()
	return args.String(0)
}
