package financing

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewIssuer(t *testing.T) {
	id := NewID()
	inv := NewIssuer(id)
	require.Equal(t, id, inv.id)

	e := NewIssuerCreatedEvent(id)
	require.Equal(t, e, inv.Changes()[0])
	require.Equal(t, 0, inv.InitialVersion())
}
