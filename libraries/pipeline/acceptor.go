package pipeline

import "context"

type Acceptor[T interface{}] interface {
	Run(ctx context.Context, input <-chan T) error
}
