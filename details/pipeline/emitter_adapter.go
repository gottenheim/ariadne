package pipeline

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

func (a *emitterAdapter[K]) SetOutputChannel(output chan<- K) {
	a.output = output
}

func (a *emitterAdapter[K]) Run() error {
	return a.emitter.Run(a.output)
}
