package common

import (
	"database/sql"
	"encoding/json"
)

type NullBool struct {
	sql.NullBool
}

func (n *NullBool) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Valid = false
		return nil
	}

	err := json.Unmarshal(data, &n.Bool)
	if err != nil {
		return err
	}
	n.Valid = true
	return nil
}

func (n NullBool) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Bool)
	} else {
		return []byte("null"), nil
	}
}
