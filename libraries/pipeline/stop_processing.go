package pipeline

import "context"

type stopProcessing[T interface{}] struct {
	pipeline *Pipeline
}

func StopProcessing[T interface{}](pipeline *Pipeline) Acceptor[T] {
	return &stopProcessing[T]{
		pipeline: pipeline,
	}
}

func (a *stopProcessing[T]) Run(ctx context.Context, input <-chan T) error {
	for {
		_, ok := <-input
		if !ok {
			break
		}

		a.pipeline.Cancel()
	}
	return nil
}
