package model

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/rickb777/date/period"
)

// TimeType xs:time

func NewTimeType(t string) *TimeType {
	value := TimeType(t)
	return &value
}

func (s *TimeType) GetTime() (time.Time, error) {
	allowedFormats := []string{
		"15:04:05.999999999",
		"15:04:05.999999999Z",
		"15:04:05",
		"15:04:05Z",
		"15:04:05+07:00",
		"15:04:05-07:00",
	}

	for _, format := range allowedFormats {
		if value, err := time.Parse(format, string(*s)); err == nil {
			return value, nil
		}
	}

	return time.Time{}, errors.New("unsupported time format")
}

// DateType xs:date

func NewDateType(t string) *DateType {
	value := DateType(t)
	return &value
}

// 2001-10-26, 2001-10-26+02:00, 2001-10-26Z, 2001-10-26+00:00, -2001-10-26, or -20000-04-01
func (d *DateType) GetTime() (time.Time, error) {
	allowedFormats := []string{
		"2006-01-02",
		"2006-01-02Z",
		"2006-01-02+07:00",
	}

	for _, format := range allowedFormats {
		if value, err := time.Parse(format, string(*d)); err == nil {
			return value, nil
		}
	}

	return time.Time{}, errors.New("unsupported date format")
}

// DateTimeType xs:datetime

func NewDateTimeType(t string) *DateTimeType {
	value := DateTimeType(t)
	return &value
}

func (d *DateTimeType) GetTime() (time.Time, error) {
	allowedFormats := []string{
		"2006-01-02T15:04:05.999999999",
		"2006-01-02T15:04:05.999999999Z",
		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05+07:00",
		"2006-01-02T15:04:05-07:00",
		time.RFC3339,
	}

	for _, format := range allowedFormats {
		if value, err := time.Parse(format, string(*d)); err == nil {
			return value, nil
		}
	}

	return time.Time{}, errors.New("unsupported datetime format")
}

// DurationType

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

func (a *AbsoluteOrRelativeTimeType) GetDateTimeType() *DateTimeType {
	value := NewDateTimeType(string(*a))
	return value
}

func (a *AbsoluteOrRelativeTimeType) GetTime() (time.Time, error) {
	value := NewDateTimeType(string(*a))
	t, err := value.GetTime()
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
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
