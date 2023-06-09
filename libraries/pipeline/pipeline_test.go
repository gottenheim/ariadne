package pipeline_test

import (
	"context"
	"errors"
	"math/rand"
	"testing"

	"github.com/gottenheim/ariadne/libraries/pipeline"
)

type serialNumbers struct {
	count int
}

func generateNSerialNumbers(count int) pipeline.Emitter[int] {
	return &serialNumbers{
		count: count,
	}
}

func (f *serialNumbers) Run(ctx context.Context, output chan<- int) error {
	for i := 0; i < f.count; i++ {
		if !pipeline.WriteToChannel[int](ctx, output, i) {
			break
		}
	}
	return nil
}

type randomNumberGenerator struct {
	count int
}

func generateNRandomNumbers(count int) pipeline.Emitter[int] {
	return &randomNumberGenerator{
		count: count,
	}
}

func (f *randomNumberGenerator) Run(ctx context.Context, output chan<- int) error {
	for i := 0; i < f.count; i++ {
		if !pipeline.WriteToChannel[int](ctx, output, rand.Int()%100) {
			break
		}
	}
	return nil
}

type lessThanFiftyFilter struct {
}

func filterLessThanFifty() pipeline.Filter[int, int] {
	return &lessThanFiftyFilter{}
}

func (f *lessThanFiftyFilter) Run(ctx context.Context, input <-chan int, output chan<- int) error {
	for {
		val, ok := <-input
		if !ok {
			break
		}

		if val < 50 {
			if !pipeline.WriteToChannel[int](ctx, output, val) {
				break
			}
		}
	}
	return nil
}

type lessThanFiftyCondition struct{}

func ifLessThanFifty() pipeline.Condition[int, int, int] {
	return &lessThanFiftyCondition{}
}

func (f *lessThanFiftyCondition) Run(ctx context.Context, input <-chan int, positiveDecision chan<- int, negativeDecision chan<- int) error {
	for {
		val, ok := <-input
		if !ok {
			break
		}

		if val < 50 {
			if !pipeline.WriteToChannel[int](ctx, positiveDecision, val) {
				break
			}
		} else {
			if !pipeline.WriteToChannel[int](ctx, negativeDecision, val) {
				break
			}
		}
	}
	return nil
}

type failureGenerator struct {
}

func generateFailure() pipeline.Filter[int, int] {
	return &failureGenerator{}
}

func (f *failureGenerator) Run(ctx context.Context, input <-chan int, output chan<- int) error {
	for {
		val, ok := <-input
		if !ok {
			break
		}

		if val > 30 {
			return errors.New("Number processing failure")
		}

		if !pipeline.WriteToChannel[int](ctx, output, val) {
			break
		}
	}
	return nil
}

type limitNumberCountCondition struct {
	pipeline *pipeline.Pipeline
	limit    int
	actual   int
}

func limitNumberCount(pipeline *pipeline.Pipeline, limit int) pipeline.Filter[int, int] {
	return &limitNumberCountCondition{
		pipeline: pipeline,
		limit:    limit,
	}
}

func (f *limitNumberCountCondition) Run(ctx context.Context, input <-chan int, output chan<- int) error {
	for {
		val, ok := <-input
		if !ok {
			break
		}

		if !pipeline.WriteToChannel[int](ctx, output, val) {
			break
		}

		f.actual++

		if f.actual >= f.limit {
			f.pipeline.Cancel()
			break
		}
	}
	return nil
}

type pipelineResult struct {
	numbers []int
}

func TestGeneratorsAndFilters(t *testing.T) {
	p := pipeline.New()

	collector := pipeline.NewItemCollector[int]()

	generator := pipeline.NewEmitter(p, generateNSerialNumbers(100))
	filter := pipeline.WithFilter[int](p, generator, filterLessThanFifty())
	pipeline.WithAcceptor[int](p, filter, collector)

	p.SyncRun()

	if len(collector.Items) != 50 {
		t.Fatal("Pipeline must produce 50 numbers")
	}

	for i := 0; i < 50; i++ {
		if collector.Items[i] != i {
			t.Fatal("Pipeline contains an unexpected number")
		}
	}
}

func TestApplyingConditions(t *testing.T) {
	p := pipeline.New()

	positiveResult, negativeResult := pipeline.NewItemCollector[int](), pipeline.NewItemCollector[int]()

	randomNumbers := pipeline.NewEmitter(p, generateNRandomNumbers(1000))

	lessThanFifty := pipeline.WithCondition[int](p, randomNumbers, ifLessThanFifty())

	pipeline.WithAcceptor[int](p, pipeline.OnPositiveDecision(lessThanFifty), positiveResult)

	pipeline.WithAcceptor[int](p, pipeline.OnNegativeDecision(lessThanFifty), negativeResult)

	err := p.SyncRun()
	if err != nil {
		t.Fatal(err)
	}

	if len(positiveResult.Items)+len(negativeResult.Items) != 1000 {
		t.Fatal("Pipeline must produce 1000 numbers")
	}

	for i := range positiveResult.Items {
		if positiveResult.Items[i] >= 50 {
			t.Fatal("Number must be less than fifty")
		}
	}

	for i := range negativeResult.Items {
		if negativeResult.Items[i] < 50 {
			t.Fatal("Number must be more than fifty")
		}
	}
}

