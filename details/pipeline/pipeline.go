package pipeline

func New[T interface{}](filter Filter[interface{}, T]) *filterAdapter[interface{}, T] {
	return newFilterAdapter(filter)
}

func Join[T interface{}, K interface{}](leftAdapter producerFilterAdapter[T], rightFilter Filter[T, K]) *filterAdapter[T, K] {
	rightAdapter := newFilterAdapter(rightFilter)
	ch := make(chan T)
	leftAdapter.SetOutputChannel(ch)
	rightAdapter.SetInputChannel(ch)
	rightAdapter.SetProducerFilter(leftAdapter)
	return rightAdapter
}
