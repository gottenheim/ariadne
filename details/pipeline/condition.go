package pipeline

type Condition[T interface{}, K interface{}, L interface{}] interface {
	Run(input <-chan T, positiveDecision chan<- K, negativeDecision chan<- L) error
}
