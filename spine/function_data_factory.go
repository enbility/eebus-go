package spine

import (
	"fmt"

	"github.com/DerAndereAndi/eebus-go/spine/model"
)

func CreateFunctionData[F any](featureType model.FeatureTypeType) []F {
	switch featureType {
	case model.FeatureTypeTypeNodeManagement:
		return []F{} // NodeManagement implementation is not using function data
	case model.FeatureTypeTypeDeviceClassification:
		return []F{
			createFunctionData[model.DeviceClassificationManufacturerDataType, F](model.FunctionTypeDeviceClassificationManufacturerData),
			createFunctionData[model.DeviceClassificationUserDataType, F](model.FunctionTypeDeviceClassificationUserData),
		}
	case model.FeatureTypeTypeDeviceDiagnosis:
		return []F{
			createFunctionData[model.DeviceDiagnosisStateDataType, F](model.FunctionTypeDeviceDiagnosisStateData),
			createFunctionData[model.DeviceDiagnosisHeartbeatDataType, F](model.FunctionTypeDeviceDiagnosisHeartbeatData),
		}
	case model.FeatureTypeTypeMeasurement:
		return []F{
			createFunctionData[model.MeasurementDataType, F](model.FunctionTypeMeasurementListData),
			createFunctionData[model.MeasurementDescriptionDataType, F](model.FunctionTypeMeasurementDescriptionListData),
			createFunctionData[model.MeasurementDescriptionListDataType, F](model.FunctionTypeMeasurementDescriptionListData),
			createFunctionData[model.MeasurementConstraintsListDataType, F](model.FunctionTypeMeasurementConstraintsListData),
			createFunctionData[model.MeasurementListDataType, F](model.FunctionTypeMeasurementListData),
		}
	case model.FeatureTypeTypeDeviceConfiguration:
		return []F{
			createFunctionData[model.DeviceConfigurationKeyValueDescriptionListDataType, F](model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData),
			createFunctionData[model.DeviceConfigurationKeyValueListDataType, F](model.FunctionTypeDeviceConfigurationKeyValueListData),
		}
	case model.FeatureTypeTypeLoadControl:
		return []F{
			createFunctionData[model.LoadControlLimitConstraintsListDataType, F](model.FunctionTypeLoadControlLimitConstraintsListData),
			createFunctionData[model.LoadControlLimitDescriptionListDataType, F](model.FunctionTypeLoadControlLimitDescriptionListData),
			createFunctionData[model.LoadControlLimitListDataType, F](model.FunctionTypeLoadControlLimitListData),
		}
	case model.FeatureTypeTypeIdentification:
		return []F{
			createFunctionData[model.IdentificationListDataType, F](model.FunctionTypeIdentificationListData),
		}
	case model.FeatureTypeTypeElectricalConnection:
		return []F{
			createFunctionData[model.ElectricalConnectionDescriptionListDataType, F](model.FunctionTypeElectricalConnectionDescriptionListData),
			createFunctionData[model.ElectricalConnectionParameterDescriptionListDataType, F](model.FunctionTypeElectricalConnectionParameterDescriptionListData),
			createFunctionData[model.ElectricalConnectionPermittedValueSetListDataType, F](model.FunctionTypeElectricalConnectionPermittedValueSetListData),
		}
	case model.FeatureTypeTypeTimeSeries:
		return []F{
			createFunctionData[model.TimeSeriesDescriptionListDataType, F](model.FunctionTypeTimeSeriesDescriptionListData),
			createFunctionData[model.TimeSeriesConstraintsListDataType, F](model.FunctionTypeTimeSeriesConstraintsListData),
			createFunctionData[model.TimeSeriesListDataType, F](model.FunctionTypeTimeSeriesListData),
		}
	case model.FeatureTypeTypeIncentiveTable:
		return []F{
			createFunctionData[model.IncentiveTableDescriptionDataType, F](model.FunctionTypeIncentiveTableDescriptionData),
			createFunctionData[model.IncentiveTableConstraintsDataType, F](model.FunctionTypeIncentiveTableConstraintsData),
			createFunctionData[model.IncentiveTableDataType, F](model.FunctionTypeIncentiveTableData),
		}
	case model.FeatureTypeTypeBill:
		return []F{
			createFunctionData[model.BillDescriptionListDataType, F](model.FunctionTypeBillDescriptionListData),
			createFunctionData[model.BillConstraintsListDataType, F](model.FunctionTypeBillConstraintsListData),
			createFunctionData[model.BillListDataType, F](model.FunctionTypeBillListData),
		}
	case model.FeatureTypeTypeGeneric:
		// TODO: add the proper function data, this is reported e.g. by the SMA HM 2.0
		return []F{
			createFunctionData[model.DeviceDiagnosisHeartbeatDataType, F](model.FunctionTypeDeviceDiagnosisHeartbeatData), // Elli Charger uses this feature type for heartbeats
		}
		// TODO: Add more feature types
		// default:
		// 	return []F{}
	}

	panic(fmt.Errorf("unknown featureType '%s'", featureType))
}

func createFunctionData[T any, F any](functionType model.FunctionType) F {
	x := any(new(F))
	switch x.(type) {
	case *FunctionDataCmd:
		return any(NewFunctionDataCmd[T](functionType)).(F)
	case *FunctionData:
		return any(NewFunctionData[T](functionType)).(F)
	default:
		panic(fmt.Errorf("only FunctionData and FunctionDataCmd are supported"))
	}
}
