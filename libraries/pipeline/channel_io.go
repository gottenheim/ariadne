package pipeline

import "golang.org/x/net/context"

func WriteToChannel[T interface{}](ctx context.Context, output chan<- T, value T) bool {
	select {
	case <-ctx.Done():
		return false
	case output <- value:
		return true
	}
}
