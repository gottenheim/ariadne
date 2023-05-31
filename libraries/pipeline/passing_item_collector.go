package pipeline

import "context"

type passingItemCollector[T interface{}] struct {
	Items []T
}

func NewPassingItemCollector[T interface{}]() *passingItemCollector[T] {
	return &passingItemCollector[T]{}
}

func (c *passingItemCollector[T]) Run(ctx context.Context, input <-chan T, output chan<- T) error {
	for {
		item, ok := <-input
		if !ok {
			break
		}
		c.Items = append(c.Items, item)

		if !WriteToChannel[T](ctx, output, item) {
			break
		}
	}

	return nil
}
