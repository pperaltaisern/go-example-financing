package esrc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAggregate_Empty(t *testing.T) {
	a := NewAggregate(nil)
	_ = assert.Equal(t, 0, a.Version()) &&
		assert.Len(t, a.Events(), 0)
}

func testOnEventFunc(c *int) func(Event) {
	return func(Event) {
		*(c)++
	}
}

type testEvent struct{}

func (testEvent) EventName() string { return "TestEvent" }

func TestAggregate_Raise(t *testing.T) {
	var c int
	a := NewAggregate(testOnEventFunc(&c))
	a.Raise(testEvent{})
	a.Raise(testEvent{})

	_ = assert.Equal(t, 0, a.Version()) &&
		assert.Len(t, a.Events(), 2) &&
		assert.Equal(t, 2, c)
}

func TestAggregate_NewFromEvents(t *testing.T) {
	var c int

	events := []Event{testEvent{}, testEvent{}}

	a := NewAggregateFromEvents(events, testOnEventFunc(&c))
	_ = assert.Equal(t, 2, a.Version()) &&
		assert.Len(t, a.Events(), 0) &&
		assert.Equal(t, 2, c)
}