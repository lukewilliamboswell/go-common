package common

import (
	"database/sql"
	"encoding/json"
)

type NullInt64 struct {
	sql.NullInt64
}

func (n *NullInt64) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Valid = false
		return nil
	}

	err := json.Unmarshal(data, &n.Int64)
	if err != nil {
		return err
	}

	n.Valid = true
	return nil
}

func (n NullInt64) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Int64)
	} else {
		return []byte("null"), nil
	}
}
