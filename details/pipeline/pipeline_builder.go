package pipeline

type producer[K interface{}] interface {
	SetOutputChannel(output chan<- K)
	Run() error
}

type consumer[T interface{}] interface {
	SetInputChannel(input <-chan T)
}

func NewGenerator[T interface{}](pipeline *Pipeline, generator Generator[T]) *generatorAdapter[T] {
	return newGeneratorAdapter(pipeline, generator)
}

func WithFilter[T interface{}, K interface{}](pipeline *Pipeline, producer producer[T], filter Filter[T, K]) *filterAdapter[T, K] {
	filterAdapter := newFilterAdapter(pipeline, filter)
	ch := make(chan T)
	producer.SetOutputChannel(ch)
	filterAdapter.SetInputChannel(ch)
	return filterAdapter
}

func WithCondition[T interface{}, K interface{}](pipeline *Pipeline, producer producer[T], condition Condition[T, K]) *conditionAdapter[T, K] {
	conditionAdapter := newConditionAdapter(pipeline, condition)
	ch := make(chan T)
	producer.SetOutputChannel(ch)
	conditionAdapter.SetInputChannel(ch)
	return conditionAdapter
}

func OnPositiveDecision[T interface{}, K interface{}](condition *conditionAdapter[T, K]) producer[K] {
	return newPositiveDecisionAdapter(condition)
}

func OnNegativeDecision[T interface{}, K interface{}](condition *conditionAdapter[T, K]) producer[K] {
	return newNegativeDecisionAdapter(condition)
}
