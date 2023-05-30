package pipeline

import "context"

type aggregatorAdapter[T interface{}, K interface{}] struct {
	leftArg    <-chan T
	rightArg   <-chan T
	output     chan<- K
	aggregator Aggregator[T, K]
}

func newAggregatorAdapter[T interface{}, K interface{}](pipeline *Pipeline, aggregator Aggregator[T, K]) *aggregatorAdapter[T, K] {
	adapter := &aggregatorAdapter[T, K]{
		aggregator: aggregator,
	}

	pipeline.attach(adapter)

	return adapter
}

func (f *aggregatorAdapter[T, K]) Name() string {
	return "Aggregator"
}

func (f *aggregatorAdapter[T, K]) SetLeftArgChannel(leftArg <-chan T) {
	f.leftArg = leftArg
}

func (f *aggregatorAdapter[T, K]) SetRightArgChannel(rightArg <-chan T) {
	f.rightArg = rightArg
}

func (f *aggregatorAdapter[T, K]) SetOutputChannel(output chan<- K) {
	f.output = output
}

func (f *aggregatorAdapter[T, K]) Run(ctx context.Context) error {
	defer func() {
		close(f.output)
	}()

	return f.aggregator.Run(ctx, f.leftArg, f.rightArg, f.output)
}
