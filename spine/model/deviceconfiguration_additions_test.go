package model_test

import (
	"testing"

	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func TestDeviceConfigurationKeyValueListDataType_Update(t *testing.T) {
	sut := model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
			{
				KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				Value: &model.DeviceConfigurationKeyValueValueType{
					Boolean: util.Ptr(true),
				},
			},
			{
				KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(1)),
				Value: &model.DeviceConfigurationKeyValueValueType{
					Boolean: util.Ptr(true),
				},
			},
		},
	}

	newData := model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
			{
				KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(1)),
				Value: &model.DeviceConfigurationKeyValueValueType{
					Boolean: util.Ptr(false),
				},
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

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
	sut := model.DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []model.DeviceConfigurationKeyValueDescriptionDataType{
			{
				KeyId:     util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				ValueType: util.Ptr(model.DeviceConfigurationKeyValueTypeTypeBoolean),
			},
			{
				KeyId:     util.Ptr(model.DeviceConfigurationKeyIdType(1)),
				ValueType: util.Ptr(model.DeviceConfigurationKeyValueTypeTypeBoolean),
			},
		},
	}

	newData := model.DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []model.DeviceConfigurationKeyValueDescriptionDataType{
			{
				KeyId:     util.Ptr(model.DeviceConfigurationKeyIdType(1)),
				ValueType: util.Ptr(model.DeviceConfigurationKeyValueTypeTypeString),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.DeviceConfigurationKeyValueDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.KeyId))
	assert.Equal(t, model.DeviceConfigurationKeyValueTypeTypeBoolean, *item1.ValueType)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.KeyId))
	assert.Equal(t, model.DeviceConfigurationKeyValueTypeTypeString, *item2.ValueType)
}

func TestDeviceConfigurationKeyValueConstraintsListDataType_Update(t *testing.T) {
	sut := model.DeviceConfigurationKeyValueConstraintsListDataType{
		DeviceConfigurationKeyValueConstraintsData: []model.DeviceConfigurationKeyValueConstraintsDataType{
			{
				KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				ValueStepSize: &model.DeviceConfigurationKeyValueValueType{
					Boolean: util.Ptr(true),
				},
			},
			{
				KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(1)),
				ValueStepSize: &model.DeviceConfigurationKeyValueValueType{
					Boolean: util.Ptr(true),
				},
			},
		},
	}

	newData := model.DeviceConfigurationKeyValueConstraintsListDataType{
		DeviceConfigurationKeyValueConstraintsData: []model.DeviceConfigurationKeyValueConstraintsDataType{
			{
				KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(1)),
				ValueStepSize: &model.DeviceConfigurationKeyValueValueType{
					Boolean: util.Ptr(false),
				},
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

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
