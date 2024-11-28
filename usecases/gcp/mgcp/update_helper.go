package mgcp

import (
	"fmt"
	"github.com/enbility/eebus-go/api"
	usecaseapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"time"
)

type UpdateValueType struct {
	updateValueTypeType     usecaseapi.GcpMGCPUpdateValueTypeType
	updateTypeMeasurement   api.MeasurementDataForID
	updateTypeConfiguration model.DeviceConfigurationKeyValueDataType
}

func (u UpdateValueType) GetUpdateValueTypeType() usecaseapi.GcpMGCPUpdateValueTypeType {
	return u.updateValueTypeType
}

func (u UpdateValueType) GetUpdateValueTypeMeasurement() api.MeasurementDataForID {
	return u.updateTypeMeasurement
}

func (u UpdateValueType) GetUpdateValueTypeConfiguration() model.DeviceConfigurationKeyValueDataType {
	return u.updateTypeConfiguration
}

var _ usecaseapi.GcpMGCPUpdateValueTypeInterface = (*UpdateValueType)(nil)

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
		updateValueTypeType: usecaseapi.GcpMGCPUpdateValueTypeTypeMeasurement,
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
