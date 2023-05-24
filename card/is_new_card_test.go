package card_test

import (
	"testing"

	"github.com/gottenheim/ariadne/card"
)

func TestIsCardNew_IfLearnActivityIsNotExecuted(t *testing.T) {
	activityChain := createTestActivityChain(learnCard)

	isNew, err := card.IsNewCard(activityChain)

	if err != nil {
		t.Fatal(err)
	}

	if !isNew {
		t.Fatal("Card should be new, because learn activity has not been executed yet")
	}
}

func TestIsCardNew_IfLearnActivityIsNotExecuted_AndRemindActivityInTheEndOfChain(t *testing.T) {
	activityChain := createTestActivityChain(learnCard, remindCard)

	isNew, err := card.IsNewCard(activityChain)

	if err != nil {
		t.Fatal(err)
	}

	if !isNew {
		t.Fatal("Card should be new, because learn activity has not been executed yet")
	}
}

func TestIsCardNew_IfLearnActivityHasAlreadyBeenExecuted(t *testing.T) {
	activityChain := createTestActivityChain(learnCard | cardExecutedToday)

	isNew, err := card.IsNewCard(activityChain)

	if err != nil {
		t.Fatal(err)
	}

	if isNew {
		t.Fatal("Card should not be new, because learn activity has already been executed")
	}
}

func TestIsCardNew_IfLearnActivityHasAlreadyBeenExecuted_AndRemindActivityInTheEndOfChain(t *testing.T) {
	activityChain := createTestActivityChain(learnCard|cardExecutedToday, remindCard)

	isNew, err := card.IsNewCard(activityChain)

	if err != nil {
		t.Fatal(err)
	}

	if isNew {
		t.Fatal("Card should not be new, because learn activity has already been executed")
	}
}
