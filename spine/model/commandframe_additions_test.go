package model_test

import (
	"testing"

	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func TestFilterType_Selector_Data(t *testing.T) {
	data := &model.ElectricalConnectionDescriptionListDataSelectorsType{
		ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(1)),
		ScopeType:              util.Ptr(model.ScopeTypeTypeACPower),
	}

	sut := &model.FilterType{
		ElectricalConnectionDescriptionListDataSelectors: data,
	}

	// Act
	cmdData, err := sut.Data()
	assert.Nil(t, err)
	assert.NotNil(t, cmdData)
	assert.Equal(t, model.FunctionTypeElectricalConnectionDescriptionListData, *cmdData.Function)
	assert.Equal(t, data, cmdData.Selector)
}

func TestFilterType_Selector_SetDataForFunction(t *testing.T) {
	cmd := model.FilterType{}
	cmd.SetDataForFunction(model.EEBusTagTypeTypeSelector, model.FunctionTypeElectricalConnectionDescriptionListData, &model.ElectricalConnectionDescriptionListDataSelectorsType{})
	assert.NotNil(t, cmd.ElectricalConnectionDescriptionListDataSelectors)

	cmd = model.FilterType{}
	cmd.SetDataForFunction(model.EEBusTagTypeTypeSelector, model.FunctionTypeElectricalConnectionDescriptionListData, nil)
	assert.Nil(t, cmd.ElectricalConnectionDescriptionListDataSelectors)

	var test *model.ElectricalConnectionDescriptionListDataSelectorsType
	cmd = model.FilterType{}
	cmd.SetDataForFunction(model.EEBusTagTypeTypeSelector, model.FunctionTypeElectricalConnectionDescriptionListData, test)
	assert.NotNil(t, cmd.ElectricalConnectionDescriptionListDataSelectors)
}

func TestFilterType_Elements_Data(t *testing.T) {
	data := &model.ElectricalConnectionDescriptionDataElementsType{
		ElectricalConnectionId: util.Ptr(model.ElementTagType{}),
	}

	sut := &model.FilterType{
		ElectricalConnectionDescriptionDataElements: data,
	}

	// Act
	cmdData, err := sut.Data()
	assert.Nil(t, err)
	assert.NotNil(t, cmdData)
	assert.Equal(t, model.FunctionTypeElectricalConnectionDescriptionListData, *cmdData.Function)
	assert.Equal(t, data, cmdData.Elements)
}

func TestFilterType_Elements_SetDataForFunction(t *testing.T) {
	cmd := model.FilterType{}
	cmd.SetDataForFunction(model.EEbusTagTypeTypeElements, model.FunctionTypeElectricalConnectionDescriptionListData, &model.ElectricalConnectionDescriptionDataElementsType{})
	assert.NotNil(t, cmd.ElectricalConnectionDescriptionDataElements)

	cmd = model.FilterType{}
	cmd.SetDataForFunction(model.EEbusTagTypeTypeElements, model.FunctionTypeElectricalConnectionDescriptionListData, nil)
	assert.Nil(t, cmd.ElectricalConnectionDescriptionDataElements)

	var test *model.ElectricalConnectionDescriptionDataElementsType
	cmd = model.FilterType{}
	cmd.SetDataForFunction(model.EEbusTagTypeTypeElements, model.FunctionTypeElectricalConnectionDescriptionListData, test)
	assert.NotNil(t, cmd.ElectricalConnectionDescriptionDataElements)
}

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

func TestCmdType_SetDataForFunction(t *testing.T) {
	cmd := model.CmdType{}
	cmd.SetDataForFunction(model.FunctionTypeElectricalConnectionDescriptionListData, &model.ElectricalConnectionDescriptionListDataType{})
	assert.NotNil(t, cmd.ElectricalConnectionDescriptionListData)

	cmd = model.CmdType{}
	cmd.SetDataForFunction(model.FunctionTypeElectricalConnectionDescriptionListData, nil)
	assert.Nil(t, cmd.ElectricalConnectionDescriptionListData)

	var test *model.ElectricalConnectionDescriptionListDataType
	cmd = model.CmdType{}
	cmd.SetDataForFunction(model.FunctionTypeElectricalConnectionDescriptionListData, test)
	assert.NotNil(t, cmd.ElectricalConnectionDescriptionListData)
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
