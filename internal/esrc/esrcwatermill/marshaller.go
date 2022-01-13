package esrcwatermill

import (
	"fmt"

	"github.com/pperaltaisern/financing/internal/esrc/relay"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
)

type RelayEventMarshaler struct {
	NewUUID      func() string
	CmdMarshaler cqrs.CommandEventMarshaler
}

var _ cqrs.CommandEventMarshaler = (*RelayEventMarshaler)(nil)

func (m RelayEventMarshaler) Marshal(v interface{}) (*message.Message, error) {
	e, ok := v.(relay.RelayEvent)
	if !ok {
		return m.CmdMarshaler.Marshal(v)
	}

	msg := message.NewMessage(
		m.newUUID(),
		e.RawEvent.Data,
	)
	msg.Metadata.Set("name", e.RawEvent.Name)
	msg.Metadata.Set("aggregateID", fmt.Sprintf("%v", e.AggregateID))

	return msg, nil
}

func (m RelayEventMarshaler) newUUID() string {
	if m.NewUUID != nil {
		return m.NewUUID()
	}

	// default
	return watermill.NewUUID()
}

func (m RelayEventMarshaler) Unmarshal(msg *message.Message, v interface{}) error {
	err := m.CmdMarshaler.Unmarshal(msg, v)
	if err != nil {
		return err
	}
	if e, ok := v.(Event); ok {
		e.WithAggregateID(msg.Metadata.Get("aggregateID"))
	}
	return nil
}

func (m RelayEventMarshaler) Name(cmdOrEvent interface{}) string {
	if re, ok := cmdOrEvent.(relay.RelayEvent); ok {
		return re.RawEvent.Name
	}
	return m.CmdMarshaler.Name(cmdOrEvent)
}

func (m RelayEventMarshaler) NameFromMessage(msg *message.Message) string {
	return msg.Metadata.Get("name")
}
