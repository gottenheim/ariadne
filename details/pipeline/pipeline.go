package pipeline

import "sync"

type Pipeline struct {
	waitGroup  sync.WaitGroup
	tasks      []task
	taskErrors []error
}

func New() *Pipeline {
	return &Pipeline{}
}

func (p *Pipeline) SyncRun() error {
	p.taskErrors = make([]error, len(p.tasks))
	for i := len(p.tasks) - 1; i >= 0; i-- {
		p.waitGroup.Add(1)
		go p.runTask(i)
	}
	p.wait()
	return p.getTaskError()
}

func (p *Pipeline) Cancel() {
	// It's enough to cancel emitter
	p.tasks[0].Cancel()
}

func (p *Pipeline) attach(task task) {
	p.tasks = append(p.tasks, task)
}

func (p *Pipeline) runTask(taskIndex int) {
	err := task.Run(p.tasks[taskIndex])
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
