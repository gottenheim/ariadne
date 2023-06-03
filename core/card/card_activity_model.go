package card

const learnActivityType = "learn"
const remindActivityType = "remind"

type CardActivityModel struct {
	ActivityType     string
	Executed         bool
	ExecutionTime    string  `yaml:",omitempty"`
	ScheduledTo      string  `yaml:",omitempty"`
	EasinessFactor   float64 `yaml:",omitempty"`
	RepetitionNumber int     `yaml:",omitempty"`
	Interval         string  `yaml:",omitempty"`
}

type CardActivitiesModel struct {
	Activities []CardActivityModel
}
