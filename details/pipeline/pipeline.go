package pipeline

import (
	"context"
	"sync"
)

type Pipeline struct {
	waitGroup  sync.WaitGroup
	cancelFunc context.CancelFunc
	tasks      []task
	taskErrors []error
}

func New() *Pipeline {
	return &Pipeline{}
}

func (p *Pipeline) SyncRun() error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	p.cancelFunc = cancelFunc

	p.taskErrors = make([]error, len(p.tasks))
	for i := len(p.tasks) - 1; i >= 0; i-- {
		p.waitGroup.Add(1)
		go p.runTask(ctx, i)
	}
	p.wait()
	return p.getTaskError()
}

func (p *Pipeline) Cancel() {
	p.cancelFunc()
}

func (p *Pipeline) attach(task task) {
	p.tasks = append(p.tasks, task)
}

func (p *Pipeline) runTask(ctx context.Context, taskIndex int) {
	task := p.tasks[taskIndex]
	err := task.Run(ctx)
	p.taskErrors[taskIndex] = err
	p.waitGroup.Done()
}

func (p *Pipeline) wait() {
	p.waitGroup.Wait()
}

func (p *Pipeline) getTaskError() error {
	for _, err := range p.taskErrors {
		if err != nil {
			return err
		}
	}
	return nil
}
