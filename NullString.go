package common

import (
	"database/sql"
	"encoding/json"
)

type NullString struct {
	sql.NullString
}

func NewNullString(value string, valid bool) NullString {
	result := NullString{}

	result.String = value
	result.Valid = valid

	return result
}

func (n *NullString) Set(value string) {
	n.Valid = true
	n.String = value
}

func (n *NullString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Valid = false
		return nil
	}

	err := json.Unmarshal(data, &n.String)
	if err != nil {
		return err
	}
	n.Valid = true
	return nil
}

func (n NullString) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.String)
	} else {
		return []byte("null"), nil
	}
}
