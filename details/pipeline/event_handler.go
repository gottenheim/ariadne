package pipeline

type pipelineEventHandler struct {
	errors []error
}

func (h *pipelineEventHandler) OnError(err error) {
	h.errors = append(h.errors, err)
}
