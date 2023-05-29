package pipeline

func New[T interface{}](generator Generator[T]) *generatorAdapter[T] {
	return newGeneratorAdapter(generator)
}

func Join[T interface{}, K interface{}](leftAdapter producerAdapter[T], rightFilter Filter[T, K]) *filterAdapter[T, K] {
	rightAdapter := newFilterAdapter(rightFilter)
	ch := make(chan T)
	leftAdapter.SetOutputChannel(ch)
	rightAdapter.SetInputChannel(ch)
	rightAdapter.SetProducerFilter(leftAdapter)
	return rightAdapter
}
