package pipeline

import "context"

type Filter[T interface{}, K interface{}] interface {
	Run(ctx context.Context, input <-chan T, output chan<- K) error
}
