package pipeline

import "context"

type limitFilter[T interface{}] struct {
	limit  int
	actual int
}

func Limit[T interface{}](limit int) Filter[T, T] {
	return &limitFilter[T]{
		limit: limit,
	}
}

func (f *limitFilter[T]) Run(ctx context.Context, input <-chan T, output chan<- T) error {
	for {
		val, ok := <-input
		if !ok {
			break
		}

		if f.actual < f.limit {
			if !WriteToChannel[T](ctx, output, val) {
				break
			}

			f.actual++
		}
	}
	return nil
}
