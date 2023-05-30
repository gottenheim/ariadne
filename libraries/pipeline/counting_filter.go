package pipeline

import "context"

type countingFilter[T interface{}] struct {
	count int
}

func NewCounter[T interface{}]() *countingFilter[T] {
	return &countingFilter[T]{}
}

func (f *countingFilter[T]) Run(ctx context.Context, input <-chan T, output chan<- int) error {
	for {
		_, ok := <-input
		if !ok {
			break
		}

		f.count++

		output <- f.count
	}

	return nil
}
