package esrctesting

import (
	"fmt"
	"sync/atomic"

	"github.com/google/uuid"
	"github.com/pperaltaisern/financing/internal/esrc"
)

type IDGenerator func() esrc.ID

func UUIDGenerator() IDGenerator {
	return func() esrc.ID {
		return uuid.New()
	}
}

func IntIDGenerator() IDGenerator {
	var c uint64
	return func() esrc.ID {
		atomic.AddUint64(&c, 1)
		return c
	}
}

func StringIDGenerator() IDGenerator {
	var c uint64
	return func() esrc.ID {
		atomic.AddUint64(&c, 1)
		return fmt.Sprintf("ID-%v", c)
	}
}
