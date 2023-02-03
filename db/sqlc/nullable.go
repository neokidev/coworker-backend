package db

import (
	"database/sql"
	"encoding/json"
)

type NullString struct {
	sql.NullString
}

func (ns *NullString) UnmarshalJSON(value []byte) error {
	err := json.Unmarshal(value, &ns.String)
	ns.Valid = err == nil
	return err
}
