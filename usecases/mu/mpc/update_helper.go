package mpc

import (
	"errors"
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"time"
)

type UpdateData struct {
	supported         bool
	notSupportedError error
	measurementData   api.MeasurementDataForID
}

func (u *UpdateData) Supported() bool {
	return u.supported
}

func (u *UpdateData) NotSupportedError() error {
	return u.notSupportedError
}

func (u *UpdateData) MeasurementData() api.MeasurementDataForID {
	return u.measurementData
}

func newUpdateData(
	errorString string,
	id *model.MeasurementIdType,
	data *model.MeasurementDataType,
) *UpdateData {
	if id == nil || data == nil {
		return &UpdateData{
			supported:         false,
			notSupportedError: errors.New(errorString),
		}
	} else {
		return &UpdateData{
			supported: true,
			measurementData: api.MeasurementDataForID{
				Id:   *id,
				Data: *data,
			},
		}
	}
}

func measurementData(
	value float64,
	timestamp *time.Time,
	valueSource *model.MeasurementValueSourceType,
	valueState *model.MeasurementValueStateType,
	evaluationStart *time.Time,
	evaluationEnd *time.Time,
) *model.MeasurementDataType {
	measurement := model.MeasurementDataType{
		ValueType:   util.Ptr(model.MeasurementValueTypeTypeValue),
		Value:       model.NewScaledNumberType(value),
		ValueSource: valueSource,
		ValueState:  valueState,
	}

	if timestamp != nil {
		measurement.Timestamp = model.NewAbsoluteOrRelativeTimeTypeFromTime(*timestamp)
	}

	if evaluationStart != nil && evaluationEnd != nil {
		measurement.EvaluationPeriod = &model.TimePeriodType{
			StartTime: model.NewAbsoluteOrRelativeTimeTypeFromTime(*evaluationStart),
			EndTime:   model.NewAbsoluteOrRelativeTimeTypeFromTime(*evaluationEnd),
		}
	}

	return &measurement
}
