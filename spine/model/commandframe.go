package model

type MsgCounterType uint64

type CmdClassifierType string

const (
	CmdClassifierTypeRead   CmdClassifierType = "read"
	CmdClassifierTypeReply  CmdClassifierType = "reply"
	CmdClassifierTypeNotify CmdClassifierType = "notify"
	CmdClassifierTypeWrite  CmdClassifierType = "write"
	CmdClassifierTypeCall   CmdClassifierType = "call"
	CmdClassifierTypeResult CmdClassifierType = "result"
)

type FilterIdType uint

type FilterType struct {
	FilterId   *FilterIdType   `json:"filterId,omitempty"`
	CmdControl *CmdControlType `json:"cmdControl,omitempty"`

	// DataSelectorsChoiceGroup
	AlarmListDataSelectors                                    *AlarmListDataSelectorsType                                    `json:"alarmListDataSelectors,omitempty"`
	BillConstraintsListDataSelectors                          *BillConstraintsListDataSelectorsType                          `json:"billConstraintsListDataSelectors,omitempty"`
	BillDescriptionListDataSelectors                          *BillDescriptionListDataSelectorsType                          `json:"billDescriptionListDataSelectors,omitempty"`
	BillListDataSelectors                                     *BillListDataSelectorsType                                     `json:"billListDataSelectors,omitempty"`
	BindingManagementEntryListDataSelectors                   *BindingManagementEntryListDataSelectorsType                   `json:"bindingManagementEntryListDataSelectors,omitempty"`
	CommodityListDataSelectors                                *CommodityListDataSelectorsType                                `json:"commodityListDataSelectors,omitempty"`
	DeviceConfigurationKeyValueConstraintsListDataSelectors   *DeviceConfigurationKeyValueConstraintsListDataSelectorsType   `json:"deviceConfigurationKeyValueConstraintsListDataSelectors,omitempty"`
	DeviceConfigurationKeyValueDescriptionListDataSelectors   *DeviceConfigurationKeyValueDescriptionListDataSelectorsType   `json:"deviceConfigurationKeyValueDescriptionListDataSelectors,omitempty"`
	DeviceConfigurationKeyValueListDataSelectors              *DeviceConfigurationKeyValueListDataSelectorsType              `json:"deviceConfigurationKeyValueListDataSelectors,omitempty"`
	DirectControlActivityListDataSelectors                    *DirectControlActivityListDataType                             `json:"directControlActivityListDataSelectors,omitempty"`
	ElectricalConnectionDescriptionListDataSelectors          *ElectricalConnectionDescriptionListDataSelectorsType          `json:"electricalConnectionDescriptionListDataSelectors,omitempty"`
	ElectricalConnectionParameterDescriptionListDataSelectors *ElectricalConnectionParameterDescriptionListDataSelectorsType `json:"electricalConnectionParameterDescriptionListDataSelectors,omitempty"`
	ElectricalConnectionPermittedValueSetListDataSelectors    *ElectricalConnectionPermittedValueSetListDataSelectorsType    `json:"electricalConnectionPermittedValueSetListDataSelectors,omitempty"`
	ElectricalConnectionStateListDataSelectors                *ElectricalConnectionStateListDataSelectorsType                `json:"electricalConnectionStateListDataSelectors,omitempty"`
	HvacOperationModeDescriptionListDataSelectors             *HvacOperationModeDescriptionListDataSelectorsType             `json:"hvacOperationModeDescriptionListDataSelectors,omitempty"`
	HvacOverrunDescriptionListDataSelectors                   *HvacOverrunDescriptionListDataSelectorsType                   `json:"hvacOverrunDescriptionListDataSelectors,omitempty"`
	HvacOverrunListDataSelectors                              *HvacOverrunListDataSelectorsType                              `json:"hvacOverrunListDataSelectors,omitempty"`
	HvacSystemFunctionDescriptionListDataSelectors            *HvacSystemFunctionDescriptionListDataSelectorsType            `json:"hvacSystemFunctionDescriptionListDataSelectors,omitempty"`
	HvacSystemFunctionListDataSelectors                       *HvacSystemFunctionListDataSelectorsType                       `json:"hvacSystemFunctionListDataSelectors,omitempty"`
	HvacSystemFunctionOperationModeRelationListDataSelectors  *HvacSystemFunctionOperationModeRelationListDataSelectorsType  `json:"hvacSystemFunctionOperationModeRelationListDataSelectors,omitempty"`
	HvacSystemFunctionPowerSequenceRelationListDataSelectors  *HvacSystemFunctionPowerSequenceRelationListDataSelectorsType  `json:"hvacSystemFunctionPowerSequenceRelationListDataSelectors,omitempty"`
	HvacSystemFunctionSetpointRelationListDataSelectors       *HvacSystemFunctionSetpointRelationListDataSelectorsType       `json:"hvacSystemFunctionSetpointRelationListDataSelectors,omitempty"`
	IdentificationListDataSelectors                           *IdentificationListDataSelectorsType                           `json:"identificationListDataSelectors,omitempty"`
	IncentiveDescriptionListDataSelectors                     *IncentiveDescriptionListDataSelectorsType                     `json:"incentiveDescriptionListDataSelectors,omitempty"`
	IncentiveListDataSelectors                                *IncentiveListDataSelectorsType                                `json:"incentiveListDataSelectors,omitempty"`
	IncentiveTableConstraintsDataSelectors                    *IncentiveTableConstraintsDataSelectorsType                    `json:"incentiveTableConstraintsDataSelectors,omitempty"`
	IncentiveTableDataSelectors                               *IncentiveTableDataSelectorsType                               `json:"incentiveTableDataSelectors,omitempty"`
	IncentiveTableDescriptionDataSelectors                    *IncentiveTableDescriptionDataSelectorsType                    `json:"incentiveTableDescriptionDataSelectors,omitempty"`
	LoadControlEventListDataSelectors                         *LoadControlEventListDataSelectorsType                         `json:"loadControlEventListDataSelectors,omitempty"`
	LoadControlLimitConstraintsListDataSelectors              *LoadControlLimitConstraintsListDataSelectorsType              `json:"loadControlLimitConstraintsListDataSelectors,omitempty"`
	LoadControlLimitDescriptionListDataSelectors              *LoadControlLimitDescriptionListDataSelectorsType              `json:"loadControlLimitDescriptionListDataSelectors,omitempty"`
	LoadControlLimitListDataSelectors                         *LoadControlLimitListDataSelectorsType                         `json:"loadControlLimitListDataSelectors,omitempty"`
	LoadControlStateListDataSelectors                         *LoadControlStateListDataSelectorsType                         `json:"loadControlStateListDataSelectors,omitempty"`
	MeasurementConstraintsListDataSelectors                   *MeasurementConstraintsListDataSelectorsType                   `json:"measurementConstraintsListDataSelectors,omitempty"`
	MeasurementDescriptionListDataSelectors                   *MeasurementDescriptionListDataSelectorsType                   `json:"measurementDescriptionListDataSelectors,omitempty"`
	MeasurementListDataSelectors                              *MeasurementListDataSelectorsType                              `json:"measurementListDataSelectors,omitempty"`
	MeasurementThresholdRelationListDataSelectors             *MeasurementThresholdRelationListDataSelectorsType             `json:"measurementThresholdRelationListDataSelectors,omitempty"`
	MessagingListDataSelectors                                *MessagingListDataSelectorsType                                `json:"messagingListDataSelectors,omitempty"`
	NetworkManagementDeviceDescriptionListDataSelectors       *NetworkManagementDeviceDescriptionListDataSelectorsType       `json:"networkManagementDeviceDescriptionListDataSelectors,omitempty"`
	NetworkManagementEntityDescriptionListDataSelectors       *NetworkManagementEntityDescriptionListDataSelectorsType       `json:"networkManagementEntityDescriptionListDataSelectors,omitempty"`
	NetworkManagementFeatureDescriptionListDataSelectors      *NetworkManagementFeatureDescriptionListDataSelectorsType      `json:"networkManagementFeatureDescriptionListDataSelectors,omitempty"`
	NodeManagementBindingDataSelectors                        *NodeManagementBindingDataSelectorsType                        `json:"nodeManagementBindingDataSelectors,omitempty"`
	NodeManagementDestinationListDataSelectors                *NodeManagementDestinationListDataSelectorsType                `json:"nodeManagementDestinationListDataSelectors,omitempty"`
	NodeManagementDetailedDiscoveryDataSelectors              *NodeManagementDetailedDiscoveryDataSelectorsType              `json:"nodeManagementDetailedDiscoveryDataSelectors,omitempty"`
	NodeManagementSubscriptionDataSelectors                   *NodeManagementSubscriptionDataSelectorsType                   `json:"nodeManagementSubscriptionDataSelectors,omitempty"`
	NodeManagementUseCaseDataSelectors                        *NodeManagementUseCaseDataSelectorsType                        `json:"nodeManagementUseCaseDataSelectors,omitempty"`
	OperatingConstraintsDurationListDataSelectors             *OperatingConstraintsDurationListDataSelectorsType             `json:"operatingConstraintsDurationListDataSelectors,omitempty"`
	OperatingConstraintsInterruptListDataSelectors            *OperatingConstraintsInterruptListDataSelectorsType            `json:"operatingConstraintsInterruptListDataSelectors,omitempty"`
	OperatingConstraintsPowerDescriptionListDataSelectors     *OperatingConstraintsPowerDescriptionListDataSelectorsType     `json:"operatingConstraintsPowerDescriptionListDataSelectors,omitempty"`
	OperatingConstraintsPowerLevelListDataSelectors           *OperatingConstraintsPowerLevelListDataSelectorsType           `json:"operatingConstraintsPowerLevelListDataSelectors,omitempty"`
	OperatingConstraintsPowerRangeListDataSelectors           *OperatingConstraintsPowerRangeListDataSelectorsType           `json:"operatingConstraintsPowerRangeListDataSelectors,omitempty"`
	OperatingConstraintsResumeImplicationListDataSelectors    *OperatingConstraintsResumeImplicationListDataSelectorsType    `json:"operatingConstraintsResumeImplicationListDataSelectors,omitempty"`
	PowerSequenceAlternativesRelationListDataSelectors        *PowerSequenceAlternativesRelationListDataSelectorsType        `json:"powerSequenceAlternativesRelationListDataSelectors,omitempty"`
	PowerSequenceDescriptionListDataSelectors                 *PowerSequenceDescriptionListDataSelectorsType                 `json:"powerSequenceDescriptionListDataSelectors,omitempty"`
	PowerSequencePriceListDataSelectors                       *PowerSequencePriceListDataSelectorsType                       `json:"powerSequencePriceListDataSelectors,omitempty"`
	PowerSequenceScheduleConstraintsListDataSelectors         *PowerSequenceScheduleConstraintsListDataSelectorsType         `json:"powerSequenceScheduleConstraintsListDataSelectors,omitempty"`
	PowerSequenceScheduleListDataSelectors                    *PowerSequenceScheduleListDataSelectorsType                    `json:"powerSequenceScheduleListDataSelectors,omitempty"`
	PowerSequenceSchedulePreferenceListDataSelectors          *PowerSequenceSchedulePreferenceListDataSelectorsType          `json:"powerSequenceSchedulePreferenceListDataSelectors,omitempty"`
	PowerSequenceStateListDataSelectors                       *PowerSequenceStateListDataSelectorsType                       `json:"powerSequenceStateListDataSelectors,omitempty"`
	PowerTimeSlotScheduleConstraintsListDataSelectors         *PowerTimeSlotScheduleConstraintsListDataSelectorsType         `json:"powerTimeSlotScheduleConstraintsListDataSelectors,omitempty"`
	PowerTimeSlotScheduleListDataSelectors                    *PowerTimeSlotScheduleListDataSelectorsType                    `json:"powerTimeSlotScheduleListDataSelectors,omitempty"`
	PowerTimeSlotValueListDataSelectors                       *PowerTimeSlotValueListDataSelectorsType                       `json:"powerTimeSlotValueListDataSelectors,omitempty"`
	SensingListDataSelectors                                  *SensingListDataSelectorsType                                  `json:"sensingListDataSelectors,omitempty"`
	SetpointConstraintsListDataSelectors                      *SetpointConstraintsListDataSelectorsType                      `json:"setpointConstraintsListDataSelectors,omitempty"`
	SetpointDescriptionListDataSelectors                      *SetpointDescriptionListDataSelectorsType                      `json:"setpointDescriptionListDataSelectors,omitempty"`
	SetpointListDataSelectors                                 *SetpointListDataSelectorsType                                 `json:"setpointListDataSelectors,omitempty"`
	SmartEnergyManagementPsDataSelectors                      *SmartEnergyManagementPsDataSelectorsType                      `json:"smartEnergyManagementPsDataSelectors,omitempty"`
	SmartEnergyManagementPsPriceDataSelectors                 *SmartEnergyManagementPsPriceDataSelectorsType                 `json:"smartEnergyManagementPsPriceDataSelectors,omitempty"`
	SpecificationVersionListDataSelectors                     *SpecificationVersionListDataSelectorsType                     `json:"specificationVersionListDataSelectors,omitempty"`
	SubscriptionManagementEntryListDataSelectors              *SubscriptionManagementEntryListDataSelectorsType              `json:"subscriptionManagementEntryListDataSelectors,omitempty"`
	SupplyConditionDescriptionListDataSelectors               *SupplyConditionDescriptionListDataSelectorsType               `json:"supplyConditionDescriptionListDataSelectors,omitempty"`
	SupplyConditionListDataSelectors                          *SupplyConditionListDataSelectorsType                          `json:"supplyConditionListDataSelectors,omitempty"`
	SupplyConditionThresholdRelationListDataSelectors         *SupplyConditionThresholdRelationListDataSelectorsType         `json:"supplyConditionThresholdRelationListDataSelectors,omitempty"`
	TariffBoundaryRelationListDataSelectors                   *TariffBoundaryRelationListDataSelectorsType                   `json:"tariffBoundaryRelationListDataSelectors,omitempty"`
	TariffDescriptionListDataSelectors                        *TariffDescriptionListDataSelectorsType                        `json:"tariffDescriptionListDataSelectors,omitempty"`
	TariffListDataSelectors                                   *TariffListDataSelectorsType                                   `json:"tariffListDataSelectors,omitempty"`
	TariffTierRelationListDataSelectors                       *TariffTierRelationListDataSelectorsType                       `json:"tariffTierRelationListDataSelectors,omitempty"`
	TaskManagementJobDescriptionListDataSelectors             *TaskManagementJobDescriptionListDataSelectorsType             `json:"taskManagementJobDescriptionListDataSelectors,omitempty"`
	TaskManagementJobListDataSelectors                        *TaskManagementJobListDataSelectorsType                        `json:"taskManagementJobListDataSelectors,omitempty"`
	TaskManagementJobRelationListDataSelectors                *TaskManagementJobRelationListDataSelectorsType                `json:"taskManagementJobRelationListDataSelectors,omitempty"`
	ThresholdConstraintsListDataSelectors                     *ThresholdConstraintsListDataSelectorsType                     `json:"thresholdConstraintsListDataSelectors,omitempty"`
	ThresholdDescriptionListDataSelectors                     *ThresholdDescriptionListDataSelectorsType                     `json:"thresholdDescriptionListDataSelectors,omitempty"`
	ThresholdListDataSelectors                                *ThresholdListDataSelectorsType                                `json:"thresholdListDataSelectors,omitempty"`
	TierBoundaryDescriptionListDataSelectors                  *TierBoundaryDescriptionListDataSelectorsType                  `json:"tierBoundaryDescriptionListDataSelectors,omitempty"`
	TierBoundaryListDataSelectors                             *TierBoundaryListDataSelectorsType                             `json:"tierBoundaryListDataSelectors,omitempty"`
	TierDescriptionListDataSelectors                          *TierDescriptionListDataSelectorsType                          `json:"tierDescriptionListDataSelectors,omitempty"`
	TierIncentiveRelationListDataSelectors                    *TierIncentiveRelationListDataSelectorsType                    `json:"tierIncentiveRelationListDataSelectors,omitempty"`
	TierListDataSelectors                                     *TierListDataSelectorsType                                     `json:"tierListDataSelectors,omitempty"`
	TimeSeriesConstraintsListDataSelectors                    *TimeSeriesConstraintsListDataSelectorsType                    `json:"timeSeriesConstraintsListDataSelectors,omitempty"`
	TimeSeriesDescriptionListDataSelectors                    *TimeSeriesDescriptionListDataSelectorsType                    `json:"timeSeriesDescriptionListDataSelectors,omitempty"`
	TimeSeriesListDataSelectors                               *TimeSeriesListDataSelectorsType                               `json:"timeSeriesListDataSelectors,omitempty"`
	TimeTableConstraintsListDataSelectors                     *TimeTableConstraintsListDataSelectorsType                     `json:"timeTableConstraintsListDataSelectors,omitempty"`
	TimeTableDescriptionListDataSelectors                     *TimeTableDescriptionListDataSelectorsType                     `json:"timeTableDescriptionListDataSelectors,omitempty"`
	TimeTableListDataSelectors                                *TimeTableListDataSelectorsType                                `json:"timeTableListDataSelectors,omitempty"`
	UseCaseInformationListDataSelectors                       *UseCaseInformationListDataSelectorsType                       `json:"useCaseInformationListDataSelectors,omitempty"`

	// DataElementsChoiceGroup
	ActuatorLevelDataElements                                  *ActuatorLevelDataElementsType                                  `json:"actuatorLevelDataElements,omitempty"`
	ActuatorLevelDescriptionDataElements                       *ActuatorLevelDescriptionDataElementsType                       `json:"actuatorLevelDescriptionDataElements,omitempty"`
	ActuatorSwitchDataElements                                 *ActuatorSwitchDataElementsType                                 `json:"actuatorSwitchDataElements,omitempty"`
	ActuatorSwitchDescriptionDataElements                      *ActuatorSwitchDescriptionDataElementsType                      `json:"actuatorSwitchDescriptionDataElements,omitempty"`
	AlarmDataElements                                          *AlarmDataElementsType                                          `json:"alarmDataElements,omitempty"`
	BillConstraintsDataElements                                *BillConstraintsDataElementsType                                `json:"billConstraintsDataElements,omitempty"`
	BillDataElements                                           *BillDataElementsType                                           `json:"billDataElements,omitempty"`
	BillDescriptionDataElements                                *BillDescriptionDataElementsType                                `json:"billDescriptionDataElements,omitempty"`
	BindingManagementDeleteCallElements                        *BindingManagementDeleteCallElementsType                        `json:"bindingManagementDeleteCallElements,omitempty"`
	BindingManagementEntryDataElements                         *BindingManagementEntryDataElementsType                         `json:"bindingManagementEntryDataElements,omitempty"`
	BindingManagementRequestCallElements                       *BindingManagementRequestCallElementsType                       `json:"bindingManagementRequestCallElements,omitempty"`
	CommodityDataElements                                      *CommodityDataElementsType                                      `json:"commodityDataElements,omitempty"`
	DataTunnelingCallElements                                  *DataTunnelingCallElementsType                                  `json:"dataTunnelingCallElements,omitempty"`
	DeviceClassificationManufacturerDataElements               *DeviceClassificationManufacturerDataElementsType               `json:"deviceClassificationManufacturerDataElements,omitempty"`
	DeviceClassificationUserDataElements                       *DeviceClassificationUserDataElementsType                       `json:"deviceClassificationUserDataElements,omitempty"`
	DeviceConfigurationKeyValueConstraintsDataElements         *DeviceConfigurationKeyValueConstraintsDataElementsType         `json:"deviceConfigurationKeyValueConstraintsDataElements,omitempty"`
	DeviceConfigurationKeyValueDataElements                    *DeviceConfigurationKeyValueDataElementsType                    `json:"deviceConfigurationKeyValueDataElements,omitempty"`
	DeviceConfigurationKeyValueDescriptionDataElements         *DeviceConfigurationKeyValueDescriptionDataElementsType         `json:"deviceConfigurationKeyValueDescriptionDataElements,omitempty"`
	DeviceDiagnosisHeartbeatDataElements                       *DeviceDiagnosisHeartbeatDataElementsType                       `json:"deviceDiagnosisHeartbeatDataElements,omitempty"`
	DeviceDiagnosisServiceDataElements                         *DeviceDiagnosisServiceDataElementsType                         `json:"deviceDiagnosisServiceDataElements,omitempty"`
	DeviceDiagnosisStateDataElements                           *DeviceDiagnosisStateDataElementsType                           `json:"deviceDiagnosisStateDataElements,omitempty"`
	DirectControlActivityDataElements                          *DirectControlActivityDataElementsType                          `json:"directControlActivityDataElements,omitempty"`
	DirectControlDescriptionDataElements                       *DirectControlDescriptionDataElementsType                       `json:"directControlDescriptionDataElements,omitempty"`
	ElectricalConnectionDescriptionDataElements                *ElectricalConnectionDescriptionDataElementsType                `json:"electricalConnectionDescriptionDataElements,omitempty"`
	ElectricalConnectionParameterDescriptionDataElements       *ElectricalConnectionParameterDescriptionDataElementsType       `json:"electricalConnectionParameterDescriptionDataElements,omitempty"`
	ElectricalConnectionPermittedValueSetDataElements          *ElectricalConnectionPermittedValueSetDataElementsType          `json:"electricalConnectionPermittedValueSetDataElements,omitempty"`
	ElectricalConnectionStateDataElements                      *ElectricalConnectionStateDataElementsType                      `json:"electricalConnectionStateDataElements,omitempty"`
	HvacOperationModeDescriptionDataElements                   *HvacOperationModeDescriptionDataElementsType                   `json:"hvacOperationModeDescriptionDataElements,omitempty"`
	HvacOverrunDataElements                                    *HvacOverrunDataElementsType                                    `json:"hvacOverrunDataElements,omitempty"`
	HvacOverrunDescriptionDataElements                         *HvacOverrunDescriptionDataElementsType                         `json:"hvacOverrunDescriptionDataElements,omitempty"`
	HvacSystemFunctionDataElements                             *HvacSystemFunctionDataElementsType                             `json:"hvacSystemFunctionDataElements,omitempty"`
	HvacSystemFunctionDescriptionDataElements                  *HvacSystemFunctionDescriptionDataElementsType                  `json:"hvacSystemFunctionDescriptionDataElements,omitempty"`
	HvacSystemFunctionOperationModeRelationDataElements        *HvacSystemFunctionOperationModeRelationDataElementsType        `json:"hvacSystemFunctionOperationModeRelationDataElements,omitempty"`
	HvacSystemFunctionPowerSequenceRelationDataElements        *HvacSystemFunctionPowerSequenceRelationDataElementsType        `json:"hvacSystemFunctionPowerSequenceRelationDataElements,omitempty"`
	HvacSystemFunctionSetpointRelationDataElements             *HvacSystemFunctionSetpointRelationDataElementsType             `json:"hvacSystemFunctionSetpointRelationDataElements,omitempty"`
	IdentificationDataElements                                 *IdentificationDataElementsType                                 `json:"identificationDataElements,omitempty"`
	IncentiveDataElements                                      *IncentiveDataElementsType                                      `json:"incentiveDataElements,omitempty"`
	IncentiveDescriptionDataElements                           *IncentiveDescriptionDataElementsType                           `json:"incentiveDescriptionDataElements,omitempty"`
	IncentiveTableConstraintsDataElements                      *IncentiveTableConstraintsDataElementsType                      `json:"incentiveTableConstraintsDataElements,omitempty"`
	IncentiveTableDataElements                                 *IncentiveTableDataElementsType                                 `json:"incentiveTableDataElements,omitempty"`
	IncentiveTableDescriptionDataElements                      *IncentiveTableDescriptionDataElementsType                      `json:"incentiveTableDescriptionDataElements,omitempty"`
	LoadControlEventDataElements                               *LoadControlEventDataElementsType                               `json:"loadControlEventDataElements,omitempty"`
	LoadControlLimitConstraintsDataElements                    *LoadControlLimitConstraintsDataElementsType                    `json:"loadControlLimitConstraintsDataElements,omitempty"`
	LoadControlLimitDataElements                               *LoadControlLimitDataElementsType                               `json:"loadControlLimitDataElements,omitempty"`
	LoadControlLimitDescriptionDataElements                    *LoadControlLimitDescriptionDataElementsType                    `json:"loadControlLimitDescriptionDataElements,omitempty"`
	LoadControlNodeDataElements                                *LoadControlNodeDataElementsType                                `json:"loadControlNodeDataElements,omitempty"`
	LoadControlStateDataElements                               *LoadControlStateDataElementsType                               `json:"loadControlStateDataElements,omitempty"`
	MeasurementConstraintsDataElements                         *MeasurementConstraintsDataElementsType                         `json:"measurementConstraintsDataElements,omitempty"`
	MeasurementDataElements                                    *MeasurementDataElementsType                                    `json:"measurementDataElements,omitempty"`
	MeasurementDescriptionDataElements                         *MeasurementDescriptionDataElementsType                         `json:"measurementDescriptionDataElements,omitempty"`
	MeasurementThresholdRelationDataElements                   *MeasurementThresholdRelationDataElementsType                   `json:"measurementThresholdRelationDataElements,omitempty"`
	MessagingDataElements                                      *MessagingDataElementsType                                      `json:"messagingDataElements,omitempty"`
	NetworkManagementAbortCallElements                         *NetworkManagementAbortCallElementsType                         `json:"networkManagementAbortCallElements,omitempty"`
	NetworkManagementAddNodeCallElements                       *NetworkManagementAddNodeCallElementsType                       `json:"networkManagementAddNodeCallElements,omitempty"`
	NetworkManagementDeviceDescriptionDataElements             *NetworkManagementDeviceDescriptionDataElementsType             `json:"networkManagementDeviceDescriptionDataElements,omitempty"`
	NetworkManagementDiscoverCallElements                      *NetworkManagementDiscoverCallElementsType                      `json:"networkManagementDiscoverCallElements,omitempty"`
	NetworkManagementEntityDescriptionDataElements             *NetworkManagementEntityDescriptionDataElementsType             `json:"networkManagementEntityDescriptionDataElements,omitempty"`
	NetworkManagementFeatureDescriptionDataElements            *NetworkManagementFeatureDescriptionDataElementsType            `json:"networkManagementFeatureDescriptionDataElements,omitempty"`
	NetworkManagementJoiningModeDataElements                   *NetworkManagementJoiningModeDataElementsType                   `json:"networkManagementJoiningModeDataElements,omitempty"`
	NetworkManagementModifyNodeCallElements                    *NetworkManagementModifyNodeCallElementsType                    `json:"networkManagementModifyNodeCallElements,omitempty"`
	NetworkManagementProcessStateDataElements                  *NetworkManagementProcessStateDataElementsType                  `json:"networkManagementProcessStateDataElements,omitempty"`
	NetworkManagementRemoveNodeCallElements                    *NetworkManagementRemoveNodeCallElementsType                    `json:"networkManagementRemoveNodeCallElements,omitempty"`
	NetworkManagementReportCandidateDataElements               *NetworkManagementReportCandidateDataElementsType               `json:"networkManagementReportCandidateDataElements,omitempty"`
	NetworkManagementScanNetworkCallElements                   *NetworkManagementScanNetworkCallElementsType                   `json:"networkManagementScanNetworkCallElements,omitempty"`
	NodeManagementBindingDataElements                          *NodeManagementBindingDataElementsType                          `json:"nodeManagementBindingDataElements,omitempty"`
	NodeManagementBindingDeleteCallElements                    *NodeManagementBindingDeleteCallElementsType                    `json:"nodeManagementBindingDeleteCallElements,omitempty"`
	NodeManagementBindingRequestCallElements                   *NodeManagementBindingRequestCallElementsType                   `json:"nodeManagementBindingRequestCallElements,omitempty"`
	NodeManagementDestinationDataElements                      *NodeManagementDestinationDataElementsType                      `json:"nodeManagementDestinationDataElements,omitempty"`
	NodeManagementDetailedDiscoveryDataElements                *NodeManagementDetailedDiscoveryDataElementsType                `json:"nodeManagementDetailedDiscoveryDataElements,omitempty"`
	NodeManagementSubscriptionDataElements                     *NodeManagementSubscriptionDataElementsType                     `json:"nodeManagementSubscriptionDataElements,omitempty"`
	NodeManagementSubscriptionDeleteCallElements               *NodeManagementSubscriptionDeleteCallElementsType               `json:"nodeManagementSubscriptionDeleteCallElements,omitempty"`
	NodeManagementSubscriptionRequestCallElements              *NodeManagementSubscriptionRequestCallElementsType              `json:"nodeManagementSubscriptionRequestCallElements,omitempty"`
	NodeManagementUseCaseDataElements                          *NodeManagementUseCaseDataElementsType                          `json:"nodeManagementUseCaseDataElements,omitempty"`
	OperatingConstraintsDurationDataElements                   *OperatingConstraintsDurationDataElementsType                   `json:"operatingConstraintsDurationDataElements,omitempty"`
	OperatingConstraintsInterruptDataElements                  *OperatingConstraintsInterruptDataElementsType                  `json:"operatingConstraintsInterruptDataElements,omitempty"`
	OperatingConstraintsPowerDescriptionDataElements           *OperatingConstraintsPowerDescriptionDataElementsType           `json:"operatingConstraintsPowerDescriptionDataElements,omitempty"`
	OperatingConstraintsPowerLevelDataElements                 *OperatingConstraintsPowerLevelDataElementsType                 `json:"operatingConstraintsPowerLevelDataElements,omitempty"`
	OperatingConstraintsPowerRangeDataElements                 *OperatingConstraintsPowerRangeDataElementsType                 `json:"operatingConstraintsPowerRangeDataElements,omitempty"`
	OperatingConstraintsResumeImplicationDataElements          *OperatingConstraintsResumeImplicationDataElementsType          `json:"operatingConstraintsResumeImplicationDataElements,omitempty"`
	PowerSequenceAlternativesRelationDataElements              *PowerSequenceAlternativesRelationDataElementsType              `json:"powerSequenceAlternativesRelationDataElements,omitempty"`
	PowerSequenceDescriptionDataElements                       *PowerSequenceDescriptionDataElementsType                       `json:"powerSequenceDescriptionDataElements,omitempty"`
	PowerSequenceNodeScheduleInformationDataElements           *PowerSequenceNodeScheduleInformationDataElementsType           `json:"powerSequenceNodeScheduleInformationDataElements,omitempty"`
	PowerSequencePriceCalculationRequestCallElements           *PowerSequencePriceCalculationRequestCallElementsType           `json:"powerSequencePriceCalculationRequestCallElements,omitempty"`
	PowerSequencePriceDataElements                             *PowerSequencePriceDataElementsType                             `json:"powerSequencePriceDataElements,omitempty"`
	PowerSequenceScheduleConfigurationRequestCallElements      *PowerSequenceScheduleConfigurationRequestCallElementsType      `json:"powerSequenceScheduleConfigurationRequestCallElements,omitempty"`
	PowerSequenceScheduleConstraintsDataElements               *PowerSequenceScheduleConstraintsDataElementsType               `json:"powerSequenceScheduleConstraintsDataElements,omitempty"`
	PowerSequenceScheduleDataElements                          *PowerSequenceScheduleDataElementsType                          `json:"powerSequenceScheduleDataElements,omitempty"`
	PowerSequenceSchedulePreferenceDataElements                *PowerSequenceSchedulePreferenceDataElementsType                `json:"powerSequenceSchedulePreferenceDataElements,omitempty"`
	PowerSequenceStateDataElements                             *PowerSequenceStateDataElementsType                             `json:"powerSequenceStateDataElements,omitempty"`
	PowerTimeSlotScheduleConstraintsDataElements               *PowerTimeSlotScheduleConstraintsDataElementsType               `json:"powerTimeSlotScheduleConstraintsDataElements,omitempty"`
	PowerTimeSlotScheduleDataElements                          *PowerTimeSlotScheduleDataElementsType                          `json:"powerTimeSlotScheduleDataElements,omitempty"`
	PowerTimeSlotValueDataElements                             *PowerTimeSlotValueDataElementsType                             `json:"powerTimeSlotValueDataElements,omitempty"`
	SensingDataElements                                        *SensingDataElementsType                                        `json:"sensingDataElements,omitempty"`
	SensingDescriptionDataElements                             *SensingDescriptionDataElementsType                             `json:"sensingDescriptionDataElements,omitempty"`
	SetpointConstraintsDataElements                            *SetpointConstraintsDataElementsType                            `json:"setpointConstraintsDataElements,omitempty"`
	SetpointDataElements                                       *SetpointDataElementsType                                       `json:"setpointDataElements,omitempty"`
	SetpointDescriptionDataElements                            *SetpointDescriptionDataElementsType                            `json:"setpointDescriptionDataElements,omitempty"`
	SmartEnergyManagementPsConfigurationRequestCallElements    *SmartEnergyManagementPsConfigurationRequestCallElementsType    `json:"smartEnergyManagementPsConfigurationRequestCallElements,omitempty"`
	SmartEnergyManagementPsDataElements                        *SmartEnergyManagementPsDataElementsType                        `json:"smartEnergyManagementPsDataElements,omitempty"`
	SmartEnergyManagementPsPriceCalculationRequestCallElements *SmartEnergyManagementPsPriceCalculationRequestCallElementsType `json:"smartEnergyManagementPsPriceCalculationRequestCallElements,omitempty"`
	SmartEnergyManagementPsPriceDataElements                   *SmartEnergyManagementPsPriceDataElementsType                   `json:"smartEnergyManagementPsPriceDataElements,omitempty"`
	SpecificationVersionDataElements                           *SpecificationVersionDataElementsType                           `json:"specificationVersionDataElements,omitempty"`
	SubscriptionManagementDeleteCallElements                   *SubscriptionManagementDeleteCallElementsType                   `json:"subscriptionManagementDeleteCallElements,omitempty"`
	SubscriptionManagementEntryDataElements                    *SubscriptionManagementEntryDataElementsType                    `json:"subscriptionManagementEntryDataElements,omitempty"`
	SubscriptionManagementRequestCallElements                  *SubscriptionManagementRequestCallElementsType                  `json:"subscriptionManagementRequestCallElements,omitempty"`
	SupplyConditionDataElements                                *SupplyConditionDataElementsType                                `json:"supplyConditionDataElements,omitempty"`
	SupplyConditionDescriptionDataElements                     *SupplyConditionDescriptionDataElementsType                     `json:"supplyConditionDescriptionDataElements,omitempty"`
	SupplyConditionThresholdRelationDataElements               *SupplyConditionThresholdRelationDataElementsType               `json:"supplyConditionThresholdRelationDataElements,omitempty"`
	TariffBoundaryRelationDataElements                         *TariffBoundaryRelationDataElementsType                         `json:"tariffBoundaryRelationDataElements,omitempty"`
	TariffDataElements                                         *TariffDataElementsType                                         `json:"tariffDataElements,omitempty"`
	TariffDescriptionDataElements                              *TariffDescriptionDataElementsType                              `json:"tariffDescriptionDataElements,omitempty"`
	TariffOverallConstraintsDataElements                       *TariffOverallConstraintsDataElementsType                       `json:"tariffOverallConstraintsDataElements,omitempty"`
	TariffTierRelationDataElements                             *TariffTierRelationDataElementsType                             `json:"tariffTierRelationDataElements,omitempty"`
	TaskManagementJobDataElements                              *TaskManagementJobDataElementsType                              `json:"taskManagementJobDataElements,omitempty"`
	TaskManagementJobDescriptionDataElements                   *TaskManagementJobDescriptionDataElementsType                   `json:"taskManagementJobDescriptionDataElements,omitempty"`
	TaskManagementJobRelationDataElements                      *TaskManagementJobRelationDataElementsType                      `json:"taskManagementJobRelationDataElements,omitempty"`
	TaskManagementOverviewDataElements                         *TaskManagementOverviewDataElementsType                         `json:"taskManagementOverviewDataElements,omitempty"`
	ThresholdConstraintsDataElements                           *ThresholdConstraintsDataElementsType                           `json:"thresholdConstraintsDataElements,omitempty"`
	ThresholdDataElements                                      *ThresholdDataElementsType                                      `json:"thresholdDataElements,omitempty"`
	ThresholdDescriptionDataElements                           *ThresholdDescriptionDataElementsType                           `json:"thresholdDescriptionDataElements,omitempty"`
	TierBoundaryDataElements                                   *TierBoundaryDataElementsType                                   `json:"tierBoundaryDataElements,omitempty"`
	TierBoundaryDescriptionDataElements                        *TierBoundaryDescriptionDataElementsType                        `json:"tierBoundaryDescriptionDataElements,omitempty"`
	TierDataElements                                           *TierDataElementsType                                           `json:"tierDataElements,omitempty"`
	TierDescriptionDataElements                                *TierDescriptionDataElementsType                                `json:"tierDescriptionDataElements,omitempty"`
	TierIncentiveRelationDataElements                          *TierIncentiveRelationDataElementsType                          `json:"tierIncentiveRelationDataElements,omitempty"`
	TimeDistributorDataElements                                *TimeDistributorDataElementsType                                `json:"timeDistributorDataElements,omitempty"`
	TimeDistributorEnquiryCallElements                         *TimeDistributorEnquiryCallElementsType                         `json:"timeDistributorEnquiryCallElements,omitempty"`
	TimeInformationDataElements                                *TimeInformationDataElementsType                                `json:"timeInformationDataElements,omitempty"`
	TimePrecisionDataElements                                  *TimePrecisionDataElementsType                                  `json:"timePrecisionDataElements,omitempty"`
	TimeSeriesConstraintsDataElements                          *TimeSeriesConstraintsDataElementsType                          `json:"timeSeriesConstraintsDataElements,omitempty"`
	TimeSeriesDataElements                                     *TimeSeriesDataElementsType                                     `json:"timeSeriesDataElements,omitempty"`
	TimeSeriesDescriptionDataElements                          *TimeSeriesDescriptionDataElementsType                          `json:"timeSeriesDescriptionDataElements,omitempty"`
	TimeTableConstraintsDataElements                           *TimeTableConstraintsDataElementsType                           `json:"timeTableConstraintsDataElements,omitempty"`
	TimeTableDataElements                                      *TimeTableDataElementsType                                      `json:"timeTableDataElements,omitempty"`
	TimeTableDescriptionDataElements                           *TimeTableDescriptionDataElementsType                           `json:"timeTableDescriptionDataElements,omitempty"`
	UseCaseInformationDataElements                             *UseCaseInformationDataElementsType                             `json:"useCaseInformationDataElements,omitempty"`
}

