package esrcpg

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/pperaltaisern/financing/internal/esrc"
	"github.com/pperaltaisern/financing/internal/esrc/relay"

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

// type EventStoreOutboxOption func(*EventStoreOutbox)

// // RelayWitInterval sets the time duration that the Relayer will wait within loops
// func RelayerWithInterval(interval time.Duration) RelayerOption {
// 	return func(c *Relayer) {
// 		c.interval = interval
// 	}
// }

func (o *EventStoreOutbox) UnpublishedEvents(ctx context.Context) ([]relay.RelayEvent, error) {
	const query = "SELECT aggregate_id, version, name, data FROM events WHERE published = FALSE ORDER BY version ASC"
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
	b.WriteString(") as p(esid, v) WHERE e.aggregate_id = p.esid AND e.version = p.v")

	_, err := o.pool.Exec(ctx, b.String())
	return err
}
