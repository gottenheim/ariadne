package pipeline

import "context"

type producer[K interface{}] interface {
	SetOutputChannel(output chan<- K)
	Run(ctx context.Context) error
}

type consumer[T interface{}] interface {
	SetInputChannel(input <-chan T)
}

func NewEmitter[T interface{}](pipeline *Pipeline, emitter Emitter[T]) *emitterAdapter[T] {
	return newEmitterAdapter(pipeline, emitter)
}

func WithFilter[T interface{}, K interface{}](pipeline *Pipeline, producer producer[T], filter Filter[T, K]) *filterAdapter[T, K] {
	filterAdapter := newFilterAdapter(pipeline, filter)
	ch := make(chan T)
	producer.SetOutputChannel(ch)
	filterAdapter.SetInputChannel(ch)
	return filterAdapter
}

func WithCondition[T interface{}, K interface{}, L interface{}](pipeline *Pipeline, producer producer[T], condition Condition[T, K, L]) *conditionAdapter[T, K, L] {
	conditionAdapter := newConditionAdapter(pipeline, condition)
	ch := make(chan T)
	producer.SetOutputChannel(ch)
	conditionAdapter.SetInputChannel(ch)
	return conditionAdapter
}

func OnPositiveDecision[T interface{}, K interface{}, L interface{}](condition *conditionAdapter[T, K, L]) producer[K] {
	return newPositiveDecisionAdapter(condition)
}

func OnNegativeDecision[T interface{}, K interface{}, L interface{}](condition *conditionAdapter[T, K, L]) producer[L] {
	return newNegativeDecisionAdapter(condition)
}

func WithAcceptor[T interface{}](pipeline *Pipeline, producer producer[T], acceptor Acceptor[T]) *acceptorAdapter[T] {
	acceptorAdapter := newAcceptorAdapter(pipeline, acceptor)
	ch := make(chan T)
	producer.SetOutputChannel(ch)
	acceptorAdapter.SetInputChannel(ch)
	return acceptorAdapter
}

func WithAggregator[T interface{}, K interface{}](pipeline *Pipeline, leftProducer producer[T], rightProducer producer[T], aggregator Aggregator[T, K]) *aggregatorAdapter[T, K] {
	aggregatorAdapter := newAggregatorAdapter(pipeline, aggregator)

	leftChannel := make(chan T)
	leftProducer.SetOutputChannel(leftChannel)
	aggregatorAdapter.SetLeftArgChannel(leftChannel)

	rightChannel := make(chan T)
	rightProducer.SetOutputChannel(rightChannel)
	aggregatorAdapter.SetRightArgChannel(rightChannel)

	return aggregatorAdapter
}
