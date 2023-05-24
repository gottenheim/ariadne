package card

type CardStatus string

const (
	New       CardStatus = "New"
	Scheduled            = "Scheduled"
)

type CardProgress struct {
	status CardStatus
}

func ScheduledCard() *CardProgress {
	return &CardProgress{
		status: Scheduled,
	}
}

func (p *CardProgress) IsNew() bool {
	return p.status == New
}

func (p *CardProgress) IsScheduled() bool {
	return p.status == Scheduled
}
