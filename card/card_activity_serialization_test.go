package card_test

import (
	"reflect"
	"testing"

	"github.com/gottenheim/ariadne/card"
)

func TestSerializeCardActivityChain(t *testing.T) {
	cardActivity := createTestActivityChain(learnCard|cardExecutedMonthAgo, remindCard|remindCardScheduledToYesterday|cardExecutedToday)

	chainBinary, err := card.SerializeCardActivityChain(cardActivity)

	if err != nil {
		t.Fatal(err)
	}

	restoredActivity, err := card.DeserializeCardActivityChain(chainBinary)

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(cardActivity, restoredActivity) {
		t.Fatal("Original and restored activity must be equal")
	}
}
