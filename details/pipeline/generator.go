package pipeline

type Generator[K interface{}] interface {
	Run(output chan<- K)
}
