package pipeline

import "context"

type (
	// SequentialStep runs two steps sequentially.
	//
	// A sequential step allows partial mutation between the two steps.
	// Eg. The first step can have an input of type 'A' and an output of type 'B'
	//     The second step can have an input of type 'B' and an output of type 'C'
	//     Hence, the sequential step is a step that goes from input 'A', to output 'C' 
	//        (mutating patially into 'B' in the middle)
	//
	// If one of the steps fails, the step is halted and the error is returned
	SequentialStep[I, M, O any] struct {
		start Step[I, M]
		end   Step[M, O]
	}
)

// NewSequentialStep creates step that will run each of the steps sequentially. If one of them fails, the operation will abort immediately
func NewSequentialStep[I, M, O any](s Step[I, M], e Step[M, O]) SequentialStep[I, M, O] {
	return SequentialStep[I, M, O]{
		start: s,
		end:   e,
	}
}

// Run both steps sequentially. If one of them fails, the step is halted and the error is returned.
func (s SequentialStep[I, M, O]) Run(ctx context.Context, in I) (O, error) {
	m, err := s.start.Run(ctx, in)
	if err != nil {
		return *new(O), err
	}

	return s.end.Run(ctx, m)
}

func (s SequentialStep[I, M, O]) Draw(graph Graph) {
	s.start.Draw(graph)
	s.end.Draw(graph)
}
