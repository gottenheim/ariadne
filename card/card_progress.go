package card

type CardProgress struct {
	previousStatus CardStatus
	status         CardStatus
}

func GetNewCardProgress() *CardProgress {
	return &CardProgress{
		status: &NewCardStatus{},
	}
}

func GetScheduledCardProgress() *CardProgress {
	return &CardProgress{
		status: &ScheduledCardStatus{},
	}
}

func (p *CardProgress) IsNew() bool {
	return p.status.Name() == CardStatusNew
}

func (p *CardProgress) IsScheduled() bool {
	return p.status.Name() == CardStatusScheduled
}
