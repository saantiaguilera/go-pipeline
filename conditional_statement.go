package pipeline

// Statement is contract for evaluating a choice to take
type Statement interface {
	// A statement is named, since it can be represented (as a logical gate)
	Named

	// Evaluate the statement, returning a boolean denoting which choice to take
	Evaluate() bool
}

type statement struct {
	Label string
	Func  func() bool
}

func (s *statement) Name() string {
	return s.Label
}

func (s *statement) Evaluate() bool {
	return s.Func != nil && s.Func()
}

// CreateSimpleStatement creates a statement represented by the given name, that will evaluate to the given evaluation
func CreateSimpleStatement(name string, evaluation func() bool) Statement {
	return &statement{
		Label: name,
		Func:  evaluation,
	}
}

// CreateAnonymousStatement creates an anonymous statement with no representation, that will evaluate to the given evaluation
func CreateAnonymousStatement(evaluation func() bool) Statement {
	return &statement{
		Func: evaluation,
	}
}
