package financing

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/google/uuid"
)

type ID uuid.UUID

func NewID() ID {
	return ID(uuid.New())
}

func NewIDFromString(text string) ID {
	return ID(uuid.MustParse(text))
}

func TryNewIDFromString(text string) (ID, error) {
	id, err := uuid.Parse(text)
	return ID(id), err
}

func (id ID) String() string {
	return uuid.UUID(id).String()
}

func (id ID) MarshalJSON() ([]byte, error) {
	v := uuid.UUID(id)
	return json.Marshal(v)
}

func (id *ID) UnmarshalJSON(data []byte) error {
	v := uuid.UUID(*id)
	err := v.UnmarshalText(data)
	if err != nil {
		return err
	}
	*id = ID(v)
	return nil
}

// Scan is for sql/driver.Scaner
func (id *ID) Scan(value interface{}) error {
	baseUuid := uuid.UUID(*id)
	err := baseUuid.Scan(value)
	*id = ID(baseUuid)
	return err
}

// Value is for sql/driver.Valuer
func (id ID) Value() (driver.Value, error) {
	baseUuid := uuid.UUID(id)
	return baseUuid.Value()
}

type Money float64

// Scan is for sql/driver.Scaner
func (m *Money) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case string:
		fv, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return err
		}
		*m = Money(fv)
		return nil
	case float64:
		*m = Money(v)
		return nil
	}
	return fmt.Errorf("failed to scan Money, incoming type is: %t", value)
}

// Value is for sql/driver.Valuer
func (m Money) Value() (driver.Value, error) {
	return float64(m), nil
}
