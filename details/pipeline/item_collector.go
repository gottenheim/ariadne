package pipeline

import "context"

type ItemCollector[T interface{}] struct {
	Items []T
}

func NewItemCollector[T interface{}]() *ItemCollector[T] {
	return &ItemCollector[T]{}
}

func (c *ItemCollector[T]) Run(ctx context.Context, input <-chan T) error {
	for {
		item, ok := <-input
		if !ok {
			break
		}
		c.Items = append(c.Items, item)
	}

	return nil
}
