package types

import (
	"errors"
	"time"

	"database/sql/driver"
)

// GORM does not support DATE/TIME columns.
// We add a special type.
// It looks like MySql DATE type
type Date struct {
	Time *time.Time
}

// Scan implements the sql.Scanner interface for database deserialization.
func (d *Date) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	v, ok := value.(time.Time)
	if !ok {
		return errors.New("cannot Scan Date")
	}

	*d = Date{}
	d.Time = &v
	return nil
}

func (d *Date) String() string {
	if d == nil || d.Time == nil {
		return ""
	}
	return d.Time.Format("2006-01-02")
}

// Value implements the driver.Valuer interface for database serialization.
func (d *Date) Value() (driver.Value, error) {
	if d == nil {
		return nil, nil
	}

	val := d.String()
	if len(val) == 0 {
		return nil, nil
	}

	return val, nil
}

func (d *Date) UnmarshalJSON(data []byte) error {
	stringVal := string(data)
	if stringVal == `""` {
		ref := Date{Time: nil}
		*d = ref
		return nil
	}
	datetime, err := time.Parse(`"2006-01-02"`, stringVal)
	if err != nil {
		return err
	}
	ref := Date{Time: &datetime}
	*d = ref
	return nil
}

func (d *Date) MarshalJSON() ([]byte, error) {
	if d == nil {
		return nil, nil
	}

	return []byte(`"` + d.String() + `"`), nil
}
