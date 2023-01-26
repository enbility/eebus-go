package spine

import (
	"fmt"

	"github.com/enbility/eebus-go/spine/model"
)

func CreateFunctionData[F any](featureType model.FeatureTypeType) []F {
	if featureType == model.FeatureTypeTypeNodeManagement {
		return []F{} // NodeManagement implementation is not using function data
	}

	// Some devices use generic for everything (e.g. Vaillant Arotherm heatpump)
	// or for some things like the SMA HM 2.0 or Elli Wallbox, which uses Generic feature
	// for Heartbeats, even though that should go into FeatureTypeTypeDeviceDiagnosis
	// Hence we add everything to the Generic feature, as we don't know what might be needed

	var result []F

	if featureType == model.FeatureTypeTypeBill || featureType == model.FeatureTypeTypeGeneric {
		result = append(result, []F{
			createFunctionData[model.BillDescriptionListDataType, F](model.FunctionTypeBillDescriptionListData),
			createFunctionData[model.BillConstraintsListDataType, F](model.FunctionTypeBillConstraintsListData),
			createFunctionData[model.BillListDataType, F](model.FunctionTypeBillListData),
		}...)
	}

	if featureType == model.FeatureTypeTypeDeviceClassification || featureType == model.FeatureTypeTypeGeneric {
		result = append(result, []F{
			createFunctionData[model.DeviceClassificationManufacturerDataType, F](model.FunctionTypeDeviceClassificationManufacturerData),
			createFunctionData[model.DeviceClassificationUserDataType, F](model.FunctionTypeDeviceClassificationUserData),
		}...)
	}

	if featureType == model.FeatureTypeTypeDeviceConfiguration || featureType == model.FeatureTypeTypeGeneric {
		result = append(result, []F{
			createFunctionData[model.DeviceConfigurationKeyValueDescriptionListDataType, F](model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData),
			createFunctionData[model.DeviceConfigurationKeyValueListDataType, F](model.FunctionTypeDeviceConfigurationKeyValueListData),
		}...)
	}

	if featureType == model.FeatureTypeTypeDeviceDiagnosis || featureType == model.FeatureTypeTypeGeneric {
		result = append(result, []F{
			createFunctionData[model.DeviceDiagnosisStateDataType, F](model.FunctionTypeDeviceDiagnosisStateData),
			createFunctionData[model.DeviceDiagnosisHeartbeatDataType, F](model.FunctionTypeDeviceDiagnosisHeartbeatData),
		}...)
	}

	if featureType == model.FeatureTypeTypeElectricalConnection || featureType == model.FeatureTypeTypeGeneric {
		result = append(result, []F{
			createFunctionData[model.ElectricalConnectionDescriptionListDataType, F](model.FunctionTypeElectricalConnectionDescriptionListData),
			createFunctionData[model.ElectricalConnectionParameterDescriptionListDataType, F](model.FunctionTypeElectricalConnectionParameterDescriptionListData),
			createFunctionData[model.ElectricalConnectionPermittedValueSetListDataType, F](model.FunctionTypeElectricalConnectionPermittedValueSetListData),
		}...)
	}

	if featureType == model.FeatureTypeTypeHvac || featureType == model.FeatureTypeTypeGeneric {
		result = append(result, []F{
			createFunctionData[model.HvacOverrunDescriptionListDataType, F](model.FunctionTypeHvacOverrunDescriptionListData),
			createFunctionData[model.HvacOverrunListDataType, F](model.FunctionTypeHvacOverrunListData),
			createFunctionData[model.HvacSystemFunctionDataType, F](model.FunctionTypeHvacSystemFunctionListData),
			createFunctionData[model.HvacSystemFunctionDescriptionDataType, F](model.FunctionTypeHvacSystemFunctionDescriptionListData),
		}...)
	}

	if featureType == model.FeatureTypeTypeIdentification || featureType == model.FeatureTypeTypeGeneric {
		result = append(result, []F{
			createFunctionData[model.IdentificationListDataType, F](model.FunctionTypeIdentificationListData),
		}...)
	}

	if featureType == model.FeatureTypeTypeIncentiveTable || featureType == model.FeatureTypeTypeGeneric {
		result = append(result, []F{
			createFunctionData[model.IncentiveTableDescriptionDataType, F](model.FunctionTypeIncentiveTableDescriptionData),
			createFunctionData[model.IncentiveTableConstraintsDataType, F](model.FunctionTypeIncentiveTableConstraintsData),
			createFunctionData[model.IncentiveTableDataType, F](model.FunctionTypeIncentiveTableData),
		}...)
	}

	if featureType == model.FeatureTypeTypeLoadControl || featureType == model.FeatureTypeTypeGeneric {
		result = append(result, []F{
			createFunctionData[model.LoadControlLimitConstraintsListDataType, F](model.FunctionTypeLoadControlLimitConstraintsListData),
			createFunctionData[model.LoadControlLimitDescriptionListDataType, F](model.FunctionTypeLoadControlLimitDescriptionListData),
			createFunctionData[model.LoadControlLimitListDataType, F](model.FunctionTypeLoadControlLimitListData),
		}...)
	}

	if featureType == model.FeatureTypeTypeMeasurement || featureType == model.FeatureTypeTypeGeneric {
		result = append(result, []F{
			createFunctionData[model.MeasurementDataType, F](model.FunctionTypeMeasurementListData),
			createFunctionData[model.MeasurementDescriptionDataType, F](model.FunctionTypeMeasurementDescriptionListData),
			createFunctionData[model.MeasurementDescriptionListDataType, F](model.FunctionTypeMeasurementDescriptionListData),
			createFunctionData[model.MeasurementConstraintsListDataType, F](model.FunctionTypeMeasurementConstraintsListData),
			createFunctionData[model.MeasurementListDataType, F](model.FunctionTypeMeasurementListData),
		}...)
	}

	if featureType == model.FeatureTypeTypeTimeSeries || featureType == model.FeatureTypeTypeGeneric {
		result = append(result, []F{
			createFunctionData[model.TimeSeriesDescriptionListDataType, F](model.FunctionTypeTimeSeriesDescriptionListData),
			createFunctionData[model.TimeSeriesConstraintsListDataType, F](model.FunctionTypeTimeSeriesConstraintsListData),
			createFunctionData[model.TimeSeriesListDataType, F](model.FunctionTypeTimeSeriesListData),
		}...)
	}

	if len(result) == 0 {
		panic(fmt.Errorf("unknown featureType '%s'", featureType))
	}

	return result
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
