package common

import (
	"encoding/json"
	"fmt"
	"time"
)

type NullTime struct {
	Time  time.Time
	Valid bool
}

func (n *NullTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Valid = false
		return nil
	}

	err := json.Unmarshal(data, &n.Time)
	if err != nil {
		return err
	}
	n.Valid = true
	return nil
}

func (n NullTime) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Time)
	} else {
		return []byte("null"), nil
	}
}

func (n *NullTime) Scan(src interface{}) error {
	if src == nil {
		n.Valid = false
		return nil
	}

	if time, ok := src.(time.Time); ok {
		n.Time = time
		n.Valid = true
		return nil
	}

	return fmt.Errorf("cannot Scan src for NullTime type")
}
