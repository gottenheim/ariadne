package card

type CardProgressWriteModel struct {
	Status CardStatus
}

func (p *CardProgress) ToWriteModel() *CardProgressWriteModel {
	return &CardProgressWriteModel{
		Status: p.status,
	}
}

func (p *CardProgressWriteModel) ToCardProgress() *CardProgress {
	return &CardProgress{
		status: p.Status,
	}
}
