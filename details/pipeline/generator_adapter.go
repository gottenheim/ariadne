package pipeline

type generatorAdapter[K interface{}] struct {
	output    chan<- K
	generator Generator[K]
	pipeline  *Pipeline
}

func newGeneratorAdapter[K interface{}](pipeline *Pipeline, generator Generator[K]) *generatorAdapter[K] {
	adapter := &generatorAdapter[K]{
		generator: generator,
		pipeline:  pipeline,
	}

	pipeline.attach(adapter)
	return adapter
}

func (a *generatorAdapter[K]) SetOutputChannel(output chan<- K) {
	a.output = output
}

func (a *generatorAdapter[K]) Run() {
	a.generator.Run(a.output)
}
