package card

type CardActivity interface {
	Accept(visitor CardActivityVisitor) error
}

type CardActivityVisitor interface {
	OnLearnCard(learn *LearnCardActivity) error
	OnRemindCard(reming *RemindCardActivity) error
}
