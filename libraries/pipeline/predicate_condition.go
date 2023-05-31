package pipeline

import (
	"context"
)

type predicateCondition[T interface{}] struct {
	predicate func(T) (bool, error)
}

func NewPredicateCondition[T interface{}](predicate func(T) (bool, error)) Condition[T, T, T] {
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

		yes, err := c.predicate(val)

		if err != nil {
			return err
		}

		if yes {
			if !WriteToChannel[T](ctx, positiveDecision, val) {
				break
			}
		} else {
			if !WriteToChannel[T](ctx, negativeDecision, val) {
				break
			}
		}
	}

	return nil
}
