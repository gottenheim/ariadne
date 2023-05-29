package pipeline

import "context"

type positiveDecisionAdapter[T interface{}, K interface{}, L interface{}] struct {
	conditionAdapter *conditionAdapter[T, K, L]
}

func newPositiveDecisionAdapter[T interface{}, K interface{}, L interface{}](conditionAdapter *conditionAdapter[T, K, L]) producer[K] {
	return &positiveDecisionAdapter[T, K, L]{
		conditionAdapter: conditionAdapter,
	}
}

func (f *positiveDecisionAdapter[T, K, L]) SetOutputChannel(output chan<- K) {
	f.conditionAdapter.SetPositiveDecisionChannel(output)
}

func (f *positiveDecisionAdapter[T, K, L]) Run(ctx context.Context) error {
	return f.conditionAdapter.Run(ctx)
}

type negativeDecisionAdapter[T interface{}, K interface{}, L interface{}] struct {
	conditionAdapter *conditionAdapter[T, K, L]
}

func newNegativeDecisionAdapter[T interface{}, K interface{}, L interface{}](conditionAdapter *conditionAdapter[T, K, L]) producer[L] {
	return &negativeDecisionAdapter[T, K, L]{
		conditionAdapter: conditionAdapter,
	}
}

func (f *negativeDecisionAdapter[T, K, L]) SetOutputChannel(output chan<- L) {
	f.conditionAdapter.SetNegativeDecisionChannel(output)
}

func (f *negativeDecisionAdapter[T, K, L]) Run(ctx context.Context) error {
	return f.conditionAdapter.Run(ctx)
}
