package pipeline

type Generator[K interface{}] interface {
	Run(output chan<- K)
}

type Filter[T interface{}, K interface{}] interface {
	Run(input <-chan T, output chan<- K)
}
