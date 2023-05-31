package card

const learnActivityType = "learn"
const remindActivityType = "remind"

type CardActivityModel struct {
	ActivityType     string
	Executed         bool
	ExecutionTime    string
	ScheduledTo      string
	EasinessFactor   float32
	RepetitionNumber int
	Interval         int
}

type CardActivitiesModel struct {
	Activities []CardActivityModel
}
