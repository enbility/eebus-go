package vapd

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/client"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

// return the current photovoltaic production power (W)
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (e *VAPD) Power(entity spineapi.EntityRemoteInterface) (float64, error) {
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

// return the nominal photovoltaic peak power (W)
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (e *VAPD) PowerNominalPeak(entity spineapi.EntityRemoteInterface) (float64, error) {
	if !e.IsCompatibleEntityType(entity) {
		return 0, api.ErrNoCompatibleEntity
	}

	deviceConfiguration, err := client.NewDeviceConfiguration(e.LocalEntity, entity)
	if err != nil {
		return 0, api.ErrFunctionNotSupported
	}

	keyName := model.DeviceConfigurationKeyNameTypePeakPowerOfPVSystem
	filter := model.DeviceConfigurationKeyValueDescriptionDataType{
		KeyName: &keyName,
	}
	if _, err := deviceConfiguration.GetKeyValueDescriptionsForFilter(filter); err != nil {
		return 0, err
	}

	filter = model.DeviceConfigurationKeyValueDescriptionDataType{
		KeyName:   &keyName,
		ValueType: util.Ptr(model.DeviceConfigurationKeyValueTypeTypeScaledNumber),
	}
	data, err := deviceConfiguration.GetKeyValueDataForFilter(filter)
	if err != nil || data == nil || data.Value == nil || data.Value.ScaledNumber == nil {
		return 0, api.ErrDataNotAvailable
	}

	return data.Value.ScaledNumber.GetValue(), nil
}

// return the total photovoltaic yield (Wh)
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (e *VAPD) PVYieldTotal(entity spineapi.EntityRemoteInterface) (float64, error) {
	if !e.IsCompatibleEntityType(entity) {
		return 0, api.ErrNoCompatibleEntity
	}

	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeEnergy),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACYieldTotal),
	}
	return e.getValuesFoFilter(entity, filter)
}

// helper

func (e *VAPD) getValuesFoFilter(
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
