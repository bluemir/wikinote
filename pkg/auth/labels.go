package auth

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Labels map[string]string

func (labels *Labels) Scan(src interface{}) error {
	str, ok := src.([]byte)
	if !ok {
		return errors.New("must []byte")
	}
	err := json.Unmarshal(str, labels)
	if err != nil {
		return err
	}
	return nil
}
func (labels Labels) Value() (driver.Value, error) {
	return json.Marshal(labels)
}
