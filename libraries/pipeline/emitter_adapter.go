package pipeline

import "context"

type emitterAdapter[K interface{}] struct {
	output   chan<- K
	emitter  Emitter[K]
	pipeline *Pipeline
}

func newEmitterAdapter[K interface{}](pipeline *Pipeline, emitter Emitter[K]) *emitterAdapter[K] {
	adapter := &emitterAdapter[K]{
		emitter:  emitter,
		pipeline: pipeline,
	}

	pipeline.attach(adapter)
	return adapter
}

func (f *emitterAdapter[K]) Name() string {
	return "Emitter"
}

func (a *emitterAdapter[K]) SetOutputChannel(output chan<- K) {
	a.output = output
}

func (a *emitterAdapter[K]) Run(ctx context.Context) error {
	defer func() {
		close(a.output)
	}()

	return a.emitter.Run(ctx, a.output)
}
