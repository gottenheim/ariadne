package pipeline

type Filter[T interface{}, K interface{}] interface {
	Run(input <-chan T, output chan<- K) error
}