type CmdControlType struct {
	Delete  *ElementTagType `json:"delete,omitempty"`
	Partial *ElementTagType `json:"partial,omitempty"`
}

type CmdType struct {
	// CmdOptionGroup
	Function *FunctionType `json:"function,omitempty"`
	Filter   []FilterType  `json:"filter,omitempty"`

	// DataChoiceGroup
	ActuatorLevelData                                  *ActuatorLevelDataType                                  `json:"actuatorLevelData,omitempty"`
	ActuatorLevelDescriptionData                       *ActuatorLevelDescriptionDataType                       `json:"actuatorLevelDescriptionData,omitempty"`
	ActuatorSwitchData                                 *ActuatorSwitchDataType                                 `json:"actuatorSwitchData,omitempty"`
	ActuatorSwitchDescriptionData                      *ActuatorSwitchDescriptionDataType                      `json:"actuatorSwitchDescriptionData,omitempty"`
	AlarmListData                                      *AlarmListDataType                                      `json:"alarmListData,omitempty"`
	BillConstraintsListData                            *BillConstraintsListDataType                            `json:"billConstraintsListData,omitempty"`
	BillDescriptionListData                            *BillDescriptionListDataType                            `json:"billDescriptionListData,omitempty"`
	BillListData                                       *BillListDataType                                       `json:"billListData,omitempty"`
	BindingManagementDeleteCall                        *BindingManagementDeleteCallType                        `json:"bindingManagementDeleteCall,omitempty"`
	BindingManagementEntryListData                     *BindingManagementEntryListDataType                     `json:"bindingManagementEntryListData,omitempty"`
	BindingManagementRequestCall                       *BindingManagementRequestCallType                       `json:"bindingManagementRequestCall,omitempty"`
	CommodityListData                                  *CommodityListDataType                                  `json:"commodityListData,omitempty"`
	DataTunnelingCall                                  *DataTunnelingCallType                                  `json:"dataTunnelingCall,omitempty"`
	DeviceClassificationManufacturerData               *DeviceClassificationManufacturerDataType               `json:"deviceClassificationManufacturerData,omitempty" eebus:"fct:deviceClassificationManufacturerData"`
	DeviceClassificationUserData                       *DeviceClassificationUserDataType                       `json:"deviceClassificationUserData,omitempty"`
	DeviceConfigurationKeyValueConstraintsListData     *DeviceConfigurationKeyValueConstraintsListDataType     `json:"deviceConfigurationKeyValueConstraintsListData,omitempty"`
	DeviceConfigurationKeyValueDescriptionListData     *DeviceConfigurationKeyValueDescriptionListDataType     `json:"deviceConfigurationKeyValueDescriptionListData,omitempty" eebus:"fct:deviceConfigurationKeyValueDescriptionListData"`
	DeviceConfigurationKeyValueListData                *DeviceConfigurationKeyValueListDataType                `json:"deviceConfigurationKeyValueListData,omitempty" eebus:"fct:deviceConfigurationKeyValueListData"`
	DeviceDiagnosisHeartbeatData                       *DeviceDiagnosisHeartbeatDataType                       `json:"deviceDiagnosisHeartbeatData,omitempty"`
	DeviceDiagnosisServiceData                         *DeviceDiagnosisServiceDataType                         `json:"deviceDiagnosisServiceData,omitempty"`
	DeviceDiagnosisStateData                           *DeviceDiagnosisStateDataType                           `json:"deviceDiagnosisStateData,omitempty" eebus:"fct:deviceDiagnosisStateData"`
	DirectControlActivityListData                      *DirectControlActivityListDataType                      `json:"directControlActivityListData,omitempty"`
	DirectControlDescriptionData                       *DirectControlDescriptionDataType                       `json:"directControlDescriptionData,omitempty"`
	ElectricalConnectionDescriptionListData            *ElectricalConnectionDescriptionListDataType            `json:"electricalConnectionDescriptionListData,omitempty" eebus:"fct:electricalConnectionDescriptionListData"`
	ElectricalConnectionParameterDescriptionListData   *ElectricalConnectionParameterDescriptionListDataType   `json:"electricalConnectionParameterDescriptionListData,omitempty" eebus:"fct:electricalConnectionParameterDescriptionListData"`
	ElectricalConnectionPermittedValueSetListData      *ElectricalConnectionPermittedValueSetListDataType      `json:"electricalConnectionPermittedValueSetListData,omitempty" eebus:"fct:electricalConnectionPermittedValueSetListData"`
	ElectricalConnectionStateListData                  *ElectricalConnectionStateListDataType                  `json:"electricalConnectionStateListData,omitempty"`
	HvacOperationModeDescriptionListData               *HvacOperationModeDescriptionListDataType               `json:"hvacOperationModeDescriptionListData,omitempty"`
	HvacOverrunDescriptionListData                     *HvacOverrunDescriptionListDataType                     `json:"hvacOverrunDescriptionListData,omitempty" eebus:"fct:hvacOverrunDescriptionListData"`
	HvacOverrunListData                                *HvacOverrunListDataType                                `json:"hvacOverrunListData,omitempty" eebus:"fct:hvacOverrunListData"`
	HvacSystemFunctionDescriptionListData              *HvacSystemFunctionDescriptionListDataType              `json:"hvacSystemFunctionDescriptionListData,omitempty"`
	HvacSystemFunctionListData                         *HvacSystemFunctionListDataType                         `json:"hvacSystemFunctionListData,omitempty"`
	HvacSystemFunctionOperationModeRelationListData    *HvacSystemFunctionOperationModeRelationListDataType    `json:"hvacSystemFunctionOperationModeRelationListData,omitempty"`
	HvacSystemFunctionPowerSequenceRelationListData    *HvacSystemFunctionPowerSequenceRelationListDataType    `json:"hvacSystemFunctionPowerSequenceRelationListData,omitempty"`
	HvacSystemFunctionSetpointRelationListData         *HvacSystemFunctionSetpointRelationListDataType         `json:"hvacSystemFunctionSetpointRelationListData,omitempty"`
	IdentificationListData                             *IdentificationListDataType                             `json:"identificationListData,omitempty" eebus:"fct:identificationListData"`
	IncentiveDescriptionListData                       *IncentiveDescriptionListDataType                       `json:"incentiveDescriptionListData,omitempty"`
	IncentiveListData                                  *IncentiveListDataType                                  `json:"incentiveListData,omitempty"`
	IncentiveTableConstraintsData                      *IncentiveTableConstraintsDataType                      `json:"incentiveTableConstraintsData,omitempty"`
	IncentiveTableData                                 *IncentiveTableDataType                                 `json:"incentiveTableData,omitempty"`
	IncentiveTableDescriptionData                      *IncentiveTableDescriptionDataType                      `json:"incentiveTableDescriptionData,omitempty"`
	LoadControlEventListData                           *LoadControlEventListDataType                           `json:"loadControlEventListData,omitempty"`
	LoadControlLimitConstraintsListData                *LoadControlLimitConstraintsListDataType                `json:"loadControlLimitConstraintsListData,omitempty"`
	LoadControlLimitDescriptionListData                *LoadControlLimitDescriptionListDataType                `json:"loadControlLimitDescriptionListData,omitempty"`
	LoadControlLimitListData                           *LoadControlLimitListDataType                           `json:"loadControlLimitListData,omitempty"`
	LoadControlNodeData                                *LoadControlNodeDataType                                `json:"loadControlNodeData,omitempty"`
	LoadControlStateListData                           *LoadControlStateListDataType                           `json:"loadControlStateListData,omitempty"`
	MeasurementConstraintsListData                     *MeasurementConstraintsListDataType                     `json:"measurementConstraintsListData,omitempty" eebus:"fct:measurementConstraintsListData"`
	MeasurementDescriptionListData                     *MeasurementDescriptionListDataType                     `json:"measurementDescriptionListData,omitempty" eebus:"fct:measurementDescriptionListData"`
	MeasurementListData                                *MeasurementListDataType                                `json:"measurementListData,omitempty" eebus:"fct:measurementListData"`
	MeasurementThresholdRelationListData               *MeasurementThresholdRelationListDataType               `json:"measurementThresholdRelationListData,omitempty"`
	MessagingListData                                  *MessagingListDataType                                  `json:"messagingListData,omitempty"`
	NetworkManagementAbortCall                         *NetworkManagementAbortCallType                         `json:"networkManagementAbortCall,omitempty"`
	NetworkManagementAddNodeCall                       *NetworkManagementAddNodeCallType                       `json:"networkManagementAddNodeCall,omitempty"`
	NetworkManagementDeviceDescriptionListData         *NetworkManagementDeviceDescriptionListDataType         `json:"networkManagementDeviceDescriptionListData,omitempty"`
	NetworkManagementDiscoverCall                      *NetworkManagementDiscoverCallType                      `json:"networkManagementDiscoverCall,omitempty"`
	NetworkManagementEntityDescriptionListData         *NetworkManagementEntityDescriptionListDataType         `json:"networkManagementEntityDescriptionListData,omitempty"`
	NetworkManagementFeatureDescriptionListData        *NetworkManagementFeatureDescriptionListDataType        `json:"networkManagementFeatureDescriptionListData,omitempty"`
	NetworkManagementJoiningModeData                   *NetworkManagementJoiningModeDataType                   `json:"networkManagementJoiningModeData,omitempty"`
	NetworkManagementModifyNodeCall                    *NetworkManagementModifyNodeCallType                    `json:"networkManagementModifyNodeCall,omitempty"`
	NetworkManagementProcessStateData                  *NetworkManagementProcessStateDataType                  `json:"networkManagementProcessStateData,omitempty"`
	NetworkManagementRemoveNodeCall                    *NetworkManagementRemoveNodeCallType                    `json:"networkManagementRemoveNodeCall,omitempty"`
	NetworkManagementReportCandidateData               *NetworkManagementReportCandidateDataType               `json:"networkManagementReportCandidateData,omitempty"`
	NetworkManagementScanNetworkCall                   *NetworkManagementScanNetworkCallType                   `json:"networkManagementScanNetworkCall,omitempty"`
	NodeManagementBindingData                          *NodeManagementBindingDataType                          `json:"nodeManagementBindingData,omitempty" eebus:"fct:nodeManagementBindingData"`
	NodeManagementBindingDeleteCall                    *NodeManagementBindingDeleteCallType                    `json:"nodeManagementBindingDeleteCall,omitempty" eebus:"fct:nodeManagementBindingDeleteCall"`
	NodeManagementBindingRequestCall                   *NodeManagementBindingRequestCallType                   `json:"nodeManagementBindingRequestCall,omitempty" eebus:"fct:nodeManagementBindingRequestCall"`
	NodeManagementDestinationListData                  *NodeManagementDestinationListDataType                  `json:"nodeManagementDestinationListData,omitempty" eebus:"fct:nodeManagementDestinationListData"`
	NodeManagementDetailedDiscoveryData                *NodeManagementDetailedDiscoveryDataType                `json:"nodeManagementDetailedDiscoveryData,omitempty" eebus:"fct:nodeManagementDetailedDiscoveryData"`
	NodeManagementSubscriptionData                     *NodeManagementSubscriptionDataType                     `json:"nodeManagementSubscriptionData,omitempty" eebus:"fct:nodeManagementSubscriptionData"`
	NodeManagementSubscriptionDeleteCall               *NodeManagementSubscriptionDeleteCallType               `json:"nodeManagementSubscriptionDeleteCall,omitempty" eebus:"fct:nodeManagementSubscriptionDeleteCall"`
	NodeManagementSubscriptionRequestCall              *NodeManagementSubscriptionRequestCallType              `json:"nodeManagementSubscriptionRequestCall,omitempty" eebus:"fct:nodeManagementSubscriptionRequestCall"`
	NodeManagementUseCaseData                          *NodeManagementUseCaseDataType                          `json:"nodeManagementUseCaseData,omitempty" eebus:"fct:nodeManagementUseCaseData"`
	OperatingConstraintsDurationListData               *OperatingConstraintsDurationListDataType               `json:"operatingConstraintsDurationListData,omitempty"`
	OperatingConstraintsInterruptListData              *OperatingConstraintsInterruptListDataType              `json:"operatingConstraintsInterruptListData,omitempty"`
	OperatingConstraintsPowerDescriptionListData       *OperatingConstraintsPowerDescriptionListDataType       `json:"operatingConstraintsPowerDescriptionListData,omitempty"`
	OperatingConstraintsPowerLevelListData             *OperatingConstraintsPowerLevelListDataType             `json:"operatingConstraintsPowerLevelListData,omitempty"`
	OperatingConstraintsPowerRangeListData             *OperatingConstraintsPowerRangeListDataType             `json:"operatingConstraintsPowerRangeListData,omitempty"`
	OperatingConstraintsResumeImplicationListData      *OperatingConstraintsResumeImplicationListDataType      `json:"operatingConstraintsResumeImplicationListData,omitempty"`
	PowerSequenceAlternativesRelationListData          *PowerSequenceAlternativesRelationListDataType          `json:"powerSequenceAlternativesRelationListData,omitempty"`
	PowerSequenceDescriptionListData                   *PowerSequenceDescriptionListDataType                   `json:"powerSequenceDescriptionListData,omitempty"`
	PowerSequenceNodeScheduleInformationData           *PowerSequenceNodeScheduleInformationDataType           `json:"powerSequenceNodeScheduleInformationData,omitempty"`
	PowerSequencePriceCalculationRequestCall           *PowerSequencePriceCalculationRequestCallType           `json:"powerSequencePriceCalculationRequestCall,omitempty"`
	PowerSequencePriceListData                         *PowerSequencePriceListDataType                         `json:"powerSequencePriceListData,omitempty"`
	PowerSequenceScheduleConfigurationRequestCall      *PowerSequenceScheduleConfigurationRequestCallType      `json:"powerSequenceScheduleConfigurationRequestCall,omitempty"`
	PowerSequenceScheduleConstraintsListData           *PowerSequenceScheduleConstraintsListDataType           `json:"powerSequenceScheduleConstraintsListData,omitempty"`
	PowerSequenceScheduleListData                      *PowerSequenceScheduleListDataType                      `json:"powerSequenceScheduleListData,omitempty"`
	PowerSequenceSchedulePreferenceListData            *PowerSequenceSchedulePreferenceListDataType            `json:"powerSequenceSchedulePreferenceListData,omitempty"`
	PowerSequenceStateListData                         *PowerSequenceStateListDataType                         `json:"powerSequenceStateListData,omitempty"`
	PowerTimeSlotScheduleConstraintsListData           *PowerTimeSlotScheduleConstraintsListDataType           `json:"powerTimeSlotScheduleConstraintsListData,omitempty"`
	PowerTimeSlotScheduleListData                      *PowerTimeSlotScheduleListDataType                      `json:"powerTimeSlotScheduleListData,omitempty"`
	PowerTimeSlotValueListData                         *PowerTimeSlotValueListDataType                         `json:"powerTimeSlotValueListData,omitempty"`
	ResultData                                         *ResultDataType                                         `json:"resultData,omitempty" eebus:"fct:resultData"`
	SensingDescriptionData                             *SensingDescriptionDataType                             `json:"sensingDescriptionData,omitempty"`
	SensingListData                                    *SensingListDataType                                    `json:"sensingListData,omitempty"`
	SetpointConstraintsListData                        *SetpointConstraintsListDataType                        `json:"setpointConstraintsListData,omitempty"`
	SetpointDescriptionListData                        *SetpointDescriptionListDataType                        `json:"setpointDescriptionListData,omitempty"`
	SetpointListData                                   *SetpointListDataType                                   `json:"setpointListData,omitempty"`
	SmartEnergyManagementPsConfigurationRequestCall    *SmartEnergyManagementPsConfigurationRequestCallType    `json:"smartEnergyManagementPsConfigurationRequestCall,omitempty"`
	SmartEnergyManagementPsData                        *SmartEnergyManagementPsDataType                        `json:"smartEnergyManagementPsData,omitempty"`
	SmartEnergyManagementPsPriceCalculationRequestCall *SmartEnergyManagementPsPriceCalculationRequestCallType `json:"smartEnergyManagementPsPriceCalculationRequestCall,omitempty"`
	SmartEnergyManagementPsPriceData                   *SmartEnergyManagementPsPriceDataType                   `json:"smartEnergyManagementPsPriceData,omitempty"`
	SpecificationVersionListData                       *SpecificationVersionListDataType                       `json:"specificationVersionListData,omitempty"`
	SubscriptionManagementDeleteCall                   *SubscriptionManagementDeleteCallType                   `json:"subscriptionManagementDeleteCall,omitempty"`
	SubscriptionManagementEntryListData                *SubscriptionManagementEntryListDataType                `json:"subscriptionManagementEntryListData,omitempty"`
	SubscriptionManagementRequestCall                  *SubscriptionManagementRequestCallType                  `json:"subscriptionManagementRequestCall,omitempty"`
	SupplyConditionDescriptionListData                 *SupplyConditionDescriptionListDataType                 `json:"supplyConditionDescriptionListData,omitempty"`
	SupplyConditionListData                            *SupplyConditionListDataType                            `json:"supplyConditionListData,omitempty"`
	SupplyConditionThresholdRelationListData           *SupplyConditionThresholdRelationListDataType           `json:"supplyConditionThresholdRelationListData,omitempty"`
	TariffBoundaryRelationListData                     *TariffBoundaryRelationListDataType                     `json:"tariffBoundaryRelationListData,omitempty"`
	TariffDescriptionListData                          *TariffDescriptionListDataType                          `json:"tariffDescriptionListData,omitempty"`
	TariffListData                                     *TariffListDataType                                     `json:"tariffListData,omitempty"`
	TariffOverallConstraintsData                       *TariffOverallConstraintsDataType                       `json:"tariffOverallConstraintsData,omitempty"`
	TariffTierRelationListData                         *TariffTierRelationListDataType                         `json:"tariffTierRelationListData,omitempty"`
	TaskManagementJobDescriptionListData               *TaskManagementJobDescriptionListDataType               `json:"taskManagementJobDescriptionListData,omitempty"`
	TaskManagementJobListData                          *TaskManagementJobListDataType                          `json:"taskManagementJobListData,omitempty"`
	TaskManagementJobRelationListData                  *TaskManagementJobRelationListDataType                  `json:"taskManagementJobRelationListData,omitempty"`
	TaskManagementOverviewData                         *TaskManagementOverviewDataType                         `json:"taskManagementOverviewData,omitempty"`
	ThresholdConstraintsListData                       *ThresholdConstraintsListDataType                       `json:"thresholdConstraintsListData,omitempty"`
	ThresholdDescriptionListData                       *ThresholdDescriptionListDataType                       `json:"thresholdDescriptionListData,omitempty"`
	ThresholdListData                                  *ThresholdListDataType                                  `json:"thresholdListData,omitempty"`
	TierBoundaryDescriptionListData                    *TierBoundaryDescriptionListDataType                    `json:"tierBoundaryDescriptionListData,omitempty"`
	TierBoundaryListData                               *TierBoundaryListDataType                               `json:"tierBoundaryListData,omitempty"`
	TierDescriptionListData                            *TierDescriptionListDataType                            `json:"tierDescriptionListData,omitempty"`
	TierIncentiveRelationListData                      *TierIncentiveRelationListDataType                      `json:"tierIncentiveRelationListData,omitempty"`
	TierListData                                       *TierListDataType                                       `json:"tierListData,omitempty"`
	TimeDistributorData                                *TimeDistributorDataType                                `json:"timeDistributorData,omitempty"`
	TimeDistributorEnquiryCall                         *TimeDistributorEnquiryCallType                         `json:"timeDistributorEnquiryCall,omitempty"`
	TimeInformationData                                *TimeInformationDataType                                `json:"timeInformationData,omitempty"`
	TimePrecisionData                                  *TimePrecisionDataType                                  `json:"timePrecisionData,omitempty"`
	TimeSeriesConstraintsListData                      *TimeSeriesConstraintsListDataType                      `json:"timeSeriesConstraintsListData,omitempty"`
	TimeSeriesDescriptionListData                      *TimeSeriesDescriptionListDataType                      `json:"timeSeriesDescriptionListData,omitempty"`
	TimeSeriesListData                                 *TimeSeriesListDataType                                 `json:"timeSeriesListData,omitempty"`
	TimeTableConstraintsListData                       *TimeTableConstraintsListDataType                       `json:"timeTableConstraintsListData,omitempty"`
	TimeTableDescriptionListData                       *TimeTableDescriptionListDataType                       `json:"timeTableDescriptionListData,omitempty"`
	TimeTableListData                                  *TimeTableListDataType                                  `json:"timeTableListData,omitempty"`
	UseCaseInformationListData                         *UseCaseInformationListDataType                         `json:"useCaseInformationListData,omitempty"`

	// DataExtendGroup
	ManufacturerSpecificExtension *string                     `json:"manufacturerSpecificExtension,omitempty"`
	LastUpdateAt                  *AbsoluteOrRelativeTimeType `json:"lastUpdateAt,omitempty"`
}
