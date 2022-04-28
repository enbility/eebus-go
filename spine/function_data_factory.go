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
			createFunctionData[model.MeasurementConstraintsListDataType, F](model.FunctionTypeMeasurementConstraintsListData),
		}
	case model.FeatureTypeTypeDeviceConfiguration:
		return []F{}
	case model.FeatureTypeTypeLoadControl:
		return []F{}
	case model.FeatureTypeTypeIdentification:
		return []F{}
	case model.FeatureTypeTypeElectricalConnection:
		return []F{}
	case model.FeatureTypeTypeTimeSeries:
		return []F{}
	case model.FeatureTypeTypeIncentiveTable:
		return []F{}
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
