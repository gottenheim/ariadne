package card

import (
	"bytes"
	"time"

	"github.com/gottenheim/ariadne/libraries/config"
)

func DeserializeCardActivityChain(activitiesBinary []byte) (CardActivity, error) {
	cfg, err := config.FromYamlReader(bytes.NewReader(activitiesBinary))

	if err != nil {
		return nil, err
	}

	cardActivitiesModel := &CardActivitiesModel{}
	err = cfg.Materialize(cardActivitiesModel)

	if err != nil {
		return nil, err
	}

	var activity CardActivity

	for i := len(cardActivitiesModel.Activities) - 1; i >= 0; i-- {
		activityModel := cardActivitiesModel.Activities[i]

		if activityModel.ActivityType == learnActivityType {
			learn := CreateLearnCardActivity()
			learn.executed = activityModel.Executed
			if learn.executed {
				executionTime, err := time.ParseInLocation(time.DateTime, activityModel.ExecutionTime, time.Local)
				if err != nil {
					return nil, err
				}
				learn.executionTime = executionTime.Local()
			}
			activity = learn
		} else if activityModel.ActivityType == remindActivityType {
			remind := CreateRemindCardActivity(activity)
			remind.executed = activityModel.Executed
			if remind.executed {
				executionTime, err := time.ParseInLocation(time.DateTime, activityModel.ExecutionTime, time.Local)
				if err != nil {
					return nil, err
				}
				remind.executionTime = executionTime.Local()
			}
			if len(activityModel.ScheduledTo) > 0 {
				scheduledTo, err := time.ParseInLocation(time.DateTime, activityModel.ScheduledTo, time.Local)
				if err != nil {
					return nil, err
				}
				remind.scheduledTo = scheduledTo.Local()
			}
			remind.easinessFactor = activityModel.EasinessFactor
			remind.repetitionNumber = activityModel.RepetitionNumber
			remind.interval = activityModel.Interval
			activity = remind
		}
	}

	return activity, nil
}
