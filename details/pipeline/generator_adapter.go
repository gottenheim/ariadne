package pipeline

type generatorAdapter[K interface{}] struct {
	output    chan<- K
	generator Generator[K]
}

func newGeneratorAdapter[K interface{}](generator Generator[K]) *generatorAdapter[K] {
	return &generatorAdapter[K]{
		generator: generator,
	}
}

func (f *generatorAdapter[K]) SetOutputChannel(output chan<- K) {
	f.output = output
}

func (f *generatorAdapter[K]) Run() {
	go f.generator.Run(f.output)
}
