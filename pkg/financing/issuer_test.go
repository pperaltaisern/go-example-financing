package financing

import (
	"ledger/internal/esrc"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewIssuer(t *testing.T) {
	id := NewID()
	inv := NewIssuer(id)
	require.Equal(t, id, inv.id)

	e := NewIssuerCreatedEvent(id)
	require.Equal(t, e, inv.aggregate.Events()[0])
	require.Equal(t, 0, inv.aggregate.Version())
}

func TestNewIssuerFromEvents(t *testing.T) {
	id := NewID()
	e := NewIssuerCreatedEvent(id)

	inv := newIssuerFromEvents([]esrc.Event{e})
	require.Equal(t, id, inv.id)

	require.Empty(t, inv.aggregate.Events())
	require.Equal(t, 1, inv.aggregate.Version())
}
