package pipeline

type Acceptor[T interface{}] interface {
	Run(input <-chan T) error
}
