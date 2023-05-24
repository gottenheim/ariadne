package card

type NewCardStatus struct {
}

func (s *NewCardStatus) Name() CardStatusName {
	return CardStatusNew
}
