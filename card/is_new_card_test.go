package card_test

import (
	"testing"
	"time"

	"github.com/gottenheim/ariadne/card"
	"github.com/gottenheim/ariadne/test"
)

func TestIsCardNew_IfLearnActivityIsNotExecuted(t *testing.T) {
	learnCard := card.CreateLearnCardActivity()

	isNew, err := card.IsNewCard(learnCard)

	if err != nil {
		t.Fatal(err)
	}

	if !isNew {
		t.Fatal("Card should be new, because learn activity has not been executed yet")
	}
}

func TestIsCardNew_IfLearnActivityIsNotExecuted_AndRemindActivityInTheEndOfChain(t *testing.T) {
	learnCard := card.CreateLearnCardActivity()

	remindCard := card.CreateRemindCardActivity(test.GetLocalTestTime(), learnCard)

	isNew, err := card.IsNewCard(remindCard)

	if err != nil {
		t.Fatal(err)
	}

	if !isNew {
		t.Fatal("Card should be new, because learn activity has not been executed yet")
	}
}

func TestIsCardNew_IfLearnActivityHasAlreadyBeenExecuted(t *testing.T) {
	learnCard := card.CreateLearnCardActivity()
	learnCard.MarkAsExecuted(time.Now())

	isNew, err := card.IsNewCard(learnCard)

	if err != nil {
		t.Fatal(err)
	}

	if isNew {
		t.Fatal("Card should not be new, because learn activity has already been executed")
	}
}

func TestIsCardNew_IfLearnActivityHasAlreadyBeenExecuted_AndRemindActivityInTheEndOfChain(t *testing.T) {
	learnCard := card.CreateLearnCardActivity()
	learnCard.MarkAsExecuted(time.Now())

	remindCard := card.CreateRemindCardActivity(test.GetLocalTestTime(), learnCard)

	isNew, err := card.IsNewCard(remindCard)

	if err != nil {
		t.Fatal(err)
	}

	if isNew {
		t.Fatal("Card should not be new, because learn activity has already been executed")
	}
}
