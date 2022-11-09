package model

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/rickb777/date/period"
)

// string as DateType
func GetDateFromString(s string) (time.Time, error) {
	value, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}, err
	}

	return value, nil
}

// string as DateTimeType
func GetDateTimeFromString(s string) (time.Time, error) {
	value, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}, err
	}

	return value, nil
}

// TimeType
func NewTimeType(t string) *TimeType {
	value := TimeType(t)
	return &value
}

func GetTime(s *TimeType) (time.Time, error) {
	if s == nil {
		return time.Time{}, fmt.Errorf("invalid time pointer")
	}

	value, err := time.Parse("15:04:05.999999999", string(*s))
	if err != nil {
		return time.Time{}, err
	}

	return value, nil
}

//  DurationType

func NewDurationType(duration time.Duration) *DurationType {
	d, _ := period.NewOf(duration)
	value := DurationType(d.String())
	return &value
}

func (d *DurationType) GetTimeDuration() (time.Duration, error) {
	return getTimeDurationFromString(string(*d))
}

// helper for DurationType and AbsoluteOrRelativeTimeType
func getTimeDurationFromString(s string) (time.Duration, error) {
	p, err := period.Parse(string(s))
	if err != nil {
		return 0, err
	}

	return p.DurationApprox(), nil
}

// AbsoluteOrRelativeTimeType
// can be of type TimeType or DurationType

func (a *AbsoluteOrRelativeTimeType) GetTime() *TimeType {
	value := NewTimeType(string(*a))
	return value
}

func (a *AbsoluteOrRelativeTimeType) GetDuration() (*DurationType, error) {
	value, err := a.GetTimeDuration()
	if err != nil {
		return nil, err
	}

	return NewDurationType(value), nil
}

func (a *AbsoluteOrRelativeTimeType) GetTimeDuration() (time.Duration, error) {
	return getTimeDurationFromString(string(*a))
}

// ScaledNumberType

func (m *ScaledNumberType) GetValue() float64 {
	if m.Number == nil {
		return 0
	}
	var scale float64 = 0
	if m.Scale != nil {
		scale = float64(*m.Scale)
	}
	return float64(*m.Number) * math.Pow(10, scale)
}

func NewScaledNumberType(value float64) *ScaledNumberType {
	m := &ScaledNumberType{}

	numberOfDecimals := 0
	temp := strconv.FormatFloat(value, 'f', -1, 64)
	index := strings.IndexByte(temp, '.')
	if index > -1 {
		numberOfDecimals = len(temp) - index - 1
	}

	// We limit this to 4 digits for now
	if numberOfDecimals > 4 {
		numberOfDecimals = 4
	}

	numberValue := NumberType(math.Trunc(value * math.Pow(10, float64(numberOfDecimals))))
	m.Number = &numberValue

	if numberValue != 0 {
		scaleValue := ScaleType(-numberOfDecimals)
		m.Scale = &scaleValue
	}

	return m
}

// FeatureAddressType

func (r *FeatureAddressType) String() string {
	if r == nil {
		return ""
	}

	var result string = ""
	if r.Device != nil {
		result += string(*r.Device)
	}
	result += ":["
	for _, id := range r.Entity {
		result += fmt.Sprintf("%d,", id)
	}
	result += "]:"
	if r.Feature != nil {
		result += fmt.Sprintf("%d", *r.Feature)
	}
	return result
}
