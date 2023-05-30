package pipeline

import "context"

type devNull[T interface{}] struct {
}

func DevNull[T interface{}]() Acceptor[T] {
	return &devNull[T]{}
}

func (f *devNull[T]) Run(ctx context.Context, input <-chan T) error {
	for {
		_, ok := <-input
		if !ok {
			break
		}
	}
	return nil
}
