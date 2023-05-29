package pipeline

type Condition[T interface{}, K interface{}] interface {
	Run(input <-chan T, positiveDecision chan<- K, negativeDecision chan<- K)
}
