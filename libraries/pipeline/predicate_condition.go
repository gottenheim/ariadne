package pipeline

import (
	"context"
)

type predicateCondition[T interface{}] struct {
	predicate func(T) bool
}

func NewPredicateCondition[T interface{}](predicate func(T) bool) Condition[T, T, T] {
	return &predicateCondition[T]{
		predicate: predicate,
	}
}

func (c *predicateCondition[T]) Run(ctx context.Context, input <-chan T, positiveDecision chan<- T, negativeDecision chan<- T) error {
	for {
		val, ok := <-input

		if !ok {
			break
		}

		if c.predicate(val) {
			positiveDecision <- val
		} else {
			negativeDecision <- val
		}
	}

	return nil
}
