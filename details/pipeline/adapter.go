package pipeline

type producerAdapter[K interface{}] interface {
	SetOutputChannel(output chan<- K)
	Run()
}

type consumerAdapter[T interface{}] interface {
	SetInputChannel(input <-chan T)
}
