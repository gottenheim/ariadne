package pipeline

import "context"

type task interface {
	Run(ctx context.Context) error
}
