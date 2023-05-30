package pipeline

import (
	"context"
)

type cancellingLimitFilter struct {
	pipeline *Pipeline
	limit    int
	actual   int
}

func CancellingLimit(pipeline *Pipeline, limit int) Filter[int, int] {
	return &cancellingLimitFilter{
		pipeline: pipeline,
		limit:    limit,
	}
}

func (f *cancellingLimitFilter) Run(ctx context.Context, input <-chan int, output chan<- int) error {
	for {
		val, ok := <-input
		if !ok {
			break
		}

		if f.actual < f.limit {
			output <- val
			f.actual++
		}

		if f.actual >= f.limit {
			f.pipeline.Cancel()
		}
	}
	return nil
}
