package esrc

type Event interface {
	IsEvent()
	Name() string
}
