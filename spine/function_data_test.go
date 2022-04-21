package spine

import (
	"testing"

	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/DerAndereAndi/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func TestFunctionData(t *testing.T) {

	setData := &model.DeviceClassificationManufacturerDataType{
		DeviceName: util.Ptr(model.DeviceClassificationStringType("device name")),
	}
	functionType := model.FunctionTypeDeviceClassificationManufacturerData
	sut := NewFunctionData[model.DeviceClassificationManufacturerDataType](functionType)
	sut.SetData(setData)
	getData := sut.Data()

	assert.Equal(t, setData.DeviceName, getData.DeviceName)
	assert.Equal(t, functionType, sut.Function())
}
