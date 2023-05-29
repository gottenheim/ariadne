package pipeline

type conditionAdapter[T interface{}, K interface{}, L interface{}] struct {
	input            <-chan T
	positiveDecision chan<- K
	negativeDecision chan<- L
	condition        Condition[T, K, L]
}

func newConditionAdapter[T interface{}, K interface{}, L interface{}](pipeline *Pipeline, condition Condition[T, K, L]) *conditionAdapter[T, K, L] {
	adapter := &conditionAdapter[T, K, L]{
		condition: condition,
	}

	pipeline.attach(adapter)

	return adapter
}

func (f *conditionAdapter[T, K, L]) SetInputChannel(input <-chan T) {
	f.input = input
}

func (f *conditionAdapter[T, K, L]) SetPositiveDecisionChannel(positiveDecision chan<- K) {
	f.positiveDecision = positiveDecision
}

func (f *conditionAdapter[T, K, L]) SetNegativeDecisionChannel(negativeDecision chan<- L) {
	f.negativeDecision = negativeDecision
}

func (f *conditionAdapter[T, K, L]) Run() error {
	defer func() {
		f.closeOutputChannels()
	}()

	return f.condition.Run(f.input, f.positiveDecision, f.negativeDecision)
}

func (f *conditionAdapter[T, L, K]) Cancel() {
	f.closeOutputChannels()
}

func (f *conditionAdapter[T, L, K]) closeOutputChannels() {
	if f.positiveDecision != nil {
		positiveDecision := f.positiveDecision
		f.positiveDecision = nil
		close(positiveDecision)
	}
	if f.negativeDecision != nil {
		negativeDecision := f.negativeDecision
		f.negativeDecision = nil
		close(negativeDecision)
	}
}
