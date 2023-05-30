package pipeline

import (
	"context"
	"errors"
)

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

func (f *acceptorAdapter[T]) Name() string {
	return "Acceptor"
}

func (f *acceptorAdapter[T]) SetInputChannel(input <-chan T) {
	f.input = input
}

func (f *acceptorAdapter[T]) Run(ctx context.Context) error {
	if f.input == nil {
		return errors.New("Acceptor input channel is not set")
	}

	return f.acceptor.Run(ctx, f.input)
}
