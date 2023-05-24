package card

type ScheduledCardStatus struct {
}

func (s *ScheduledCardStatus) Name() CardStatusName {
	return CardStatusScheduled
}
