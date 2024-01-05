package model

import (
	"testing"

	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func TestDeviceConfigurationKeyValueListDataType_Update(t *testing.T) {
	sut := DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []DeviceConfigurationKeyValueDataType{
			{
				KeyId: util.Ptr(DeviceConfigurationKeyIdType(0)),
				Value: &DeviceConfigurationKeyValueValueType{
					Boolean: util.Ptr(true),
				},
			},
			{
				KeyId: util.Ptr(DeviceConfigurationKeyIdType(1)),
				Value: &DeviceConfigurationKeyValueValueType{
					Boolean: util.Ptr(true),
				},
			},
		},
	}

	newData := DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []DeviceConfigurationKeyValueDataType{
			{
				KeyId: util.Ptr(DeviceConfigurationKeyIdType(1)),
				Value: &DeviceConfigurationKeyValueValueType{
					Boolean: util.Ptr(false),
				},
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

	data := sut.DeviceConfigurationKeyValueData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.KeyId))
	assert.Equal(t, true, *item1.Value.Boolean)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.KeyId))
	assert.Equal(t, false, *item2.Value.Boolean)
}

func TestDeviceConfigurationKeyValueDescriptionListDataType_Update(t *testing.T) {
	sut := DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []DeviceConfigurationKeyValueDescriptionDataType{
			{
				KeyId:     util.Ptr(DeviceConfigurationKeyIdType(0)),
				ValueType: util.Ptr(DeviceConfigurationKeyValueTypeTypeBoolean),
			},
			{
				KeyId:     util.Ptr(DeviceConfigurationKeyIdType(1)),
				ValueType: util.Ptr(DeviceConfigurationKeyValueTypeTypeBoolean),
			},
		},
	}

	newData := DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []DeviceConfigurationKeyValueDescriptionDataType{
			{
				KeyId:     util.Ptr(DeviceConfigurationKeyIdType(1)),
				ValueType: util.Ptr(DeviceConfigurationKeyValueTypeTypeString),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

	data := sut.DeviceConfigurationKeyValueDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.KeyId))
	assert.Equal(t, DeviceConfigurationKeyValueTypeTypeBoolean, *item1.ValueType)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.KeyId))
	assert.Equal(t, DeviceConfigurationKeyValueTypeTypeString, *item2.ValueType)
}

func TestDeviceConfigurationKeyValueConstraintsListDataType_Update(t *testing.T) {
	sut := DeviceConfigurationKeyValueConstraintsListDataType{
		DeviceConfigurationKeyValueConstraintsData: []DeviceConfigurationKeyValueConstraintsDataType{
			{
				KeyId: util.Ptr(DeviceConfigurationKeyIdType(0)),
				ValueStepSize: &DeviceConfigurationKeyValueValueType{
					Boolean: util.Ptr(true),
				},
			},
			{
				KeyId: util.Ptr(DeviceConfigurationKeyIdType(1)),
				ValueStepSize: &DeviceConfigurationKeyValueValueType{
					Boolean: util.Ptr(true),
				},
			},
		},
	}

	newData := DeviceConfigurationKeyValueConstraintsListDataType{
		DeviceConfigurationKeyValueConstraintsData: []DeviceConfigurationKeyValueConstraintsDataType{
			{
				KeyId: util.Ptr(DeviceConfigurationKeyIdType(1)),
				ValueStepSize: &DeviceConfigurationKeyValueValueType{
					Boolean: util.Ptr(false),
				},
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

	data := sut.DeviceConfigurationKeyValueConstraintsData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.KeyId))
	assert.Equal(t, true, *item1.ValueStepSize.Boolean)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.KeyId))
	assert.Equal(t, false, *item2.ValueStepSize.Boolean)
}
