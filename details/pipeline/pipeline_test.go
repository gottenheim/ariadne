package pipeline_test

import (
	"errors"
	"math/rand"
	"testing"

	"github.com/gottenheim/ariadne/details/pipeline"
)

type serialNumbers struct {
	count int
}

func generateNSerialNumbers(count int) pipeline.Generator[int] {
	return &serialNumbers{
		count: count,
	}
}

func (f *serialNumbers) Run(output chan<- int) error {
	for i := 0; i < f.count; i++ {
		output <- i
	}
	close(output)
	return nil
}

type randomNumberGenerator struct {
	count int
}

func generateNRandomNumbers(count int) pipeline.Generator[int] {
	return &randomNumberGenerator{
		count: count,
	}
}

func (f *randomNumberGenerator) Run(output chan<- int) error {
	for i := 0; i < f.count; i++ {
		output <- rand.Int() % 100
	}
	close(output)
	return nil
}

type lessThanFiftyFilter struct {
}

func filterLessThanFifty() pipeline.Filter[int, int] {
	return &lessThanFiftyFilter{}
}

func (f *lessThanFiftyFilter) Run(input <-chan int, output chan<- int) error {
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
	return nil
}

type lessThanFiftyCondition struct{}

func ifLessThanFifty() pipeline.Condition[int, int] {
	return &lessThanFiftyCondition{}
}

func (f *lessThanFiftyCondition) Run(input <-chan int, positiveDecision chan<- int, negativeDecision chan<- int) error {
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
	return nil
}

type failureGenerator struct {
}

func generateFailure() pipeline.Generator[int] {
	return &failureGenerator{}
}

func (f *failureGenerator) Run(output chan<- int) error {
	return errors.New("Failed to generate numbers")
}

type pipelineResult struct {
	numbers []int
}

type numberCollector struct {
	result *pipelineResult
}

func collectNumbers(result *pipelineResult) pipeline.Filter[int, int] {
	return &numberCollector{
		result: result,
	}
}

func (f *numberCollector) Run(input <-chan int, output chan<- int) error {
	for {
		val, ok := <-input
		if !ok {
			break
		}

		f.result.numbers = append(f.result.numbers, val)
	}
	return nil
}

func TestGeneratorsAndFilters(t *testing.T) {
	p := pipeline.New()

	result := &pipelineResult{}

	generator := pipeline.NewGenerator(p, generateNSerialNumbers(100))
	filter := pipeline.WithFilter[int](p, generator, filterLessThanFifty())
	pipeline.WithFilter[int](p, filter, collectNumbers(result))

	p.SyncRun()

	if len(result.numbers) != 50 {
		t.Fatal("Pipeline must produce 50 numbers")
	}

	for i := 0; i < 50; i++ {
		if result.numbers[i] != i {
			t.Fatal("Pipeline contains an unexpected number")
		}
	}
}

func TestApplyingConditions(t *testing.T) {
	p := pipeline.New()

	positiveResult, negativeResult := &pipelineResult{}, &pipelineResult{}

	randomNumbers := pipeline.NewGenerator(p, generateNRandomNumbers(1000))

	lessThanFifty := pipeline.WithCondition[int](p, randomNumbers, ifLessThanFifty())

	pipeline.WithFilter(p, pipeline.OnPositiveDecision(lessThanFifty), collectNumbers(positiveResult))

	pipeline.WithFilter(p, pipeline.OnNegativeDecision(lessThanFifty), collectNumbers(negativeResult))

	p.SyncRun()

	if len(positiveResult.numbers)+len(negativeResult.numbers) != 1000 {
		t.Fatal("Pipeline must produce 1000 numbers")
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

func TestPipelineFailure(t *testing.T) {
	p := pipeline.New()

	pipeline.NewGenerator(p, generateFailure())

	err := p.SyncRun()

	if err == nil {
		t.Fatal("Pipeline should return error")
	}

	if err.Error() != "Failed to generate numbers" {
		t.Fatal("Unexpected error")
	}
}
