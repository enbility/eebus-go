package model_test

import (
	"testing"

	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/stretchr/testify/assert"
)

func TestCmdType_Data(t *testing.T) {
	data := &model.NodeManagementDetailedDiscoveryDataType{
		SpecificationVersionList: &model.NodeManagementSpecificationVersionListType{[]model.SpecificationVersionDataType{model.SpecificationVersionDataType("dummy")}},
	}

	sut := &model.CmdType{
		NodeManagementDetailedDiscoveryData: data,
	}

	// Act
	cmdData, err := sut.Data()
	assert.Nil(t, err)
	assert.NotNil(t, cmdData)
	assert.Equal(t, "NodeManagementDetailedDiscoveryData", cmdData.FieldName)
	assert.Equal(t, model.FunctionTypeNodeManagementDetailedDiscoveryData, *cmdData.Function)
	assert.Equal(t, data, cmdData.Value)
}
