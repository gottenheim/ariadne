package pipeline

import (
	"context"
)

type countingFilter[T interface{}] struct {
}

func NewCounter[T interface{}]() *countingFilter[T] {
	return &countingFilter[T]{}
}

func (f *countingFilter[T]) Run(ctx context.Context, input <-chan T, output chan<- int) error {
	count := 0

	for {
		_, ok := <-input
		if !ok {
			break
		}

		count++

		if !WriteToChannel[int](ctx, output, count) {
			break
		}
	}

	return nil
}
