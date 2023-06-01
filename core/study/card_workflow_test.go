package study_test

import (
	"testing"
	"time"

	"github.com/gottenheim/ariadne/core/card"
	"github.com/gottenheim/ariadne/core/study"
	"github.com/gottenheim/ariadne/libraries/datetime"
)

const (
	Again int = 0
	Hard      = 3
	Good      = 4
	Easy      = 5
)

func TestCardWorkflow_ShouldAssignTenMinutesInterval_ForNewCard_IfAnyButtonIsPressed(t *testing.T) {
	grades := []int{0, 3, 4, 5}

	for _, grade := range grades {
		newCard := card.NewFakeCard().
			WithSection("languages/cpp").
			WithEntry("1").
			WithActivityChain(card.GenerateActivityChain(card.LearnCard)).
			Build()

		repetitionParams := pressButtonByGrade(t, newCard, grade)

		if repetitionParams.Interval != 10*time.Minute {
			t.Fatal("Ten minutes interval should be assigned")
		}
	}
}

func TestCardWorkflow_ShouldAssignOneDayInterval_ForLearnedNewCard(t *testing.T) {
	newCard := card.NewFakeCard().
		WithSection("languages/cpp").
		WithEntry("1").
		WithActivityChain(card.GenerateActivityChain(card.LearnCard)).
		Build()

	// New card shown for the first time
	pressGood(t, newCard)

	// Confirming knowledge after ten minutes
	repetitionParams := pressGood(t, newCard)

	if repetitionParams.Interval != 1*study.Day {
		t.Fatal("One day interval should be assigned for learned new card")
	}
}

func TestCardWorkflow_ShouldAssignSixDaysInterval_ForNewCard_ShownNextDayAfterLearning(t *testing.T) {
	newCard := card.NewFakeCard().
		WithSection("languages/cpp").
		WithEntry("1").
		WithActivityChain(card.GenerateActivityChain(card.LearnCard)).
		Build()

	// New card shown for the first time
	pressGood(t, newCard)

	// Confirming knowledge after ten minutes
	pressGood(t, newCard)

	// Confirming knowledge next day
	repetitionParams := pressGood(t, newCard)

	if repetitionParams.Interval != 6*study.Day {
		t.Fatal("Six days interval should be assigned for new card shown next day after learning")
	}
}

func TestCardWorkflow_ShouldAssignFifteenDaysInterval_IfCardCanBeStillRemembered_AfterSixDaysInterval(t *testing.T) {
	newCard := card.NewFakeCard().
		WithSection("languages/cpp").
		WithEntry("1").
		WithActivityChain(card.GenerateActivityChain(card.LearnCard)).
		Build()

	// New card shown for the first time
	pressGood(t, newCard)

	// Confirming knowledge after ten minutes
	pressGood(t, newCard)

	// Confirming knowledge next day
	pressGood(t, newCard)

	// Confirming knowledge after six days
	repetitionParams := pressGood(t, newCard)

	// Grade 4 makes easiness factor stable and always equal to 2.5
	// So to calculate interval we need to multiply previous interval to 2.5 and round to nearest integer
	// 6 * 2.5 = 15
	if repetitionParams.Interval != 15*study.Day {
		t.Fatal("Fifteen days interval should be assigned for card that can be remembered after six days interval")
	}
}

func TestCardWorkflow_ShouldAssignTenMinutesInterval_IfCardWasForgotten_AfterSixDaysInterval(t *testing.T) {
	newCard := card.NewFakeCard().
		WithSection("languages/cpp").
		WithEntry("1").
		WithActivityChain(card.GenerateActivityChain(card.LearnCard)).
		Build()

	// New card shown for the first time
	pressGood(t, newCard)

	// Confirming knowledge after ten minutes
	pressGood(t, newCard)

	// Confirming knowledge next day
	pressGood(t, newCard)

	// Failed to confirm knowledge after six days
	repetitionParams := pressAgain(t, newCard)

	if repetitionParams.Interval != 10*time.Minute {
		t.Fatal("Should start remembering card from scratch")
	}
}

func TestCardWorkflow_ShouldAssignTenDaysInterval_IfCardWasForgotten_ButThenKnowledgeWasConfirmedAfterSixDaysInterval(t *testing.T) {
	newCard := card.NewFakeCard().
		WithSection("languages/cpp").
		WithEntry("1").
		WithActivityChain(card.GenerateActivityChain(card.LearnCard)).
		Build()

	// New card shown for the first time
	pressGood(t, newCard)

	// Confirming knowledge after ten minutes
	pressGood(t, newCard)

	// Confirming knowledge next day
	pressGood(t, newCard)

	// Failed to confirm knowledge after six days
	pressAgain(t, newCard)

	// Confirming knowledge after ten minutes
	pressGood(t, newCard)

	// Confirming knowledge next day
	pressGood(t, newCard)

	// Confirming knowledge after six days
	repetitionParams := pressGood(t, newCard)

	// Since card was forgotten, easiness factor dropped to 1.7
	// Grade 4 makes easiness factor again stable, but now on new lower level
	// So to calculate interval we need to multiply previous interval to 1.7 and round to nearest integer
	// 6 * 1.7 = 10
	if repetitionParams.Interval != 10*study.Day {
		t.Fatal("Ten days interval should be assigned for card that was forgotten but can be remembered after six days interval after that")
	}
}

func pressAgain(t *testing.T, c *card.Card) *study.CardRepetitionParams {
	return pressButtonByGrade(t, c, Again)
}

func pressHard(t *testing.T, c *card.Card) *study.CardRepetitionParams {
	return pressButtonByGrade(t, c, Hard)
}

func pressGood(t *testing.T, c *card.Card) *study.CardRepetitionParams {
	return pressButtonByGrade(t, c, Good)
}

func pressEasy(t *testing.T, c *card.Card) *study.CardRepetitionParams {
	return pressButtonByGrade(t, c, Easy)
}

func pressButtonByGrade(t *testing.T, c *card.Card, grade int) *study.CardRepetitionParams {
	timeSource := datetime.NewFakeTimeSource()

	workflow := study.NewCardWorkflow(timeSource, c)

	nextStates, err := workflow.GetNextStates()

	if err != nil {
		t.Fatal(err)
	}

	state := getStateByGrade(nextStates, grade)

	err = workflow.TransitTo(state)

	if err != nil {
		t.Fatal(err)
	}

	repetitionParams, err := study.GetCardRepetitionParams(c.Activities())

	if err != nil {
		t.Fatal(err)
	}

	return repetitionParams
}

func getStateByGrade(states []*study.CardState, grade int) *study.CardState {
	for _, state := range states {
		if state.Grade == grade {
			return state
		}
	}
	return nil
}
