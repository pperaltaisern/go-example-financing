package financing

import (
	"encoding/json"

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

type Money float64
