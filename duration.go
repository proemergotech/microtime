package microtime

import (
	"strconv"
	"time"

	"github.com/pkg/errors"
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

func (d *Duration) MarshalJSON() ([]byte, error) {
	if d == nil {
		return []byte("null"), nil
	}
	s := time.Duration(*d).String()

	return []byte(strconv.Quote(s)), nil
}

func (d *Duration) UnmarshalJSON(data []byte) error {
	s := string(data)
	if s == "null" {
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
