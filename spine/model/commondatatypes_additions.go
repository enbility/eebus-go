package model

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/rickb777/date/period"
)

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

func NewISO8601Duration(duration time.Duration) *string {
	d, _ := period.NewOf(duration)
	value := d.String()
	return &value
}
