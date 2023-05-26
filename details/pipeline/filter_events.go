package pipeline

type FilterEvents interface {
	OnStart()
	OnFinish()
}
