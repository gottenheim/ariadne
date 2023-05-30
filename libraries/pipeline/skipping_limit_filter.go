package pipeline

import "context"

type skippingLimitFilter struct {
	limit  int
	actual int
}

func SkippingLimit(limit int) Filter[int, int] {
	return &skippingLimitFilter{
		limit: limit,
	}
}

func (f *skippingLimitFilter) Run(ctx context.Context, input <-chan int, output chan<- int) error {
	for {
		val, ok := <-input
		if !ok {
			break
		}

		if f.actual < f.limit {
			output <- val
			f.actual++
		}
	}
	return nil
}
