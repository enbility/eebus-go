package model

import (
	"testing"

	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func TestFilterType_Selector_Data(t *testing.T) {
	data := &ElectricalConnectionDescriptionListDataSelectorsType{
		ElectricalConnectionId: util.Ptr(ElectricalConnectionIdType(1)),
		ScopeType:              util.Ptr(ScopeTypeTypeACPower),
	}

	sut := &FilterType{
		ElectricalConnectionDescriptionListDataSelectors: data,
	}

	// Act
	cmdData, err := sut.Data()
	assert.Nil(t, err)
	assert.NotNil(t, cmdData)
	assert.Equal(t, FunctionTypeElectricalConnectionDescriptionListData, *cmdData.Function)
	assert.Equal(t, data, cmdData.Selector)
}

func TestFilterType_Selector_SetDataForFunction(t *testing.T) {
	cmd := FilterType{}
	cmd.SetDataForFunction(EEBusTagTypeTypeSelector, FunctionTypeElectricalConnectionDescriptionListData, &ElectricalConnectionDescriptionListDataSelectorsType{})
	assert.NotNil(t, cmd.ElectricalConnectionDescriptionListDataSelectors)

	cmd = FilterType{}
	cmd.SetDataForFunction(EEBusTagTypeTypeSelector, FunctionTypeElectricalConnectionDescriptionListData, nil)
	assert.Nil(t, cmd.ElectricalConnectionDescriptionListDataSelectors)

	var test *ElectricalConnectionDescriptionListDataSelectorsType
	cmd = FilterType{}
	cmd.SetDataForFunction(EEBusTagTypeTypeSelector, FunctionTypeElectricalConnectionDescriptionListData, test)
	assert.NotNil(t, cmd.ElectricalConnectionDescriptionListDataSelectors)
}

func TestFilterType_Elements_Data(t *testing.T) {
	data := &ElectricalConnectionDescriptionDataElementsType{
		ElectricalConnectionId: util.Ptr(ElementTagType{}),
	}

	sut := &FilterType{
		ElectricalConnectionDescriptionDataElements: data,
	}

	// Act
	cmdData, err := sut.Data()
	assert.Nil(t, err)
	assert.NotNil(t, cmdData)
	assert.Equal(t, FunctionTypeElectricalConnectionDescriptionListData, *cmdData.Function)
	assert.Equal(t, data, cmdData.Elements)
}

func TestFilterType_Elements_SetDataForFunction(t *testing.T) {
	cmd := FilterType{}
	cmd.SetDataForFunction(EEbusTagTypeTypeElements, FunctionTypeElectricalConnectionDescriptionListData, &ElectricalConnectionDescriptionDataElementsType{})
	assert.NotNil(t, cmd.ElectricalConnectionDescriptionDataElements)

	cmd = FilterType{}
	cmd.SetDataForFunction(EEbusTagTypeTypeElements, FunctionTypeElectricalConnectionDescriptionListData, nil)
	assert.Nil(t, cmd.ElectricalConnectionDescriptionDataElements)

	var test *ElectricalConnectionDescriptionDataElementsType
	cmd = FilterType{}
	cmd.SetDataForFunction(EEbusTagTypeTypeElements, FunctionTypeElectricalConnectionDescriptionListData, test)
	assert.NotNil(t, cmd.ElectricalConnectionDescriptionDataElements)
}

func TestCmdType_Data(t *testing.T) {
	data := &NodeManagementDetailedDiscoveryDataType{
		SpecificationVersionList: &NodeManagementSpecificationVersionListType{[]SpecificationVersionDataType{SpecificationVersionDataType("dummy")}},
	}

	sut := &CmdType{
		NodeManagementDetailedDiscoveryData: data,
	}

	// Act
	cmdData, err := sut.Data()
	assert.Nil(t, err)
	assert.NotNil(t, cmdData)
	assert.Equal(t, "NodeManagementDetailedDiscoveryData", cmdData.FieldName)
	assert.Equal(t, FunctionTypeNodeManagementDetailedDiscoveryData, *cmdData.Function)
	assert.Equal(t, data, cmdData.Value)
}

func TestCmdType_SetDataForFunction(t *testing.T) {
	cmd := CmdType{}
	cmd.SetDataForFunction(FunctionTypeElectricalConnectionDescriptionListData, &ElectricalConnectionDescriptionListDataType{})
	assert.NotNil(t, cmd.ElectricalConnectionDescriptionListData)

	cmd = CmdType{}
	cmd.SetDataForFunction(FunctionTypeElectricalConnectionDescriptionListData, nil)
	assert.Nil(t, cmd.ElectricalConnectionDescriptionListData)

	var test *ElectricalConnectionDescriptionListDataType
	cmd = CmdType{}
	cmd.SetDataForFunction(FunctionTypeElectricalConnectionDescriptionListData, test)
	assert.NotNil(t, cmd.ElectricalConnectionDescriptionListData)
}

func TestCmdType_ExtractFilter_NoFilter(t *testing.T) {
	sut := &CmdType{
		NodeManagementDetailedDiscoveryData: &NodeManagementDetailedDiscoveryDataType{},
	}

	// Act
	filterPartial, filterDelete := sut.ExtractFilter()
	assert.Nil(t, filterPartial)
	assert.Nil(t, filterDelete)
}

func TestCmdType_ExtractFilter_FilterPartialDelete(t *testing.T) {

	filterP := FilterType{
		CmdControl: &CmdControlType{Partial: &ElementTagType{}},
		NodeManagementDetailedDiscoveryDataSelectors: &NodeManagementDetailedDiscoveryDataSelectorsType{},
	}
	filterD := FilterType{
		CmdControl: &CmdControlType{Delete: &ElementTagType{}},
		NodeManagementDetailedDiscoveryDataSelectors: &NodeManagementDetailedDiscoveryDataSelectorsType{},
	}

	sut := &CmdType{
		Filter:                              []FilterType{filterD, filterP},
		NodeManagementDetailedDiscoveryData: &NodeManagementDetailedDiscoveryDataType{},
	}

	// Act
	filterPartial, filterDelete := sut.ExtractFilter()
	assert.NotNil(t, filterPartial)
	assert.Equal(t, &filterP, filterPartial)
	assert.NotNil(t, filterDelete)
	assert.Equal(t, &filterD, filterDelete)
}
