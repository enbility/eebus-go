package model

import (
	"testing"
)

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
