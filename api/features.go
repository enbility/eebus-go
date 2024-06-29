package api

import (
	"time"

	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// Feature client interface were the local feature role is client and the remote feature role is server
type FeatureClientInterface interface {
	// check if there is a subscription to the remote feature
	HasSubscription() bool

	// subscribe to the feature of the entity
	Subscribe() (*model.MsgCounterType, error)

	// check if there is a binding to the remote feature
	HasBinding() bool

	// bind to the feature of the entity
	Bind() (*model.MsgCounterType, error)

	// add a callback function to be invoked once a result or reply message for a msgCounter came in
	AddResponseCallback(msgCounterReference model.MsgCounterType, function func(msg api.ResponseMessage)) error

	// add a callback function to be invoked once a result came in
	AddResultCallback(function func(msg api.ResponseMessage))
}

// Feature server interface were the local feature role is a server
type FeatureServerInterface interface {
}

// Common interface for DeviceClassificationClientInterface and DeviceClassificationServerInterface
type DeviceClassificationCommonInterface interface {
	// get the current manufacturer details for a remote device entity
	GetManufacturerDetails() (*model.DeviceClassificationManufacturerDataType, error)
}

// Common interface for DeviceConfigurationClientInterface and DeviceConfigurationServerInterface
type DeviceConfigurationCommonInterface interface {
	// check if spine.EventPayload Data contains data for a given filter
	//
	// data type will be checked for model.DeviceConfigurationKeyValueListDataType,
	// filter type will be checked for model.DeviceConfigurationKeyValueDescriptionDataType
	CheckEventPayloadDataForFilter(payloadData any, filter any) bool

	// Get the description for a given keyId
	//
	// Will return nil if no matching description was found
	GetKeyValueDescriptionFoKeyId(keyId model.DeviceConfigurationKeyIdType) (*model.DeviceConfigurationKeyValueDescriptionDataType, error)

	// Get the description for a given value combination
	//
	// Returns an error if no matching description was found
	GetKeyValueDescriptionsForFilter(filter model.DeviceConfigurationKeyValueDescriptionDataType) ([]model.DeviceConfigurationKeyValueDescriptionDataType, error)

	// Get the key value data for a given keyId
	//
	// Will return nil if no matching data was found
	GetKeyValueDataForKeyId(keyId model.DeviceConfigurationKeyIdType) (*model.DeviceConfigurationKeyValueDataType, error)

	// Get key value data for a given filter
	//
	// Will return nil if no matching data was found
	GetKeyValueDataForFilter(filter model.DeviceConfigurationKeyValueDescriptionDataType) (*model.DeviceConfigurationKeyValueDataType, error)
}

// Common interface for DeviceDiagnosisClientInterface and DeviceDiagnosisServerInterface
type DeviceDiagnosisCommonInterface interface {
	// get the current diagnosis state for an device entity
	GetState() (*model.DeviceDiagnosisStateDataType, error)

	// check if the currently available heartbeat data is within a time duration
	IsHeartbeatWithinDuration(duration time.Duration) bool
}

// Common interface for ElectricalConnectionClientInterface and ElectricalConnectionServerInterface
type ElectricalConnectionCommonInterface interface {
	// check if spine.EventPayload Data contains data for a given filter
	//
	// data type will be checked for model.ElectricalConnectionPermittedValueSetListDataType,
	// filter type will be checked for model.ElectricalConnectionParameterDescriptionDataType
	CheckEventPayloadDataForFilter(payloadData any, filter any) bool

	// Get the description for a given filter
	//
	// Returns an error if no matching description is found
	GetDescriptionsForFilter(
		filter model.ElectricalConnectionDescriptionDataType,
	) ([]model.ElectricalConnectionDescriptionDataType, error)

	// Get the description for a given parameter description
	//
	// Returns an error if no matching description is found
	GetDescriptionForParameterDescriptionFilter(
		filter model.ElectricalConnectionParameterDescriptionDataType) (
		*model.ElectricalConnectionDescriptionDataType, error)

	// Get the description for a given filter
	//
	// Returns an error if no matching description is found
	GetParameterDescriptionsForFilter(
		filter model.ElectricalConnectionParameterDescriptionDataType,
	) ([]model.ElectricalConnectionParameterDescriptionDataType, error)

	// return permitted values for all Electrical Connections
	GetPermittedValueSetForFilter(filter model.ElectricalConnectionPermittedValueSetDataType) (
		[]model.ElectricalConnectionPermittedValueSetDataType, error)

	// returns minimum, maximum, default/pause limit values
	GetPermittedValueDataForFilter(filter model.ElectricalConnectionPermittedValueSetDataType) (
		float64, float64, float64, error)

	// Get the min, max, default current limits for each phase
	GetPhaseCurrentLimits(measDesc []model.MeasurementDescriptionDataType) (
		resultMin []float64, resultMax []float64, resultDefault []float64, resultErr error)

	// Adjust a value to be within the permitted value range
	AdjustValueToBeWithinPermittedValuesForParameterId(
		value float64, parameterId model.ElectricalConnectionParameterIdType) float64

	// Get the characteristics for a given filter
	//
	// Returns an error if no matching description is found
	GetCharacteristicsForFilter(
		filter model.ElectricalConnectionCharacteristicDataType,
	) ([]model.ElectricalConnectionCharacteristicDataType, error)
}

// Common interface for LoadControlClientInterface and LoadControlServerInterface
type LoadControlCommonInterface interface {
	// check if spine.EventPayload Data contains data for a given filter
	//
	// data type will be checked for model.LoadControlLimitListDataType,
	// filter type will be checked for model.LoadControlLimitDescriptionDataType
	CheckEventPayloadDataForFilter(payloadData any, filter any) bool

	// Get the description for a given limitId
	//
	// Will return nil if no matching description is found
	GetLimitDescriptionForId(limitId model.LoadControlLimitIdType) (
		*model.LoadControlLimitDescriptionDataType, error)

	// Get the description for a given filter
	//
	// Returns an error if no matching description is found
	GetLimitDescriptionsForFilter(
		filter model.LoadControlLimitDescriptionDataType,
	) ([]model.LoadControlLimitDescriptionDataType, error)

	// Get the description for a given limitId
	//
	// Will return nil if no data is available
	GetLimitDataForId(limitId model.LoadControlLimitIdType) (*model.LoadControlLimitDataType, error)

	// Get limit data for a given filter
	//
	// Will return nil if no data is available
	GetLimitDataForFilter(filter model.LoadControlLimitDescriptionDataType) ([]model.LoadControlLimitDataType, error)
}

// Common interface for MeasurementClientInterface and MeasurementServerInterface
type MeasurementCommonInterface interface {
	// check if spine.EventPayload Data contains data for a given filter
	//
	// data type will be checked for model.MeasurementListDataType,
	// filter type will be checked for model.MeasurementDescriptionDataType
	CheckEventPayloadDataForFilter(payloadData any, filter any) bool

	// Get the description for a given id
	//
	// Returns an error if no matching description is found
	GetDescriptionForId(
		measurementId model.MeasurementIdType,
	) (*model.MeasurementDescriptionDataType, error)

	// Get the description for a given filter
	//
	// Returns an error if no matching description is found
	GetDescriptionsForFilter(
		filter model.MeasurementDescriptionDataType,
	) ([]model.MeasurementDescriptionDataType, error)

	// Get the constraints for a given filter
	//
	// Returns an error if no matching constraint is found
	GetConstraintsForFilter(
		filter model.MeasurementConstraintsDataType,
	) ([]model.MeasurementConstraintsDataType, error)

	// Get the measuement data for a given measurementId
	//
	// Will return nil if no data is available
	GetDataForId(measurementId model.MeasurementIdType) (*model.MeasurementDataType, error)

	// Get measuement data for a given filter
	//
	// Will return nil if no data is available
	GetDataForFilter(filter model.MeasurementDescriptionDataType) (
		[]model.MeasurementDataType, error)
}

// Common interface for IdentificationClientInterface and IdentificationServerInterface
type IdentificationCommonInterface interface {
	// check if spine.EventPayload Data contains identification data
	//
	// data type will be checked for model.IdentificationListDataType
	CheckEventPayloadDataForFilter(payloadData any) bool

	// return current values for Identification
	GetDataForFilter(filter model.IdentificationDataType) ([]model.IdentificationDataType, error)
}

// Common interface for IncentiveTableClientInterface and IncentiveTableServerInterface
type IncentiveTableCommonInterface interface {
	// return list of descriptions for a given filter
	GetDescriptionsForFilter(filter model.TariffDescriptionDataType) ([]model.IncentiveTableDescriptionType, error)

	// return list of constraints
	GetConstraints() ([]model.IncentiveTableConstraintsType, error)

	// return current data for Time Series
	GetData() ([]model.IncentiveTableType, error)
}

// Common interface for SmartEnergyManagementPsClientInterface and SmartEnergyManagementPsServerInterface
type SmartEnergyManagementPsCommonInterface interface {
	// return current data for FunctionTypeSmartEnergyManagementPsData
	GetData() (*model.SmartEnergyManagementPsDataType, error)
}

// Common interface for TimeSeriesClientInterface and TimeSeriesServerInterface
type TimeSeriesCommonInterface interface {
	// return list of descriptions for a given filter
	GetDescriptionsForFilter(filter model.TimeSeriesDescriptionDataType) ([]model.TimeSeriesDescriptionDataType, error)

	// return current constraints for Time Series
	GetConstraints() ([]model.TimeSeriesConstraintsDataType, error)

	// return current data for Time Series for a given filter
	GetDataForFilter(filter model.TimeSeriesDescriptionDataType) ([]model.TimeSeriesDataType, error)
}
