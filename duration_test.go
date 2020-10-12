package microtime

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestDurationMarshalJson(t *testing.T) {
	for index, test := range []struct {
		duration      Duration
		expectedValue []byte
	}{
		{
			duration:      Duration(0),
			expectedValue: []byte("null"),
		},
		{
			duration:      Duration(-1 << 63),
			expectedValue: []byte(`"-2562047h47m16.854775808s"`),
		},
		{
			duration:      Duration(1<<63 - 1),
			expectedValue: []byte(`"2562047h47m16.854775807s"`),
		},
	} {
		t.Run(fmt.Sprintf("Case %d", index+1), func(t *testing.T) {
			result, err := json.Marshal(test.duration)
			if err != nil {
				t.Fatalf("expected: nil, got: %v", err)
			}
			resultPtr, err := json.Marshal(&test.duration)
			if err != nil {
				t.Fatalf("expected: nil, got: %v", err)
			}
			assert.Equal(t, test.expectedValue, result)
			assert.Equal(t, test.expectedValue, resultPtr)
		})
	}

}

func TestDurationUnmarshalJson(t *testing.T) {
	for index, test := range []struct {
		duration      []byte
		expectedValue Duration
	}{
		{
			duration:      []byte("null"),
			expectedValue: Duration(-1),
		},
		{
			duration:      []byte(`"0"`),
			expectedValue: Duration(0),
		},
		{
			duration:      []byte(`"10s"`),
			expectedValue: Second * 10,
		},
		{
			duration:      []byte(`"10m"`),
			expectedValue: Minute * 10,
		},
		{
			duration:      []byte(`"10h"`),
			expectedValue: Hour * 10,
		},
	} {
		t.Run(fmt.Sprintf("Case %d", index+1), func(t *testing.T) {
			result := Duration(-1)
			err := json.Unmarshal(test.duration, &result)
			if err != nil {
				t.Fatalf("expected: nil, got: %v", err)
			}
			assert.Equal(t, test.expectedValue, result)
		})
	}
}

func TestDurationUnmarshalJsonErrorCases(t *testing.T) {
	for index, test := range []struct {
		duration    []byte
		expectedErr string
	}{
		{
			duration:    nil,
			expectedErr: "unexpected end of JSON input",
		},
		{
			duration:    []byte("-9223372036854775808"),
			expectedErr: "duration must be valid json string: invalid syntax",
		},
		{
			duration:    []byte(`"9223372036854775807"`),
			expectedErr: "time: missing unit in duration \"9223372036854775807\"",
		},
	} {
		t.Run(fmt.Sprintf("Case %d", index+1), func(t *testing.T) {
			result := Duration(-1)
			err := json.Unmarshal(test.duration, &result)
			assert.Equal(t, test.expectedErr, err.Error())
		})
	}

}

func TestDurationMarshalBinary(t *testing.T) {
	for index, test := range []struct {
		duration      Duration
		expectedValue []byte
	}{
		{
			duration:      Duration(0),
			expectedValue: nil,
		},
		{
			duration:      Duration(-1 << 63),
			expectedValue: []byte("-2562047h47m16.854775808s"),
		},
		{
			duration:      Duration(1<<63 - 1),
			expectedValue: []byte("2562047h47m16.854775807s"),
		},
	} {
		t.Run(fmt.Sprintf("Case %d", index+1), func(t *testing.T) {
			result, err := test.duration.MarshalBinary()
			if err != nil {
				t.Fatalf("expected: nil, got: %v", err)
			}
			assert.Equal(t, test.expectedValue, result)
		})
	}

}

func TestDurationUnmarshalBinary(t *testing.T) {
	for index, test := range []struct {
		duration      []byte
		expectedValue Duration
	}{
		{
			duration:      nil,
			expectedValue: Duration(-1),
		},
		{
			duration:      []byte("0"),
			expectedValue: Duration(0),
		},
		{
			duration:      []byte("10s"),
			expectedValue: Second * 10,
		},
		{
			duration:      []byte("10m"),
			expectedValue: Minute * 10,
		},
		{
			duration:      []byte("10h"),
			expectedValue: Hour * 10,
		},
	} {
		t.Run(fmt.Sprintf("Case %d", index+1), func(t *testing.T) {
			result := Duration(-1)
			err := result.UnmarshalBinary(test.duration)
			if err != nil {
				t.Fatalf("expected: nil, got: %v", err)
			}
			assert.Equal(t, test.expectedValue, result)
		})
	}
}
