package card

type CardStatusName string

const (
	CardStatusNew       = "New"
	CardStatusScheduled = "Scheduled"
)

type CardStatus interface {
	Name() CardStatusName
}
