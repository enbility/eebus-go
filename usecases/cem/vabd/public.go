package vabd

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/client"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

// return the current battery (dis-)charge power (W)
//
//   - positive values charge power
//   - negative values discharge power
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (e *VABD) Power(entity spineapi.EntityRemoteInterface) (float64, error) {
	if !e.IsCompatibleEntityType(entity) {
		return 0, api.ErrNoCompatibleEntity
	}

	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACPowerTotal),
	}
	return e.getValuesFoFilter(entity, filter)
}

// return the total charge energy (Wh)
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (e *VABD) EnergyCharged(entity spineapi.EntityRemoteInterface) (float64, error) {
	if !e.IsCompatibleEntityType(entity) {
		return 0, api.ErrNoCompatibleEntity
	}

	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeEnergy),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		ScopeType:       util.Ptr(model.ScopeTypeTypeCharge),
	}
	return e.getValuesFoFilter(entity, filter)
}

// return the total discharge energy (Wh)
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (e *VABD) EnergyDischarged(entity spineapi.EntityRemoteInterface) (float64, error) {
	if !e.IsCompatibleEntityType(entity) {
		return 0, api.ErrNoCompatibleEntity
	}

	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeEnergy),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		ScopeType:       util.Ptr(model.ScopeTypeTypeDischarge),
	}
	return e.getValuesFoFilter(entity, filter)
}

// return the current state of charge in %
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (e *VABD) StateOfCharge(entity spineapi.EntityRemoteInterface) (float64, error) {
	if !e.IsCompatibleEntityType(entity) {
		return 0, api.ErrNoCompatibleEntity
	}

	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypePercentage),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		ScopeType:       util.Ptr(model.ScopeTypeTypeStateOfCharge),
	}
	return e.getValuesFoFilter(entity, filter)
}

// helper

func (e *VABD) getValuesFoFilter(
	entity spineapi.EntityRemoteInterface,
	filter model.MeasurementDescriptionDataType,
) (float64, error) {
	if entity == nil {
		return 0, api.ErrDeviceDisconnected
	}

	measurementF, err := client.NewMeasurement(e.LocalEntity, entity)
	if err != nil {
		return 0, api.ErrFunctionNotSupported
	}

	result, err := measurementF.GetDataForFilter(filter)
	if err != nil || len(result) == 0 || result[0].Value == nil {
		return 0, api.ErrDataNotAvailable
	}
	return result[0].Value.GetValue(), nil
}
