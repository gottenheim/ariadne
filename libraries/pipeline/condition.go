package pipeline

import "context"

type Condition[T interface{}, K interface{}, L interface{}] interface {
	Run(ctx context.Context, input <-chan T, positiveDecision chan<- K, negativeDecision chan<- L) error
}