func TestPipelineFailure(t *testing.T) {
	p := pipeline.New()

	collector := pipeline.NewItemCollector[int]()

	emitter := pipeline.NewEmitter(p, generateNSerialNumbers(100))
	filter := pipeline.WithFilter[int](p, emitter, generateFailure())
	pipeline.WithAcceptor[int](p, filter, collector)

	err := p.SyncRun()

	if err == nil {
		t.Fatal("Pipeline should return error")
	}

	if err.Error() != "Number processing failure" {
		t.Fatal("Unexpected error")
	}
}

func TestPipelineCancellation(t *testing.T) {
	p := pipeline.New()

	collector := pipeline.NewItemCollector[int]()

	generator := pipeline.NewEmitter(p, generateNSerialNumbers(100))
	filter := pipeline.WithFilter[int](p, generator, limitNumberCount(p, 30))
	pipeline.WithAcceptor[int](p, filter, collector)

	p.SyncRun()

	if len(collector.Items) != 30 {
		t.Fatal("Pipeline must produce 30 numbers")
	}
}

func TestLimitFilter(t *testing.T) {
	p := pipeline.New()

	complete := pipeline.NewPassingItemCollector[int]()
	partial := pipeline.NewItemCollector[int]()

	generator := pipeline.NewEmitter(p, generateNSerialNumbers(100))
	passingStep := pipeline.WithFilter[int, int](p, generator, complete)
	skippingStep := pipeline.WithFilter[int](p, passingStep, pipeline.Limit[int](30))
	pipeline.WithAcceptor[int](p, skippingStep, partial)

	p.SyncRun()

	if len(partial.Items) != 30 {
		t.Fatal("Skipping limit filter must pass only 30 records")
	}

	if len(complete.Items) != 100 {
		t.Fatal("Generator must produce 100 records")
	}
}

func TestPredicateCondition(t *testing.T) {
	p := pipeline.New()

	positiveResult, negativeResult := pipeline.NewItemCollector[int](), pipeline.NewItemCollector[int]()

	randomNumbers := pipeline.NewEmitter(p, generateNRandomNumbers(1000))

	predicate := func(val int) (bool, error) {
		return val < 50, nil
	}

	lessThanFifty := pipeline.WithCondition[int](p, randomNumbers, pipeline.NewPredicateCondition(predicate))

	pipeline.WithAcceptor[int](p, pipeline.OnPositiveDecision(lessThanFifty), positiveResult)

	pipeline.WithAcceptor[int](p, pipeline.OnNegativeDecision(lessThanFifty), negativeResult)

	err := p.SyncRun()

	if err != nil {
		t.Fatal(err)
	}

	if len(positiveResult.Items)+len(negativeResult.Items) != 1000 {
		t.Fatal("Pipeline must produce 1000 numbers")
	}

	for i := range positiveResult.Items {
		if positiveResult.Items[i] >= 50 {
			t.Fatal("Number must be less than fifty")
		}
	}

	for i := range negativeResult.Items {
		if negativeResult.Items[i] < 50 {
			t.Fatal("Number must be more than fifty")
		}
	}
}

func TestCountingFilter(t *testing.T) {
	p := pipeline.New()

	valueStore := pipeline.NewValueStore[int]()

	generator := pipeline.NewEmitter(p, generateNSerialNumbers(100))
	counting := pipeline.WithFilter[int, int](p, generator, pipeline.NewCounter[int]())
	pipeline.WithAcceptor[int](p, counting, valueStore)

	p.SyncRun()

	if valueStore.Value() != 100 {
		t.Fatal("Counter must count 100 elements")
	}
}

func TestSumCalculation(t *testing.T) {
	p := pipeline.New()

	valueStore := pipeline.NewValueStore[int]()

	generator1 := pipeline.NewEmitter(p, generateNSerialNumbers(100))
	counting1 := pipeline.WithFilter[int, int](p, generator1, pipeline.NewCounter[int]())

	generator2 := pipeline.NewEmitter(p, generateNSerialNumbers(50))
	counting2 := pipeline.WithFilter[int, int](p, generator2, pipeline.NewCounter[int]())

	sumCalculator := pipeline.WithAggregator[int](p, counting1, counting2, pipeline.SumCalculator())

	pipeline.WithAcceptor[int](p, sumCalculator, valueStore)

	err := p.SyncRun()

	if err != nil {
		t.Fatal(err)
	}

	if valueStore.Value() != 150 {
		t.Fatalf("Wrong sum calculated. Expected 150, actual %d", valueStore.Value())
	}
}

func TestDevNull(t *testing.T) {
	p := pipeline.New()

	generator := pipeline.NewEmitter(p, generateNSerialNumbers(100))

	pipeline.WithAcceptor[int](p, generator, pipeline.Skip[int]())

	err := p.SyncRun()

	if err != nil {
		t.Fatal(err)
	}
}
