package card

import (
	"time"

	"github.com/gottenheim/ariadne/libraries/config"
)

type cardActivitySerializer struct {
	activities []CardActivityModel
}

func (s *cardActivitySerializer) OnLearnCard(learn *LearnCardActivity) error {
	activityModel := CardActivityModel{
		ActivityType: learnActivityType,
	}
	activityModel.Executed = learn.executed
	if learn.executed {
		activityModel.ExecutionTime = learn.executionTime.Format(time.DateTime)
	}

	s.activities = append(s.activities, activityModel)

	return nil
}

func (s *cardActivitySerializer) OnRemindCard(remind *RemindCardActivity) error {
	activityModel := CardActivityModel{
		ActivityType: remindActivityType,
	}
	activityModel.Executed = remind.executed

	if remind.executed {
		activityModel.ExecutionTime = remind.executionTime.Format(time.DateTime)
	}
	activityModel.ScheduledTo = remind.scheduledTo.Format(time.DateTime)

	s.activities = append(s.activities, activityModel)

	if remind.previousActivity == nil {
		return nil
	}

	return remind.previousActivity.Accept(s)
}

func SerializeCardActivityChain(cardActivity CardActivity) ([]byte, error) {
	serializer := &cardActivitySerializer{}
	err := cardActivity.Accept(serializer)

	if err != nil {
		return nil, err
	}

	activities := &CardActivitiesModel{
		Activities: serializer.activities,
	}

	return config.SerializeToYaml(activities)
}
