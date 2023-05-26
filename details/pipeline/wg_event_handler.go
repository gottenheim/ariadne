package pipeline

import (
	"sync"
)

type WaitGroupEventHandler struct {
	waitGroup sync.WaitGroup
	errors    []error
}

func (h *WaitGroupEventHandler) OnStart() {
	h.waitGroup.Add(1)
}

func (h *WaitGroupEventHandler) OnFinish() {
	h.waitGroup.Done()
}

func (h *WaitGroupEventHandler) OnError(err error) {
	h.errors = append(h.errors, err)
}

func (h *WaitGroupEventHandler) Wait() {
	h.waitGroup.Wait()
}
