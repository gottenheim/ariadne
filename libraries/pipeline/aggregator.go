package pipeline

import "context"

type Aggregator[T interface{}, K interface{}] interface {
	Run(ctx context.Context, leftArg <-chan T, rightArg <-chan T, output chan<- K) error
}
