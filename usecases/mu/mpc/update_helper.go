package mpc

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

type UpdateData struct {
	supported         bool
	notSupportedError error
	measurementData   api.MeasurementDataForID
}

type serUpdateData struct {
	Supported         bool
	NotSupportedError string
	MeasurementData   api.MeasurementDataForID
}

func (r *UpdateData) UnmarshalJSON(data []byte) error {
	aux := serUpdateData{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	r.supported = aux.Supported
	if aux.NotSupportedError != "" {
		r.notSupportedError = errors.New(aux.NotSupportedError)
	}
	r.measurementData = aux.MeasurementData

	return nil
}

func (r *UpdateData) MarshalJSON() ([]byte, error) {
	aux := serUpdateData{
		Supported:       r.supported,
		MeasurementData: r.measurementData,
	}
	if r.notSupportedError != nil {
		aux.NotSupportedError = r.notSupportedError.Error()
	}

	return json.Marshal(aux)
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
