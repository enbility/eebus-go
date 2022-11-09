package model

import (
	"testing"
	"time"
)

func TestTimeType(t *testing.T) {
	tc := []struct {
		in    string
		parse string
	}{
		{"21:32:52.12679", "15:04:05.999999999"},
		{"21:32:52.12679Z", "15:04:05.999999999Z"},
		{"21:32:52", "15:04:05"},
		{"19:32:52Z", "15:04:05Z"},
		{"19:32:52+07:00", "15:04:05+07:00"},
		{"19:32:52-07:00", "15:04:05-07:00"},
	}

	for _, tc := range tc {
		got := NewTimeType(tc.in)
		expect, err := time.Parse(tc.parse, tc.in)
		if err != nil {
			t.Errorf("Parsing failure with %s and parser %s: %s", tc.in, tc.parse, err)
			continue
		}
		value, err := got.GetTime()
		if err != nil {
			t.Errorf("Test Failure with %s and parser %s: %s", tc.in, tc.parse, err)
			continue
		}

		if value.UTC() != expect.UTC() {
			t.Errorf("Test failure for %s, expected %s and got %s", tc.in, value.String(), expect.String())
		}
	}
}

func TestDateTimeType(t *testing.T) {
	tc := []struct {
		in    string
		parse string
	}{
		{"2022-02-01T21:32:52.12679", "2006-01-02T15:04:05.999999999"},
		{"2022-02-01T21:32:52.12679Z", "2006-01-02T15:04:05.999999999Z"},
		{"2022-02-01T21:32:52", "2006-01-02T15:04:05"},
		{"2022-02-01T19:32:52Z", "2006-01-02T15:04:05Z"},
		{"2022-02-01T19:32:52+07:00", "2006-01-02T15:04:05+07:00"},
		{"2022-02-01T19:32:52-07:00", "2006-01-02T15:04:05-07:00"},
	}

	for _, tc := range tc {
		got := NewDateTimeType(tc.in)
		expect, err := time.Parse(tc.parse, tc.in)
		if err != nil {
			t.Errorf("Parsing failure with %s and parser %s: %s", tc.in, tc.parse, err)
			continue
		}
		value, err := got.GetTime()
		if err != nil {
			t.Errorf("Test Failure with %s and parser %s: %s", tc.in, tc.parse, err)
			continue
		}

		if value.UTC() != expect.UTC() {
			t.Errorf("Test failure for %s, expected %s and got %s", tc.in, value.String(), expect.String())
		}
	}
}

func TestNewScaledNumberType(t *testing.T) {
	tc := []struct {
		in     float64
		number int64
		scale  int
	}{
		{0.0, 0, 0},
		{0.1, 1, -1},
		{1.0, 1, 0},
		{6.25, 625, -2},
		{10.0, 10, 0},
		{12.5952, 125952, -4},
		{13.1637, 131637, -4},
	}

	for _, tc := range tc {
		got := NewScaledNumberType(tc.in)
		number := int64(*got.Number)
		scale := 0
		if got.Scale != nil {
			scale = int(*got.Scale)
		}
		if number != tc.number || scale != tc.scale {
			t.Errorf("NewScaledNumberType(%v) = %d %d, want %d %d", tc.in, got.Number, got.Scale, tc.number, tc.scale)
		}

		val := got.GetValue()
		if val != tc.in {
			t.Errorf("GetValue(%d %d) = %f, want %f", tc.number, tc.scale, val, tc.in)
		}
	}
}
