package card_test

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	"github.com/gottenheim/ariadne/core/card"
	"github.com/gottenheim/ariadne/libraries/config"
	"github.com/gottenheim/ariadne/libraries/datetime"
)

func TestSerializeCardActivityChain(t *testing.T) {
	cardActivity := card.GenerateActivityChain(card.LearnCard|card.CardExecutedMonthAgo, card.RemindCard|card.RemindCardScheduledToYesterday|card.CardExecutedToday)

	remindCard := cardActivity.(*card.RemindCardActivity)
	remindCard.SetEasinessFactor(2.5)
	remindCard.SetRepetitionNumber(4)
	remindCard.SetInterval(time.Hour * 2)

	chainBinary, err := card.SerializeCardActivityChain(cardActivity)

	if err != nil {
		t.Fatal(err)
	}

	cfg, err := config.FromYamlReader(bytes.NewBuffer(chainBinary))

	if err != nil {
		t.Fatal(err)
	}

	cardActivitiesModel := &card.CardActivitiesModel{}
	err = cfg.Materialize(cardActivitiesModel)

	if err != nil {
		t.Fatal(err)
	}

	if len(cardActivitiesModel.Activities) != 2 {
		t.Fatal("There must be two activities")
	}

	remind := cardActivitiesModel.Activities[0]

	if remind.ActivityType != "remind" {
		t.Error("Remind activity expected")
	}

	if !remind.Executed {
		t.Error("Remind activity is expected to be executed")
	}

	today := datetime.FakeNow()

	if remind.ExecutionTime != today.Format(time.DateTime) {
		t.Error("Remind activity execution time must be today")
	}

	yesterday := today.AddDate(0, 0, -1)

	if remind.ScheduledTo != yesterday.Format(time.DateTime) {
		t.Error("Remind activity schedule time must be yesterday")
	}

	if remind.EasinessFactor != 2.5 {
		t.Error("Remind activity should have easiness factor 2.5")
	}

	if remind.RepetitionNumber != 4 {
		t.Error("Remind activity should have repetition number 4")
	}

	if remind.Interval != "2h0m0s" {
		t.Error("Remind activity should have interval 12")
	}

	learn := cardActivitiesModel.Activities[1]

	if learn.ActivityType != "learn" {
		t.Error("Learn activity expected")
	}

	if !learn.Executed {
		t.Error("Learn activity is expected to be executed")
	}

	monthAgo := today.AddDate(0, -1, 0)

	if learn.ExecutionTime != monthAgo.Format(time.DateTime) {
		t.Error("Learn activity execution time must be month ago")
	}
}

type testCardActivityVisitor struct {
	activities []card.CardActivityModel
}

func (v *testCardActivityVisitor) OnLearnCard(learn *card.LearnCardActivity) error {
	learnModel := card.CardActivityModel{
		ActivityType:  "learn",
		Executed:      learn.IsExecuted(),
		ExecutionTime: learn.ExecutionTime().Format(time.DateTime),
	}

	v.activities = append(v.activities, learnModel)

	return nil
}

func (v *testCardActivityVisitor) OnRemindCard(remind *card.RemindCardActivity) error {
	remindModel := card.CardActivityModel{
		ActivityType:     "remind",
		Executed:         remind.IsExecuted(),
		ExecutionTime:    remind.ExecutionTime().Format(time.DateTime),
		ScheduledTo:      remind.ScheduledTo().Format(time.DateTime),
		EasinessFactor:   2.5,
		RepetitionNumber: 4,
		Interval:         "2h0m0s",
	}

	v.activities = append(v.activities, remindModel)

	return remind.PreviousActivity().Accept(v)
}

func TestDeserializeCardActivityChain(t *testing.T) {
	today := datetime.FakeNow()
	monthAgo := today.AddDate(0, -1, 0)
	yesterday := today.AddDate(0, 0, -1)

	initialCardActivitiesModel := &card.CardActivitiesModel{
		Activities: []card.CardActivityModel{
			{
				ActivityType:     "remind",
				Executed:         true,
				ExecutionTime:    today.Format(time.DateTime),
				ScheduledTo:      yesterday.Format(time.DateTime),
				EasinessFactor:   2.5,
				RepetitionNumber: 4,
				Interval:         "2h0m0s",
			},
			{
				ActivityType:  "learn",
				Executed:      true,
				ExecutionTime: monthAgo.Format(time.DateTime),
			},
		},
	}

	activitiesBinary, err := config.SerializeToYaml(initialCardActivitiesModel)

	if err != nil {
		t.Fatal(err)
	}

	activity, err := card.DeserializeCardActivityChain(activitiesBinary)

	if err != nil {
		t.Fatal(err)
	}

	visitor := &testCardActivityVisitor{}
	activity.Accept(visitor)

	actualCardActivitiesModel := &card.CardActivitiesModel{
		Activities: visitor.activities,
	}

	if !reflect.DeepEqual(actualCardActivitiesModel, initialCardActivitiesModel) {
		t.Fatal("Initial activities don't match with deserialized ones")
	}
}
