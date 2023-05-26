package pipeline

type FilterEvents interface {
	OnStart()
	OnFinish()
	OnError(error)
}
