package mgcp

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"time"
)

type SupportedUpdateValueType int

const (
	SupportedUpdateValueTypeMeasurement   SupportedUpdateValueType = 0
	SupportedUpdateValueTypeConfiguration SupportedUpdateValueType = 1
)

type UpdateValueType struct {
	updateValueType         SupportedUpdateValueType
	updateTypeMeasurement   api.MeasurementDataForID
	updateTypeConfiguration model.DeviceConfigurationKeyValueDataType
}

func measuredValue(value float64) model.MeasurementDataType {
	return model.MeasurementDataType{
		ValueType:   util.Ptr(model.MeasurementValueTypeTypeValue),
		Timestamp:   model.NewAbsoluteOrRelativeTimeTypeFromTime(time.Now()),
		Value:       model.NewScaledNumberType(value),
		ValueSource: util.Ptr(model.MeasurementValueSourceTypeMeasuredValue),
		ValueState:  util.Ptr(model.MeasurementValueStateTypeNormal),
	}
}
