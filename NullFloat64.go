package common

import (
	"database/sql"
	"encoding/json"
)

type NullFloat64 struct {
	sql.NullFloat64
}

func (n *NullFloat64) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Valid = false
		return nil
	}

	err := json.Unmarshal(data, &n.Float64)
	if err != nil {
		return err
	}
	n.Valid = true
	return nil
}

func (n NullFloat64) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Float64)
	} else {
		return []byte("null"), nil
	}
}
