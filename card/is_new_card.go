package card

type isNewCard struct {
	result bool
}

func (s *isNewCard) OnLearnCard(learn *LearnCardActivity) error {
	s.result = !learn.executed
	return nil
}

func (s *isNewCard) OnRemindCard(remind *RemindCardActivity) error {
	if remind.previousActivity == nil {
		return nil
	}

	return remind.previousActivity.Accept(s)
}

func IsNewCard(activity CardActivity) (bool, error) {
	isNew := &isNewCard{
		result: true,
	}
	err := activity.Accept(isNew)
	if err != nil {
		return false, err
	}
	return isNew.result, nil
}
