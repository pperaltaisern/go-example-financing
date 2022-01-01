package financing

import (
	"encoding/json"
	"reflect"
	"testing"
)

var testID = NewID()

func TestJSON(t *testing.T) {
	type S struct {
		ID1 ID
		ID2 ID
	}
	s1 := S{ID1: testID}
	data, err := json.Marshal(&s1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("marshalled data: %s", data)

	var s2 S
	if err := json.Unmarshal(data, &s2); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(&s1, &s2) {
		t.Errorf("got %#v, want %#v", s2, s1)
	}
}
