package pipeline

type Emitter[K interface{}] interface {
	Run(output chan<- K) error
}
