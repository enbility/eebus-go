package spine

import (
	"testing"

	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/stretchr/testify/assert"
)

func TestFunctionDataFactory_FunctionData(t *testing.T) {
	result := CreateFunctionData[FunctionData](model.FeatureTypeEnumTypeDeviceClassification)
	assert.Equal(t, 2, len(result))
	assert.IsType(t, &FunctionDataImpl[model.DeviceClassificationManufacturerDataType]{}, result[0])
	assert.IsType(t, &FunctionDataImpl[model.DeviceClassificationUserDataType]{}, result[1])
}

func TestFunctionDataFactory_FunctionDataCmd(t *testing.T) {
	result := CreateFunctionData[FunctionDataCmd](model.FeatureTypeEnumTypeDeviceClassification)
	assert.Equal(t, 2, len(result))
	assert.IsType(t, &FunctionDataCmdImpl[model.DeviceClassificationManufacturerDataType]{}, result[0])
	assert.IsType(t, &FunctionDataCmdImpl[model.DeviceClassificationUserDataType]{}, result[1])
}

func TestFunctionDataFactory_unknownFeatureType(t *testing.T) {
	assert.PanicsWithError(t, "unknown featureType 'Alarm'",
		func() { CreateFunctionData[FunctionDataCmd](model.FeatureTypeEnumTypeAlarm) })
}

func TestFunctionDataFactory_unknownFunctionDataType(t *testing.T) {
	assert.PanicsWithError(t, "only FunctionData and FunctionDataCmd are supported",
		func() { CreateFunctionData[int](model.FeatureTypeEnumTypeDeviceClassification) })
}
