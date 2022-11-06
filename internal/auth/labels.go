package auth

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/pkg/errors"
	//"github.com/sirupsen/logrus"
)

type Labels map[string]string

func (labels *Labels) Scan(src interface{}) error {
	//logrus.Tracef("src type: %T", src)
	switch str := src.(type) {
	case []byte:
		if err := json.Unmarshal(str, labels); err != nil {
			return err
		}
	case string:
		if err := json.Unmarshal([]byte(str), labels); err != nil {
			return err
		}
	default:
		return errors.Errorf("must []byte was '%T'", src)
	}

	return nil
}
func (labels Labels) Value() (driver.Value, error) {
	return json.Marshal(labels)
}

type List []string

func (list *List) Scan(src interface{}) error {
	//logrus.Tracef("src type: %T", src)
	switch str := src.(type) {
	case []byte:
		if err := json.Unmarshal(str, list); err != nil {
			return err
		}
	case string:
		if err := json.Unmarshal([]byte(str), list); err != nil {
			return err
		}
	default:
		return errors.Errorf("must []byte was '%T'", src)
	}

	return nil

}
func (list List) Value() (driver.Value, error) {
	return json.Marshal(list)
}
