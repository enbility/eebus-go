package mgcp

import (
	"fmt"
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

func measurementUpdateValueType(
	errorName string,
	id *model.MeasurementIdType,
	valueSource *model.MeasurementValueSourceType,
	value float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
	evaluationStart *time.Time,
	evaluationEnd *time.Time,
) UpdateValueType {
	if id == nil {
		panic(fmt.Sprintf("%s is not supported by the use case MGCP, please check the MGCP configuration", errorName))
	}

	updateValueType := UpdateValueType{
		updateValueType: SupportedUpdateValueTypeMeasurement,
		updateTypeMeasurement: api.MeasurementDataForID{
			Id: *id,
			Data: model.MeasurementDataType{
				ValueType:   util.Ptr(model.MeasurementValueTypeTypeValue),
				ValueSource: valueSource,
				Value:       model.NewScaledNumberType(value),
			},
		},
	}

	if timestamp != nil {
		updateValueType.updateTypeMeasurement.Data.Timestamp = model.NewAbsoluteOrRelativeTimeTypeFromTime(*timestamp)
	}

	if valueState != nil {
		updateValueType.updateTypeMeasurement.Data.ValueState = valueState
	}

	if evaluationStart != nil && evaluationEnd != nil {
		updateValueType.updateTypeMeasurement.Data.EvaluationPeriod = &model.TimePeriodType{
			StartTime: model.NewAbsoluteOrRelativeTimeTypeFromTime(*evaluationStart),
			EndTime:   model.NewAbsoluteOrRelativeTimeTypeFromTime(*evaluationEnd),
		}
	}

	return updateValueType
}
