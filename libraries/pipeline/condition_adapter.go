package pipeline

import "context"

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

func (f *conditionAdapter[T, K, L]) Name() string {
	return "Condition"
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

func (f *conditionAdapter[T, K, L]) Run(ctx context.Context) error {
	defer func() {
		close(f.positiveDecision)
		close(f.negativeDecision)
	}()

	return f.condition.Run(ctx, f.input, f.positiveDecision, f.negativeDecision)
}
