package pipeline_test

import (
	"testing"

	"github.com/gottenheim/ariadne/details/pipeline"
)

type serialNumbers struct {
	events pipeline.FilterEvents
}

func newSerialNumbers(events pipeline.FilterEvents) pipeline.Filter[interface{}, int] {
	events.OnStart()

	return &serialNumbers{
		events: events,
	}
}

func (f *serialNumbers) Run(input <-chan interface{}, output chan<- int) {
	for i := 0; i < 100; i++ {
		output <- i
	}
	close(output)
	f.events.OnFinish()
}

type lessThanFifty struct {
	events pipeline.FilterEvents
}

func newLessThanFifty(events pipeline.FilterEvents) pipeline.Filter[int, int] {
	events.OnStart()

	return &lessThanFifty{
		events: events,
	}
}

func (f *lessThanFifty) Run(input <-chan int, output chan<- int) {
	for {
		val, ok := <-input
		if !ok {
			break
		}

		if val < 50 {
			output <- val
		}
	}
	close(output)
	f.events.OnFinish()
}

type pipelineResult struct {
	numbers []int
}

type collectNumbers struct {
	events pipeline.FilterEvents
	result *pipelineResult
}

func newCollectNumbers(events pipeline.FilterEvents, result *pipelineResult) pipeline.Filter[int, int] {
	events.OnStart()

	return &collectNumbers{
		events: events,
		result: result,
	}
}

func (f *collectNumbers) Run(input <-chan int, output chan<- int) {
	for {
		val, ok := <-input
		if !ok {
			break
		}

		f.result.numbers = append(f.result.numbers, val)
	}
	f.events.OnFinish()
}

func TestFilterProducingAndCollectingFilteredIntegers(t *testing.T) {
	wgEvents := &pipeline.WaitGroupEventHandler{}

	result := &pipelineResult{}

	pipeLine := pipeline.Join[int](
		pipeline.Join[int](
			pipeline.New(
				newSerialNumbers(wgEvents)),
			newLessThanFifty(wgEvents)),
		newCollectNumbers(wgEvents, result))

	pipeLine.Run()

	wgEvents.Wait()

	if len(result.numbers) != 50 {
		t.Fatal("Pipeline must produce 50 numbers")
	}

	for i := 0; i < 50; i++ {
		if result.numbers[i] != i {
			t.Fatal("Pipeline contains a wrong number")
		}
	}
}
