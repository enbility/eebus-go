package spine

import (
	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
)

type FunctionDataCmd interface {
	FunctionData
	ReadCmdType(partialSelector any, elements any) model.CmdType
	ReplyCmdType(partial bool) model.CmdType
	NotifyCmdType(deleteSelector, partialSelector any, partialWithoutSelector bool, deleteElements any) model.CmdType
	WriteCmdType(deleteSelector, partialSelector any, deleteElements any) model.CmdType
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

func (r *FunctionDataCmdImpl[T]) ReadCmdType(partialSelector any, elements any) model.CmdType {
	cmd := createCmd[T](r.functionType, nil)

	var filters []model.FilterType
	filters = filtersForSelectorsElements(r.functionType, filters, nil, partialSelector, nil, elements)
	if len(filters) > 0 {
		cmd.Filter = filters
	}

	return cmd
}

func (r *FunctionDataCmdImpl[T]) ReplyCmdType(partial bool) model.CmdType {
	cmd := createCmd(r.functionType, r.data)
	if partial {
		cmd.Filter = filterEmptyPartial()
	}
	return cmd
}

func (r *FunctionDataCmdImpl[T]) NotifyCmdType(deleteSelector, partialSelector any, partialWithoutSelector bool, deleteElements any) model.CmdType {
	cmd := createCmd(r.functionType, r.data)
	cmd.Function = util.Ptr(model.FunctionType(r.functionType))

	if partialWithoutSelector {
		cmd.Filter = filterEmptyPartial()
		return cmd
	}
	var filters []model.FilterType
	if filters := filtersForSelectorsElements(r.functionType, filters, deleteSelector, partialSelector, deleteElements, nil); len(filters) > 0 {
		cmd.Filter = filters
	}

	return cmd
}

func (r *FunctionDataCmdImpl[T]) WriteCmdType(deleteSelector, partialSelector any, deleteElements any) model.CmdType {
	cmd := createCmd(r.functionType, r.data)

	var filters []model.FilterType
	if filters := filtersForSelectorsElements(r.functionType, filters, deleteSelector, partialSelector, deleteElements, nil); len(filters) > 0 {
		cmd.Filter = filters
	}

	return cmd
}

func filtersForSelectorsElements(functionType model.FunctionType, filters []model.FilterType, deleteSelector, partialSelector any, deleteElements, readElements any) []model.FilterType {
	if deleteSelector != nil || deleteElements != nil {
		filter := model.FilterType{CmdControl: &model.CmdControlType{Delete: &model.ElementTagType{}}}
		if deleteSelector != nil {
			filter = addSelectorToFilter(filter, functionType, &deleteSelector)
		}
		if deleteElements != nil {
			filter = addElementToFilter(filter, functionType, &deleteElements)
		}
		filters = append(filters, filter)
	}

	if partialSelector != nil || readElements != nil {
		filter := model.FilterType{CmdControl: &model.CmdControlType{Partial: &model.ElementTagType{}}}
		if partialSelector != nil {
			filter = addSelectorToFilter(filter, functionType, &partialSelector)
		}
		if readElements != nil {
			filter = addElementToFilter(filter, functionType, &readElements)
		}
		filters = append(filters, filter)
	}

	return filters
}

// simple helper for adding a single filterType without any selectors
func filterEmptyPartial() []model.FilterType {
	return []model.FilterType{{CmdControl: &model.CmdControlType{Partial: &model.ElementTagType{}}}}
}

func addSelectorToFilter[T any](filter model.FilterType, function model.FunctionType, data *T) model.FilterType {
	result := filter

	switch function {
	case model.FunctionTypeAlarmListData:
		result.AlarmListDataSelectors = castData[model.AlarmListDataSelectorsType](data)
	case model.FunctionTypeBillConstraintsListData:
		result.BillConstraintsListDataSelectors = castData[model.BillConstraintsListDataSelectorsType](data)
	case model.FunctionTypeBillDescriptionListData:
		result.BillDescriptionListDataSelectors = castData[model.BillDescriptionListDataSelectorsType](data)
	case model.FunctionTypeBillListData:
		result.BillListDataSelectors = castData[model.BillListDataSelectorsType](data)
	case model.FunctionTypeBindingManagementEntryListData:
		result.BindingManagementEntryListDataSelectors = castData[model.BindingManagementEntryListDataSelectorsType](data)
	case model.FunctionTypeCommodityListData:
		result.CommodityListDataSelectors = castData[model.CommodityListDataSelectorsType](data)
	case model.FunctionTypeDeviceConfigurationKeyValueConstraintsListData:
		result.DeviceConfigurationKeyValueConstraintsListDataSelectors = castData[model.DeviceConfigurationKeyValueConstraintsListDataSelectorsType](data)
	case model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData:
		result.DeviceConfigurationKeyValueDescriptionListDataSelectors = castData[model.DeviceConfigurationKeyValueDescriptionListDataSelectorsType](data)
	case model.FunctionTypeDeviceConfigurationKeyValueListData:
		result.DeviceConfigurationKeyValueListDataSelectors = castData[model.DeviceConfigurationKeyValueListDataSelectorsType](data)
	case model.FunctionTypeDirectControlActivityListData:
		result.DirectControlActivityListDataSelectors = castData[model.DirectControlActivityListDataSelectorsType](data)
	case model.FunctionTypeElectricalConnectionDescriptionListData:
		filter.ElectricalConnectionDescriptionListDataSelectors = castData[model.ElectricalConnectionDescriptionListDataSelectorsType](data)
	case model.FunctionTypeElectricalConnectionParameterDescriptionListData:
		filter.ElectricalConnectionParameterDescriptionListDataSelectors = castData[model.ElectricalConnectionParameterDescriptionListDataSelectorsType](data)
	case model.FunctionTypeElectricalConnectionPermittedValueSetListData:
		result.ElectricalConnectionPermittedValueSetListDataSelectors = castData[model.ElectricalConnectionPermittedValueSetListDataSelectorsType](data)
	case model.FunctionTypeElectricalConnectionStateListData:
		result.ElectricalConnectionStateListDataSelectors = castData[model.ElectricalConnectionStateListDataSelectorsType](data)
	case model.FunctionTypeElectricalConnectionCharacteristicListData:
		result.ElectricalConnectionCharacteristicListDataSelectors = castData[model.ElectricalConnectionCharacteristicListDataSelectorsType](data)
	case model.FunctionTypeHvacOperationModeDescriptionListData:
		result.HvacOperationModeDescriptionListDataSelectors = castData[model.HvacOperationModeDescriptionListDataSelectorsType](data)
	case model.FunctionTypeHvacOverrunDescriptionListData:
		result.HvacOverrunDescriptionListDataSelectors = castData[model.HvacOverrunDescriptionListDataSelectorsType](data)
	case model.FunctionTypeHvacOverrunListData:
		result.HvacOverrunListDataSelectors = castData[model.HvacOverrunListDataSelectorsType](data)
	case model.FunctionTypeHvacSystemFunctionDescriptionListData:
		result.HvacSystemFunctionDescriptionListDataSelectors = castData[model.HvacSystemFunctionDescriptionListDataSelectorsType](data)
	case model.FunctionTypeHvacSystemFunctionListData:
		result.HvacSystemFunctionListDataSelectors = castData[model.HvacSystemFunctionListDataSelectorsType](data)
	case model.FunctionTypeHvacSystemFunctionOperationModeRelationListData:
		result.HvacSystemFunctionOperationModeRelationListDataSelectors = castData[model.HvacSystemFunctionOperationModeRelationListDataSelectorsType](data)
	case model.FunctionTypeHvacSystemFunctionPowerSequenceRelationListData:
		result.HvacSystemFunctionPowerSequenceRelationListDataSelectors = castData[model.HvacSystemFunctionPowerSequenceRelationListDataSelectorsType](data)
	case model.FunctionTypeHvacSystemFunctionSetPointRelationListData:
		result.HvacSystemFunctionSetpointRelationListDataSelectors = castData[model.HvacSystemFunctionSetpointRelationListDataSelectorsType](data)
	case model.FunctionTypeIdentificationListData:
		result.IdentificationListDataSelectors = castData[model.IdentificationListDataSelectorsType](data)
	case model.FunctionTypeIncentiveDescriptionListData:
		result.IncentiveDescriptionListDataSelectors = castData[model.IncentiveDescriptionListDataSelectorsType](data)
	case model.FunctionTypeIncentiveListData:
		result.IncentiveListDataSelectors = castData[model.IncentiveListDataSelectorsType](data)
	case model.FunctionTypeIncentiveTableConstraintsData:
		result.IncentiveTableConstraintsDataSelectors = castData[model.IncentiveTableConstraintsDataSelectorsType](data)
	case model.FunctionTypeIncentiveTableData:
		result.IncentiveTableDataSelectors = castData[model.IncentiveTableDataSelectorsType](data)
	case model.FunctionTypeIncentiveTableDescriptionData:
		result.IncentiveTableDescriptionDataSelectors = castData[model.IncentiveTableDescriptionDataSelectorsType](data)
	case model.FunctionTypeLoadControlEventListData:
		result.LoadControlEventListDataSelectors = castData[model.LoadControlEventListDataSelectorsType](data)
	case model.FunctionTypeLoadControlLimitConstraintsListData:
		result.LoadControlLimitConstraintsListDataSelectors = castData[model.LoadControlLimitConstraintsListDataSelectorsType](data)
	case model.FunctionTypeLoadControlLimitDescriptionListData:
		result.LoadControlLimitDescriptionListDataSelectors = castData[model.LoadControlLimitDescriptionListDataSelectorsType](data)
	case model.FunctionTypeLoadControlLimitListData:
		result.LoadControlLimitListDataSelectors = castData[model.LoadControlLimitListDataSelectorsType](data)
	case model.FunctionTypeLoadControlStateListData:
		result.LoadControlStateListDataSelectors = castData[model.LoadControlStateListDataSelectorsType](data)
	case model.FunctionTypeMeasurementConstraintsListData:
		result.MeasurementConstraintsListDataSelectors = castData[model.MeasurementConstraintsListDataSelectorsType](data)
	case model.FunctionTypeMeasurementDescriptionListData:
		result.MeasurementDescriptionListDataSelectors = castData[model.MeasurementDescriptionListDataSelectorsType](data)
	case model.FunctionTypeMeasurementListData:
		result.MeasurementListDataSelectors = castData[model.MeasurementListDataSelectorsType](data)
	case model.FunctionTypeMeasurementSeriesListData:
		result.MeasurementSeriesListDataSelectors = castData[model.MeasurementSeriesListDataSelectorsType](data)
	case model.FunctionTypeMeasurementThresholdRelationListData:
		result.MeasurementThresholdRelationListDataSelectors = castData[model.MeasurementThresholdRelationListDataSelectorsType](data)
	case model.FunctionTypeMessagingListData:
		result.MessagingListDataSelectors = castData[model.MessagingListDataSelectorsType](data)
	case model.FunctionTypeNetworkManagementDeviceDescriptionListData:
		result.NetworkManagementDeviceDescriptionListDataSelectors = castData[model.NetworkManagementDeviceDescriptionListDataSelectorsType](data)
	case model.FunctionTypeNetworkManagementEntityDescriptionListData:
		result.NetworkManagementEntityDescriptionListDataSelectors = castData[model.NetworkManagementEntityDescriptionListDataSelectorsType](data)
	case model.FunctionTypeNetworkManagementFeatureDescriptionListData:
		result.NetworkManagementFeatureDescriptionListDataSelectors = castData[model.NetworkManagementFeatureDescriptionListDataSelectorsType](data)
	case model.FunctionTypeNodeManagementBindingData:
		result.NodeManagementBindingDataSelectors = castData[model.NodeManagementBindingDataSelectorsType](data)
	case model.FunctionTypeNodeManagementDestinationListData:
		result.NodeManagementDestinationListDataSelectors = castData[model.NodeManagementDestinationListDataSelectorsType](data)
	case model.FunctionTypeNodeManagementDetailedDiscoveryData:
		result.NodeManagementDetailedDiscoveryDataSelectors = castData[model.NodeManagementDetailedDiscoveryDataSelectorsType](data)
	case model.FunctionTypeNodeManagementSubscriptionData:
		result.NodeManagementSubscriptionDataSelectors = castData[model.NodeManagementSubscriptionDataSelectorsType](data)
	case model.FunctionTypeNodeManagementUseCaseData:
		result.NodeManagementUseCaseDataSelectors = castData[model.NodeManagementUseCaseDataSelectorsType](data)
	case model.FunctionTypeOperatingConstraintsDurationListData:
		result.OperatingConstraintsDurationListDataSelectors = castData[model.OperatingConstraintsDurationListDataSelectorsType](data)
	case model.FunctionTypeOperatingConstraintsInterruptListData:
		result.OperatingConstraintsInterruptListDataSelectors = castData[model.OperatingConstraintsInterruptListDataSelectorsType](data)
	case model.FunctionTypeOperatingConstraintsPowerDescriptionListData:
		result.OperatingConstraintsPowerDescriptionListDataSelectors = castData[model.OperatingConstraintsPowerDescriptionListDataSelectorsType](data)
	case model.FunctionTypeOperatingConstraintsPowerLevelListData:
		result.OperatingConstraintsPowerLevelListDataSelectors = castData[model.OperatingConstraintsPowerLevelListDataSelectorsType](data)
	case model.FunctionTypeOperatingConstraintsPowerRangeListData:
		result.OperatingConstraintsPowerRangeListDataSelectors = castData[model.OperatingConstraintsPowerRangeListDataSelectorsType](data)
	case model.FunctionTypeOperatingConstraintsResumeImplicationListData:
		result.OperatingConstraintsResumeImplicationListDataSelectors = castData[model.OperatingConstraintsResumeImplicationListDataSelectorsType](data)
	case model.FunctionTypePowerSequenceAlternativesRelationListData:
		result.PowerSequenceAlternativesRelationListDataSelectors = castData[model.PowerSequenceAlternativesRelationListDataSelectorsType](data)
	case model.FunctionTypePowerSequenceDescriptionListData:
		result.PowerSequenceDescriptionListDataSelectors = castData[model.PowerSequenceDescriptionListDataSelectorsType](data)
	case model.FunctionTypePowerSequencePriceListData:
		result.PowerSequencePriceListDataSelectors = castData[model.PowerSequencePriceListDataSelectorsType](data)
	case model.FunctionTypePowerSequenceScheduleConstraintsListData:
		result.PowerSequenceScheduleConstraintsListDataSelectors = castData[model.PowerSequenceScheduleConstraintsListDataSelectorsType](data)
	case model.FunctionTypePowerSequenceScheduleListData:
		result.PowerSequenceScheduleListDataSelectors = castData[model.PowerSequenceScheduleListDataSelectorsType](data)
	case model.FunctionTypePowerSequenceSchedulePreferenceListData:
		result.PowerSequenceSchedulePreferenceListDataSelectors = castData[model.PowerSequenceSchedulePreferenceListDataSelectorsType](data)
	case model.FunctionTypePowerSequenceStateListData:
		result.PowerSequenceStateListDataSelectors = castData[model.PowerSequenceStateListDataSelectorsType](data)
	case model.FunctionTypePowerTimeSlotScheduleConstraintsListData:
		result.PowerTimeSlotScheduleConstraintsListDataSelectors = castData[model.PowerTimeSlotScheduleConstraintsListDataSelectorsType](data)
	case model.FunctionTypePowerTimeSlotScheduleListData:
		result.PowerTimeSlotScheduleListDataSelectors = castData[model.PowerTimeSlotScheduleListDataSelectorsType](data)
	case model.FunctionTypePowerTimeSlotValueListData:
		result.PowerTimeSlotValueListDataSelectors = castData[model.PowerTimeSlotValueListDataSelectorsType](data)
	case model.FunctionTypeSensingListData:
		result.SensingListDataSelectors = castData[model.SensingListDataSelectorsType](data)
	case model.FunctionTypeSessionIdentificationListData:
		result.SessionIdentificationListDataSelectors = castData[model.SessionIdentificationListDataSelectorsType](data)
	case model.FunctionTypeSessionMeasurementRelationListData:
		result.SessionMeasurementRelationListDataSelectors = castData[model.SessionMeasurementRelationListDataSelectorsType](data)
	case model.FunctionTypeSetpointConstraintsListData:
		result.SetpointConstraintsListDataSelectors = castData[model.SetpointConstraintsListDataSelectorsType](data)
	case model.FunctionTypeSetpointDescriptionListData:
		result.SetpointDescriptionListDataSelectors = castData[model.SetpointDescriptionListDataSelectorsType](data)
	case model.FunctionTypeSetpointListData:
		result.SetpointListDataSelectors = castData[model.SetpointListDataSelectorsType](data)
	case model.FunctionTypeSmartEnergyManagementPsData:
		result.SmartEnergyManagementPsDataSelectors = castData[model.SmartEnergyManagementPsDataSelectorsType](data)
	case model.FunctionTypeSmartEnergyManagementPsPriceData:
		result.SmartEnergyManagementPsPriceDataSelectors = castData[model.SmartEnergyManagementPsPriceDataSelectorsType](data)
	case model.FunctionTypeSpecificationVersionListData:
		result.SpecificationVersionListDataSelectors = castData[model.SpecificationVersionListDataSelectorsType](data)
	case model.FunctionTypeStateInformationListData:
		result.StateInformationListDataSelectors = castData[model.StateInformationListDataSelectorsType](data)
	case model.FunctionTypeSupplyConditionListData:
		result.SupplyConditionListDataSelectors = castData[model.SupplyConditionListDataSelectorsType](data)
	case model.FunctionTypeSupplyConditionThresholdRelationListData:
		result.SupplyConditionThresholdRelationListDataSelectors = castData[model.SupplyConditionThresholdRelationListDataSelectorsType](data)
	case model.FunctionTypeTariffBoundaryRelationListData:
		result.TariffBoundaryRelationListDataSelectors = castData[model.TariffBoundaryRelationListDataSelectorsType](data)
	case model.FunctionTypeTariffDescriptionListData:
		result.TariffDescriptionListDataSelectors = castData[model.TariffDescriptionListDataSelectorsType](data)
	case model.FunctionTypeTariffListData:
		result.TariffListDataSelectors = castData[model.TariffListDataSelectorsType](data)
	case model.FunctionTypeTariffTierRelationListData:
		result.TariffTierRelationListDataSelectors = castData[model.TariffTierRelationListDataSelectorsType](data)
	case model.FunctionTypeTaskManagementJobDescriptionListData:
		result.TaskManagementJobDescriptionListDataSelectors = castData[model.TaskManagementJobDescriptionListDataSelectorsType](data)
	case model.FunctionTypeTaskManagementJobListData:
		result.TaskManagementJobListDataSelectors = castData[model.TaskManagementJobListDataSelectorsType](data)
	case model.FunctionTypeTaskManagementJobRelationListData:
		result.TaskManagementJobRelationListDataSelectors = castData[model.TaskManagementJobRelationListDataSelectorsType](data)
	case model.FunctionTypeThresholdConstraintsListData:
		result.ThresholdConstraintsListDataSelectors = castData[model.ThresholdConstraintsListDataSelectorsType](data)
	case model.FunctionTypeThresholdDescriptionListData:
		result.ThresholdDescriptionListDataSelectors = castData[model.ThresholdDescriptionListDataSelectorsType](data)
	case model.FunctionTypeThresholdListData:
		result.ThresholdListDataSelectors = castData[model.ThresholdListDataSelectorsType](data)
	case model.FunctionTypeTierBoundaryDescriptionListData:
		result.TierBoundaryDescriptionListDataSelectors = castData[model.TierBoundaryDescriptionListDataSelectorsType](data)
	case model.FunctionTypeTierBoundaryListData:
		result.TierBoundaryListDataSelectors = castData[model.TierBoundaryListDataSelectorsType](data)
	case model.FunctionTypeTierDescriptionListData:
		result.TierDescriptionListDataSelectors = castData[model.TierDescriptionListDataSelectorsType](data)
	case model.FunctionTypeTierIncentiveRelationListData:
		result.TierIncentiveRelationListDataSelectors = castData[model.TierIncentiveRelationListDataSelectorsType](data)
	case model.FunctionTypeTierListData:
		result.TierListDataSelectors = castData[model.TierListDataSelectorsType](data)
	case model.FunctionTypeTimeSeriesConstraintsListData:
		result.TimeSeriesConstraintsListDataSelectors = castData[model.TimeSeriesConstraintsListDataSelectorsType](data)
	case model.FunctionTypeTimeSeriesDescriptionListData:
		result.TimeSeriesDescriptionListDataSelectors = castData[model.TimeSeriesDescriptionListDataSelectorsType](data)
	case model.FunctionTypeTimeSeriesListData:
		result.TimeSeriesListDataSelectors = castData[model.TimeSeriesListDataSelectorsType](data)
	case model.FunctionTypeTimeTableConstraintsListData:
		result.TimeTableConstraintsListDataSelectors = castData[model.TimeTableConstraintsListDataSelectorsType](data)
	case model.FunctionTypeTimeTableDescriptionListData:
		result.TimeTableDescriptionListDataSelectors = castData[model.TimeTableDescriptionListDataSelectorsType](data)
	case model.FunctionTypeTimeTableListData:
		result.TimeTableListDataSelectors = castData[model.TimeTableListDataSelectorsType](data)
	case model.FunctionTypeUseCaseInformationListData:
		result.UseCaseInformationListDataSelectors = castData[model.UseCaseInformationListDataSelectorsType](data)
	}

	return result
}

func addElementToFilter[T any](filter model.FilterType, function model.FunctionType, data *T) model.FilterType {
	result := filter

	switch function {
	case model.FunctionTypeActuatorLevelData:
		result.ActuatorLevelDataElements = castData[model.ActuatorLevelDataElementsType](data)
	case model.FunctionTypeActuatorLevelDescriptionData:
		result.ActuatorLevelDescriptionDataElements = castData[model.ActuatorLevelDescriptionDataElementsType](data)
	case model.FunctionTypeActuatorSwitchData:
		result.ActuatorSwitchDataElements = castData[model.ActuatorSwitchDataElementsType](data)
	case model.FunctionTypeActuatorSwitchDescriptionData:
		result.ActuatorSwitchDescriptionDataElements = castData[model.ActuatorSwitchDescriptionDataElementsType](data)
	case model.FunctionTypeAlarmListData:
		result.AlarmDataElements = castData[model.AlarmDataElementsType](data)
	case model.FunctionTypeBillConstraintsListData:
		result.BillConstraintsDataElements = castData[model.BillConstraintsDataElementsType](data)
	case model.FunctionTypeBillDescriptionListData:
		result.BillDescriptionDataElements = castData[model.BillDescriptionDataElementsType](data)
	case model.FunctionTypeBillListData:
		result.BillDataElements = castData[model.BillDataElementsType](data)
	case model.FunctionTypeBindingManagementDeleteCall:
		result.BindingManagementDeleteCallElements = castData[model.BindingManagementDeleteCallElementsType](data)
	case model.FunctionTypeBindingManagementEntryListData:
		result.BindingManagementEntryDataElements = castData[model.BindingManagementEntryDataElementsType](data)
	case model.FunctionTypeBindingManagementRequestCall:
		result.BindingManagementRequestCallElements = castData[model.BindingManagementRequestCallElementsType](data)
	case model.FunctionTypeCommodityListData:
		result.CommodityDataElements = castData[model.CommodityDataElementsType](data)
	case model.FunctionTypeDataTunnelingCall:
		result.DataTunnelingCallElements = castData[model.DataTunnelingCallElementsType](data)
	case model.FunctionTypeDeviceClassificationManufacturerData:
		result.DeviceClassificationManufacturerDataElements = castData[model.DeviceClassificationManufacturerDataElementsType](data)
	case model.FunctionTypeDeviceClassificationUserData:
		result.DeviceClassificationUserDataElements = castData[model.DeviceClassificationUserDataElementsType](data)
	case model.FunctionTypeDeviceConfigurationKeyValueConstraintsListData:
		result.DeviceConfigurationKeyValueConstraintsDataElements = castData[model.DeviceConfigurationKeyValueConstraintsDataElementsType](data)
	case model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData:
		result.DeviceConfigurationKeyValueDescriptionDataElements = castData[model.DeviceConfigurationKeyValueDescriptionDataElementsType](data)
	case model.FunctionTypeDeviceConfigurationKeyValueListData:
		result.DeviceConfigurationKeyValueDataElements = castData[model.DeviceConfigurationKeyValueDataElementsType](data)
	case model.FunctionTypeDeviceDiagnosisHeartbeatData:
		result.DeviceDiagnosisHeartbeatDataElements = castData[model.DeviceDiagnosisHeartbeatDataElementsType](data)
	case model.FunctionTypeDeviceDiagnosisServiceData:
		result.DeviceDiagnosisServiceDataElements = castData[model.DeviceDiagnosisServiceDataElementsType](data)
	case model.FunctionTypeDeviceDiagnosisStateData:
		result.DeviceDiagnosisStateDataElements = castData[model.DeviceDiagnosisStateDataElementsType](data)
	case model.FunctionTypeDirectControlActivityListData:
		result.DirectControlActivityDataElements = castData[model.DirectControlActivityDataElementsType](data)
	case model.FunctionTypeDirectControlDescriptionData:
		result.DirectControlDescriptionDataElements = castData[model.DirectControlDescriptionDataElementsType](data)
	case model.FunctionTypeElectricalConnectionDescriptionListData:
		result.ElectricalConnectionDescriptionDataElements = castData[model.ElectricalConnectionDescriptionDataElementsType](data)
	case model.FunctionTypeElectricalConnectionParameterDescriptionListData:
		result.ElectricalConnectionParameterDescriptionDataElements = castData[model.ElectricalConnectionParameterDescriptionDataElementsType](data)
	case model.FunctionTypeElectricalConnectionPermittedValueSetListData:
		result.ElectricalConnectionPermittedValueSetDataElements = castData[model.ElectricalConnectionPermittedValueSetDataElementsType](data)
	case model.FunctionTypeElectricalConnectionStateListData:
		result.ElectricalConnectionStateDataElements = castData[model.ElectricalConnectionStateDataElementsType](data)
	case model.FunctionTypeElectricalConnectionCharacteristicListData:
		result.ElectricalConnectionCharacteristicDataElements = castData[model.ElectricalConnectionCharacteristicDataElementsType](data)
	case model.FunctionTypeHvacOperationModeDescriptionListData:
		result.HvacOperationModeDescriptionDataElements = castData[model.HvacOperationModeDescriptionDataElementsType](data)
	case model.FunctionTypeHvacOverrunDescriptionListData:
		result.HvacOverrunDescriptionDataElements = castData[model.HvacOverrunDescriptionDataElementsType](data)
	case model.FunctionTypeHvacOverrunListData:
		result.HvacOverrunDataElements = castData[model.HvacOverrunDataElementsType](data)
	case model.FunctionTypeHvacSystemFunctionDescriptionListData:
		result.HvacSystemFunctionDescriptionDataElements = castData[model.HvacSystemFunctionDescriptionDataElementsType](data)
	case model.FunctionTypeHvacSystemFunctionListData:
		result.HvacSystemFunctionDataElements = castData[model.HvacSystemFunctionDataElementsType](data)
	case model.FunctionTypeHvacSystemFunctionOperationModeRelationListData:
		result.HvacSystemFunctionOperationModeRelationDataElements = castData[model.HvacSystemFunctionOperationModeRelationDataElementsType](data)
	case model.FunctionTypeHvacSystemFunctionPowerSequenceRelationListData:
		result.HvacSystemFunctionPowerSequenceRelationDataElements = castData[model.HvacSystemFunctionPowerSequenceRelationDataElementsType](data)
	case model.FunctionTypeHvacSystemFunctionSetPointRelationListData:
		result.HvacSystemFunctionSetpointRelationDataElements = castData[model.HvacSystemFunctionSetpointRelationDataElementsType](data)
	case model.FunctionTypeIdentificationListData:
		result.IdentificationDataElements = castData[model.IdentificationDataElementsType](data)
	case model.FunctionTypeIncentiveDescriptionListData:
		result.IncentiveDescriptionDataElements = castData[model.IncentiveDescriptionDataElementsType](data)
	case model.FunctionTypeIncentiveListData:
		result.IncentiveDataElements = castData[model.IncentiveDataElementsType](data)
	case model.FunctionTypeIncentiveTableConstraintsData:
		result.IncentiveTableConstraintsDataElements = castData[model.IncentiveTableConstraintsDataElementsType](data)
	case model.FunctionTypeIncentiveTableData:
		result.IncentiveTableDataElements = castData[model.IncentiveTableDataElementsType](data)
	case model.FunctionTypeIncentiveTableDescriptionData:
		result.IncentiveTableDescriptionDataElements = castData[model.IncentiveTableDescriptionDataElementsType](data)
	case model.FunctionTypeLoadControlEventListData:
		result.LoadControlEventDataElements = castData[model.LoadControlEventDataElementsType](data)
	case model.FunctionTypeLoadControlLimitConstraintsListData:
		result.LoadControlLimitConstraintsDataElements = castData[model.LoadControlLimitConstraintsDataElementsType](data)
	case model.FunctionTypeLoadControlLimitDescriptionListData:
		result.LoadControlLimitDescriptionDataElements = castData[model.LoadControlLimitDescriptionDataElementsType](data)
	case model.FunctionTypeLoadControlLimitListData:
		result.LoadControlLimitDataElements = castData[model.LoadControlLimitDataElementsType](data)
	case model.FunctionTypeLoadControlNodeData:
		result.LoadControlNodeDataElements = castData[model.LoadControlNodeDataElementsType](data)
	case model.FunctionTypeLoadControlStateListData:
		result.LoadControlStateDataElements = castData[model.LoadControlStateDataElementsType](data)
	case model.FunctionTypeMeasurementConstraintsListData:
		result.MeasurementConstraintsDataElements = castData[model.MeasurementConstraintsDataElementsType](data)
	case model.FunctionTypeMeasurementDescriptionListData:
		result.MeasurementDescriptionDataElements = castData[model.MeasurementDescriptionDataElementsType](data)
	case model.FunctionTypeMeasurementListData:
		result.MeasurementDataElements = castData[model.MeasurementDataElementsType](data)
	case model.FunctionTypeMeasurementSeriesListData:
		result.MeasurementSeriesDataElements = castData[model.MeasurementSeriesDataElementsType](data)
	case model.FunctionTypeMeasurementThresholdRelationListData:
		result.MeasurementThresholdRelationDataElements = castData[model.MeasurementThresholdRelationDataElementsType](data)
	case model.FunctionTypeMessagingListData:
		result.MessagingDataElements = castData[model.MessagingDataElementsType](data)
	case model.FunctionTypeNetworkManagementAbortCall:
		result.NetworkManagementAbortCallElements = castData[model.NetworkManagementAbortCallElementsType](data)
	case model.FunctionTypeNetworkManagementAddNodeCall:
		result.NetworkManagementAddNodeCallElements = castData[model.NetworkManagementAddNodeCallElementsType](data)
	case model.FunctionTypeNetworkManagementDeviceDescriptionListData:
		result.NetworkManagementDeviceDescriptionDataElements = castData[model.NetworkManagementDeviceDescriptionDataElementsType](data)
	case model.FunctionTypeNetworkManagementDiscoverCall:
		result.NetworkManagementDiscoverCallElements = castData[model.NetworkManagementDiscoverCallElementsType](data)
	case model.FunctionTypeNetworkManagementEntityDescriptionListData:
		result.NetworkManagementEntityDescriptionDataElements = castData[model.NetworkManagementEntityDescriptionDataElementsType](data)
	case model.FunctionTypeNetworkManagementFeatureDescriptionListData:
		result.NetworkManagementFeatureDescriptionDataElements = castData[model.NetworkManagementFeatureDescriptionDataElementsType](data)
	case model.FunctionTypeNetworkManagementJoiningModeData:
		result.NetworkManagementJoiningModeDataElements = castData[model.NetworkManagementJoiningModeDataElementsType](data)
	case model.FunctionTypeNetworkManagementModifyNodeCall:
		result.NetworkManagementModifyNodeCallElements = castData[model.NetworkManagementModifyNodeCallElementsType](data)
	case model.FunctionTypeNetworkManagementProcessStateData:
		result.NetworkManagementProcessStateDataElements = castData[model.NetworkManagementProcessStateDataElementsType](data)
	case model.FunctionTypeNetworkManagementRemoveNodeCall:
		result.NetworkManagementRemoveNodeCallElements = castData[model.NetworkManagementRemoveNodeCallElementsType](data)
	case model.FunctionTypeNetworkManagementReportCandidateData:
		result.NetworkManagementReportCandidateDataElements = castData[model.NetworkManagementReportCandidateDataElementsType](data)
	case model.FunctionTypeNetworkManagementScanNetworkCall:
		result.NetworkManagementScanNetworkCallElements = castData[model.NetworkManagementScanNetworkCallElementsType](data)
	case model.FunctionTypeNodeManagementBindingData:
		result.NodeManagementBindingDataElements = castData[model.NodeManagementBindingDataElementsType](data)
	case model.FunctionTypeNodeManagementBindingDeleteCall:
		result.NodeManagementBindingDeleteCallElements = castData[model.NodeManagementBindingDeleteCallElementsType](data)
	case model.FunctionTypeNodeManagementBindingRequestCall:
		result.NodeManagementBindingRequestCallElements = castData[model.NodeManagementBindingRequestCallElementsType](data)
	case model.FunctionTypeNodeManagementDestinationListData:
		result.NodeManagementDestinationDataElements = castData[model.NodeManagementDestinationDataElementsType](data)
	case model.FunctionTypeNodeManagementDetailedDiscoveryData:
		result.NodeManagementDetailedDiscoveryDataElements = castData[model.NodeManagementDetailedDiscoveryDataElementsType](data)
	case model.FunctionTypeNodeManagementSubscriptionData:
		result.NodeManagementSubscriptionDataElements = castData[model.NodeManagementSubscriptionDataElementsType](data)
	case model.FunctionTypeNodeManagementSubscriptionDeleteCall:
		result.NodeManagementSubscriptionDeleteCallElements = castData[model.NodeManagementSubscriptionDeleteCallElementsType](data)
	case model.FunctionTypeNodeManagementSubscriptionRequestCall:
		result.NodeManagementSubscriptionRequestCallElements = castData[model.NodeManagementSubscriptionRequestCallElementsType](data)
	case model.FunctionTypeNodeManagementUseCaseData:
		result.NodeManagementUseCaseDataElements = castData[model.NodeManagementUseCaseDataElementsType](data)
	case model.FunctionTypeOperatingConstraintsDurationListData:
		result.OperatingConstraintsDurationDataElements = castData[model.OperatingConstraintsDurationDataElementsType](data)
	case model.FunctionTypeOperatingConstraintsInterruptListData:
		result.OperatingConstraintsInterruptDataElements = castData[model.OperatingConstraintsInterruptDataElementsType](data)
	case model.FunctionTypeOperatingConstraintsPowerDescriptionListData:
		result.OperatingConstraintsPowerDescriptionDataElements = castData[model.OperatingConstraintsPowerDescriptionDataElementsType](data)
	case model.FunctionTypeOperatingConstraintsPowerLevelListData:
		result.OperatingConstraintsPowerLevelDataElements = castData[model.OperatingConstraintsPowerLevelDataElementsType](data)
	case model.FunctionTypeOperatingConstraintsPowerRangeListData:
		result.OperatingConstraintsPowerRangeDataElements = castData[model.OperatingConstraintsPowerRangeDataElementsType](data)
	case model.FunctionTypeOperatingConstraintsResumeImplicationListData:
		result.OperatingConstraintsResumeImplicationDataElements = castData[model.OperatingConstraintsResumeImplicationDataElementsType](data)
	case model.FunctionTypePowerSequenceAlternativesRelationListData:
		result.PowerSequenceAlternativesRelationDataElements = castData[model.PowerSequenceAlternativesRelationDataElementsType](data)
	case model.FunctionTypePowerSequenceDescriptionListData:
		result.PowerSequenceDescriptionDataElements = castData[model.PowerSequenceDescriptionDataElementsType](data)
	case model.FunctionTypePowerSequenceNodeScheduleInformationData:
		result.PowerSequenceNodeScheduleInformationDataElements = castData[model.PowerSequenceNodeScheduleInformationDataElementsType](data)
	case model.FunctionTypePowerSequencePriceCalculationRequestCall:
		result.PowerSequencePriceCalculationRequestCallElements = castData[model.PowerSequencePriceCalculationRequestCallElementsType](data)
	case model.FunctionTypePowerSequencePriceListData:
		result.PowerSequencePriceDataElements = castData[model.PowerSequencePriceDataElementsType](data)
	case model.FunctionTypePowerSequenceScheduleConfigurationRequestCall:
		result.PowerSequenceScheduleConfigurationRequestCallElements = castData[model.PowerSequenceScheduleConfigurationRequestCallElementsType](data)
	case model.FunctionTypePowerSequenceScheduleConstraintsListData:
		result.PowerSequenceScheduleConstraintsDataElements = castData[model.PowerSequenceScheduleConstraintsDataElementsType](data)
	case model.FunctionTypePowerSequenceScheduleListData:
		result.PowerSequenceScheduleDataElements = castData[model.PowerSequenceScheduleDataElementsType](data)
	case model.FunctionTypePowerSequenceSchedulePreferenceListData:
		result.PowerSequenceSchedulePreferenceDataElements = castData[model.PowerSequenceSchedulePreferenceDataElementsType](data)
	case model.FunctionTypePowerSequenceStateListData:
		result.PowerSequenceStateDataElements = castData[model.PowerSequenceStateDataElementsType](data)
	case model.FunctionTypePowerTimeSlotScheduleConstraintsListData:
		result.PowerTimeSlotScheduleConstraintsDataElements = castData[model.PowerTimeSlotScheduleConstraintsDataElementsType](data)
	case model.FunctionTypePowerTimeSlotScheduleListData:
		result.PowerTimeSlotScheduleDataElements = castData[model.PowerTimeSlotScheduleDataElementsType](data)
	case model.FunctionTypePowerTimeSlotValueListData:
		result.PowerTimeSlotValueDataElements = castData[model.PowerTimeSlotValueDataElementsType](data)
	case model.FunctionTypeSensingListData:
		result.SensingDataElements = castData[model.SensingDataElementsType](data)
	case model.FunctionTypeSessionIdentificationListData:
		result.SessionIdentificationDataElements = castData[model.SessionIdentificationDataElementsType](data)
	case model.FunctionTypeSessionMeasurementRelationListData:
		result.SessionMeasurementRelationDataElements = castData[model.SessionMeasurementRelationDataElementsType](data)
	case model.FunctionTypeSetpointConstraintsListData:
		result.SetpointConstraintsDataElements = castData[model.SetpointConstraintsDataElementsType](data)
	case model.FunctionTypeSetpointDescriptionListData:
		result.SensingDescriptionDataElements = castData[model.SensingDescriptionDataElementsType](data)
	case model.FunctionTypeSetpointListData:
		result.SetpointDataElements = castData[model.SetpointDataElementsType](data)
	case model.FunctionTypeSmartEnergyManagementPsConfigurationRequestCall:
		result.SmartEnergyManagementPsConfigurationRequestCallElements = castData[model.SmartEnergyManagementPsConfigurationRequestCallElementsType](data)
	case model.FunctionTypeSmartEnergyManagementPsData:
		result.SmartEnergyManagementPsDataElements = castData[model.SmartEnergyManagementPsDataElementsType](data)
	case model.FunctionTypeSmartEnergyManagementPsPriceCalculationRequestCall:
		result.SmartEnergyManagementPsPriceCalculationRequestCallElements = castData[model.SmartEnergyManagementPsPriceCalculationRequestCallElementsType](data)
	case model.FunctionTypeSmartEnergyManagementPsPriceData:
		result.SmartEnergyManagementPsPriceDataElements = castData[model.SmartEnergyManagementPsPriceDataElementsType](data)
	case model.FunctionTypeSpecificationVersionListData:
		result.SpecificationVersionDataElements = castData[model.SpecificationVersionDataElementsType](data)
	case model.FunctionTypeStateInformationListData:
		result.StateInformationDataElements = castData[model.StateInformationDataElementsType](data)
	case model.FunctionTypeSubscriptionManagementDeleteCall:
		result.SubscriptionManagementDeleteCallElements = castData[model.SubscriptionManagementDeleteCallElementsType](data)
	case model.FunctionTypeSubscriptionManagementEntryListData:
		result.SubscriptionManagementEntryDataElements = castData[model.SubscriptionManagementEntryDataElementsType](data)
	case model.FunctionTypeSubscriptionManagementRequestCall:
		result.SubscriptionManagementRequestCallElements = castData[model.SubscriptionManagementRequestCallElementsType](data)
	case model.FunctionTypeSupplyConditionListData:
		result.SupplyConditionDataElements = castData[model.SupplyConditionDataElementsType](data)
	case model.FunctionTypeSupplyConditionDescriptionListData:
		result.SupplyConditionDescriptionDataElements = castData[model.SupplyConditionDescriptionDataElementsType](data)
	case model.FunctionTypeSupplyConditionThresholdRelationListData:
		result.SupplyConditionThresholdRelationDataElements = castData[model.SupplyConditionThresholdRelationDataElementsType](data)
	case model.FunctionTypeTariffBoundaryRelationListData:
		result.TariffBoundaryRelationDataElements = castData[model.TariffBoundaryRelationDataElementsType](data)
	case model.FunctionTypeTariffDescriptionListData:
		result.TariffDescriptionDataElements = castData[model.TariffDescriptionDataElementsType](data)
	case model.FunctionTypeTariffListData:
		result.TariffDataElements = castData[model.TariffDataElementsType](data)
	case model.FunctionTypeTariffOverallConstraintsData:
		result.TariffOverallConstraintsDataElements = castData[model.TariffOverallConstraintsDataElementsType](data)
	case model.FunctionTypeTariffTierRelationListData:
		result.TariffTierRelationDataElements = castData[model.TariffTierRelationDataElementsType](data)
	case model.FunctionTypeTaskManagementJobDescriptionListData:
		result.TaskManagementJobDescriptionDataElements = castData[model.TaskManagementJobDescriptionDataElementsType](data)
	case model.FunctionTypeTaskManagementJobListData:
		result.TaskManagementJobDataElements = castData[model.TaskManagementJobDataElementsType](data)
	case model.FunctionTypeTaskManagementJobRelationListData:
		result.TaskManagementJobRelationDataElements = castData[model.TaskManagementJobRelationDataElementsType](data)
	case model.FunctionTypeTaskManagementOverviewData:
		result.TaskManagementOverviewDataElements = castData[model.TaskManagementOverviewDataElementsType](data)
	case model.FunctionTypeThresholdConstraintsListData:
		result.ThresholdConstraintsDataElements = castData[model.ThresholdConstraintsDataElementsType](data)
	case model.FunctionTypeThresholdDescriptionListData:
		result.ThresholdDescriptionDataElements = castData[model.ThresholdDescriptionDataElementsType](data)
	case model.FunctionTypeThresholdListData:
		result.ThresholdDataElements = castData[model.ThresholdDataElementsType](data)
	case model.FunctionTypeTierBoundaryDescriptionListData:
		result.TierBoundaryDescriptionDataElements = castData[model.TierBoundaryDescriptionDataElementsType](data)
	case model.FunctionTypeTierBoundaryListData:
		result.TierBoundaryDataElements = castData[model.TierBoundaryDataElementsType](data)
	case model.FunctionTypeTierDescriptionListData:
		result.TierDescriptionDataElements = castData[model.TierDescriptionDataElementsType](data)
	case model.FunctionTypeTierIncentiveRelationListData:
		result.TierIncentiveRelationDataElements = castData[model.TierIncentiveRelationDataElementsType](data)
	case model.FunctionTypeTierListData:
		result.TierDataElements = castData[model.TierDataElementsType](data)
	case model.FunctionTypeTimeDistributorData:
		result.TimeDistributorDataElements = castData[model.TimeDistributorDataElementsType](data)
	case model.FunctionTypeTimeDistributorEnquiryCall:
		result.TimeDistributorEnquiryCallElements = castData[model.TimeDistributorEnquiryCallElementsType](data)
	case model.FunctionTypeTimeInformationData:
		result.TimeInformationDataElements = castData[model.TimeInformationDataElementsType](data)
	case model.FunctionTypeTimePrecisionData:
		result.TimePrecisionDataElements = castData[model.TimePrecisionDataElementsType](data)
	case model.FunctionTypeTimeSeriesConstraintsListData:
		result.TimeSeriesConstraintsDataElements = castData[model.TimeSeriesConstraintsDataElementsType](data)
	case model.FunctionTypeTimeSeriesDescriptionListData:
		result.TimeSeriesDescriptionDataElements = castData[model.TimeSeriesDescriptionDataElementsType](data)
	case model.FunctionTypeTimeSeriesListData:
		result.TimeSeriesDataElements = castData[model.TimeSeriesDataElementsType](data)
	case model.FunctionTypeTimeTableConstraintsListData:
		result.TimeTableConstraintsDataElements = castData[model.TimeTableConstraintsDataElementsType](data)
	case model.FunctionTypeTimeTableDescriptionListData:
		result.TimeTableDescriptionDataElements = castData[model.TimeTableDescriptionDataElementsType](data)
	case model.FunctionTypeTimeTableListData:
		result.TimeTableDataElements = castData[model.TimeTableDataElementsType](data)
	case model.FunctionTypeUseCaseInformationListData:
		result.UseCaseInformationDataElements = castData[model.UseCaseInformationDataElementsType](data)
	}

	return result
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
	case model.FunctionTypeElectricalConnectionCharacteristicListData:
		result.ElectricalConnectionCharacteristicListData = castData[model.ElectricalConnectionCharacteristicListDataType](data)
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
	case model.FunctionTypeMeasurementSeriesListData:
		result.MeasurementSeriesListData = castData[model.MeasurementSeriesListDataType](data)
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
	case model.FunctionTypeSessionIdentificationListData:
		result.SessionIdentificationListData = castData[model.SessionIdentificationListDataType](data)
	case model.FunctionTypeSessionMeasurementRelationListData:
		result.SessionMeasurementRelationListData = castData[model.SessionMeasurementRelationListDataType](data)
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
	case model.FunctionTypeStateInformationListData:
		result.StateInformationListData = castData[model.StateInformationListDataType](data)
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
