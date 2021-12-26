package es

type Event interface {
	IsEvent()
	Name() string
}
