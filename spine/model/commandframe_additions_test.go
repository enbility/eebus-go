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

func TestCmdType_ExtractFilter_NoFilter(t *testing.T) {
	sut := &model.CmdType{
		NodeManagementDetailedDiscoveryData: &model.NodeManagementDetailedDiscoveryDataType{},
	}

	// Act
	filterPartial, filterDelete := sut.ExtractFilter()
	assert.Nil(t, filterPartial)
	assert.Nil(t, filterDelete)
}

func TestCmdType_ExtractFilter_FilterPartialDelete(t *testing.T) {

	filterP := model.FilterType{
		CmdControl: &model.CmdControlType{Partial: &model.ElementTagType{}},
		NodeManagementDetailedDiscoveryDataSelectors: &model.NodeManagementDetailedDiscoveryDataSelectorsType{},
	}
	filterD := model.FilterType{
		CmdControl: &model.CmdControlType{Delete: &model.ElementTagType{}},
		NodeManagementDetailedDiscoveryDataSelectors: &model.NodeManagementDetailedDiscoveryDataSelectorsType{},
	}

	sut := &model.CmdType{
		Filter:                              []model.FilterType{filterD, filterP},
		NodeManagementDetailedDiscoveryData: &model.NodeManagementDetailedDiscoveryDataType{},
	}

	// Act
	filterPartial, filterDelete := sut.ExtractFilter()
	assert.NotNil(t, filterPartial)
	assert.Equal(t, &filterP, filterPartial)
	assert.NotNil(t, filterDelete)
	assert.Equal(t, &filterD, filterDelete)
}
