package spine

import (
	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/DerAndereAndi/eebus-go/util"
)

type FunctionDataCmd interface {
	FunctionData
	ReadCmdType() model.CmdType
	ReplyCmdType() model.CmdType
	NotifyCmdType(partial bool) model.CmdType
	WriteCmdType() model.CmdType
}

var _ FunctionDataCmd = (*FunctionDataCmdImpl[int])(nil)

type FunctionDataCmdImpl[T any] struct {
	*FunctionDataImpl[T]
}

func NewFunctionDataCmd[T any](function model.FunctionType) *FunctionDataCmdImpl[T] {
	return &FunctionDataCmdImpl[T]{
		FunctionDataImpl: NewFunctionData[T](function),
	}
}

func (r *FunctionDataCmdImpl[T]) ReadCmdType() model.CmdType {
	cmd := createCmd[T](r.functionType, nil)
	return cmd
}

func (r *FunctionDataCmdImpl[T]) ReplyCmdType() model.CmdType {
	cmd := createCmd(r.functionType, r.data)
	return cmd
}

func (r *FunctionDataCmdImpl[T]) NotifyCmdType(partial bool) model.CmdType {
	cmd := createCmd(r.functionType, r.data)
	cmd.Function = util.Ptr(model.FunctionType(r.functionType))
	cmd.Filter = filterType(partial)
	return cmd
}

func (r *FunctionDataCmdImpl[T]) WriteCmdType() model.CmdType {
	cmd := createCmd(r.functionType, r.data)
	return cmd
}

func filterType(partial bool) []model.FilterType {
	if partial {
		return []model.FilterType{{CmdControl: &model.CmdControlType{Partial: &model.ElementTagType{}}}}
	}
	return nil
}

