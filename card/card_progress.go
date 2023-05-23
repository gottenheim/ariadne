package card

type CardStatus string

const (
	New       CardStatus = "New"
	Scheduled            = "Scheduled"
)

type CardProgress struct {
	Status CardStatus
}

func ScheduledCard() *CardProgress {
	return &CardProgress{
		Status: Scheduled,
	}
}
