package pipeline

import "context"

type task interface {
	Name() string
	Run(ctx context.Context) error
}
