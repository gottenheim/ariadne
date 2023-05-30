package card

const learnActivityType = "learn"
const remindActivityType = "remind"

type CardActivityModel struct {
	ActivityType  string
	Executed      bool
	ExecutionTime string
	ScheduledTo   string
}

type CardActivitiesModel struct {
	Activities []CardActivityModel
}
