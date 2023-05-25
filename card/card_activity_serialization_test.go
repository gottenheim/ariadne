package card_test

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	"github.com/gottenheim/ariadne/card"
	"github.com/gottenheim/ariadne/details/config"
	"github.com/gottenheim/ariadne/details/datetime"
)

func TestSerializeCardActivityChain(t *testing.T) {
	cardActivity := card.GenerateActivityChain(card.LearnCard|card.CardExecutedMonthAgo, card.RemindCard|card.RemindCardScheduledToYesterday|card.CardExecutedToday)

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

	today := datetime.GetLocalTestTime()

	if remind.ExecutionTime != today.Format(time.DateTime) {
		t.Error("Remind activity execution time must be today")
	}

	yesterday := today.AddDate(0, 0, -1)

	if remind.ScheduledTo != yesterday.Format(time.DateTime) {
		t.Error("Remind activity schedule time must be yesterday")
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
		ActivityType:  "remind",
		Executed:      remind.IsExecuted(),
		ExecutionTime: remind.ExecutionTime().Format(time.DateTime),
		ScheduledTo:   remind.ScheduledTo().Format(time.DateTime),
	}

	v.activities = append(v.activities, remindModel)

	return remind.PreviousActivity().Accept(v)
}

func TestDeserializeCardActivityChain(t *testing.T) {
	today := datetime.GetLocalTestTime()
	monthAgo := today.AddDate(0, -1, 0)
	yesterday := today.AddDate(0, 0, -1)

	initialCardActivitiesModel := &card.CardActivitiesModel{
		Activities: []card.CardActivityModel{
			{
				ActivityType:  "remind",
				Executed:      true,
				ExecutionTime: today.Format(time.DateTime),
				ScheduledTo:   yesterday.Format(time.DateTime),
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
