package esrcwatermill

import (
	"encoding/json"
	"testing"

	"github.com/pperaltaisern/financing/internal/esrc"
	"github.com/pperaltaisern/financing/internal/esrc/relay"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestEvent struct {
	ID    string `json:"-"`
	Field int
}

func (e *TestEvent) WithAggregateID(id string) {
	e.ID = id
}

var eventToMarshal = &TestEvent{
	ID:    watermill.NewULID(),
	Field: 10,
}

func TestRelayEventMarshaler(t *testing.T) {
	marshaler := RelayEventMarshaler{
		CmdMarshaler: cqrs.JSONMarshaler{},
	}

	re, err := relayEventFromTestEvent(eventToMarshal)
	require.NoError(t, err)

	msg, err := marshaler.Marshal(re)
	require.NoError(t, err)

	eventToUnmarshal := TestEvent{}
	err = marshaler.Unmarshal(msg, &eventToUnmarshal)
	require.NoError(t, err)

	assert.EqualValues(t, eventToUnmarshal, eventToUnmarshal)
}

func TestRelayEventMarshaler_Marshal_generate_name(t *testing.T) {
	marshaler := RelayEventMarshaler{
		CmdMarshaler: cqrs.JSONMarshaler{},
	}

	re, err := relayEventFromTestEvent(eventToMarshal)
	require.NoError(t, err)

	name := marshaler.Name(re)
	require.NoError(t, err)

	assert.Equal(t, "TestEvent", name)
}

func relayEventFromTestEvent(eventToMarshal *TestEvent) (relay.RelayEvent, error) {
	b, err := json.Marshal(eventToMarshal)
	if err != nil {
		return relay.RelayEvent{}, err
	}

	return relay.NewRelayEvent(eventToMarshal.ID, 0, esrc.RawEvent{Name: "TestEvent", Data: b}), nil
}
