package spine

import (
	"testing"

	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestFunctionDataCmdSuite(t *testing.T) {
	suite.Run(t, new(FctDataCmdSuite))
}

type FctDataCmdSuite struct {
	suite.Suite
	function model.FunctionType
	data     *model.DeviceClassificationManufacturerDataType
	sut      *FunctionDataCmdImpl[model.DeviceClassificationManufacturerDataType]
}

func (suite *FctDataCmdSuite) SetupSuite() {
	suite.function = model.FunctionTypeDeviceClassificationManufacturerData
	suite.data = &model.DeviceClassificationManufacturerDataType{
		DeviceName: util.Ptr(model.DeviceClassificationStringType("device name")),
	}
	suite.sut = NewFunctionDataCmd[model.DeviceClassificationManufacturerDataType](suite.function)
	suite.sut.UpdateData(suite.data, nil, nil)
}

func (suite *FctDataCmdSuite) TestFunctionDataCmd_ReadCmd() {
	readCmd := suite.sut.ReadCmdType(nil, nil)
	assert.NotNil(suite.T(), readCmd.DeviceClassificationManufacturerData)
	assert.Nil(suite.T(), readCmd.DeviceClassificationManufacturerData.DeviceName)

	partialS := model.NewFilterTypePartial()
	readCmd = suite.sut.ReadCmdType(partialS, nil)
	assert.NotNil(suite.T(), readCmd.DeviceClassificationManufacturerData)
	assert.Nil(suite.T(), readCmd.DeviceClassificationManufacturerData.DeviceName)
}

func (suite *FctDataCmdSuite) TestFunctionDataCmd_ReplyCmd() {
	readCmd := suite.sut.ReplyCmdType(false)
	assert.NotNil(suite.T(), readCmd.DeviceClassificationManufacturerData)
	assert.Equal(suite.T(), suite.data.DeviceName, readCmd.DeviceClassificationManufacturerData.DeviceName)

	readCmd = suite.sut.ReplyCmdType(true)
	assert.NotNil(suite.T(), readCmd.DeviceClassificationManufacturerData)
	assert.Equal(suite.T(), suite.data.DeviceName, readCmd.DeviceClassificationManufacturerData.DeviceName)
}

func (suite *FctDataCmdSuite) TestFunctionDataCmd_NotifyCmd() {
	readCmd := suite.sut.NotifyCmdType(nil, nil, false, nil)
	assert.NotNil(suite.T(), readCmd.DeviceClassificationManufacturerData)
	assert.Equal(suite.T(), suite.data.DeviceName, readCmd.DeviceClassificationManufacturerData.DeviceName)

	readCmd = suite.sut.NotifyCmdType(nil, nil, true, nil)
	assert.NotNil(suite.T(), readCmd.DeviceClassificationManufacturerData)
	assert.Equal(suite.T(), suite.data.DeviceName, readCmd.DeviceClassificationManufacturerData.DeviceName)

	deleteS := model.NewFilterTypePartial()
	readCmd = suite.sut.NotifyCmdType(deleteS, nil, false, nil)
	assert.NotNil(suite.T(), readCmd.DeviceClassificationManufacturerData)
	assert.Equal(suite.T(), suite.data.DeviceName, readCmd.DeviceClassificationManufacturerData.DeviceName)
}

func (suite *FctDataCmdSuite) TestFunctionDataCmd_WriteCmd() {
	readCmd := suite.sut.WriteCmdType(nil, nil, nil)
	assert.NotNil(suite.T(), readCmd.DeviceClassificationManufacturerData)
	assert.Equal(suite.T(), suite.data.DeviceName, readCmd.DeviceClassificationManufacturerData.DeviceName)

	partialS := model.NewFilterTypePartial()
	readCmd = suite.sut.WriteCmdType(nil, partialS, nil)
	assert.NotNil(suite.T(), readCmd.DeviceClassificationManufacturerData)
	assert.Equal(suite.T(), suite.data.DeviceName, readCmd.DeviceClassificationManufacturerData.DeviceName)
}

func (suite *FctDataCmdSuite) Test_AddSelectorToFilter() {
	filter := model.FilterType{CmdControl: &model.CmdControlType{Partial: &model.ElementTagType{}}}

	result := addSelectorToFilter(filter, model.FunctionTypeAlarmListData, &model.AlarmListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeBillConstraintsListData, &model.BillConstraintsListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeBillDescriptionListData, &model.BillDescriptionListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeBillListData, &model.BillListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeBindingManagementEntryListData, &model.BindingManagementEntryListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeCommodityListData, &model.CommodityListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeDeviceConfigurationKeyValueConstraintsListData, &model.DeviceConfigurationKeyValueConstraintsListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, &model.DeviceConfigurationKeyValueDescriptionListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeDeviceConfigurationKeyValueListData, &model.DeviceConfigurationKeyValueListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeDirectControlActivityListData, &model.DirectControlActivityListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeElectricalConnectionDescriptionListData, &model.ElectricalConnectionDescriptionListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeElectricalConnectionParameterDescriptionListData, &model.ElectricalConnectionParameterDescriptionListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeElectricalConnectionPermittedValueSetListData, &model.ElectricalConnectionPermittedValueSetListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeElectricalConnectionStateListData, &model.ElectricalConnectionStateListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeHvacOperationModeDescriptionListData, &model.HvacOperationModeDescriptionListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeHvacOverrunDescriptionListData, &model.HvacOverrunDescriptionListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeHvacOverrunListData, &model.HvacOverrunListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeHvacSystemFunctionDescriptionListData, &model.HvacSystemFunctionDescriptionListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeHvacSystemFunctionListData, &model.HvacSystemFunctionListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeHvacSystemFunctionOperationModeRelationListData, &model.HvacSystemFunctionOperationModeRelationListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeHvacSystemFunctionPowerSequenceRelationListData, &model.HvacSystemFunctionPowerSequenceRelationListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeHvacSystemFunctionSetPointRelationListData, &model.HvacSystemFunctionSetpointRelationListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeIdentificationListData, &model.IdentificationListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeIncentiveDescriptionListData, &model.IncentiveDescriptionListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeIncentiveListData, &model.IncentiveListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeIncentiveTableConstraintsData, &model.IncentiveTableConstraintsDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeIncentiveTableData, &model.IncentiveTableDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeIncentiveTableDescriptionData, &model.IncentiveTableDescriptionDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeLoadControlEventListData, &model.LoadControlEventListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeLoadControlLimitConstraintsListData, &model.LoadControlLimitConstraintsListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeLoadControlLimitDescriptionListData, &model.LoadControlLimitDescriptionListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeLoadControlLimitListData, &model.LoadControlLimitListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeLoadControlStateListData, &model.LoadControlStateListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeMeasurementConstraintsListData, &model.MeasurementConstraintsListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeMeasurementDescriptionListData, &model.MeasurementDescriptionListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeMeasurementListData, &model.MeasurementListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeMeasurementThresholdRelationListData, &model.MeasurementThresholdRelationListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeMessagingListData, &model.MessagingListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeNetworkManagementDeviceDescriptionListData, &model.NetworkManagementDeviceDescriptionListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeNetworkManagementEntityDescriptionListData, &model.NetworkManagementEntityDescriptionListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeNetworkManagementFeatureDescriptionListData, &model.NetworkManagementFeatureDescriptionListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeNodeManagementBindingData, &model.NodeManagementBindingDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeNodeManagementDestinationListData, &model.NodeManagementDestinationListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeNodeManagementDetailedDiscoveryData, &model.NodeManagementDetailedDiscoveryDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeNodeManagementSubscriptionData, &model.NodeManagementSubscriptionDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeNodeManagementUseCaseData, &model.NodeManagementUseCaseDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeOperatingConstraintsDurationListData, &model.OperatingConstraintsDurationListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeOperatingConstraintsInterruptListData, &model.OperatingConstraintsInterruptListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeOperatingConstraintsPowerDescriptionListData, &model.OperatingConstraintsPowerDescriptionListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeOperatingConstraintsPowerLevelListData, &model.OperatingConstraintsPowerLevelListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeOperatingConstraintsPowerRangeListData, &model.OperatingConstraintsPowerRangeListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeOperatingConstraintsResumeImplicationListData, &model.OperatingConstraintsResumeImplicationListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypePowerSequenceAlternativesRelationListData, &model.PowerSequenceAlternativesRelationListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypePowerSequenceDescriptionListData, &model.PowerSequenceDescriptionListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypePowerSequencePriceListData, &model.PowerSequencePriceListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypePowerSequenceScheduleConstraintsListData, &model.PowerSequenceScheduleConstraintsListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypePowerSequenceScheduleListData, &model.PowerSequenceScheduleListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypePowerSequenceSchedulePreferenceListData, &model.PowerSequenceSchedulePreferenceListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypePowerSequenceStateListData, &model.PowerSequenceStateListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypePowerTimeSlotScheduleConstraintsListData, &model.PowerTimeSlotScheduleConstraintsListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypePowerTimeSlotScheduleListData, &model.PowerTimeSlotScheduleListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypePowerTimeSlotValueListData, &model.PowerTimeSlotValueListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeSensingListData, &model.SensingListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeSetpointConstraintsListData, &model.SetpointConstraintsListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeSetpointDescriptionListData, &model.SetpointDescriptionListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeSetpointListData, &model.SetpointListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeSmartEnergyManagementPsData, &model.SmartEnergyManagementPsDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeSmartEnergyManagementPsPriceData, &model.SmartEnergyManagementPsPriceDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeSpecificationVersionListData, &model.SpecificationVersionListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeSupplyConditionListData, &model.SupplyConditionListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeSupplyConditionThresholdRelationListData, &model.SupplyConditionThresholdRelationListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeTariffBoundaryRelationListData, &model.TariffBoundaryRelationListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeTariffDescriptionListData, &model.TariffDescriptionListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeTariffListData, &model.TariffListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeTariffTierRelationListData, &model.TariffTierRelationListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeTaskManagementJobDescriptionListData, &model.TaskManagementJobDescriptionListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeTaskManagementJobListData, &model.TaskManagementJobListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeTaskManagementJobRelationListData, &model.TaskManagementJobRelationListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeThresholdConstraintsListData, &model.ThresholdConstraintsListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeThresholdDescriptionListData, &model.ThresholdDescriptionListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeThresholdListData, &model.ThresholdListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeTierBoundaryDescriptionListData, &model.TierBoundaryDescriptionListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeTierBoundaryListData, &model.TierBoundaryListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeTierDescriptionListData, &model.TierDescriptionListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeTierIncentiveRelationListData, &model.TierIncentiveRelationListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeTierListData, &model.TierListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeTimeSeriesConstraintsListData, &model.TimeSeriesConstraintsListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeTimeSeriesDescriptionListData, &model.TimeSeriesDescriptionListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeTimeSeriesListData, &model.TimeSeriesListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeTimeTableConstraintsListData, &model.TimeTableConstraintsListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeTimeTableDescriptionListData, &model.TimeTableDescriptionListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeTimeTableListData, &model.TimeTableListDataSelectorsType{})
	assert.NotNil(suite.T(), result)

	result = addSelectorToFilter(filter, model.FunctionTypeUseCaseInformationListData, &model.UseCaseInformationListDataSelectorsType{})
	assert.NotNil(suite.T(), result)
}

func (suite *FctDataCmdSuite) Test_AddElementsToFilter() {
	filter := model.FilterType{CmdControl: &model.CmdControlType{Partial: &model.ElementTagType{}}}

	result := addElementToFilter(filter, model.FunctionTypeActuatorLevelData, &model.ActuatorLevelDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeActuatorLevelDescriptionData, &model.ActuatorLevelDescriptionDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeActuatorSwitchData, &model.ActuatorSwitchDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeActuatorSwitchDescriptionData, &model.ActuatorSwitchDescriptionDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeAlarmListData, &model.AlarmDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeBillConstraintsListData, &model.BillConstraintsDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeBillDescriptionListData, &model.BillDescriptionDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeBillListData, &model.BillDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeBindingManagementDeleteCall, &model.BindingManagementDeleteCallElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeBindingManagementEntryListData, &model.BindingManagementEntryDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeBindingManagementRequestCall, &model.BindingManagementRequestCallElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeCommodityListData, &model.CommodityDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeDataTunnelingCall, &model.DataTunnelingCallElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeDeviceClassificationManufacturerData, &model.DeviceClassificationManufacturerDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeDeviceClassificationUserData, &model.DeviceClassificationUserDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeDeviceConfigurationKeyValueConstraintsListData, &model.DeviceConfigurationKeyValueConstraintsDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, &model.DeviceConfigurationKeyValueDescriptionDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeDeviceConfigurationKeyValueListData, &model.DeviceConfigurationKeyValueDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeDeviceDiagnosisHeartbeatData, &model.DeviceDiagnosisHeartbeatDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeDeviceDiagnosisServiceData, &model.DeviceDiagnosisServiceDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeDeviceDiagnosisStateData, &model.DeviceDiagnosisStateDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeDirectControlActivityListData, &model.DirectControlActivityDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeDirectControlDescriptionData, &model.DirectControlDescriptionDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeElectricalConnectionDescriptionListData, &model.ElectricalConnectionDescriptionDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeElectricalConnectionParameterDescriptionListData, &model.ElectricalConnectionParameterDescriptionDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeElectricalConnectionPermittedValueSetListData, &model.ElectricalConnectionPermittedValueSetDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeElectricalConnectionStateListData, &model.ElectricalConnectionStateDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeHvacOperationModeDescriptionListData, &model.HvacOperationModeDescriptionDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeHvacOverrunDescriptionListData, &model.HvacOverrunDescriptionDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeHvacOverrunListData, &model.HvacOverrunDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeHvacSystemFunctionDescriptionListData, &model.HvacSystemFunctionDescriptionDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeHvacSystemFunctionListData, &model.HvacSystemFunctionDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeHvacSystemFunctionOperationModeRelationListData, &model.HvacSystemFunctionOperationModeRelationDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeHvacSystemFunctionPowerSequenceRelationListData, &model.HvacSystemFunctionPowerSequenceRelationDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeHvacSystemFunctionSetPointRelationListData, &model.HvacSystemFunctionSetpointRelationDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeIdentificationListData, &model.IdentificationDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeIncentiveDescriptionListData, &model.IncentiveDescriptionDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeIncentiveListData, &model.IncentiveDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeIncentiveTableConstraintsData, &model.IncentiveTableConstraintsDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeIncentiveTableData, &model.IncentiveTableDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeIncentiveTableDescriptionData, &model.IncentiveTableDescriptionDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeLoadControlEventListData, &model.LoadControlEventDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeLoadControlLimitConstraintsListData, &model.LoadControlLimitConstraintsDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeLoadControlLimitDescriptionListData, &model.LoadControlLimitDescriptionDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeLoadControlLimitListData, &model.LoadControlLimitDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeLoadControlNodeData, &model.LoadControlNodeDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeLoadControlStateListData, &model.LoadControlStateDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeMeasurementConstraintsListData, &model.MeasurementConstraintsDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeMeasurementDescriptionListData, &model.MeasurementDescriptionDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeMeasurementListData, &model.MeasurementDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeMeasurementThresholdRelationListData, &model.MeasurementThresholdRelationDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeMessagingListData, &model.MessagingDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeNetworkManagementAbortCall, &model.NetworkManagementAbortCallElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeNetworkManagementAddNodeCall, &model.NetworkManagementAddNodeCallElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeNetworkManagementDeviceDescriptionListData, &model.NetworkManagementDeviceDescriptionDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeNetworkManagementDiscoverCall, &model.NetworkManagementDiscoverCallElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeNetworkManagementEntityDescriptionListData, &model.NetworkManagementEntityDescriptionDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeNetworkManagementFeatureDescriptionListData, &model.NetworkManagementFeatureDescriptionDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeNetworkManagementJoiningModeData, &model.NetworkManagementJoiningModeDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeNetworkManagementModifyNodeCall, &model.NetworkManagementModifyNodeCallElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeNetworkManagementProcessStateData, &model.NetworkManagementProcessStateDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeNetworkManagementRemoveNodeCall, &model.NetworkManagementRemoveNodeCallElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeNetworkManagementReportCandidateData, &model.NetworkManagementReportCandidateDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeNetworkManagementScanNetworkCall, &model.NetworkManagementScanNetworkCallElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeNodeManagementBindingData, &model.NodeManagementBindingDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeNodeManagementBindingDeleteCall, &model.NodeManagementBindingDeleteCallElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeNodeManagementBindingRequestCall, &model.NodeManagementBindingRequestCallElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeNodeManagementDestinationListData, &model.NodeManagementDestinationDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeNodeManagementDetailedDiscoveryData, &model.NodeManagementDetailedDiscoveryDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeNodeManagementSubscriptionData, &model.NodeManagementSubscriptionDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeNodeManagementSubscriptionDeleteCall, &model.NodeManagementSubscriptionDeleteCallElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeNodeManagementSubscriptionRequestCall, &model.NodeManagementSubscriptionRequestCallElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeNodeManagementUseCaseData, &model.NodeManagementUseCaseDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeOperatingConstraintsDurationListData, &model.OperatingConstraintsDurationDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeOperatingConstraintsInterruptListData, &model.OperatingConstraintsInterruptDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeOperatingConstraintsPowerDescriptionListData, &model.OperatingConstraintsPowerDescriptionDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeOperatingConstraintsPowerLevelListData, &model.OperatingConstraintsPowerLevelDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeOperatingConstraintsPowerRangeListData, &model.OperatingConstraintsPowerRangeDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeOperatingConstraintsResumeImplicationListData, &model.OperatingConstraintsResumeImplicationDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypePowerSequenceAlternativesRelationListData, &model.PowerSequenceAlternativesRelationDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypePowerSequenceDescriptionListData, &model.PowerSequenceDescriptionDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypePowerSequenceNodeScheduleInformationData, &model.PowerSequenceNodeScheduleInformationDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypePowerSequencePriceCalculationRequestCall, &model.PowerSequencePriceCalculationRequestCallElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypePowerSequencePriceListData, &model.PowerSequencePriceDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypePowerSequenceScheduleConfigurationRequestCall, &model.PowerSequenceScheduleConfigurationRequestCallElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypePowerSequenceScheduleConstraintsListData, &model.PowerSequenceScheduleConstraintsDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypePowerSequenceScheduleListData, &model.PowerSequenceScheduleDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypePowerSequenceSchedulePreferenceListData, &model.PowerSequenceSchedulePreferenceDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypePowerSequenceStateListData, &model.PowerSequenceStateDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypePowerTimeSlotScheduleConstraintsListData, &model.PowerTimeSlotScheduleConstraintsDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypePowerTimeSlotScheduleListData, &model.PowerTimeSlotScheduleDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypePowerTimeSlotValueListData, &model.PowerTimeSlotValueDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeSensingListData, &model.SensingDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeSetpointConstraintsListData, &model.SetpointConstraintsDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeSetpointDescriptionListData, &model.SensingDescriptionDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeSetpointListData, &model.SetpointDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeSmartEnergyManagementPsConfigurationRequestCall, &model.SmartEnergyManagementPsConfigurationRequestCallElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeSmartEnergyManagementPsData, &model.SmartEnergyManagementPsDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeSmartEnergyManagementPsPriceCalculationRequestCall, &model.SmartEnergyManagementPsPriceCalculationRequestCallElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeSmartEnergyManagementPsPriceData, &model.SmartEnergyManagementPsPriceDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeSpecificationVersionListData, &model.SpecificationVersionDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeSubscriptionManagementDeleteCall, &model.SubscriptionManagementDeleteCallElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeSubscriptionManagementEntryListData, &model.SubscriptionManagementEntryDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeSubscriptionManagementRequestCall, &model.SubscriptionManagementRequestCallElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeSupplyConditionListData, &model.SupplyConditionDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeSupplyConditionDescriptionListData, &model.SupplyConditionDescriptionDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeSupplyConditionThresholdRelationListData, &model.SupplyConditionThresholdRelationDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeTariffBoundaryRelationListData, &model.TariffBoundaryRelationDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeTariffDescriptionListData, &model.TariffDescriptionDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeTariffListData, &model.TariffDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeTariffOverallConstraintsData, &model.TariffOverallConstraintsDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeTariffTierRelationListData, &model.TariffTierRelationDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeTaskManagementJobDescriptionListData, &model.TaskManagementJobDescriptionDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeTaskManagementJobListData, &model.TaskManagementJobDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeTaskManagementJobRelationListData, &model.TaskManagementJobRelationDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeTaskManagementOverviewData, &model.TaskManagementOverviewDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeThresholdConstraintsListData, &model.ThresholdConstraintsDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeThresholdDescriptionListData, &model.ThresholdDescriptionDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeThresholdListData, &model.ThresholdDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeTierBoundaryDescriptionListData, &model.TierBoundaryDescriptionDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeTierBoundaryListData, &model.TierBoundaryDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeTierDescriptionListData, &model.TierDescriptionDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeTierIncentiveRelationListData, &model.TierIncentiveRelationDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeTierListData, &model.TierDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeTimeDistributorData, &model.TimeDistributorDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeTimeDistributorEnquiryCall, &model.TimeDistributorEnquiryCallElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeTimeInformationData, &model.TimeInformationDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeTimePrecisionData, &model.TimePrecisionDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeTimeSeriesConstraintsListData, &model.TimeSeriesConstraintsDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeTimeSeriesDescriptionListData, &model.TimeSeriesDescriptionDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeTimeSeriesListData, &model.TimeSeriesDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeTimeTableConstraintsListData, &model.TimeTableConstraintsDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeTimeTableDescriptionListData, &model.TimeTableDescriptionDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeTimeTableListData, &model.TimeTableDataElementsType{})
	assert.NotNil(suite.T(), result)

	result = addElementToFilter(filter, model.FunctionTypeUseCaseInformationListData, &model.UseCaseInformationDataElementsType{})
	assert.NotNil(suite.T(), result)
}

func (suite *FctDataCmdSuite) Test_CreateCmd() {
	result := createCmd(model.FunctionTypeActuatorLevelData, &model.ActuatorLevelDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeActuatorLevelDescriptionData, &model.ActuatorLevelDescriptionDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeActuatorSwitchData, &model.ActuatorSwitchDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeActuatorSwitchDescriptionData, &model.ActuatorSwitchDescriptionDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeAlarmListData, &model.AlarmListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeBillConstraintsListData, &model.BillConstraintsListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeBillDescriptionListData, &model.BillDescriptionListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeBillListData, &model.BillListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeBindingManagementDeleteCall, &model.BindingManagementDeleteCallType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeBindingManagementEntryListData, &model.BindingManagementEntryListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeBindingManagementRequestCall, &model.BindingManagementRequestCallType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeCommodityListData, &model.CommodityListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeDataTunnelingCall, &model.DataTunnelingCallType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeDeviceClassificationManufacturerData, &model.DeviceClassificationManufacturerDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeDeviceClassificationUserData, &model.DeviceClassificationUserDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeDeviceConfigurationKeyValueConstraintsListData, &model.DeviceConfigurationKeyValueConstraintsListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, &model.DeviceConfigurationKeyValueDescriptionListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeDeviceConfigurationKeyValueListData, &model.DeviceConfigurationKeyValueListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeDeviceDiagnosisHeartbeatData, &model.DeviceDiagnosisHeartbeatDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeDeviceDiagnosisServiceData, &model.DeviceDiagnosisServiceDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeDeviceDiagnosisStateData, &model.DeviceDiagnosisStateDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeDirectControlActivityListData, &model.DirectControlActivityListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeDirectControlDescriptionData, &model.DirectControlDescriptionDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeElectricalConnectionDescriptionListData, &model.ElectricalConnectionDescriptionListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeElectricalConnectionParameterDescriptionListData, &model.ElectricalConnectionParameterDescriptionListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeElectricalConnectionPermittedValueSetListData, &model.ElectricalConnectionPermittedValueSetListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeElectricalConnectionStateListData, &model.ElectricalConnectionStateListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeHvacOperationModeDescriptionListData, &model.HvacOperationModeDescriptionListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeHvacOverrunDescriptionListData, &model.HvacOverrunDescriptionListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeHvacOverrunListData, &model.HvacOverrunListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeHvacSystemFunctionDescriptionListData, &model.HvacSystemFunctionDescriptionListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeHvacSystemFunctionListData, &model.HvacSystemFunctionListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeHvacSystemFunctionOperationModeRelationListData, &model.HvacSystemFunctionOperationModeRelationListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeHvacSystemFunctionPowerSequenceRelationListData, &model.HvacSystemFunctionPowerSequenceRelationListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeHvacSystemFunctionSetPointRelationListData, &model.HvacSystemFunctionSetpointRelationListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeIdentificationListData, &model.IdentificationListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeIncentiveDescriptionListData, &model.IncentiveDescriptionListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeIncentiveListData, &model.IncentiveListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeIncentiveTableConstraintsData, &model.IncentiveTableConstraintsDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeIncentiveTableData, &model.IncentiveTableDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeIncentiveTableDescriptionData, &model.IncentiveTableDescriptionDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeLoadControlEventListData, &model.LoadControlEventListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeLoadControlLimitConstraintsListData, &model.LoadControlLimitConstraintsListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeLoadControlLimitDescriptionListData, &model.LoadControlLimitDescriptionListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeLoadControlLimitListData, &model.LoadControlLimitListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeLoadControlNodeData, &model.LoadControlNodeDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeLoadControlStateListData, &model.LoadControlStateListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeMeasurementConstraintsListData, &model.MeasurementConstraintsListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeMeasurementDescriptionListData, &model.MeasurementDescriptionListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeMeasurementListData, &model.MeasurementListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeMeasurementThresholdRelationListData, &model.MeasurementThresholdRelationListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeMessagingListData, &model.MessagingListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeNetworkManagementAbortCall, &model.NetworkManagementAbortCallType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeNetworkManagementAddNodeCall, &model.NetworkManagementAddNodeCallType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeNetworkManagementDeviceDescriptionListData, &model.NetworkManagementDeviceDescriptionListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeNetworkManagementDiscoverCall, &model.NetworkManagementDiscoverCallType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeNetworkManagementEntityDescriptionListData, &model.NetworkManagementEntityDescriptionListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeNetworkManagementFeatureDescriptionListData, &model.NetworkManagementFeatureDescriptionListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeNetworkManagementJoiningModeData, &model.NetworkManagementJoiningModeDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeNetworkManagementModifyNodeCall, &model.NetworkManagementModifyNodeCallType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeNetworkManagementProcessStateData, &model.NetworkManagementProcessStateDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeNetworkManagementRemoveNodeCall, &model.NetworkManagementRemoveNodeCallType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeNetworkManagementReportCandidateData, &model.NetworkManagementReportCandidateDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeNetworkManagementScanNetworkCall, &model.NetworkManagementScanNetworkCallType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeOperatingConstraintsDurationListData, &model.OperatingConstraintsDurationListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeOperatingConstraintsInterruptListData, &model.OperatingConstraintsInterruptListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeOperatingConstraintsPowerDescriptionListData, &model.OperatingConstraintsPowerDescriptionListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeOperatingConstraintsPowerLevelListData, &model.OperatingConstraintsPowerLevelListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeOperatingConstraintsPowerRangeListData, &model.OperatingConstraintsPowerRangeListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeOperatingConstraintsResumeImplicationListData, &model.OperatingConstraintsResumeImplicationListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypePowerSequenceAlternativesRelationListData, &model.PowerSequenceAlternativesRelationListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypePowerSequenceDescriptionListData, &model.PowerSequenceDescriptionListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypePowerSequenceNodeScheduleInformationData, &model.PowerSequenceNodeScheduleInformationDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypePowerSequencePriceCalculationRequestCall, &model.PowerSequencePriceCalculationRequestCallType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypePowerSequencePriceListData, &model.PowerSequencePriceListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypePowerSequenceScheduleConfigurationRequestCall, &model.PowerSequenceScheduleConfigurationRequestCallType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypePowerSequenceScheduleConstraintsListData, &model.PowerSequenceScheduleConstraintsListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypePowerSequenceScheduleListData, &model.PowerSequenceScheduleListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypePowerSequenceSchedulePreferenceListData, &model.PowerSequenceSchedulePreferenceListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypePowerSequenceStateListData, &model.PowerSequenceStateListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypePowerTimeSlotScheduleConstraintsListData, &model.PowerTimeSlotScheduleConstraintsListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypePowerTimeSlotScheduleListData, &model.PowerTimeSlotScheduleListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypePowerTimeSlotValueListData, &model.PowerTimeSlotValueListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeResultData, &model.ResultDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeSensingDescriptionData, &model.SensingDescriptionDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeSensingListData, &model.SensingListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeSetpointConstraintsListData, &model.SetpointConstraintsListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeSetpointDescriptionListData, &model.SetpointDescriptionListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeSetpointListData, &model.SetpointListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeSmartEnergyManagementPsConfigurationRequestCall, &model.SmartEnergyManagementPsConfigurationRequestCallType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeSmartEnergyManagementPsData, &model.SmartEnergyManagementPsDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeSmartEnergyManagementPsPriceCalculationRequestCall, &model.SmartEnergyManagementPsPriceCalculationRequestCallType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeSmartEnergyManagementPsPriceData, &model.SmartEnergyManagementPsPriceDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeSpecificationVersionListData, &model.SpecificationVersionListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeSupplyConditionListData, &model.SupplyConditionListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeSupplyConditionThresholdRelationListData, &model.SupplyConditionThresholdRelationListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeTariffBoundaryRelationListData, &model.TariffBoundaryRelationListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeTariffDescriptionListData, &model.TariffDescriptionListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeTariffListData, &model.TariffListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeTariffOverallConstraintsData, &model.TariffOverallConstraintsDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeTariffTierRelationListData, &model.TariffTierRelationListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeTaskManagementJobDescriptionListData, &model.TaskManagementJobDescriptionListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeTaskManagementJobListData, &model.TaskManagementJobListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeTaskManagementJobRelationListData, &model.TaskManagementJobRelationListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeTaskManagementOverviewData, &model.TaskManagementOverviewDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeThresholdConstraintsListData, &model.ThresholdConstraintsListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeThresholdDescriptionListData, &model.ThresholdDescriptionListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeThresholdListData, &model.ThresholdListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeTierBoundaryDescriptionListData, &model.TierBoundaryDescriptionListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeTierBoundaryListData, &model.TierBoundaryListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeTierDescriptionListData, &model.TierDescriptionListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeTierIncentiveRelationListData, &model.TierIncentiveRelationListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeTierListData, &model.TierListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeTimeDistributorData, &model.TimeDistributorDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeTimeDistributorEnquiryCall, &model.TimeDistributorEnquiryCallType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeTimeInformationData, &model.TimeInformationDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeTimePrecisionData, &model.TimePrecisionDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeTimeSeriesConstraintsListData, &model.TimeSeriesConstraintsListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeTimeSeriesDescriptionListData, &model.TimeSeriesDescriptionListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeTimeSeriesListData, &model.TimeSeriesListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeTimeTableConstraintsListData, &model.TimeTableConstraintsListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeTimeTableDescriptionListData, &model.TimeTableDescriptionListDataType{})
	assert.NotNil(suite.T(), result)

	result = createCmd(model.FunctionTypeTimeTableListData, &model.TimeTableListDataType{})
	assert.NotNil(suite.T(), result)
}
