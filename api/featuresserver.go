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
	// set the local device diagnosis state of the device
	SetLocalState(statetate *model.DeviceDiagnosisStateDataType)

	// set the local device diagnosis operating state
	SetLocalOperatingState(operatingState model.DeviceDiagnosisOperatingStateType)
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
	// Elements provided in deleteElements will be removed from the data set before the update
	//
	// Will return an error if the data set could not be updated
	UpdateLimitDataForId(
		data model.LoadControlLimitDataType,
		deleteElements *model.LoadControlLimitDataElementsType,
		limitId model.LoadControlLimitIdType,
	) error

	// Set or update data set for a filter
	// Elements provided in deleteElements will be removed from the data set before the update
	//
	// Will return an error if the data set could not be updated
	UpdateLimitDataForFilter(
		data model.LoadControlLimitDataType,
		deleteElements *model.LoadControlLimitDataElementsType,
		filter model.LoadControlLimitDescriptionDataType,
	) error
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
	// Elements provided in deleteElements will be removed from the data set before the update
	//
	// Will return an error if the data set could not be updated
	UpdateDataForId(
		data model.MeasurementDataType,
		deleteElements *model.MeasurementDataElementsType,
		measurementId model.MeasurementIdType,
	) error

	// Set or update data set for a filter
	// Elements provided in deleteElements will be removed from the data set before the update
	//
	// Will return an error if the data set could not be updated
	UpdateDataForFilter(
		data model.MeasurementDataType,
		deleteElements *model.MeasurementDataElementsType,
		filter model.MeasurementDescriptionDataType,
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
