package pipeline

import "context"

type Emitter[K interface{}] interface {
	Run(ctx context.Context, output chan<- K) error
}
