package pipeline

type acceptorAdapter[T interface{}] struct {
	input    <-chan T
	acceptor Acceptor[T]
}

func newAcceptorAdapter[T interface{}](p *Pipeline, acceptor Acceptor[T]) *acceptorAdapter[T] {
	adapter := &acceptorAdapter[T]{
		acceptor: acceptor,
	}

	p.attach(adapter)

	return adapter
}

func (f *acceptorAdapter[T]) SetInputChannel(input <-chan T) {
	f.input = input
}

func (f *acceptorAdapter[T]) Run() error {
	return f.acceptor.Run(f.input)
}
