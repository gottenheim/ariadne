package pipeline

import "context"

type filterAdapter[T interface{}, K interface{}] struct {
	input  <-chan T
	output chan<- K
	filter Filter[T, K]
}

func newFilterAdapter[T interface{}, K interface{}](p *Pipeline, filter Filter[T, K]) *filterAdapter[T, K] {
	adapter := &filterAdapter[T, K]{
		filter: filter,
	}

	p.attach(adapter)

	return adapter
}

func (f *filterAdapter[T, K]) SetInputChannel(input <-chan T) {
	f.input = input
}

func (f *filterAdapter[T, K]) SetOutputChannel(output chan<- K) {
	f.output = output
}

func (f *filterAdapter[T, K]) Run(ctx context.Context) error {
	defer func() {
		close(f.output)
	}()

	return f.filter.Run(ctx, f.input, f.output)
}
