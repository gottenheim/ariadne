package pipeline

type conditionAdapter[T interface{}, K interface{}] struct {
	input            <-chan T
	positiveDecision chan<- K
	negativeDecision chan<- K
	condition        Condition[T, K]
	producer         producerAdapter[T]
	running          bool
}

func newConditionAdapter[T interface{}, K interface{}](condition Condition[T, K]) *conditionAdapter[T, K] {
	return &conditionAdapter[T, K]{
		condition: condition,
	}
}

func (f *conditionAdapter[T, K]) SetInputChannel(input <-chan T) {
	f.input = input
}

func (f *conditionAdapter[T, K]) SetPositiveDecisionChannel(positiveDecision chan<- K) {
	f.positiveDecision = positiveDecision
}

func (f *conditionAdapter[T, K]) SetNegativeDecisionChannel(negativeDecision chan<- K) {
	f.negativeDecision = negativeDecision
}

func (f *conditionAdapter[T, K]) SetProducerFilter(producer producerAdapter[T]) {
	f.producer = producer
}

func (f *conditionAdapter[T, K]) Run() {
	if !f.running {
		go f.condition.Run(f.input, f.positiveDecision, f.negativeDecision)

		if f.producer != nil {
			f.producer.Run()
		}

		f.running = true
	}
}
