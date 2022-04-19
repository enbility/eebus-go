package spine

import (
	"fmt"

	"github.com/DerAndereAndi/eebus-go/spine/model"
)

func CreateFunctionData[F any](featureType model.FeatureTypeEnumType) []F {
	switch featureType {
	case model.FeatureTypeEnumTypeDeviceClassification:
		return []F{
			createFunctionData[model.DeviceClassificationManufacturerDataType, F](model.FunctionEnumTypeDeviceClassificationManufacturerData),
			createFunctionData[model.DeviceClassificationUserDataType, F](model.FunctionEnumTypeDeviceClassificationUserData),
		}
		// TODO: Add more feature types
	}

	panic(fmt.Errorf("unknown featureType '%s'", featureType))
}

func createFunctionData[T any, F any](functionType model.FunctionEnumType) F {
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
