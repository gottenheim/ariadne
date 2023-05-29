package card_test

import (
	"testing"

	"github.com/gottenheim/ariadne/card"
)

func TestIsCardNew_IfLearnActivityIsNotExecuted(t *testing.T) {
	activityChain := card.GenerateActivityChain(card.LearnCard)

	isNew, err := card.IsNewCardActivities(activityChain)

	if err != nil {
		t.Fatal(err)
	}

	if !isNew {
		t.Fatal("Card should be new, because learn activity has not been executed yet")
	}
}

func TestIsCardNew_IfLearnActivityIsNotExecuted_AndRemindActivityInTheEndOfChain(t *testing.T) {
	activityChain := card.GenerateActivityChain(card.LearnCard, card.RemindCard)

	isNew, err := card.IsNewCardActivities(activityChain)

	if err != nil {
		t.Fatal(err)
	}

	if !isNew {
		t.Fatal("Card should be new, because learn activity has not been executed yet")
	}
}

func TestIsCardNew_IfLearnActivityHasAlreadyBeenExecuted(t *testing.T) {
	activityChain := card.GenerateActivityChain(card.LearnCard | card.CardExecutedToday)

	isNew, err := card.IsNewCardActivities(activityChain)

	if err != nil {
		t.Fatal(err)
	}

	if isNew {
		t.Fatal("Card should not be new, because learn activity has already been executed")
	}
}

func TestIsCardNew_IfLearnActivityHasAlreadyBeenExecuted_AndRemindActivityInTheEndOfChain(t *testing.T) {
	activityChain := card.GenerateActivityChain(card.LearnCard|card.CardExecutedToday, card.RemindCard)

	isNew, err := card.IsNewCardActivities(activityChain)

	if err != nil {
		t.Fatal(err)
	}

	if isNew {
		t.Fatal("Card should not be new, because learn activity has already been executed")
	}
}
