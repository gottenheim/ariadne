package pipeline

import "sync"

type Pipeline struct {
	waitGroup    sync.WaitGroup
	eventHandler *pipelineEventHandler
	elements     []elementAdapter
}

func New() *Pipeline {
	return &Pipeline{}
}

func (p *Pipeline) SyncRun() {
	for i := len(p.elements) - 1; i >= 0; i-- {
		element := p.elements[i]
		p.waitGroup.Add(1)
		go p.runElement(element)
	}
	p.wait()
}

func (p *Pipeline) attach(element elementAdapter) {
	p.elements = append(p.elements, element)
}

func (p *Pipeline) runElement(element elementAdapter) {
	element.Run()
	p.waitGroup.Done()
}

func (p *Pipeline) wait() {
	p.waitGroup.Wait()
}
