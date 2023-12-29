package spine

import (
	"fmt"

	"github.com/enbility/eebus-go/spine/model"
)

func CreateFunctionData[F any](featureType model.FeatureTypeType) []F {
	if featureType == model.FeatureTypeTypeNodeManagement {
		return []F{} // NodeManagement implementation is not using function data
	}

	// Some devices use generic for everything (e.g. Vaillant Arotherm heatpump)
	// or for some things like the SMA HM 2.0 or Elli Wallbox, which uses Generic feature
	// for Heartbeats, even though that should go into FeatureTypeTypeDeviceDiagnosis
	// Hence we add everything to the Generic feature, as we don't know what might be needed

	var result []F

	if featureType == model.FeatureTypeTypeActuatorLevel || featureType == model.FeatureTypeTypeGeneric {
		result = append(result, []F{
			createFunctionData[model.ActuatorLevelDataType, F](model.FunctionTypeActuatorLevelData),
			createFunctionData[model.ActuatorLevelDescriptionDataType, F](model.FunctionTypeActuatorLevelDescriptionData),
		}...)
	}

	if featureType == model.FeatureTypeTypeActuatorSwitch || featureType == model.FeatureTypeTypeGeneric {
		result = append(result, []F{
			createFunctionData[model.ActuatorSwitchDataType, F](model.FunctionTypeActuatorSwitchData),
			createFunctionData[model.ActuatorSwitchDescriptionDataType, F](model.FunctionTypeActuatorSwitchDescriptionData),
		}...)
	}

	if featureType == model.FeatureTypeTypeAlarm || featureType == model.FeatureTypeTypeGeneric {
		result = append(result, []F{
			createFunctionData[model.AlarmListDataType, F](model.FunctionTypeAlarmListData),
		}...)
	}

	if featureType == model.FeatureTypeTypeBill || featureType == model.FeatureTypeTypeGeneric {
		result = append(result, []F{
			createFunctionData[model.BillDescriptionListDataType, F](model.FunctionTypeBillDescriptionListData),
			createFunctionData[model.BillConstraintsListDataType, F](model.FunctionTypeBillConstraintsListData),
			createFunctionData[model.BillListDataType, F](model.FunctionTypeBillListData),
		}...)
	}

	if featureType == model.FeatureTypeTypeDataTunneling || featureType == model.FeatureTypeTypeGeneric {
		result = append(result, []F{
			createFunctionData[model.DataTunnelingCallType, F](model.FunctionTypeDataTunnelingCall),
		}...)
	}

	if featureType == model.FeatureTypeTypeDeviceClassification || featureType == model.FeatureTypeTypeGeneric {
		result = append(result, []F{
			createFunctionData[model.DeviceClassificationManufacturerDataType, F](model.FunctionTypeDeviceClassificationManufacturerData),
			createFunctionData[model.DeviceClassificationUserDataType, F](model.FunctionTypeDeviceClassificationUserData),
		}...)
	}

	if featureType == model.FeatureTypeTypeDeviceConfiguration || featureType == model.FeatureTypeTypeGeneric {
		result = append(result, []F{
			createFunctionData[model.DeviceConfigurationKeyValueConstraintsListDataType, F](model.FunctionTypeDeviceConfigurationKeyValueConstraintsListData),
			createFunctionData[model.DeviceConfigurationKeyValueDescriptionListDataType, F](model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData),
			createFunctionData[model.DeviceConfigurationKeyValueListDataType, F](model.FunctionTypeDeviceConfigurationKeyValueListData),
		}...)
	}

	if featureType == model.FeatureTypeTypeDeviceDiagnosis || featureType == model.FeatureTypeTypeGeneric {
		result = append(result, []F{
			createFunctionData[model.DeviceDiagnosisStateDataType, F](model.FunctionTypeDeviceDiagnosisStateData),
			createFunctionData[model.DeviceDiagnosisHeartbeatDataType, F](model.FunctionTypeDeviceDiagnosisHeartbeatData),
			createFunctionData[model.DeviceDiagnosisServiceDataType, F](model.FunctionTypeDeviceDiagnosisServiceData),
		}...)
	}

	if featureType == model.FeatureTypeTypeDirectControl || featureType == model.FeatureTypeTypeGeneric {
		result = append(result, []F{
			createFunctionData[model.DirectControlActivityListDataType, F](model.FunctionTypeDirectControlActivityListData),
			createFunctionData[model.DirectControlDescriptionDataType, F](model.FunctionTypeDirectControlDescriptionData),
		}...)
	}

	if featureType == model.FeatureTypeTypeElectricalConnection || featureType == model.FeatureTypeTypeGeneric {
		result = append(result, []F{
			createFunctionData[model.ElectricalConnectionDescriptionListDataType, F](model.FunctionTypeElectricalConnectionDescriptionListData),
			createFunctionData[model.ElectricalConnectionParameterDescriptionListDataType, F](model.FunctionTypeElectricalConnectionParameterDescriptionListData),
			createFunctionData[model.ElectricalConnectionPermittedValueSetListDataType, F](model.FunctionTypeElectricalConnectionPermittedValueSetListData),
			createFunctionData[model.ElectricalConnectionStateListDataType, F](model.FunctionTypeElectricalConnectionStateListData),
			createFunctionData[model.ElectricalConnectionCharacteristicListDataType, F](model.FunctionTypeElectricalConnectionCharacteristicListData),
		}...)
	}

	if featureType == model.FeatureTypeTypeHvac || featureType == model.FeatureTypeTypeGeneric {
		result = append(result, []F{
			createFunctionData[model.HvacOperationModeDescriptionDataType, F](model.FunctionTypeHvacOperationModeDescriptionListData),
			createFunctionData[model.HvacOverrunDescriptionListDataType, F](model.FunctionTypeHvacOverrunDescriptionListData),
			createFunctionData[model.HvacOverrunListDataType, F](model.FunctionTypeHvacOverrunListData),
			createFunctionData[model.HvacSystemFunctionDescriptionDataType, F](model.FunctionTypeHvacSystemFunctionDescriptionListData),
			createFunctionData[model.HvacSystemFunctionListDataType, F](model.FunctionTypeHvacSystemFunctionListData),
			createFunctionData[model.HvacSystemFunctionOperationModeRelationListDataType, F](model.FunctionTypeHvacSystemFunctionOperationModeRelationListData),
			createFunctionData[model.HvacSystemFunctionPowerSequenceRelationListDataType, F](model.FunctionTypeHvacSystemFunctionPowerSequenceRelationListData),
			createFunctionData[model.HvacSystemFunctionSetpointRelationListDataType, F](model.FunctionTypeHvacSystemFunctionSetPointRelationListData),
		}...)
	}

	if featureType == model.FeatureTypeTypeIdentification || featureType == model.FeatureTypeTypeGeneric {
		result = append(result, []F{
			createFunctionData[model.IdentificationListDataType, F](model.FunctionTypeIdentificationListData),
			createFunctionData[model.SessionIdentificationListDataType, F](model.FunctionTypeSessionIdentificationListData),
			createFunctionData[model.SessionMeasurementRelationListDataType, F](model.FunctionTypeSessionMeasurementRelationListData),
		}...)
	}

	if featureType == model.FeatureTypeTypeIncentiveTable || featureType == model.FeatureTypeTypeGeneric {
		result = append(result, []F{
			createFunctionData[model.IncentiveTableDescriptionDataType, F](model.FunctionTypeIncentiveTableDescriptionData),
			createFunctionData[model.IncentiveTableConstraintsDataType, F](model.FunctionTypeIncentiveTableConstraintsData),
			createFunctionData[model.IncentiveTableDataType, F](model.FunctionTypeIncentiveTableData),
		}...)
	}

	if featureType == model.FeatureTypeTypeLoadControl || featureType == model.FeatureTypeTypeGeneric {
		result = append(result, []F{
			createFunctionData[model.LoadControlEventListDataType, F](model.FunctionTypeLoadControlEventListData),
			createFunctionData[model.LoadControlLimitConstraintsListDataType, F](model.FunctionTypeLoadControlLimitConstraintsListData),
			createFunctionData[model.LoadControlLimitDescriptionListDataType, F](model.FunctionTypeLoadControlLimitDescriptionListData),
			createFunctionData[model.LoadControlLimitListDataType, F](model.FunctionTypeLoadControlLimitListData),
			createFunctionData[model.LoadControlNodeDataType, F](model.FunctionTypeLoadControlNodeData),
			createFunctionData[model.LoadControlStateListDataType, F](model.FunctionTypeLoadControlStateListData),
		}...)
	}

	if featureType == model.FeatureTypeTypeMeasurement || featureType == model.FeatureTypeTypeGeneric {
		result = append(result, []F{
			createFunctionData[model.MeasurementListDataType, F](model.FunctionTypeMeasurementListData),
			createFunctionData[model.MeasurementDescriptionListDataType, F](model.FunctionTypeMeasurementDescriptionListData),
			createFunctionData[model.MeasurementConstraintsListDataType, F](model.FunctionTypeMeasurementConstraintsListData),
			createFunctionData[model.MeasurementThresholdRelationListDataType, F](model.FunctionTypeMeasurementThresholdRelationListData),
			createFunctionData[model.MeasurementSeriesListDataType, F](model.FunctionTypeMeasurementSeriesListData),
		}...)
	}

	if featureType == model.FeatureTypeTypeMessaging || featureType == model.FeatureTypeTypeGeneric {
		result = append(result, []F{
			createFunctionData[model.MessagingListDataType, F](model.FunctionTypeMessagingListData),
		}...)
	}

	if featureType == model.FeatureTypeTypeNetworkManagement || featureType == model.FeatureTypeTypeGeneric {
		result = append(result, []F{
			createFunctionData[model.NetworkManagementAbortCallType, F](model.FunctionTypeNetworkManagementAbortCall),
			createFunctionData[model.NetworkManagementAddNodeCallType, F](model.FunctionTypeNetworkManagementAddNodeCall),
			createFunctionData[model.NetworkManagementDeviceDescriptionListDataType, F](model.FunctionTypeNetworkManagementDeviceDescriptionListData),
			createFunctionData[model.NetworkManagementDiscoverCallType, F](model.FunctionTypeNetworkManagementDiscoverCall),
			createFunctionData[model.NetworkManagementEntityDescriptionListDataType, F](model.FunctionTypeNetworkManagementEntityDescriptionListData),
			createFunctionData[model.NetworkManagementFeatureDescriptionListDataType, F](model.FunctionTypeNetworkManagementFeatureDescriptionListData),
			createFunctionData[model.NetworkManagementJoiningModeDataType, F](model.FunctionTypeNetworkManagementJoiningModeData),
			createFunctionData[model.NetworkManagementModifyNodeCallType, F](model.FunctionTypeNetworkManagementModifyNodeCall),
			createFunctionData[model.NetworkManagementProcessStateDataType, F](model.FunctionTypeNetworkManagementProcessStateData),
			createFunctionData[model.NetworkManagementRemoveNodeCallType, F](model.FunctionTypeNetworkManagementRemoveNodeCall),
			createFunctionData[model.NetworkManagementReportCandidateDataType, F](model.FunctionTypeNetworkManagementReportCandidateData),
			createFunctionData[model.NetworkManagementScanNetworkCallType, F](model.FunctionTypeNetworkManagementScanNetworkCall),
		}...)
	}

	if featureType == model.FeatureTypeTypeOperatingConstraints || featureType == model.FeatureTypeTypeGeneric {
		result = append(result, []F{
			createFunctionData[model.OperatingConstraintsDurationListDataType, F](model.FunctionTypeOperatingConstraintsDurationListData),
			createFunctionData[model.OperatingConstraintsInterruptListDataType, F](model.FunctionTypeOperatingConstraintsInterruptListData),
			createFunctionData[model.OperatingConstraintsPowerDescriptionListDataType, F](model.FunctionTypeOperatingConstraintsPowerDescriptionListData),
			createFunctionData[model.OperatingConstraintsPowerLevelListDataType, F](model.FunctionTypeOperatingConstraintsPowerLevelListData),
			createFunctionData[model.OperatingConstraintsPowerRangeListDataType, F](model.FunctionTypeOperatingConstraintsPowerRangeListData),
			createFunctionData[model.OperatingConstraintsResumeImplicationListDataType, F](model.FunctionTypeOperatingConstraintsResumeImplicationListData),
		}...)
	}

	if featureType == model.FeatureTypeTypePowerSequences || featureType == model.FeatureTypeTypeGeneric {
		result = append(result, []F{
			createFunctionData[model.PowerSequenceAlternativesRelationListDataType, F](model.FunctionTypePowerSequenceAlternativesRelationListData),
			createFunctionData[model.PowerSequenceDescriptionListDataType, F](model.FunctionTypePowerSequenceDescriptionListData),
			createFunctionData[model.PowerSequenceNodeScheduleInformationDataType, F](model.FunctionTypePowerSequenceNodeScheduleInformationData),
			createFunctionData[model.PowerSequencePriceCalculationRequestCallType, F](model.FunctionTypePowerSequencePriceCalculationRequestCall),
			createFunctionData[model.PowerSequencePriceListDataType, F](model.FunctionTypePowerSequencePriceListData),
			createFunctionData[model.PowerSequenceScheduleConfigurationRequestCallType, F](model.FunctionTypePowerSequenceScheduleConfigurationRequestCall),
			createFunctionData[model.PowerSequenceScheduleConstraintsListDataType, F](model.FunctionTypePowerSequenceScheduleConstraintsListData),
			createFunctionData[model.PowerSequenceScheduleListDataType, F](model.FunctionTypePowerSequenceScheduleListData),
			createFunctionData[model.PowerSequenceSchedulePreferenceListDataType, F](model.FunctionTypePowerSequenceSchedulePreferenceListData),
			createFunctionData[model.PowerSequenceStateListDataType, F](model.FunctionTypePowerSequenceStateListData),
			createFunctionData[model.PowerTimeSlotScheduleConstraintsListDataType, F](model.FunctionTypePowerTimeSlotScheduleConstraintsListData),
			createFunctionData[model.PowerTimeSlotScheduleListDataType, F](model.FunctionTypePowerTimeSlotScheduleListData),
			createFunctionData[model.PowerTimeSlotValueListDataType, F](model.FunctionTypePowerTimeSlotValueListData),
		}...)
	}

	if featureType == model.FeatureTypeTypeSensing || featureType == model.FeatureTypeTypeGeneric {
		result = append(result, []F{
			createFunctionData[model.SensingDescriptionDataType, F](model.FunctionTypeSensingDescriptionData),
			createFunctionData[model.SensingListDataType, F](model.FunctionTypeSensingListData),
		}...)
	}

	if featureType == model.FeatureTypeTypeSetpoint || featureType == model.FeatureTypeTypeGeneric {
		result = append(result, []F{
			createFunctionData[model.SetpointConstraintsListDataType, F](model.FunctionTypeSetpointConstraintsListData),
			createFunctionData[model.SetpointDescriptionListDataType, F](model.FunctionTypeSetpointDescriptionListData),
			createFunctionData[model.SetpointListDataType, F](model.FunctionTypeSetpointListData),
		}...)
	}

	if featureType == model.FeatureTypeTypeSmartEnergyManagementPs || featureType == model.FeatureTypeTypeGeneric {
		result = append(result, []F{
			createFunctionData[model.SmartEnergyManagementPsConfigurationRequestCallType, F](model.FunctionTypeSmartEnergyManagementPsConfigurationRequestCall),
			createFunctionData[model.SmartEnergyManagementPsDataType, F](model.FunctionTypeSmartEnergyManagementPsData),
			createFunctionData[model.SmartEnergyManagementPsPriceCalculationRequestCallType, F](model.FunctionTypeSmartEnergyManagementPsPriceCalculationRequestCall),
			createFunctionData[model.SmartEnergyManagementPsPriceDataType, F](model.FunctionTypeSmartEnergyManagementPsPriceData),
		}...)
	}

	if featureType == model.FeatureTypeTypeStateInformation || featureType == model.FeatureTypeTypeGeneric {
		result = append(result, []F{
			createFunctionData[model.StateInformationListDataType, F](model.FunctionTypeStateInformationListData),
		}...)
	}

	if featureType == model.FeatureTypeTypeSupplyCondition || featureType == model.FeatureTypeTypeGeneric {
		result = append(result, []F{
			createFunctionData[model.SupplyConditionDescriptionListDataType, F](model.FunctionTypeSupplyConditionDescriptionListData),
			createFunctionData[model.SupplyConditionListDataType, F](model.FunctionTypeSupplyConditionListData),
			createFunctionData[model.SupplyConditionThresholdRelationListDataType, F](model.FunctionTypeSupplyConditionThresholdRelationListData),
		}...)
	}

	if featureType == model.FeatureTypeTypeTariffInformation || featureType == model.FeatureTypeTypeGeneric {
		result = append(result, []F{
			createFunctionData[model.IncentiveDescriptionListDataType, F](model.FunctionTypeIncentiveDescriptionListData),
			createFunctionData[model.IncentiveListDataType, F](model.FunctionTypeIncentiveListData),
			createFunctionData[model.TariffBoundaryRelationListDataType, F](model.FunctionTypeTariffBoundaryRelationListData),
			createFunctionData[model.TariffDescriptionListDataType, F](model.FunctionTypeTariffDescriptionListData),
			createFunctionData[model.TariffListDataType, F](model.FunctionTypeTariffListData),
			createFunctionData[model.TariffOverallConstraintsDataType, F](model.FunctionTypeTariffOverallConstraintsData),
			createFunctionData[model.TariffTierRelationListDataType, F](model.FunctionTypeTariffTierRelationListData),
			createFunctionData[model.TierBoundaryDescriptionListDataType, F](model.FunctionTypeTierBoundaryDescriptionListData),
			createFunctionData[model.TierBoundaryListDataType, F](model.FunctionTypeTierBoundaryListData),
			createFunctionData[model.TierDescriptionListDataType, F](model.FunctionTypeTierDescriptionListData),
			createFunctionData[model.TierIncentiveRelationListDataType, F](model.FunctionTypeTierIncentiveRelationListData),
			createFunctionData[model.TierListDataType, F](model.FunctionTypeTierListData),
		}...)
	}

	if featureType == model.FeatureTypeTypeTaskManagement || featureType == model.FeatureTypeTypeGeneric {
		result = append(result, []F{
			createFunctionData[model.TaskManagementJobDescriptionListDataType, F](model.FunctionTypeTaskManagementJobDescriptionListData),
			createFunctionData[model.TaskManagementJobListDataType, F](model.FunctionTypeTaskManagementJobListData),
			createFunctionData[model.TaskManagementJobRelationListDataType, F](model.FunctionTypeTaskManagementJobRelationListData),
			createFunctionData[model.TaskManagementOverviewDataType, F](model.FunctionTypeTaskManagementOverviewData),
		}...)
	}

	if featureType == model.FeatureTypeTypeThreshold || featureType == model.FeatureTypeTypeGeneric {
		result = append(result, []F{
			createFunctionData[model.ThresholdConstraintsListDataType, F](model.FunctionTypeThresholdConstraintsListData),
			createFunctionData[model.ThresholdDescriptionListDataType, F](model.FunctionTypeThresholdDescriptionListData),
			createFunctionData[model.ThresholdListDataType, F](model.FunctionTypeThresholdListData),
		}...)
	}

	if featureType == model.FeatureTypeTypeTimeInformation || featureType == model.FeatureTypeTypeGeneric {
		result = append(result, []F{
			createFunctionData[model.TimeDistributorDataType, F](model.FunctionTypeTimeDistributorData),
			createFunctionData[model.TimeDistributorEnquiryCallType, F](model.FunctionTypeTimeDistributorEnquiryCall),
			createFunctionData[model.TimeInformationDataType, F](model.FunctionTypeTimeInformationData),
			createFunctionData[model.TimePrecisionDataType, F](model.FunctionTypeTimePrecisionData),
		}...)
	}

	if featureType == model.FeatureTypeTypeTimeSeries || featureType == model.FeatureTypeTypeGeneric {
		result = append(result, []F{
			createFunctionData[model.TimeSeriesDescriptionListDataType, F](model.FunctionTypeTimeSeriesDescriptionListData),
			createFunctionData[model.TimeSeriesConstraintsListDataType, F](model.FunctionTypeTimeSeriesConstraintsListData),
			createFunctionData[model.TimeSeriesListDataType, F](model.FunctionTypeTimeSeriesListData),
		}...)
	}

	if featureType == model.FeatureTypeTypeTimeTable || featureType == model.FeatureTypeTypeGeneric {
		result = append(result, []F{
			createFunctionData[model.TimeTableConstraintsListDataType, F](model.FunctionTypeTimeTableConstraintsListData),
			createFunctionData[model.TimeTableDescriptionListDataType, F](model.FunctionTypeTimeTableDescriptionListData),
			createFunctionData[model.TimeTableListDataType, F](model.FunctionTypeTimeTableListData),
		}...)
	}

	if len(result) == 0 {
		panic(fmt.Errorf("unknown featureType '%s'", featureType))
	}

	return result
}

func createFunctionData[T any, F any](functionType model.FunctionType) F {
	x := any(new(F))
	switch x.(type) {
	case *FunctionDataCmd:
		return any(NewFunctionDataCmd[T](functionType)).(F)
	case *FunctionData:
		return any(NewFunctionData[T](functionType)).(F)
	default:
		panic(fmt.Errorf("only FunctionData and FunctionDataCmd are supported"))
	}
}
