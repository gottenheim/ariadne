package pipeline

func NewGenerator[T interface{}](pipeline *Pipeline, generator Generator[T]) *generatorAdapter[T] {
	return newGeneratorAdapter(pipeline, generator)
}

func WithFilter[T interface{}, K interface{}](pipeline *Pipeline, leftAdapter producerAdapter[T], rightFilter Filter[T, K]) *filterAdapter[T, K] {
	rightAdapter := newFilterAdapter(pipeline, rightFilter)
	ch := make(chan T)
	leftAdapter.SetOutputChannel(ch)
	rightAdapter.SetInputChannel(ch)
	return rightAdapter
}

func WithCondition[T interface{}, K interface{}](pipeline *Pipeline, leftAdapter producerAdapter[T], rightCondition Condition[T, K]) *conditionAdapter[T, K] {
	rightAdapter := newConditionAdapter(pipeline, rightCondition)
	ch := make(chan T)
	leftAdapter.SetOutputChannel(ch)
	rightAdapter.SetInputChannel(ch)
	return rightAdapter
}

func OnPositiveDecision[T interface{}, K interface{}](condition *conditionAdapter[T, K]) producerAdapter[K] {
	return newPositiveDecisionAdapter(condition)
}

func OnNegativeDecision[T interface{}, K interface{}](condition *conditionAdapter[T, K]) producerAdapter[K] {
	return newNegativeDecisionAdapter(condition)
}
