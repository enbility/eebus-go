package api

import "github.com/enbility/spine-go/model"

type DeviceClassificationClientInterface interface {
	// request DeviceClassificationManufacturerData from a remote device entity
	RequestManufacturerDetails() (*model.MsgCounterType, error)
}

type DeviceConfigurationClientInterface interface {
	DeviceConfigurationCommonInterface

	// request DeviceConfigurationKeyValueDescriptionDataType from a remote entity
	RequestKeyValueDescriptions(
		selector *model.DeviceConfigurationKeyValueDescriptionListDataSelectorsType,
		elements *model.DeviceConfigurationKeyValueDescriptionDataElementsType,
	) (*model.MsgCounterType, error)

	// request DeviceConfigurationKeyValueListData from a remote entity
	RequestKeyValues(
		selector *model.DeviceConfigurationKeyValueListDataSelectorsType,
		elements *model.DeviceConfigurationKeyValueDataElementsType,
	) (*model.MsgCounterType, error)

	// write key values
	// returns an error if this failed
	WriteKeyValues(data []model.DeviceConfigurationKeyValueDataType) (*model.MsgCounterType, error)
}

type DeviceDiagnosisClientInterface interface {
	// request DeviceDiagnosisStateData from a remote entity
	RequestState() (*model.MsgCounterType, error)

	// request FunctionTypeDeviceDiagnosisHeartbeatData from a remote device
	RequestHeartbeat() (*model.MsgCounterType, error)
}

type ElectricalConnectionClientInterface interface {
	// request ElectricalConnectionDescriptionListDataType from a remote entity
	RequestDescriptions(
		selector *model.ElectricalConnectionDescriptionListDataSelectorsType,
		elements *model.ElectricalConnectionDescriptionDataElementsType,
	) (*model.MsgCounterType, error)

	// request FunctionTypeElectricalConnectionParameterDescriptionListData from a remote entity
	RequestParameterDescriptions(
		selector *model.ElectricalConnectionParameterDescriptionListDataSelectorsType,
		elements *model.ElectricalConnectionParameterDescriptionDataElementsType,
	) (*model.MsgCounterType, error)

	// request FunctionTypeElectricalConnectionPermittedValueSetListData from a remote entity
	RequestPermittedValueSets(
		selector *model.ElectricalConnectionPermittedValueSetListDataSelectorsType,
		elements *model.ElectricalConnectionPermittedValueSetDataElementsType,
	) (*model.MsgCounterType, error)

	// request FunctionTypeElectricalConnectionCharacteristicListData from a remote entity
	RequestCharacteristics(
		selector *model.ElectricalConnectionCharacteristicListDataSelectorsType,
		elements *model.ElectricalConnectionCharacteristicDataElementsType,
	) (*model.MsgCounterType, error)
}

type IdentificationClientInterface interface {
	// request FunctionTypeIdentificationListData from a remote entity
	RequestValues() (*model.MsgCounterType, error)
}

type IncentiveTableClientInterface interface {
	// request FunctionTypeIncentiveTableDescriptionData from a remote entity
	RequestDescriptions() (*model.MsgCounterType, error)

	// request FunctionTypeIncentiveTableConstraintsData from a remote entity
	RequestConstraints() (*model.MsgCounterType, error)

	// request FunctionTypeIncentiveTableData from a remote entity
	RequestValues() (*model.MsgCounterType, error)

	// write incentivetable descriptions
	// returns an error if this failed
	WriteDescriptions(data []model.IncentiveTableDescriptionType) (*model.MsgCounterType, error)

	// write incentivetable descriptions
	// returns an error if this failed
	WriteValues(data []model.IncentiveTableType) (*model.MsgCounterType, error)
}

type LoadControlClientInterface interface {
	// request FunctionTypeLoadControlLimitDescriptionListData from a remote device
	RequestLimitDescriptions(
		selector *model.LoadControlLimitDescriptionListDataSelectorsType,
		elements *model.LoadControlLimitDescriptionDataElementsType,
	) (*model.MsgCounterType, error)

	// request FunctionTypeLoadControlLimitConstraintsListData from a remote device
	RequestLimitConstraints(
		selector *model.LoadControlLimitConstraintsListDataSelectorsType,
		elements *model.LoadControlLimitConstraintsDataElementsType,
	) (*model.MsgCounterType, error)

	// request FunctionTypeLoadControlLimitListData from a remote device
	RequestLimitData(
		selector *model.LoadControlLimitListDataSelectorsType,
		elements *model.LoadControlLimitDataElementsType,
	) (*model.MsgCounterType, error)

	// write load control limits
	// returns an error if this failed
	WriteLimitData(
		data []model.LoadControlLimitDataType,
		deleteSelectors *model.LoadControlLimitListDataSelectorsType,
		deleteElements *model.LoadControlLimitDataElementsType,
	) (*model.MsgCounterType, error)
}

type MeasurementClientInterface interface {
	// request FunctionTypeMeasurementDescriptionListData from a remote device
	RequestDescriptions(
		selector *model.MeasurementDescriptionListDataSelectorsType,
		elements *model.MeasurementDescriptionDataElementsType,
	) (*model.MsgCounterType, error)

	// request FunctionTypeMeasurementConstraintsListData from a remote entity
	RequestConstraints(
		selector *model.MeasurementConstraintsListDataSelectorsType,
		elements *model.MeasurementConstraintsDataElementsType,
	) (*model.MsgCounterType, error)

	// request FunctionTypeMeasurementListData from a remote entity
	RequestData(
		selector *model.MeasurementListDataSelectorsType,
		elements *model.MeasurementDataElementsType,
	) (*model.MsgCounterType, error)
}

type SmartEnergyManagementPsClientInterface interface {
	// request FunctionTypeSmartEnergyManagementPsData from a remote entity
	RequestData() (*model.MsgCounterType, error)

	// write SmartEnergyManagementPsData
	// returns an error if this failed
	WriteData(data *model.SmartEnergyManagementPsDataType) (*model.MsgCounterType, error)
}

type TimeSeriesClientInterface interface {
	// request FunctionTypeTimeSeriesDescriptionListData from a remote entity
	RequestDescriptions(
		selector *model.TimeSeriesDescriptionListDataSelectorsType,
		elements *model.TimeSeriesDescriptionDataElementsType,
	) (*model.MsgCounterType, error)

	// request FunctionTypeTimeSeriesConstraintsListData from a remote entity
	RequestConstraints(
		selector *model.TimeSeriesConstraintsListDataSelectorsType,
		elements *model.TimeSeriesConstraintsDataElementsType,
	) (*model.MsgCounterType, error)

	// request FunctionTypeTimeSeriesListData from a remote device
	RequestData(
		selector *model.TimeSeriesListDataSelectorsType,
		elements *model.TimeSeriesDataElementsType,
	) (*model.MsgCounterType, error)

	// write Time Series values
	// returns an error if this failed
	WriteData(data []model.TimeSeriesDataType) (*model.MsgCounterType, error)
}
