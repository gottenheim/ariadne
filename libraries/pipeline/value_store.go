package pipeline

import "context"

type valueStore[T interface{}] struct {
	value T
}

func NewValueStore[T interface{}]() *valueStore[T] {
	return &valueStore[T]{}
}

func (s *valueStore[T]) Run(ctx context.Context, input <-chan T) error {
	for {
		value, ok := <-input
		if !ok {
			break
		}
		s.value = value
	}

	return nil
}

func (s *valueStore[T]) Value() T {
	return s.value
}
