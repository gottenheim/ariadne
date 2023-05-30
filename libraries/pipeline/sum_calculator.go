package pipeline

import (
	"context"
)

type sumCalculator struct {
}

func SumCalculator() Aggregator[int, int] {
	return &sumCalculator{}
}

func (f *sumCalculator) Run(ctx context.Context, leftArg <-chan int, rightArg <-chan int, output chan<- int) error {
	leftVal := 0
	rightVal := 0
	leftOk, rightOk := true, true

	for {
		select {
		case val, ok := <-leftArg:
			if ok {
				leftVal = val
			} else {
				leftOk = false
			}
		case val, ok := <-rightArg:
			if ok {
				rightVal = val
			} else {
				rightOk = false
			}
		}

		if !leftOk && !rightOk {
			break
		}

		output <- leftVal + rightVal
	}

	return nil
}
