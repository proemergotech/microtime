package microtime

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var now = time.Now()

func TestSub(t *testing.T) {
	for index, test := range []struct {
		first    Time
		second   Time
		expected Duration
	}{
		{
			first:    Time{now.Add(1000)},
			second:   Time{now},
			expected: Duration(1000),
		},
		{
			first:    Time{now.Add(330000)},
			second:   Time{now},
			expected: Duration(330000),
		},
	} {
		t.Run(fmt.Sprintf("Case %d: %v sub %v", index+1, test.first, test.second), func(t *testing.T) {
			result := test.first.Sub(test.second)
			assert.Equal(t, test.expected, result)
		})
	}

}

func TestString(t *testing.T) {
	for index, test := range []struct {
		time     Time
		expected string
	}{
		{
			time:     Time{time.Date(1987, time.February, 4, 16, 52, 19, 330000, time.UTC)},
			expected: "1987-02-04T16:52:19.00033Z",
		},
		{
			time:     Time{time.Date(2019, time.January, 28, 16, 54, 58, 99000000, time.UTC)},
			expected: "2019-01-28T16:54:58.099Z",
		},
	} {
		t.Run(fmt.Sprintf("Case %d: %v -> %v", index+1, test.time, test.expected), func(t *testing.T) {
			result := test.time.String()
			assert.Equal(t, test.expected, result)
		})
	}

}

func TestFromString(t *testing.T) {
	for index, test := range []struct {
		time          string
		expectedValue Time
		expectedErr   error
	}{
		{
			time:          "1987-02-04T16:52:19.00033Z",
			expectedValue: Time{time.Date(1987, time.February, 4, 16, 52, 19, 330000, time.UTC)},
		},
		{
			time:          "2019-01-28T16:54:58.099Z",
			expectedValue: Time{time.Date(2019, time.January, 28, 16, 54, 58, 99000000, time.UTC)},
		},
		{
			time:          "1987-02-04T16:52:19.00033000Z",
			expectedValue: Time{time.Date(1987, time.February, 4, 16, 52, 19, 330000, time.UTC)},
		},
		{
			time:          "2019-01-28T16:54:58.099000Z",
			expectedValue: Time{time.Date(2019, time.January, 28, 16, 54, 58, 99000000, time.UTC)},
		},
		{
			time:          "2019-01-28 16:54:58.099",
			expectedValue: Time{time.Date(2019, time.January, 28, 16, 54, 58, 99000000, time.UTC)},
		},
		{
			time:          "2019-01-28 16:54:58.099Z",
			expectedValue: Time{time.Date(2019, time.January, 28, 16, 54, 58, 99000000, time.UTC)},
		},
		{
			time:          "2019-01-28T16:54:58.099",
			expectedValue: Time{time.Date(2019, time.January, 28, 16, 54, 58, 99000000, time.UTC)},
		},
		{
			time:          "2019-01-28 16:54:58.09900000",
			expectedValue: Time{time.Date(2019, time.January, 28, 16, 54, 58, 99000000, time.UTC)},
		},
		{
			time:          "2019012816545809900000",
			expectedValue: Time{},
			expectedErr:   fmt.Errorf("Could not find format for %q", "2019012816545809900000"),
		},
		{
			time:          "-20190128165458099",
			expectedValue: Time{},
			expectedErr:   fmt.Errorf("Could not find format for %q", "-20190128165458099"),
		},
	} {
		t.Run(fmt.Sprintf("Case %d: %v -> %v", index+1, test.time, test.expectedValue), func(t *testing.T) {
			result, err := FromString(test.time)
			assert.Equal(t, test.expectedErr, err)
			assert.Equal(t, test.expectedValue, result)
		})
	}

}

func TestRound(t *testing.T) {
	for index, test := range []struct {
		dur      Duration
		rounder  Duration
		expected Duration
	}{
		{
			dur:      Duration(123456789101112),
			rounder:  Hour,
			expected: Duration(122400000000000),
		},
		{
			dur:      Duration(123456789101112),
			rounder:  Minute,
			expected: Duration(123480000000000),
		},
		{
			dur:      Duration(123456789101112),
			rounder:  Second,
			expected: Duration(123457000000000),
		},
		{
			dur:      Duration(123456789101112),
			rounder:  Millisecond,
			expected: Duration(123456789000000),
		},
		{
			dur:      Duration(123456789101112),
			rounder:  Microsecond,
			expected: Duration(123456789101000),
		},
		{
			dur:      Duration(123456789101112),
			rounder:  Nanosecond,
			expected: Duration(123456789101112),
		},
	} {
		t.Run(fmt.Sprintf("Case %d: round %v -> %v", index+1, test.dur, test.expected), func(t *testing.T) {
			result := test.dur.Round(test.rounder)
			assert.Equal(t, test.expected, result)
		})
	}

}

func TestMarshalJson(t *testing.T) {
	for index, test := range []struct {
		time          Time
		expectedValue []byte
		expectedErr   error
	}{
		{
			time:          Time{time.Time{}},
			expectedValue: []byte("null"),
		},
		{
			time:          Time{time.Date(1987, time.February, 4, 16, 52, 19, 330000, time.UTC)},
			expectedValue: []byte("\"1987-02-04T16:52:19.00033Z\""),
		},
	} {
		t.Run(fmt.Sprintf("Case %d", index+1), func(t *testing.T) {
			result, err := test.time.MarshalJSON()
			assert.Equal(t, test.expectedErr, err)
			assert.Equal(t, test.expectedValue, result)
		})
	}

}

func TestUnmarshalJson(t *testing.T) {
	for index, test := range []struct {
		time          []byte
		expectedValue Time
		expectedErr   error
	}{
		{
			time:          []byte("null"),
			expectedValue: Time{time.Time{}},
		},
		{
			time:          []byte("\"1987-02-04T16:52:19.00033Z\""),
			expectedValue: Time{time.Date(1987, time.February, 4, 16, 52, 19, 330000, time.UTC)},
		},
		{
			time:          []byte("1987-02-04T16:52:19.00033Z"),
			expectedValue: Time{},
			expectedErr:   fmt.Errorf("invalid syntax"),
		},
		{
			time:          []byte("\"19870204x165n2190v0033\""),
			expectedValue: Time{},
			expectedErr:   fmt.Errorf("Could not find format for \"19870204x165n2190v0033\""),
		},
	} {
		t.Run(fmt.Sprintf("Case %d", index+1), func(t *testing.T) {
			result := Time{}
			err := result.UnmarshalJSON(test.time)
			assert.Equal(t, test.expectedErr, err)
			assert.Equal(t, test.expectedValue, result)
		})
	}

}