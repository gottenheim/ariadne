package pipeline_test

import (
	"math/rand"
	"testing"

	"github.com/gottenheim/ariadne/details/pipeline"
)

type serialNumbers struct {
	events pipeline.FilterEvents
}

func newSerialNumbers(events pipeline.FilterEvents) pipeline.Generator[int] {
	events.OnStart()

	return &serialNumbers{
		events: events,
	}
}

func (f *serialNumbers) Run(output chan<- int) {
	for i := 0; i < 100; i++ {
		output <- i
	}
	close(output)
	f.events.OnFinish()
}

type randomNumbers struct {
	events pipeline.FilterEvents
}

func newRandomNumbers(events pipeline.FilterEvents) pipeline.Generator[int] {
	events.OnStart()

	return &randomNumbers{
		events: events,
	}
}

func (f *randomNumbers) Run(output chan<- int) {
	for i := 0; i < 100; i++ {
		output <- rand.Int() % 100
	}
	close(output)
	f.events.OnFinish()
}

type lessThanFiftyFilter struct {
	events pipeline.FilterEvents
}

func filterLessThanFifty(events pipeline.FilterEvents) pipeline.Filter[int, int] {
	events.OnStart()

	return &lessThanFiftyFilter{
		events: events,
	}
}

func (f *lessThanFiftyFilter) Run(input <-chan int, output chan<- int) {
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

type lessThanFiftyCondition struct {
	events pipeline.FilterEvents
}

func ifLessThanFifty(events pipeline.FilterEvents) pipeline.Condition[int, int] {
	events.OnStart()

	return &lessThanFiftyCondition{
		events: events,
	}
}

func (f *lessThanFiftyCondition) Run(input <-chan int, positiveDecision chan<- int, negativeDecision chan<- int) {
	for {
		val, ok := <-input
		if !ok {
			break
		}

		if val < 50 {
			positiveDecision <- val
		} else {
			negativeDecision <- val
		}
	}
	close(positiveDecision)
	close(negativeDecision)
	f.events.OnFinish()
}

type pipelineResult struct {
	numbers []int
}

type collectNumbers struct {
	events pipeline.FilterEvents
	result *pipelineResult
}

func newNumberCollector(events pipeline.FilterEvents, result *pipelineResult) pipeline.Filter[int, int] {
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

	pipeLine := pipeline.WithFilter[int](
		pipeline.WithFilter[int](
			pipeline.NewGenerator(
				newSerialNumbers(wgEvents)),
			filterLessThanFifty(wgEvents)),
		newNumberCollector(wgEvents, result))

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

func TestCondition(t *testing.T) {
	wgEvents := &pipeline.WaitGroupEventHandler{}

	positiveResult := &pipelineResult{}
	negativeResult := &pipelineResult{}

	randomNumbers := pipeline.NewGenerator(newRandomNumbers(wgEvents))

	lessThanFifty := pipeline.WithCondition[int](randomNumbers, ifLessThanFifty(wgEvents))

	positiveCollector := pipeline.WithFilter(pipeline.OnPositiveDecision(lessThanFifty), newNumberCollector(wgEvents, positiveResult))

	negativeCollector := pipeline.WithFilter(pipeline.OnNegativeDecision(lessThanFifty), newNumberCollector(wgEvents, negativeResult))

	positiveCollector.Run()
	negativeCollector.Run()

	wgEvents.Wait()

	if len(positiveResult.numbers)+len(negativeResult.numbers) != 100 {
		t.Fatal("Pipeline must produce 100 numbers")
	}

	for i := range positiveResult.numbers {
		if positiveResult.numbers[i] >= 50 {
			t.Fatal("Number must be less than fifty")
		}
	}

	for i := range negativeResult.numbers {
		if negativeResult.numbers[i] < 50 {
			t.Fatal("Number must be more than fifty")
		}
	}
}
