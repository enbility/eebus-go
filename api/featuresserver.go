package api

import "github.com/enbility/spine-go/model"

type DeviceClassificationServerInterface interface {
}

type DeviceConfigurationServerInterface interface {
	DeviceConfigurationCommonInterface

	// Add a new description data set and return the keyId
	//
	// will return nil if the data set could not be added
	AddKeyValueDescription(description model.DeviceConfigurationKeyValueDescriptionDataType) *model.DeviceConfigurationKeyIdType

	// Set or update data set for a keyId
	// Elements provided in deleteElements will be removed from the data set before the update
	//
	// Will return an error if the data set could not be updated
	UpdateKeyValueDataForKeyId(
		data model.DeviceConfigurationKeyValueDataType,
		deleteElements *model.DeviceConfigurationKeyValueDataElementsType,
		keyId model.DeviceConfigurationKeyIdType,
	) error

	// Set or update data set for a filter
	// Elements provided in deleteElements will be removed from the data set before the update
	//
	// Will return an error if the data set could not be updated
	UpdateKeyValueDataForFilter(
		data model.DeviceConfigurationKeyValueDataType,
		deleteElements *model.DeviceConfigurationKeyValueDataElementsType,
		filter model.DeviceConfigurationKeyValueDescriptionDataType,
	) error
}

type DeviceDiagnosisServerInterface interface {
	// set the local diagnosis state of the device
	SetLocalState(operatingState *model.DeviceDiagnosisStateDataType)
}

type ElectricalConnectionServerInterface interface {
	ElectricalConnectionCommonInterface

	// Add a new description data set
	//
	// NOTE: the electricalConnectionId has to be provided
	//
	// will return nil if the data set could not be added
	AddDescription(
		description model.ElectricalConnectionDescriptionDataType,
	) error

	// Add a new parameter description data sett and return the parameterId
	//
	// NOTE: the electricalConnectionId has to be provided, parameterId may not be provided
	//
	// will return nil if the data set could not be added
	AddParameterDescription(
		description model.ElectricalConnectionParameterDescriptionDataType,
	) *model.ElectricalConnectionParameterIdType

	// Add a new characteristic data set
	//
	// Note: ElectricalConnectionId and ParameterId must be set, CharacteristicId will be set automatically
	//
	// Will return an error if the data set could not be added
	AddCharacteristic(data model.ElectricalConnectionCharacteristicDataType) (*model.ElectricalConnectionCharacteristicIdType, error)

	// Update data set for a filter
	// Elements provided in deleteElements will be removed from the data set before the update
	//
	// // ElectricalConnectionId, ParameterId and CharacteristicId must be set
	//
	// Will return an error if the data set could not be updated
	UpdateCharacteristic(
		data model.ElectricalConnectionCharacteristicDataType,
		deleteElements *model.ElectricalConnectionCharacteristicDataElementsType,
	) error
}

type LoadControlLimitDataForID struct {
	Data model.LoadControlLimitDataType
	Id   model.LoadControlLimitIdType
}

type LoadControlLimitDataForFilter struct {
	Data   model.LoadControlLimitDataType
	Filter model.LoadControlLimitDescriptionDataType
}

type LoadControlServerInterface interface {
	// Add a new description data set and return the limitId
	//
	// NOTE: the limitId may not be provided
	//
	// will return nil if the data set could not be added
	AddLimitDescription(
		description model.LoadControlLimitDescriptionDataType,
	) *model.LoadControlLimitIdType

	// Set or update data set for a limitId
	// Id provided in deleteId will trigger removal of matching items from the data set before the update
	// Elements provided in deleteElement will limit the fields to be removed using Id
	//
	// Will return an error if the data set could not be updated
	UpdateLimitDataForIds(
		data []LoadControlLimitDataForID,
		deleteId *model.LoadControlLimitIdType,
		deleteElements *model.LoadControlLimitDataElementsType,
	) error

	// Set or update data set for a filter
	// Id provided in deleteId will trigger removal of matching items from the data set before the update
	// Elements provided in deleteElement will limit the fields to be removed using Id
	//
	// Will return an error if the data set could not be updated
	UpdateLimitDataForFilter(
		data []LoadControlLimitDataForFilter,
		deleteSelector *model.LoadControlLimitListDataSelectorsType,
		deleteElements *model.LoadControlLimitDataElementsType,
	) error
}

type MeasurementDataForID struct {
	Data model.MeasurementDataType
	Id   model.MeasurementIdType
}

type MeasurementDataForFilter struct {
	Data   model.MeasurementDataType
	Filter model.MeasurementDescriptionDataType
}

type MeasurementServerInterface interface {
	// Add a new parameter description data sett and return the measurementId
	//
	// NOTE: the measurementId may not be provided
	//
	// will return nil if the data set could not be added
	AddDescription(
		description model.MeasurementDescriptionDataType,
	) *model.MeasurementIdType

	// Set or update data set for a measurementId
	// Id provided in deleteId will trigger removal of matching items from the data set before the update
	// Elements provided in deleteElement will limit the fields to be removed using Id
	//
	// Will return an error if the data set could not be updated
	UpdateDataForIds(
		data []MeasurementDataForID,
		deleteId *model.MeasurementIdType,
		deleteElements *model.MeasurementDataElementsType,
	) error

	// Set or update data set for a filter
	// Selector provided in deleteSelector will trigger removal of matching items from the data set before the update
	// Elements provided in deleteElement will limit the fields to be removed using Selector
	//
	// Will return an error if the data set could not be updated
	UpdateDataForFilters(
		data []MeasurementDataForFilter,
		deleteSelector *model.MeasurementListDataSelectorsType,
		deleteElements *model.MeasurementDataElementsType,
	) error
}

type IdentificationServerInterface interface {
}

type IncentiveTableServerInterface interface {
}

type SmartEnergyManagementPsServerInterface interface {
}

type TimeSeriesServerInterface interface {
}
