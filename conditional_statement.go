package pipeline

// Statement is contract for evaluating a choice to take
type Statement interface {
	// A statement is named, since it can be represented (as a logical gate)
	Named

	// Evaluate the statement with a provided context of the transaction, returning a boolean denoting which choice to take
	Evaluate(ctx Context) bool
}

type statement struct {
	Label string
	Func  func(ctx Context) bool
}

func (s *statement) Name() string {
	return s.Label
}

func (s *statement) Evaluate(ctx Context) bool {
	return s.Func != nil && s.Func(ctx)
}

// CreateSimpleStatement creates a statement represented by the given name, that will evaluate to the given evaluation
func CreateSimpleStatement(name string, evaluation func(ctx Context) bool) Statement {
	return &statement{
		Label: name,
		Func:  evaluation,
	}
}

// CreateAnonymousStatement creates an anonymous statement with no representation, that will evaluate to the given evaluation
func CreateAnonymousStatement(evaluation func(ctx Context) bool) Statement {
	return &statement{
		Func: evaluation,
	}
}
