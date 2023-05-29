package pipeline

func NewGenerator[T interface{}](generator Generator[T]) *generatorAdapter[T] {
	return newGeneratorAdapter(generator)
}

func WithFilter[T interface{}, K interface{}](leftAdapter producerAdapter[T], rightFilter Filter[T, K]) *filterAdapter[T, K] {
	rightAdapter := newFilterAdapter(rightFilter)
	ch := make(chan T)
	leftAdapter.SetOutputChannel(ch)
	rightAdapter.SetInputChannel(ch)
	rightAdapter.SetProducerFilter(leftAdapter)
	return rightAdapter
}

func WithCondition[T interface{}, K interface{}](leftAdapter producerAdapter[T], rightCondition Condition[T, K]) *conditionAdapter[T, K] {
	rightAdapter := newConditionAdapter(rightCondition)
	ch := make(chan T)
	leftAdapter.SetOutputChannel(ch)
	rightAdapter.SetInputChannel(ch)
	rightAdapter.SetProducerFilter(leftAdapter)
	return rightAdapter
}

func OnPositiveDecision[T interface{}, K interface{}](condition *conditionAdapter[T, K]) producerAdapter[K] {
	return newPositiveDecisionAdapter(condition)
}

func OnNegativeDecision[T interface{}, K interface{}](condition *conditionAdapter[T, K]) producerAdapter[K] {
	return newNegativeDecisionAdapter(condition)
}
