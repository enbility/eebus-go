package mgcp

import (
	"fmt"
	"github.com/enbility/eebus-go/api"
	usecaseapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"time"
)

type UpdateData struct {
	supported         bool
	notSupportedError error
}

type UpdateMeasurementData struct {
	UpdateData

	measurementData api.MeasurementDataForID
}

type UpdateConfigurationData struct {
	UpdateData

	configurationData model.DeviceConfigurationKeyValueDataType
}

func (u *UpdateData) Supported() bool {
	return u.supported
}

func (u *UpdateData) NotSupportedError() error {
	return u.notSupportedError
}

func (u *UpdateMeasurementData) MeasurementData() api.MeasurementDataForID {
	return u.measurementData
}

func (u UpdateConfigurationData) ConfigurationData() model.DeviceConfigurationKeyValueDataType {
	return u.configurationData
}

func updateMeasurementData(
	errorName string,
	id *model.MeasurementIdType,
	valueSource *model.MeasurementValueSourceType,
	value float64,
	timestamp *time.Time,
	valueState *model.MeasurementValueStateType,
	evaluationStart *time.Time,
	evaluationEnd *time.Time,
) usecaseapi.UpdateData {
	if id == nil {
		return &UpdateData{
			supported:         false,
			notSupportedError: fmt.Errorf("id is nil: %s, please check the mgcp configuration", errorName),
		}
	}

	updateValueType := UpdateMeasurementData{
		UpdateData: UpdateData{
			supported: true,
		},
		measurementData: api.MeasurementDataForID{
			Id: *id,
			Data: model.MeasurementDataType{
				ValueType:   util.Ptr(model.MeasurementValueTypeTypeValue),
				ValueSource: valueSource,
				Value:       model.NewScaledNumberType(value),
			},
		},
	}

	if timestamp != nil {
		updateValueType.measurementData.Data.Timestamp = model.NewAbsoluteOrRelativeTimeTypeFromTime(*timestamp)
	}

	if valueState != nil {
		updateValueType.measurementData.Data.ValueState = valueState
	}

	if evaluationStart != nil && evaluationEnd != nil {
		updateValueType.measurementData.Data.EvaluationPeriod = &model.TimePeriodType{
			StartTime: model.NewAbsoluteOrRelativeTimeTypeFromTime(*evaluationStart),
			EndTime:   model.NewAbsoluteOrRelativeTimeTypeFromTime(*evaluationEnd),
		}
	}

	return &updateValueType
}
