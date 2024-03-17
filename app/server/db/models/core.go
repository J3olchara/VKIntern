package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"time"
)

type Search struct {
	Search    string `json:"search"`
	ActorName string `json:"actor"`
	Ordering  string `json:"ordering"`
	Field     string `json:"field"`
}

type Date time.Time

func (date *Date) Scan(value interface{}) error {
	nullTime := &sql.NullTime{}
	err := nullTime.Scan(value)
	*date = Date(nullTime.Time)
	return err
}

func (date Date) Value() (driver.Value, error) {
	y, m, d := time.Time(date).Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.Time(date).Location()), nil
}

func (date Date) GormDataType() string {
	return "date"
}

func (date Date) GobEncode() ([]byte, error) {
	return time.Time(date).GobEncode()
}

func (date *Date) GobDecode(b []byte) error {
	return (*time.Time)(date).GobDecode(b)
}

func (date Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(date).Format("2006-01-02"))
}

func (date *Date) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	dt, err := time.Parse("2006-01-02", str)
	if err != nil {
		return err
	}
	*date = Date(dt)
	return nil
}
