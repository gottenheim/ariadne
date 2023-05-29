package pipeline

type conditionAdapter[T interface{}, K interface{}] struct {
	input            <-chan T
	positiveDecision chan<- K
	negativeDecision chan<- K
	condition        Condition[T, K]
}

func newConditionAdapter[T interface{}, K interface{}](pipeline *Pipeline, condition Condition[T, K]) *conditionAdapter[T, K] {
	adapter := &conditionAdapter[T, K]{
		condition: condition,
	}

	pipeline.attach(adapter)

	return adapter
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

func (f *conditionAdapter[T, K]) Run() {
	f.condition.Run(f.input, f.positiveDecision, f.negativeDecision)
}
