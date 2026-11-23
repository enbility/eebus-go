package mdt

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/client"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

// Scenario 1

// return the domestic hot water temperature converted to the requested unit
//
// the DHW circuit may announce the temperature in degC, degF or K; the value is
// converted to the requested unit (one of degC, degF, K)
//
// possible errors:
//   - ErrDataNotAvailable if no such value is (yet) available
//   - ErrDataInvalid if the currently available data is invalid and should be ignored
//   - ErrNotSupported if the requested unit is not a supported temperature unit
//   - and others
func (e *MDT) Temperature(entity spineapi.EntityRemoteInterface, unit model.UnitOfMeasurementType) (float64, error) {
	if !e.IsCompatibleEntityType(entity) {
		return 0, api.ErrNoCompatibleEntity
	}

	measurement, err := client.NewMeasurement(e.LocalEntity, entity)
	if err != nil {
		return 0, err
	}

	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeTemperature),
		CommodityType:   util.Ptr(model.CommodityTypeTypeDomestichotwater),
		ScopeType:       util.Ptr(model.ScopeTypeTypeDhwTemperature),
	}
	data, err := measurement.GetDataForFilter(filter)
	if err != nil || len(data) == 0 || data[0].Value == nil || data[0].MeasurementId == nil {
		return 0, api.ErrDataNotAvailable
	}

	// if the value state is set and not normal, the value is not valid and should be ignored
	if data[0].ValueState != nil && *data[0].ValueState != model.MeasurementValueStateTypeNormal {
		return 0, api.ErrDataInvalid
	}

	description, err := measurement.GetDescriptionForId(*data[0].MeasurementId)
	if err != nil || description.Unit == nil {
		return 0, api.ErrDataNotAvailable
	}

	return convertTemperature(data[0].Value.GetValue(), *description.Unit, unit)
}

// convertTemperature converts a temperature value between degC, degF and K.
// It returns ErrDataInvalid if the announced unit is not a supported temperature
// unit, and ErrNotSupported if the requested unit is not.
func convertTemperature(value float64, from, to model.UnitOfMeasurementType) (float64, error) {
	var celsius float64
	switch from {
	case model.UnitOfMeasurementTypedegC:
		celsius = value
	case model.UnitOfMeasurementTypedegF:
		celsius = (value - 32) * 5 / 9
	case model.UnitOfMeasurementTypeK:
		celsius = value - 273.15
	default:
		return 0, api.ErrDataInvalid
	}

	switch to {
	case model.UnitOfMeasurementTypedegC:
		return celsius, nil
	case model.UnitOfMeasurementTypedegF:
		return celsius*9/5 + 32, nil
	case model.UnitOfMeasurementTypeK:
		return celsius + 273.15, nil
	default:
		return 0, api.ErrNotSupported
	}
}
