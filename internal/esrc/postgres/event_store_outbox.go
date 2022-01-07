package postgres

import (
	"context"
	"errors"
	"fmt"
	"ledger/internal/esrc"
	"ledger/internal/esrc/relay"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type EventStoreOutbox struct {
	pool     *pgxpool.Pool
	encodeID func(esrc.ID) string
}

var _ (relay.EventStoreOutbox) = (*EventStoreOutbox)(nil)

func NewEventStoreOutbox(pool *pgxpool.Pool) *EventStoreOutbox {
	return &EventStoreOutbox{
		pool: pool,
	}
}

func (o *EventStoreOutbox) UnpublishedEvents(ctx context.Context) ([]relay.RelayEvent, error) {
	const query = "SELECT event_source_id, version, name, data FROM events WHERE published = FALSE ORDER BY version ASC"
	rows, err := o.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []relay.RelayEvent
	for rows.Next() {
		var re relay.RelayEvent
		var id []byte
		err := rows.Scan(&id, &re.Sequence, &re.RawEvent.Name, &re.RawEvent.Data)
		if err != nil {
			return nil, err
		}
		re.AggregateID, err = uuid.FromBytes(id)
		if err != nil {
			return nil, err
		}
		events = append(events, re)
	}
	return events, nil
}

func (o *EventStoreOutbox) MarkEventsAsPublised(ctx context.Context, events []relay.RelayEvent) error {
	if len(events) == 0 {
		return errors.New("no events to mark as published")
	}

	b := &strings.Builder{}
	b.WriteString("UPDATE events AS e SET published = TRUE FROM (values")
	for i := range events {
		b.WriteString("('")
		b.WriteString(fmt.Sprintf("%v", events[i].AggregateID))
		// TODO: UUID should be configurable
		b.WriteString("'::UUID, ")
		b.WriteString(strconv.FormatUint(events[i].Sequence, 10))
		b.WriteString(")")
		if i+1 != len(events) {
			b.WriteByte(',')
		}
	}
	b.WriteString(") as p(esid, v) WHERE e.event_source_id = p.esid AND e.version = p.v")

	_, err := o.pool.Exec(ctx, b.String())
	return err
}
