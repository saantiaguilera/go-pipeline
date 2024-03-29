package pipeline_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/saantiaguilera/go-pipeline"
)

func TestUMLGraph_GivenAGraph_WhenAddingEmptyConcurrentCases_ThenPlantUMLForksAreAdded(t *testing.T) {
	diagram := pipeline.NewUMLGraph()
	diagram.AddConcurrency(
		func(graph pipeline.Graph) {
		},
		func(graph pipeline.Graph) {
		},
		func(graph pipeline.Graph) {
		},
	)

	content := diagram.String()
	expectedContent := "\nfork\nfork again\nfork again\nend fork\n"

	assert.Contains(t, content, expectedContent)
}

func TestUMLGraph_GivenAGraph_WhenAddingConcurrentCases_ThenPlantUMLForksAreAdded(t *testing.T) {
	diagram := pipeline.NewUMLGraph()
	diagram.AddConcurrency(
		func(graph pipeline.Graph) {
			graph.AddActivity("1")
		},
		func(graph pipeline.Graph) {
			graph.AddActivity("2")
		},
		func(graph pipeline.Graph) {
			graph.AddActivity("3")
		},
	)

	content := diagram.String()
	expectedContent := "\nfork\n:1;\nfork again\n:2;\nfork again\n:3;\nend fork\n"

	assert.Contains(t, content, expectedContent)
}

func TestUMLGraph_GivenAGraph_WhenAddingZeroConcurrentCases_ThenNothingHappens(t *testing.T) {
	diagram := pipeline.NewUMLGraph()
	diagram.AddConcurrency()

	content := diagram.String()
	notExpectedContent := "fork"

	assert.NotContains(t, content, notExpectedContent)
}

func TestUMLGraph_GivenAGraph_WhenAddingActivities_ThenPlantUMLActivitiesAreAdded(t *testing.T) {
	diagram := pipeline.NewUMLGraph()
	diagram.AddActivity("1")
	diagram.AddActivity("1 2")
	diagram.AddActivity("1 2 3")
	diagram.AddActivity("1 2 3 4")
	diagram.AddActivity("1 2 3 4 5")

	content := diagram.String()
	expectedContent := "\n:1;\n:1 2;\n:1 2 3;\n:1 2 3 4;\n:1 2 3 4 5;\n"

	assert.Contains(t, content, expectedContent)
}

func TestUMLGraph_GivenAGraph_WhenAddingADecision_ThenPlantUMLChoiceIsAdded(t *testing.T) {
	diagram := pipeline.NewUMLGraph()
	diagram.AddDecision("is this a test?", func(graph pipeline.Graph) {
		graph.AddActivity("yes, this is a test")
	}, func(graph pipeline.Graph) {
		graph.AddActivity("seems this isn't a test")
	})

	content := diagram.String()
	expectedContent := "\nif (is this a test?) then (yes)\n:yes, this is a test;\nelse (no)\n:seems this isn't a test;\n"

	assert.Contains(t, content, expectedContent)
}

func TestUMLGraph_GivenAGraph_WhenStringRepresentationIsAsked_ThenCompletePlantUMLIsRepresented(t *testing.T) {
	diagram := pipeline.NewUMLGraph()
	diagram.AddActivity("beginning")
	diagram.AddConcurrency(func(graph pipeline.Graph) {
		graph.AddActivity("branch 1")
	}, func(graph pipeline.Graph) {
		diagram.AddDecision("is this a test?", func(graph pipeline.Graph) {
			graph.AddActivity("yes, this is a test")
		}, func(graph pipeline.Graph) {
			graph.AddActivity("seems this isn't a test")
		})
	})

	content := diagram.String()
	expectedContent := "@startuml\nstart\n:beginning;\nfork\n:branch 1;\nfork again\nif (is this a test?) then (yes)" +
		"\n:yes, this is a test;\nelse (no)\n:seems this isn't a test;\nendif\nend fork\nstop\n@enduml\n"

	assert.Equal(t, expectedContent, content)
}
