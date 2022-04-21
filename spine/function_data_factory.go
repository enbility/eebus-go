package spine

import (
	"fmt"

	"github.com/DerAndereAndi/eebus-go/spine/model"
)

func CreateFunctionData[F any](featureType model.FeatureTypeType) []F {
	switch featureType {
	case model.FeatureTypeTypeDeviceClassification:
		return []F{
			createFunctionData[model.DeviceClassificationManufacturerDataType, F](model.FunctionTypeDeviceClassificationManufacturerData),
			createFunctionData[model.DeviceClassificationUserDataType, F](model.FunctionTypeDeviceClassificationUserData),
		}
		// TODO: Add more feature types
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
