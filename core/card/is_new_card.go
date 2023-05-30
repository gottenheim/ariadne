package card

type newCardVisitor struct {
	result bool
}

func (s *newCardVisitor) OnLearnCard(learn *LearnCardActivity) error {
	s.result = !learn.IsExecuted()
	return nil
}

func (s *newCardVisitor) OnRemindCard(remind *RemindCardActivity) error {
	if remind.PreviousActivity() == nil {
		return nil
	}

	return remind.PreviousActivity().Accept(s)
}

func IsNewCardActivities(activities CardActivity) (bool, error) {
	isNew := &newCardVisitor{
		result: true,
	}
	err := activities.Accept(isNew)
	if err != nil {
		return false, err
	}
	return isNew.result, nil
}
