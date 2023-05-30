package study

import "github.com/gottenheim/ariadne/core/card"

type newCardVisitor struct {
	result bool
}

func (s *newCardVisitor) OnLearnCard(learn *card.LearnCardActivity) error {
	s.result = !learn.IsExecuted()
	return nil
}

func (s *newCardVisitor) OnRemindCard(remind *card.RemindCardActivity) error {
	if remind.PreviousActivity() == nil {
		return nil
	}

	return remind.PreviousActivity().Accept(s)
}

func IsNewCardActivities(activities card.CardActivity) (bool, error) {
	isNew := &newCardVisitor{
		result: true,
	}
	err := activities.Accept(isNew)
	if err != nil {
		return false, err
	}
	return isNew.result, nil
}
