package spine

import (
	"testing"

	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/DerAndereAndi/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func TestFunctionData_UpdateData(t *testing.T) {
	newData := &model.DeviceClassificationManufacturerDataType{
		DeviceName: util.Ptr(model.DeviceClassificationStringType("device name")),
	}
	functionType := model.FunctionTypeDeviceClassificationManufacturerData
	sut := NewFunctionData[model.DeviceClassificationManufacturerDataType](functionType)
	sut.UpdateData(newData, nil, nil)
	getData := sut.Data()

	assert.Equal(t, newData.DeviceName, getData.DeviceName)
	assert.Equal(t, functionType, sut.Function())
}

func TestFunctionData_UpdateDataPartial(t *testing.T) {
	newData := &model.ElectricalConnectionPermittedValueSetListDataType{
		ElectricalConnectionPermittedValueSetData: []model.ElectricalConnectionPermittedValueSetDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(1)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(1)),
				PermittedValueSet: []model.ScaledNumberSetType{
					{
						Range: []model.ScaledNumberRangeType{
							{
								Min: &model.ScaledNumberType{
									Number: util.Ptr(model.NumberType(6)),
									Scale:  util.Ptr(model.ScaleType(0)),
								},
							},
						},
					},
				},
			},
		},
	}
	functionType := model.FunctionTypeElectricalConnectionPermittedValueSetListData
	sut := NewFunctionData[model.ElectricalConnectionPermittedValueSetListDataType](functionType)

	err := sut.UpdateData(newData, &model.FilterType{CmdControl: &model.CmdControlType{Partial: &model.ElementTagType{}}}, nil)
	if assert.Nil(t, err) {
		getData := sut.Data()
		assert.Equal(t, 1, len(getData.ElectricalConnectionPermittedValueSetData))
	}
}

func TestFunctionData_UpdateDataPartial_NotSupported(t *testing.T) {
	newData := &model.HvacOverrunListDataType{
		HvacOverrunData: []model.HvacOverrunDataType{
			{
				OverrunId: util.Ptr(model.HvacOverrunIdType(1)),
			},
		},
	}
	functionType := model.FunctionTypeHvacOverrunListData
	sut := NewFunctionData[model.HvacOverrunListDataType](functionType)

	err := sut.UpdateData(newData, &model.FilterType{CmdControl: &model.CmdControlType{Partial: &model.ElementTagType{}}}, nil)
	assert.NotNil(t, err)
	assert.Equal(t, "partial updates are not supported for type 'HvacOverrunListDataType'", string(*err.Description))
}
