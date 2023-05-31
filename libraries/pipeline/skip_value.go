package pipeline

import "context"

type skipValue[T interface{}] struct {
}

func Skip[T interface{}]() Acceptor[T] {
	return &skipValue[T]{}
}

func (f *skipValue[T]) Run(ctx context.Context, input <-chan T) error {
	for {
		_, ok := <-input
		if !ok {
			break
		}
	}
	return nil
}
