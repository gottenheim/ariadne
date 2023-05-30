package pipeline

import (
	"context"
)

type stopOnValueGreaterThan struct {
	pipeline *Pipeline
	value    int
}

func StopOnValueGreaterThan(pipeline *Pipeline, value int) Filter[int, int] {
	return &stopOnValueGreaterThan{
		pipeline: pipeline,
		value:    value,
	}
}

func (f *stopOnValueGreaterThan) Run(ctx context.Context, input <-chan int, output chan<- int) error {
	for {
		val, ok := <-input
		if !ok {
			break
		}

		if val <= f.value {
			output <- val
		} else {
			f.pipeline.Cancel()
		}
	}
	return nil
}
