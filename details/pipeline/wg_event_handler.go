package pipeline

import "sync"

type WaitGroupEventHandler struct {
	waitGroup sync.WaitGroup
}

func (h *WaitGroupEventHandler) OnStart() {
	h.waitGroup.Add(1)
}

func (h *WaitGroupEventHandler) OnFinish() {
	h.waitGroup.Done()
}

func (h *WaitGroupEventHandler) Wait() {
	h.waitGroup.Wait()
}
