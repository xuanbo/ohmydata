package entity

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// Entity 实体
type Entity struct {
	ID        string `json:"id" gorm:"primaryKey;type:string;size:30"`
	CreatedAt *Time  `json:"createdAt" gorm:"<-:create"`
	UpdatedAt *Time  `json:"updatedAt" gorm:"<-:create;<-:update"`
	CreatedBy string `json:"createdBy" gorm:"<-:create"`
	UpdatedBy string `json:"updatedBy" gorm:"<-:create;<-:update"`
}

const defaultFmt = "2006-01-02 15:04:05.000"

// Time json time
type Time struct {
	time.Time
}

// MarshalJSON makes Time implements json.Marshaler interface
func (t Time) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf(`"%s"`, t.Format(defaultFmt))
	return []byte(formatted), nil
}

// UnmarshalJSON makes Time implements json.Unmarshaler interface
func (t *Time) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" || string(data) == "NULL" {
		return nil
	}
	// Fractional seconds are handled implicitly by Parse.
	value, err := time.Parse(`"`+defaultFmt+`"`, string(data))
	if err != nil {
		return err
	}
	*t = Time{value}
	return nil
}

// Value makes Time implements drive.Valuer interface
func (t Time) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan makes Time implements sql.Scanner interface
func (t *Time) Scan(v interface{}) error {
	if value, ok := v.(time.Time); ok {
		*t = Time{value}
		return nil
	}
	return fmt.Errorf("can not convert %v to time.Time", v)
}

// Now returns the current local time
func Now() Time {
	return Time{Time: time.Now()}
}
