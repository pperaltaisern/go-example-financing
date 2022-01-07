package esrcwatermill

type Event interface {
	WithAggregateID(id string)
}
