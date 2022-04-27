package model

import "fmt"

func (r *MsgCounterType) String() string {
	if r == nil {
		return ""
	}
	return fmt.Sprintf("%d", *r)
}

func (cmd CmdType) DataName() string {
	switch {
	case cmd.DeviceClassificationManufacturerData != nil:
		return "DeviceClassificationManufacturerData"
	case cmd.DeviceConfigurationKeyValueDescriptionListData != nil:
		return "DeviceConfigurationKeyValueDescriptionListData"
	case cmd.DeviceConfigurationKeyValueListData != nil:
		return "DeviceConfigurationKeyValueListData"
	case cmd.DeviceDiagnosisHeartbeatData != nil:
		return "DeviceDiagnosisHeartbeatData"
	case cmd.DeviceDiagnosisStateData != nil:
		return "DeviceDiagnosisStateData"
	case cmd.ElectricalConnectionDescriptionListData != nil:
		return "ElectricalConnectionDescriptionListData"
	case cmd.ElectricalConnectionParameterDescriptionListData != nil:
		return "ElectricalConnectionParameterDescriptionListData"
	case cmd.ElectricalConnectionPermittedValueSetListData != nil:
		return "ElectricalConnectionPermittedValueSetListData"
	case cmd.IdentificationListData != nil:
		return "IdentificationListData"
	case cmd.IncentiveTableDescriptionData != nil:
		return "IncentiveTableDescriptionData"
	case cmd.IncentiveTableConstraintsData != nil:
		return "IncentiveTableConstraintsData"
	case cmd.IncentiveTableData != nil:
		return "IncentiveTableData"
	case cmd.LoadControlLimitDescriptionListData != nil:
		return "LoadControlLimitDescriptionListData"
	case cmd.LoadControlLimitListData != nil:
		return "LoadControlLimitListData"
	case cmd.NodeManagementBindingRequestCall != nil:
		return "NodeManagementBindingRequestCall"
	case cmd.NodeManagementDetailedDiscoveryData != nil:
		return "NodeManagementDetailedDiscoveryData"
	case cmd.NodeManagementSubscriptionData != nil:
		return "NodeManagementSubscriptionData"
	case cmd.NodeManagementSubscriptionRequestCall != nil:
		return "NodeManagementSubscriptionRequestCall"
	case cmd.NodeManagementSubscriptionDeleteCall != nil:
		return "NodeManagementSubscriptionDeleteCall"
	case cmd.NodeManagementUseCaseData != nil:
		return "NodeManagementUseCaseData"
	case cmd.MeasurementConstraintsListData != nil:
		return "MeasurementConstraintsListData"
	case cmd.MeasurementDescriptionListData != nil:
		return "MeasurementDescriptionListData"
	case cmd.MeasurementListData != nil:
		return "MeasurementListData"
	case cmd.TimeSeriesConstraintsListData != nil:
		return "TimeSeriesConstraintsListData"
	case cmd.TimeSeriesDescriptionListData != nil:
		return "TimeSeriesDescriptionListData"
	case cmd.TimeSeriesListData != nil:
		return "TimeSeriesListData"
	case cmd.ResultData != nil:
		return "ResultData"
	}

	return "unknown"
}
