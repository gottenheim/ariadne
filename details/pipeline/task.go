package pipeline

type task interface {
	Run() error
}
