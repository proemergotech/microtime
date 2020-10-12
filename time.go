package microtime

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"time"

	"github.com/araddon/dateparse"
	"github.com/pkg/errors"
)

type Time struct {
	time.Time
}

func Now() Time {
	return Time{Time: time.Now().UTC()}
}

func (t Time) Sub(u Time) Duration {
	return Duration(t.Time.Sub(u.Time))
}

func (t Time) String() string {
	return t.Time.UTC().Round(time.Microsecond).Format(time.RFC3339Nano)
}

func FromString(str string) (Time, error) {
	tim, err := dateparse.ParseStrict(str)
	if err != nil {
		return Time{}, err
	}

	return Time{tim.UTC().Round(time.Microsecond)}, nil
}

func (t Time) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte("null"), nil
	}

	return t.UTC().Round(time.Microsecond).MarshalJSON()
}

func (t *Time) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	unquoted, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}

	t.Time, err = dateparse.ParseStrict(unquoted)
	if err != nil {
		return err
	}
	t.Time = t.Time.UTC()

	return nil
}

func (t *Time) UnmarshalParam(data string) error {
	quotedData := data
	if _, err := strconv.Unquote(data); err != nil {
		quotedData = strconv.Quote(data)
	}

	return t.UnmarshalJSON([]byte(quotedData))
}

func (t Time) MarshalBinary() (data []byte, err error) {
	if t.Time.IsZero() {
		return nil, nil
	}
	return []byte(t.String()), nil
}

func (t *Time) UnmarshalBinary(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	var err error
	*t, err = FromString(string(data))

	return err
}

func (t Time) RedisArg() interface{} {
	return strconv.FormatInt(t.Unix(), 10)
}

func (t *Time) RedisScan(src interface{}) error {
	if src == nil {
		return nil
	}

	var str string
	switch val := src.(type) {
	case []byte:
		str = string(val)
	case string:
		str = val
	default:
		return errors.Errorf("schema.RedisScan: invalid time: %v", src)
	}

	unixTime, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return errors.Errorf("schema.RedisScan: invalid time: %v", str)
	}

	*t = Time{time.Unix(unixTime, 0)}

	return nil
}

func (t Time) Value() (driver.Value, error) {
	if t.Time.IsZero() {
		return nil, nil
	}
	return t.UTC().Round(time.Microsecond), nil
}

func (t *Time) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	if srcTime, ok := src.(time.Time); ok {
		*t = Time{srcTime.UTC()}
	} else {
		return fmt.Errorf("microtime: cannot convert value '%v(%T)' to microtime", src, src)
	}

	return nil
}
