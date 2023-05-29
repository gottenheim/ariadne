package pipeline

type positiveDecisionAdapter[T interface{}, K interface{}] struct {
	conditionAdapter *conditionAdapter[T, K]
}

func newPositiveDecisionAdapter[T interface{}, K interface{}](conditionAdapter *conditionAdapter[T, K]) producer[K] {
	return &positiveDecisionAdapter[T, K]{
		conditionAdapter: conditionAdapter,
	}
}

func (f *positiveDecisionAdapter[T, K]) SetOutputChannel(output chan<- K) {
	f.conditionAdapter.SetPositiveDecisionChannel(output)
}

func (f *positiveDecisionAdapter[T, K]) Run() error {
	return f.conditionAdapter.Run()
}

type negativeDecisionAdapter[T interface{}, K interface{}] struct {
	conditionAdapter *conditionAdapter[T, K]
}

func newNegativeDecisionAdapter[T interface{}, K interface{}](conditionAdapter *conditionAdapter[T, K]) producer[K] {
	return &negativeDecisionAdapter[T, K]{
		conditionAdapter: conditionAdapter,
	}
}

func (f *negativeDecisionAdapter[T, K]) SetOutputChannel(output chan<- K) {
	f.conditionAdapter.SetNegativeDecisionChannel(output)
}

func (f *negativeDecisionAdapter[T, K]) Run() error {
	return f.conditionAdapter.Run()
}