func createCmd[T any](function model.FunctionType, data *T) model.CmdType {
	result := model.CmdType{}

	switch function {
	case model.FunctionTypeActuatorLevelData:
		result.ActuatorLevelData = castData[model.ActuatorLevelDataType](data)
	case model.FunctionTypeActuatorLevelDescriptionData:
		result.ActuatorLevelDescriptionData = castData[model.ActuatorLevelDescriptionDataType](data)
	case model.FunctionTypeActuatorSwitchData:
		result.ActuatorSwitchData = castData[model.ActuatorSwitchDataType](data)
	case model.FunctionTypeActuatorSwitchDescriptionData:
		result.ActuatorSwitchDescriptionData = castData[model.ActuatorSwitchDescriptionDataType](data)
	case model.FunctionTypeAlarmListData:
		result.AlarmListData = castData[model.AlarmListDataType](data)
	case model.FunctionTypeBillConstraintsListData:
		result.BillConstraintsListData = castData[model.BillConstraintsListDataType](data)
	case model.FunctionTypeBillDescriptionListData:
		result.BillDescriptionListData = castData[model.BillDescriptionListDataType](data)
	case model.FunctionTypeBillListData:
		result.BillListData = castData[model.BillListDataType](data)
	case model.FunctionTypeBindingManagementDeleteCall:
		result.BindingManagementDeleteCall = castData[model.BindingManagementDeleteCallType](data)
	case model.FunctionTypeBindingManagementEntryListData:
		result.BindingManagementEntryListData = castData[model.BindingManagementEntryListDataType](data)
	case model.FunctionTypeBindingManagementRequestCall:
		result.BindingManagementRequestCall = castData[model.BindingManagementRequestCallType](data)
	case model.FunctionTypeCommodityListData:
		result.CommodityListData = castData[model.CommodityListDataType](data)
	case model.FunctionTypeDataTunnelingCall:
		result.DataTunnelingCall = castData[model.DataTunnelingCallType](data)
	case model.FunctionTypeDeviceClassificationManufacturerData:
		result.DeviceClassificationManufacturerData = castData[model.DeviceClassificationManufacturerDataType](data)
	case model.FunctionTypeDeviceClassificationUserData:
		result.DeviceClassificationUserData = castData[model.DeviceClassificationUserDataType](data)
	case model.FunctionTypeDeviceConfigurationKeyValueConstraintsListData:
		result.DeviceConfigurationKeyValueConstraintsListData = castData[model.DeviceConfigurationKeyValueConstraintsListDataType](data)
	case model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData:
		result.DeviceConfigurationKeyValueDescriptionListData = castData[model.DeviceConfigurationKeyValueDescriptionListDataType](data)
	case model.FunctionTypeDeviceConfigurationKeyValueListData:
		result.DeviceConfigurationKeyValueListData = castData[model.DeviceConfigurationKeyValueListDataType](data)
	case model.FunctionTypeDeviceDiagnosisHeartbeatData:
		result.DeviceDiagnosisHeartbeatData = castData[model.DeviceDiagnosisHeartbeatDataType](data)
	case model.FunctionTypeDeviceDiagnosisServiceData:
		result.DeviceDiagnosisServiceData = castData[model.DeviceDiagnosisServiceDataType](data)
	case model.FunctionTypeDeviceDiagnosisStateData:
		result.DeviceDiagnosisStateData = castData[model.DeviceDiagnosisStateDataType](data)
	case model.FunctionTypeDirectControlActivityListData:
		result.DirectControlActivityListData = castData[model.DirectControlActivityListDataType](data)
	case model.FunctionTypeDirectControlDescriptionData:
		result.DirectControlDescriptionData = castData[model.DirectControlDescriptionDataType](data)
	case model.FunctionTypeElectricalConnectionDescriptionListData:
		result.ElectricalConnectionDescriptionListData = castData[model.ElectricalConnectionDescriptionListDataType](data)
	case model.FunctionTypeElectricalConnectionParameterDescriptionListData:
		result.ElectricalConnectionParameterDescriptionListData = castData[model.ElectricalConnectionParameterDescriptionListDataType](data)
	case model.FunctionTypeElectricalConnectionPermittedValueSetListData:
		result.ElectricalConnectionPermittedValueSetListData = castData[model.ElectricalConnectionPermittedValueSetListDataType](data)
	case model.FunctionTypeElectricalConnectionStateListData:
		result.ElectricalConnectionStateListData = castData[model.ElectricalConnectionStateListDataType](data)
	case model.FunctionTypeHvacOperationModeDescriptionListData:
		result.HvacOperationModeDescriptionListData = castData[model.HvacOperationModeDescriptionListDataType](data)
	case model.FunctionTypeHvacOverrunDescriptionListData:
		result.HvacOverrunDescriptionListData = castData[model.HvacOverrunDescriptionListDataType](data)
	case model.FunctionTypeHvacOverrunListData:
		result.HvacOverrunListData = castData[model.HvacOverrunListDataType](data)
	case model.FunctionTypeHvacSystemFunctionDescriptionListData:
		result.HvacSystemFunctionDescriptionListData = castData[model.HvacSystemFunctionDescriptionListDataType](data)
	case model.FunctionTypeHvacSystemFunctionListData:
		result.HvacSystemFunctionListData = castData[model.HvacSystemFunctionListDataType](data)
	case model.FunctionTypeHvacSystemFunctionOperationModeRelationListData:
		result.HvacSystemFunctionOperationModeRelationListData = castData[model.HvacSystemFunctionOperationModeRelationListDataType](data)
	case model.FunctionTypeHvacSystemFunctionPowerSequenceRelationListData:
		result.HvacSystemFunctionPowerSequenceRelationListData = castData[model.HvacSystemFunctionPowerSequenceRelationListDataType](data)
	case model.FunctionTypeHvacSystemFunctionSetPointRelationListData:
		result.HvacSystemFunctionSetPointRelationListData = castData[model.HvacSystemFunctionSetpointRelationListDataType](data)
	case model.FunctionTypeIdentificationListData:
		result.IdentificationListData = castData[model.IdentificationListDataType](data)
	case model.FunctionTypeIncentiveDescriptionListData:
		result.IncentiveDescriptionListData = castData[model.IncentiveDescriptionListDataType](data)
	case model.FunctionTypeIncentiveListData:
		result.IncentiveListData = castData[model.IncentiveListDataType](data)
	case model.FunctionTypeIncentiveTableConstraintsData:
		result.IncentiveTableConstraintsData = castData[model.IncentiveTableConstraintsDataType](data)
	case model.FunctionTypeIncentiveTableData:
		result.IncentiveTableData = castData[model.IncentiveTableDataType](data)
	case model.FunctionTypeIncentiveTableDescriptionData:
		result.IncentiveTableDescriptionData = castData[model.IncentiveTableDescriptionDataType](data)
	case model.FunctionTypeLoadControlEventListData:
		result.LoadControlEventListData = castData[model.LoadControlEventListDataType](data)
	case model.FunctionTypeLoadControlLimitConstraintsListData:
		result.LoadControlLimitConstraintsListData = castData[model.LoadControlLimitConstraintsListDataType](data)
	case model.FunctionTypeLoadControlLimitDescriptionListData:
		result.LoadControlLimitDescriptionListData = castData[model.LoadControlLimitDescriptionListDataType](data)
	case model.FunctionTypeLoadControlLimitListData:
		result.LoadControlLimitListData = castData[model.LoadControlLimitListDataType](data)
	case model.FunctionTypeLoadControlNodeData:
		result.LoadControlNodeData = castData[model.LoadControlNodeDataType](data)
	case model.FunctionTypeLoadControlStateListData:
		result.LoadControlStateListData = castData[model.LoadControlStateListDataType](data)
	case model.FunctionTypeMeasurementConstraintsListData:
		result.MeasurementConstraintsListData = castData[model.MeasurementConstraintsListDataType](data)
	case model.FunctionTypeMeasurementDescriptionListData:
		result.MeasurementDescriptionListData = castData[model.MeasurementDescriptionListDataType](data)
	case model.FunctionTypeMeasurementListData:
		result.MeasurementListData = castData[model.MeasurementListDataType](data)
	case model.FunctionTypeMeasurementThresholdRelationListData:
		result.MeasurementThresholdRelationListData = castData[model.MeasurementThresholdRelationListDataType](data)
	case model.FunctionTypeMessagingListData:
		result.MessagingListData = castData[model.MessagingListDataType](data)
	case model.FunctionTypeNetworkManagementAbortCall:
		result.NetworkManagementAbortCall = castData[model.NetworkManagementAbortCallType](data)
	case model.FunctionTypeNetworkManagementAddNodeCall:
		result.NetworkManagementAddNodeCall = castData[model.NetworkManagementAddNodeCallType](data)
	case model.FunctionTypeNetworkManagementDeviceDescriptionListData:
		result.NetworkManagementDeviceDescriptionListData = castData[model.NetworkManagementDeviceDescriptionListDataType](data)
	case model.FunctionTypeNetworkManagementDiscoverCall:
		result.NetworkManagementDiscoverCall = castData[model.NetworkManagementDiscoverCallType](data)
	case model.FunctionTypeNetworkManagementEntityDescriptionListData:
		result.NetworkManagementEntityDescriptionListData = castData[model.NetworkManagementEntityDescriptionListDataType](data)
	case model.FunctionTypeNetworkManagementFeatureDescriptionListData:
		result.NetworkManagementFeatureDescriptionListData = castData[model.NetworkManagementFeatureDescriptionListDataType](data)
	case model.FunctionTypeNetworkManagementJoiningModeData:
		result.NetworkManagementJoiningModeData = castData[model.NetworkManagementJoiningModeDataType](data)
	case model.FunctionTypeNetworkManagementModifyNodeCall:
		result.NetworkManagementModifyNodeCall = castData[model.NetworkManagementModifyNodeCallType](data)
	case model.FunctionTypeNetworkManagementProcessStateData:
		result.NetworkManagementProcessStateData = castData[model.NetworkManagementProcessStateDataType](data)
	case model.FunctionTypeNetworkManagementRemoveNodeCall:
		result.NetworkManagementRemoveNodeCall = castData[model.NetworkManagementRemoveNodeCallType](data)
	case model.FunctionTypeNetworkManagementReportCandidateData:
		result.NetworkManagementReportCandidateData = castData[model.NetworkManagementReportCandidateDataType](data)
	case model.FunctionTypeNetworkManagementScanNetworkCall:
		result.NetworkManagementScanNetworkCall = castData[model.NetworkManagementScanNetworkCallType](data)
	case model.FunctionTypeOperatingConstraintsDurationListData:
		result.OperatingConstraintsDurationListData = castData[model.OperatingConstraintsDurationListDataType](data)
	case model.FunctionTypeOperatingConstraintsInterruptListData:
		result.OperatingConstraintsInterruptListData = castData[model.OperatingConstraintsInterruptListDataType](data)
	case model.FunctionTypeOperatingConstraintsPowerDescriptionListData:
		result.OperatingConstraintsPowerDescriptionListData = castData[model.OperatingConstraintsPowerDescriptionListDataType](data)
	case model.FunctionTypeOperatingConstraintsPowerLevelListData:
		result.OperatingConstraintsPowerLevelListData = castData[model.OperatingConstraintsPowerLevelListDataType](data)
	case model.FunctionTypeOperatingConstraintsPowerRangeListData:
		result.OperatingConstraintsPowerRangeListData = castData[model.OperatingConstraintsPowerRangeListDataType](data)
	case model.FunctionTypeOperatingConstraintsResumeImplicationListData:
		result.OperatingConstraintsResumeImplicationListData = castData[model.OperatingConstraintsResumeImplicationListDataType](data)
	case model.FunctionTypePowerSequenceAlternativesRelationListData:
		result.PowerSequenceAlternativesRelationListData = castData[model.PowerSequenceAlternativesRelationListDataType](data)
	case model.FunctionTypePowerSequenceDescriptionListData:
		result.PowerSequenceDescriptionListData = castData[model.PowerSequenceDescriptionListDataType](data)
	case model.FunctionTypePowerSequenceNodeScheduleInformationData:
		result.PowerSequenceNodeScheduleInformationData = castData[model.PowerSequenceNodeScheduleInformationDataType](data)
	case model.FunctionTypePowerSequencePriceCalculationRequestCall:
		result.PowerSequencePriceCalculationRequestCall = castData[model.PowerSequencePriceCalculationRequestCallType](data)
	case model.FunctionTypePowerSequencePriceListData:
		result.PowerSequencePriceListData = castData[model.PowerSequencePriceListDataType](data)
	case model.FunctionTypePowerSequenceScheduleConfigurationRequestCall:
		result.PowerSequenceScheduleConfigurationRequestCall = castData[model.PowerSequenceScheduleConfigurationRequestCallType](data)
	case model.FunctionTypePowerSequenceScheduleConstraintsListData:
		result.PowerSequenceScheduleConstraintsListData = castData[model.PowerSequenceScheduleConstraintsListDataType](data)
	case model.FunctionTypePowerSequenceScheduleListData:
		result.PowerSequenceScheduleListData = castData[model.PowerSequenceScheduleListDataType](data)
	case model.FunctionTypePowerSequenceSchedulePreferenceListData:
		result.PowerSequenceSchedulePreferenceListData = castData[model.PowerSequenceSchedulePreferenceListDataType](data)
	case model.FunctionTypePowerSequenceStateListData:
		result.PowerSequenceStateListData = castData[model.PowerSequenceStateListDataType](data)
	case model.FunctionTypePowerTimeSlotScheduleConstraintsListData:
		result.PowerTimeSlotScheduleConstraintsListData = castData[model.PowerTimeSlotScheduleConstraintsListDataType](data)
	case model.FunctionTypePowerTimeSlotScheduleListData:
		result.PowerTimeSlotScheduleListData = castData[model.PowerTimeSlotScheduleListDataType](data)
	case model.FunctionTypePowerTimeSlotValueListData:
		result.PowerTimeSlotValueListData = castData[model.PowerTimeSlotValueListDataType](data)
	case model.FunctionTypeResultData:
		result.ResultData = castData[model.ResultDataType](data)
	case model.FunctionTypeSensingDescriptionData:
		result.SensingDescriptionData = castData[model.SensingDescriptionDataType](data)
	case model.FunctionTypeSensingListData:
		result.SensingListData = castData[model.SensingListDataType](data)
	case model.FunctionTypeSetpointConstraintsListData:
		result.SetpointConstraintsListData = castData[model.SetpointConstraintsListDataType](data)
	case model.FunctionTypeSetpointDescriptionListData:
		result.SetpointDescriptionListData = castData[model.SetpointDescriptionListDataType](data)
	case model.FunctionTypeSetpointListData:
		result.SetpointListData = castData[model.SetpointListDataType](data)
	case model.FunctionTypeSmartEnergyManagementPsConfigurationRequestCall:
		result.SmartEnergyManagementPsConfigurationRequestCall = castData[model.SmartEnergyManagementPsConfigurationRequestCallType](data)
	case model.FunctionTypeSmartEnergyManagementPsData:
		result.SmartEnergyManagementPsData = castData[model.SmartEnergyManagementPsDataType](data)
	case model.FunctionTypeSmartEnergyManagementPsPriceCalculationRequestCall:
		result.SmartEnergyManagementPsPriceCalculationRequestCall = castData[model.SmartEnergyManagementPsPriceCalculationRequestCallType](data)
	case model.FunctionTypeSmartEnergyManagementPsPriceData:
		result.SmartEnergyManagementPsPriceData = castData[model.SmartEnergyManagementPsPriceDataType](data)
	case model.FunctionTypeSpecificationVersionListData:
		result.SpecificationVersionListData = castData[model.SpecificationVersionListDataType](data)
	case model.FunctionTypeSupplyConditionListData:
		result.SupplyConditionListData = castData[model.SupplyConditionListDataType](data)
	case model.FunctionTypeSupplyConditionThresholdRelationListData:
		result.SupplyConditionThresholdRelationListData = castData[model.SupplyConditionThresholdRelationListDataType](data)
	case model.FunctionTypeTariffBoundaryRelationListData:
		result.TariffBoundaryRelationListData = castData[model.TariffBoundaryRelationListDataType](data)
	case model.FunctionTypeTariffDescriptionListData:
		result.TariffDescriptionListData = castData[model.TariffDescriptionListDataType](data)
	case model.FunctionTypeTariffListData:
		result.TariffListData = castData[model.TariffListDataType](data)
	case model.FunctionTypeTariffOverallConstraintsData:
		result.TariffOverallConstraintsData = castData[model.TariffOverallConstraintsDataType](data)
	case model.FunctionTypeTariffTierRelationListData:
		result.TariffTierRelationListData = castData[model.TariffTierRelationListDataType](data)
	case model.FunctionTypeTaskManagementJobDescriptionListData:
		result.TaskManagementJobDescriptionListData = castData[model.TaskManagementJobDescriptionListDataType](data)
	case model.FunctionTypeTaskManagementJobListData:
		result.TaskManagementJobListData = castData[model.TaskManagementJobListDataType](data)
	case model.FunctionTypeTaskManagementJobRelationListData:
		result.TaskManagementJobRelationListData = castData[model.TaskManagementJobRelationListDataType](data)
	case model.FunctionTypeTaskManagementOverviewData:
		result.TaskManagementOverviewData = castData[model.TaskManagementOverviewDataType](data)
	case model.FunctionTypeThresholdConstraintsListData:
		result.ThresholdConstraintsListData = castData[model.ThresholdConstraintsListDataType](data)
	case model.FunctionTypeThresholdDescriptionListData:
		result.ThresholdDescriptionListData = castData[model.ThresholdDescriptionListDataType](data)
	case model.FunctionTypeThresholdListData:
		result.ThresholdListData = castData[model.ThresholdListDataType](data)
	case model.FunctionTypeTierBoundaryDescriptionListData:
		result.TierBoundaryDescriptionListData = castData[model.TierBoundaryDescriptionListDataType](data)
	case model.FunctionTypeTierBoundaryListData:
		result.TierBoundaryListData = castData[model.TierBoundaryListDataType](data)
	case model.FunctionTypeTierDescriptionListData:
		result.TierDescriptionListData = castData[model.TierDescriptionListDataType](data)
	case model.FunctionTypeTierIncentiveRelationListData:
		result.TierIncentiveRelationListData = castData[model.TierIncentiveRelationListDataType](data)
	case model.FunctionTypeTierListData:
		result.TierListData = castData[model.TierListDataType](data)
	case model.FunctionTypeTimeDistributorData:
		result.TimeDistributorData = castData[model.TimeDistributorDataType](data)
	case model.FunctionTypeTimeDistributorEnquiryCall:
		result.TimeDistributorEnquiryCall = castData[model.TimeDistributorEnquiryCallType](data)
	case model.FunctionTypeTimeInformationData:
		result.TimeInformationData = castData[model.TimeInformationDataType](data)
	case model.FunctionTypeTimePrecisionData:
		result.TimePrecisionData = castData[model.TimePrecisionDataType](data)
	case model.FunctionTypeTimeSeriesConstraintsListData:
		result.TimeSeriesConstraintsListData = castData[model.TimeSeriesConstraintsListDataType](data)
	case model.FunctionTypeTimeSeriesDescriptionListData:
		result.TimeSeriesDescriptionListData = castData[model.TimeSeriesDescriptionListDataType](data)
	case model.FunctionTypeTimeSeriesListData:
		result.TimeSeriesListData = castData[model.TimeSeriesListDataType](data)
	case model.FunctionTypeTimeTableConstraintsListData:
		result.TimeTableConstraintsListData = castData[model.TimeTableConstraintsListDataType](data)
	case model.FunctionTypeTimeTableDescriptionListData:
		result.TimeTableDescriptionListData = castData[model.TimeTableDescriptionListDataType](data)
	case model.FunctionTypeTimeTableListData:
		result.TimeTableListData = castData[model.TimeTableListDataType](data)
		// add more model types here
	}

	return result
}

func castData[D, S any](data *S) *D {
	if data == nil {
		return new(D)
	}
	return any(data).(*D)
}
