package dbutil

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/seerx/gpa/engine/sql/types"
)

// FormatTime format time as column type
func FormatTime(fmtTimeStampz string, sqlTypeName string, t time.Time) (v interface{}) {
	switch sqlTypeName {
	case types.Time:
		s := t.Format("2006-01-02 15:04:05") // time.RFC3339
		v = s[11:19]
	case types.Date:
		v = t.Format("2006-01-02")
	case types.DateTime, types.TimeStamp, types.Varchar: // !DarthPestilane! format time when sqlTypeName is schemas.Varchar.
		v = t.Format("2006-01-02 15:04:05")
	case types.TimeStampz:
		if fmtTimeStampz != "" {
			// dialect.URI().DBType == types.MSSQL ? "2006-01-02T15:04:05.9999999Z07:00"
			v = t.Format(fmtTimeStampz)
		} else {
			// if dialect.URI().DBType == types.MSSQL {
			// 	v = t.Format("2006-01-02T15:04:05.9999999Z07:00")
			// } else {
			v = t.Format(time.RFC3339Nano)
		}
	case types.BigInt, types.Int:
		v = t.Unix()
	default:
		v = t
	}
	return
}

type TimeProp struct {
	TypeName string
	Nullable bool
	TimeZone *time.Location
}

type TimePropDesc struct {
	TypeName string
	Nullable bool
	TimeZone string
}

func NewTimeProp(typeName string, nullable bool, timeZoneName string) (*TimeProp, error) {
	tz, err := time.LoadLocation(timeZoneName)
	if err != nil {
		return nil, err
	}
	return &TimeProp{
		TypeName: typeName,
		Nullable: nullable,
		TimeZone: tz,
	}, nil
}

func FormatColumnTime(fmtTimeStampz string, defaultTimeZone *time.Location, prop *TimeProp, t time.Time) (v interface{}) {
	if t.IsZero() {
		if prop.Nullable {
			return nil
		}
		return ""
	}

	if prop.TimeZone != nil {
		return FormatTime(fmtTimeStampz, prop.TypeName, t.In(prop.TimeZone))
	}
	return FormatTime(fmtTimeStampz, prop.TypeName, t.In(defaultTimeZone))
}

type NullTime time.Time

var (
	_ driver.Valuer = NullTime{}
)

func (ns *NullTime) Time() time.Time {
	return time.Time(*ns)
}

func (ns *NullTime) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	return convertTime(ns, value)
}

// Value implements the driver Valuer interface.
func (ns NullTime) Value() (driver.Value, error) {
	if (time.Time)(ns).IsZero() {
		return nil, nil
	}
	return (time.Time)(ns).Format("2006-01-02 15:04:05"), nil
}

func convertTime(dest *NullTime, src interface{}) error {
	// Common cases, without reflect.
	switch s := src.(type) {
	case string:
		t, err := time.Parse("2006-01-02 15:04:05", s)
		if err != nil {
			return err
		}
		*dest = NullTime(t)
		return nil
	case []uint8:
		t, err := time.Parse("2006-01-02 15:04:05", string(s))
		if err != nil {
			return err
		}
		*dest = NullTime(t)
		return nil
	case time.Time:
		*dest = NullTime(s)
		return nil
	case nil:
	default:
		return fmt.Errorf("unsupported driver -> Scan pair: %T -> %T", src, dest)
	}
	return nil
}
