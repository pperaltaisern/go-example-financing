package esrc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEventRaiserAggregateShouldCreateAggregateWithVersion0AndNoEvents(t *testing.T) {
	a := NewEventRaiserAggregate(nil)
	_ = assert.Equal(t, 0, a.InitialVersion()) &&
		assert.Len(t, a.Changes(), 0)
}

func testOnEventFunc(c *int) func(Event) {
	return func(Event) {
		*(c)++
	}
}

type testEvent struct{}

func (testEvent) EventName() string { return "TestEvent" }

func TestNewEventRaiserAggregateShouldAppendRaisedEventsAndExecuteTheOnEventFunctionEachTimeAnEventIsRaised(t *testing.T) {
	var c int
	a := NewEventRaiserAggregate(testOnEventFunc(&c))
	a.Raise(testEvent{})
	a.Raise(testEvent{})

	_ = assert.Equal(t, 0, a.InitialVersion()) &&
		assert.Len(t, a.Changes(), 2) &&
		assert.Equal(t, 2, c)
}

func TestNewEventRaiserAggregateFromEventsShouldCreateAnAggregateWithAllEventsReplayedAndIncrementHisVersion(t *testing.T) {
	var c int
	events := []Event{testEvent{}, testEvent{}}

	a := NewEventRaiserAggregateFromEvents(1, events, testOnEventFunc(&c))
	_ = assert.Equal(t, 3, a.InitialVersion()) &&
		assert.Len(t, a.Changes(), 0) &&
		assert.Equal(t, 2, c)
}
