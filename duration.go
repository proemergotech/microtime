package microtime

import (
	"strconv"
	"time"

	"github.com/proemergotech/errors/v2"
)

const (
	Nanosecond  = Duration(time.Nanosecond)
	Microsecond = Duration(time.Microsecond)
	Millisecond = Duration(time.Millisecond)
	Second      = Duration(time.Second)
	Minute      = Duration(time.Minute)
	Hour        = Duration(time.Hour)
)

type Duration time.Duration

func (d Duration) Round(m Duration) Duration {
	return Duration(time.Duration(d).Round(time.Duration(m)))
}

func (d Duration) RedisArg() interface{} {
	return strconv.FormatInt(int64(d), 10)
}

func (d *Duration) RedisScan(src interface{}) error {
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
		return errors.Errorf("schema.RedisScan: invalid duration: %v", src)
	}

	dur, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return errors.Errorf("schema.RedisScan: invalid time: %v", str)
	}

	*d = Duration(dur)

	return nil
}

func (d Duration) MarshalJSON() ([]byte, error) {
	if d == Duration(0) {
		return []byte("null"), nil
	}

	return []byte(strconv.Quote(d.String())), nil
}

func (d *Duration) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	unquoted, err := strconv.Unquote(string(data))
	if err != nil {
		return errors.Wrap(err, "duration must be valid json string")
	}

	duration, err := time.ParseDuration(unquoted)
	if err != nil {
		return errors.WithStack(err)
	}

	*d = Duration(duration)

	return nil
}

func (d *Duration) UnmarshalParam(data string) error {
	quotedData := data
	if _, err := strconv.Unquote(data); err != nil {
		quotedData = strconv.Quote(data)
	}

	return d.UnmarshalJSON([]byte(quotedData))
}

func (d Duration) MarshalBinary() (data []byte, err error) {
	if d == Duration(0) {
		return nil, nil
	}
	return []byte(d.String()), nil
}

func (d *Duration) UnmarshalBinary(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	duration, err := time.ParseDuration(string(data))
	if err != nil {
		return errors.WithStack(err)
	}

	*d = Duration(duration)

	return nil
}

func (d Duration) String() string {
	return time.Duration(d).String()
}
