package pipeline

type filterAdapter[T interface{}, K interface{}] struct {
	input    <-chan T
	output   chan<- K
	filter   Filter[T, K]
	producer producerAdapter[T]
}

func newFilterAdapter[T interface{}, K interface{}](filter Filter[T, K]) *filterAdapter[T, K] {
	return &filterAdapter[T, K]{
		filter: filter,
	}
}

func (f *filterAdapter[T, K]) SetInputChannel(input <-chan T) {
	f.input = input
}

func (f *filterAdapter[T, K]) SetOutputChannel(output chan<- K) {
	f.output = output
}

func (f *filterAdapter[T, K]) SetProducerFilter(producer producerAdapter[T]) {
	f.producer = producer
}

func (f *filterAdapter[T, K]) Run() {
	go f.filter.Run(f.input, f.output)

	if f.producer != nil {
		f.producer.Run()
	}
}
