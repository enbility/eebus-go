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
	ActuatorLevelData                                  *ActuatorLevelDataType                                  `json:"actuatorLevelData,omitempty" eebus:"fct:actuatorLevelData"`
	ActuatorLevelDescriptionData                       *ActuatorLevelDescriptionDataType                       `json:"actuatorLevelDescriptionData,omitempty" eebus:"fct:actuatorLevelDescriptionData"`
	ActuatorSwitchData                                 *ActuatorSwitchDataType                                 `json:"actuatorSwitchData,omitempty" eebus:"fct:actuatorSwitchData"`
	ActuatorSwitchDescriptionData                      *ActuatorSwitchDescriptionDataType                      `json:"actuatorSwitchDescriptionData,omitempty" eebus:"fct:actuatorSwitchDescriptionData"`
	AlarmListData                                      *AlarmListDataType                                      `json:"alarmListData,omitempty" eebus:"fct:alarmListData"`
	BillConstraintsListData                            *BillConstraintsListDataType                            `json:"billConstraintsListData,omitempty" eebus:"fct:billConstraintsListData"`
	BillDescriptionListData                            *BillDescriptionListDataType                            `json:"billDescriptionListData,omitempty" eebus:"fct:billDescriptionListData"`
	BillListData                                       *BillListDataType                                       `json:"billListData,omitempty" eebus:"fct:billListData"`
	BindingManagementDeleteCall                        *BindingManagementDeleteCallType                        `json:"bindingManagementDeleteCall,omitempty" eebus:"fct:bindingManagementDeleteCall"`
	BindingManagementEntryListData                     *BindingManagementEntryListDataType                     `json:"bindingManagementEntryListData,omitempty" eebus:"fct:bindingManagementEntryListData"`
	BindingManagementRequestCall                       *BindingManagementRequestCallType                       `json:"bindingManagementRequestCall,omitempty" eebus:"fct:bindingManagementRequestCall"`
	CommodityListData                                  *CommodityListDataType                                  `json:"commodityListData,omitempty" eebus:"fct:commodityListData"`
	DataTunnelingCall                                  *DataTunnelingCallType                                  `json:"dataTunnelingCall,omitempty" eebus:"fct:dataTunnelingCall"`
	DeviceClassificationManufacturerData               *DeviceClassificationManufacturerDataType               `json:"deviceClassificationManufacturerData,omitempty" eebus:"fct:deviceClassificationManufacturerData"`
	DeviceClassificationUserData                       *DeviceClassificationUserDataType                       `json:"deviceClassificationUserData,omitempty" eebus:"fct:deviceClassificationUserData"`
	DeviceConfigurationKeyValueConstraintsListData     *DeviceConfigurationKeyValueConstraintsListDataType     `json:"deviceConfigurationKeyValueConstraintsListData,omitempty"`
	DeviceConfigurationKeyValueDescriptionListData     *DeviceConfigurationKeyValueDescriptionListDataType     `json:"deviceConfigurationKeyValueDescriptionListData,omitempty" eebus:"fct:deviceConfigurationKeyValueDescriptionListData"`
	DeviceConfigurationKeyValueListData                *DeviceConfigurationKeyValueListDataType                `json:"deviceConfigurationKeyValueListData,omitempty" eebus:"fct:deviceConfigurationKeyValueListData"`
	DeviceDiagnosisHeartbeatData                       *DeviceDiagnosisHeartbeatDataType                       `json:"deviceDiagnosisHeartbeatData,omitempty" eebus:"fct:deviceDiagnosisHeartbeatData"`
	DeviceDiagnosisServiceData                         *DeviceDiagnosisServiceDataType                         `json:"deviceDiagnosisServiceData,omitempty" eebus:"fct:deviceDiagnosisServiceData"`
	DeviceDiagnosisStateData                           *DeviceDiagnosisStateDataType                           `json:"deviceDiagnosisStateData,omitempty" eebus:"fct:deviceDiagnosisStateData"`
	DirectControlActivityListData                      *DirectControlActivityListDataType                      `json:"directControlActivityListData,omitempty" eebus:"fct:directControlActivityListData"`
	DirectControlDescriptionData                       *DirectControlDescriptionDataType                       `json:"directControlDescriptionData,omitempty" eebus:"fct:directControlDescriptionData"`
	ElectricalConnectionDescriptionListData            *ElectricalConnectionDescriptionListDataType            `json:"electricalConnectionDescriptionListData,omitempty" eebus:"fct:electricalConnectionDescriptionListData"`
	ElectricalConnectionParameterDescriptionListData   *ElectricalConnectionParameterDescriptionListDataType   `json:"electricalConnectionParameterDescriptionListData,omitempty" eebus:"fct:electricalConnectionParameterDescriptionListData"`
	ElectricalConnectionPermittedValueSetListData      *ElectricalConnectionPermittedValueSetListDataType      `json:"electricalConnectionPermittedValueSetListData,omitempty" eebus:"fct:electricalConnectionPermittedValueSetListData"`
	ElectricalConnectionStateListData                  *ElectricalConnectionStateListDataType                  `json:"electricalConnectionStateListData,omitempty" eebus:"fct:electricalConnectionStateListData"`
	HvacOperationModeDescriptionListData               *HvacOperationModeDescriptionListDataType               `json:"hvacOperationModeDescriptionListData,omitempty" eebus:"fct:hvacOperationModeDescriptionListData"`
	HvacOverrunDescriptionListData                     *HvacOverrunDescriptionListDataType                     `json:"hvacOverrunDescriptionListData,omitempty"`
	HvacOverrunListData                                *HvacOverrunListDataType                                `json:"hvacOverrunListData,omitempty"`
	HvacSystemFunctionDescriptionListData              *HvacSystemFunctionDescriptionListDataType              `json:"hvacSystemFunctionDescriptionListData,omitempty" eebus:"fct:hvacSystemFunctionDescriptionListData"`
	HvacSystemFunctionListData                         *HvacSystemFunctionListDataType                         `json:"hvacSystemFunctionListData,omitempty" eebus:"fct:hvacSystemFunctionListData"`
	HvacSystemFunctionOperationModeRelationListData    *HvacSystemFunctionOperationModeRelationListDataType    `json:"hvacSystemFunctionOperationModeRelationListData,omitempty" eebus:"fct:hvacSystemFunctionOperationModeRelationListData"`
	HvacSystemFunctionPowerSequenceRelationListData    *HvacSystemFunctionPowerSequenceRelationListDataType    `json:"hvacSystemFunctionPowerSequenceRelationListData,omitempty" eebus:"fct:hvacSystemFunctionPowerSequenceRelationListData"`
	HvacSystemFunctionSetpointRelationListData         *HvacSystemFunctionSetpointRelationListDataType         `json:"hvacSystemFunctionSetpointRelationListData,omitempty" eebus:"fct:hvacSystemFunctionSetpointRelationListData"`
	IdentificationListData                             *IdentificationListDataType                             `json:"identificationListData,omitempty" eebus:"fct:identificationListData"`
	IncentiveDescriptionListData                       *IncentiveDescriptionListDataType                       `json:"incentiveDescriptionListData,omitempty" eebus:"fct:incentiveDescriptionListData"`
	IncentiveListData                                  *IncentiveListDataType                                  `json:"incentiveListData,omitempty" eebus:"fct:incentiveListData"`
	IncentiveTableConstraintsData                      *IncentiveTableConstraintsDataType                      `json:"incentiveTableConstraintsData,omitempty" eebus:"fct:incentiveTableConstraintsData"`
	IncentiveTableData                                 *IncentiveTableDataType                                 `json:"incentiveTableData,omitempty" eebus:"fct:incentiveTableData"`
	IncentiveTableDescriptionData                      *IncentiveTableDescriptionDataType                      `json:"incentiveTableDescriptionData,omitempty" eebus:"fct:incentiveTableDescriptionData"`
	LoadControlEventListData                           *LoadControlEventListDataType                           `json:"loadControlEventListData,omitempty" eebus:"fct:loadControlEventListData"`
	LoadControlLimitConstraintsListData                *LoadControlLimitConstraintsListDataType                `json:"loadControlLimitConstraintsListData,omitempty" eebus:"fct:loadControlLimitConstraintsListData"`
	LoadControlLimitDescriptionListData                *LoadControlLimitDescriptionListDataType                `json:"loadControlLimitDescriptionListData,omitempty" eebus:"fct:loadControlLimitDescriptionListData"`
	LoadControlLimitListData                           *LoadControlLimitListDataType                           `json:"loadControlLimitListData,omitempty" eebus:"fct:loadControlLimitListData"`
	LoadControlNodeData                                *LoadControlNodeDataType                                `json:"loadControlNodeData,omitempty" eebus:"fct:loadControlNodeData"`
	LoadControlStateListData                           *LoadControlStateListDataType                           `json:"loadControlStateListData,omitempty" eebus:"fct:loadControlStateListData"`
	MeasurementConstraintsListData                     *MeasurementConstraintsListDataType                     `json:"measurementConstraintsListData,omitempty" eebus:"fct:measurementConstraintsListData"`
	MeasurementDescriptionListData                     *MeasurementDescriptionListDataType                     `json:"measurementDescriptionListData,omitempty" eebus:"fct:measurementDescriptionListData"`
	MeasurementListData                                *MeasurementListDataType                                `json:"measurementListData,omitempty" eebus:"fct:measurementListData"`
	MeasurementThresholdRelationListData               *MeasurementThresholdRelationListDataType               `json:"measurementThresholdRelationListData,omitempty" eebus:"fct:measurementThresholdRelationListData"`
	MessagingListData                                  *MessagingListDataType                                  `json:"messagingListData,omitempty" eebus:"fct:messagingListData"`
	NetworkManagementAbortCall                         *NetworkManagementAbortCallType                         `json:"networkManagementAbortCall,omitempty" eebus:"fct:networkManagementAbortCall"`
	NetworkManagementAddNodeCall                       *NetworkManagementAddNodeCallType                       `json:"networkManagementAddNodeCall,omitempty" eebus:"fct:networkManagementAddNodeCall"`
	NetworkManagementDeviceDescriptionListData         *NetworkManagementDeviceDescriptionListDataType         `json:"networkManagementDeviceDescriptionListData,omitempty" eebus:"fct:networkManagementDeviceDescriptionListData"`
	NetworkManagementDiscoverCall                      *NetworkManagementDiscoverCallType                      `json:"networkManagementDiscoverCall,omitempty" eebus:"fct:networkManagementDiscoverCall"`
	NetworkManagementEntityDescriptionListData         *NetworkManagementEntityDescriptionListDataType         `json:"networkManagementEntityDescriptionListData,omitempty" eebus:"fct:networkManagementEntityDescriptionListData"`
	NetworkManagementFeatureDescriptionListData        *NetworkManagementFeatureDescriptionListDataType        `json:"networkManagementFeatureDescriptionListData,omitempty" eebus:"fct:networkManagementFeatureDescriptionListData"`
	NetworkManagementJoiningModeData                   *NetworkManagementJoiningModeDataType                   `json:"networkManagementJoiningModeData,omitempty" eebus:"fct:networkManagementJoiningModeData"`
	NetworkManagementModifyNodeCall                    *NetworkManagementModifyNodeCallType                    `json:"networkManagementModifyNodeCall,omitempty" eebus:"fct:networkManagementModifyNodeCall"`
	NetworkManagementProcessStateData                  *NetworkManagementProcessStateDataType                  `json:"networkManagementProcessStateData,omitempty" eebus:"fct:networkManagementProcessStateData"`
	NetworkManagementRemoveNodeCall                    *NetworkManagementRemoveNodeCallType                    `json:"networkManagementRemoveNodeCall,omitempty" eebus:"fct:networkManagementRemoveNodeCall"`
	NetworkManagementReportCandidateData               *NetworkManagementReportCandidateDataType               `json:"networkManagementReportCandidateData,omitempty" eebus:"fct:networkManagementReportCandidateData"`
	NetworkManagementScanNetworkCall                   *NetworkManagementScanNetworkCallType                   `json:"networkManagementScanNetworkCall,omitempty" eebus:"fct:networkManagementScanNetworkCall"`
	NodeManagementBindingData                          *NodeManagementBindingDataType                          `json:"nodeManagementBindingData,omitempty" eebus:"fct:nodeManagementBindingData"`
	NodeManagementBindingDeleteCall                    *NodeManagementBindingDeleteCallType                    `json:"nodeManagementBindingDeleteCall,omitempty" eebus:"fct:nodeManagementBindingDeleteCall"`
	NodeManagementBindingRequestCall                   *NodeManagementBindingRequestCallType                   `json:"nodeManagementBindingRequestCall,omitempty" eebus:"fct:nodeManagementBindingRequestCall"`
	NodeManagementDestinationListData                  *NodeManagementDestinationListDataType                  `json:"nodeManagementDestinationListData,omitempty" eebus:"fct:nodeManagementDestinationListData"`
	NodeManagementDetailedDiscoveryData                *NodeManagementDetailedDiscoveryDataType                `json:"nodeManagementDetailedDiscoveryData,omitempty" eebus:"fct:nodeManagementDetailedDiscoveryData"`
	NodeManagementSubscriptionData                     *NodeManagementSubscriptionDataType                     `json:"nodeManagementSubscriptionData,omitempty" eebus:"fct:nodeManagementSubscriptionData"`
	NodeManagementSubscriptionDeleteCall               *NodeManagementSubscriptionDeleteCallType               `json:"nodeManagementSubscriptionDeleteCall,omitempty" eebus:"fct:nodeManagementSubscriptionDeleteCall"`
	NodeManagementSubscriptionRequestCall              *NodeManagementSubscriptionRequestCallType              `json:"nodeManagementSubscriptionRequestCall,omitempty" eebus:"fct:nodeManagementSubscriptionRequestCall"`
	NodeManagementUseCaseData                          *NodeManagementUseCaseDataType                          `json:"nodeManagementUseCaseData,omitempty" eebus:"fct:nodeManagementUseCaseData"`
	OperatingConstraintsDurationListData               *OperatingConstraintsDurationListDataType               `json:"operatingConstraintsDurationListData,omitempty" eebus:"fct:operatingConstraintsDurationListData"`
	OperatingConstraintsInterruptListData              *OperatingConstraintsInterruptListDataType              `json:"operatingConstraintsInterruptListData,omitempty" eebus:"fct:operatingConstraintsInterruptListData"`
	OperatingConstraintsPowerDescriptionListData       *OperatingConstraintsPowerDescriptionListDataType       `json:"operatingConstraintsPowerDescriptionListData,omitempty" eebus:"fct:operatingConstraintsPowerDescriptionListData"`
	OperatingConstraintsPowerLevelListData             *OperatingConstraintsPowerLevelListDataType             `json:"operatingConstraintsPowerLevelListData,omitempty" eebus:"fct:operatingConstraintsPowerLevelListData"`
	OperatingConstraintsPowerRangeListData             *OperatingConstraintsPowerRangeListDataType             `json:"operatingConstraintsPowerRangeListData,omitempty" eebus:"fct:operatingConstraintsPowerRangeListData"`
	OperatingConstraintsResumeImplicationListData      *OperatingConstraintsResumeImplicationListDataType      `json:"operatingConstraintsResumeImplicationListData,omitempty" eebus:"fct:operatingConstraintsResumeImplicationListData"`
	PowerSequenceAlternativesRelationListData          *PowerSequenceAlternativesRelationListDataType          `json:"powerSequenceAlternativesRelationListData,omitempty" eebus:"fct:powerSequenceAlternativesRelationListData"`
	PowerSequenceDescriptionListData                   *PowerSequenceDescriptionListDataType                   `json:"powerSequenceDescriptionListData,omitempty" eebus:"fct:powerSequenceDescriptionListData"`
	PowerSequenceNodeScheduleInformationData           *PowerSequenceNodeScheduleInformationDataType           `json:"powerSequenceNodeScheduleInformationData,omitempty" eebus:"fct:powerSequenceNodeScheduleInformationData"`
	PowerSequencePriceCalculationRequestCall           *PowerSequencePriceCalculationRequestCallType           `json:"powerSequencePriceCalculationRequestCall,omitempty" eebus:"fct:powerSequencePriceCalculationRequestCall"`
	PowerSequencePriceListData                         *PowerSequencePriceListDataType                         `json:"powerSequencePriceListData,omitempty" eebus:"fct:powerSequencePriceListData"`
	PowerSequenceScheduleConfigurationRequestCall      *PowerSequenceScheduleConfigurationRequestCallType      `json:"powerSequenceScheduleConfigurationRequestCall,omitempty" eebus:"fct:powerSequenceScheduleConfigurationRequestCall"`
	PowerSequenceScheduleConstraintsListData           *PowerSequenceScheduleConstraintsListDataType           `json:"powerSequenceScheduleConstraintsListData,omitempty" eebus:"fct:powerSequenceScheduleConstraintsListData"`
	PowerSequenceScheduleListData                      *PowerSequenceScheduleListDataType                      `json:"powerSequenceScheduleListData,omitempty" eebus:"fct:powerSequenceScheduleListData"`
	PowerSequenceSchedulePreferenceListData            *PowerSequenceSchedulePreferenceListDataType            `json:"powerSequenceSchedulePreferenceListData,omitempty" eebus:"fct:powerSequenceSchedulePreferenceListData"`
	PowerSequenceStateListData                         *PowerSequenceStateListDataType                         `json:"powerSequenceStateListData,omitempty" eebus:"fct:powerSequenceStateListData"`
	PowerTimeSlotScheduleConstraintsListData           *PowerTimeSlotScheduleConstraintsListDataType           `json:"powerTimeSlotScheduleConstraintsListData,omitempty" eebus:"fct:powerTimeSlotScheduleConstraintsListData"`
	PowerTimeSlotScheduleListData                      *PowerTimeSlotScheduleListDataType                      `json:"powerTimeSlotScheduleListData,omitempty" eebus:"fct:powerTimeSlotScheduleListData"`
	PowerTimeSlotValueListData                         *PowerTimeSlotValueListDataType                         `json:"powerTimeSlotValueListData,omitempty" eebus:"fct:powerTimeSlotValueListData"`
	ResultData                                         *ResultDataType                                         `json:"resultData,omitempty" eebus:"fct:resultData"`
	SensingDescriptionData                             *SensingDescriptionDataType                             `json:"sensingDescriptionData,omitempty" eebus:"fct:sensingDescriptionData"`
	SensingListData                                    *SensingListDataType                                    `json:"sensingListData,omitempty" eebus:"fct:sensingListData"`
	SetpointConstraintsListData                        *SetpointConstraintsListDataType                        `json:"setpointConstraintsListData,omitempty" eebus:"fct:setpointConstraintsListData"`
	SetpointDescriptionListData                        *SetpointDescriptionListDataType                        `json:"setpointDescriptionListData,omitempty" eebus:"fct:setpointDescriptionListData"`
	SetpointListData                                   *SetpointListDataType                                   `json:"setpointListData,omitempty" eebus:"fct:setpointListData"`
	SmartEnergyManagementPsConfigurationRequestCall    *SmartEnergyManagementPsConfigurationRequestCallType    `json:"smartEnergyManagementPsConfigurationRequestCall,omitempty" eebus:"fct:smartEnergyManagementPsConfigurationRequestCall"`
	SmartEnergyManagementPsData                        *SmartEnergyManagementPsDataType                        `json:"smartEnergyManagementPsData,omitempty" eebus:"fct:smartEnergyManagementPsData"`
	SmartEnergyManagementPsPriceCalculationRequestCall *SmartEnergyManagementPsPriceCalculationRequestCallType `json:"smartEnergyManagementPsPriceCalculationRequestCall,omitempty" eebus:"fct:smartEnergyManagementPsPriceCalculationRequestCall"`
	SmartEnergyManagementPsPriceData                   *SmartEnergyManagementPsPriceDataType                   `json:"smartEnergyManagementPsPriceData,omitempty" eebus:"fct:smartEnergyManagementPsPriceData"`
	SpecificationVersionListData                       *SpecificationVersionListDataType                       `json:"specificationVersionListData,omitempty" eebus:"fct:specificationVersionListData"`
	SubscriptionManagementDeleteCall                   *SubscriptionManagementDeleteCallType                   `json:"subscriptionManagementDeleteCall,omitempty" eebus:"fct:subscriptionManagementDeleteCall"`
	SubscriptionManagementEntryListData                *SubscriptionManagementEntryListDataType                `json:"subscriptionManagementEntryListData,omitempty" eebus:"fct:subscriptionManagementEntryListData"`
	SubscriptionManagementRequestCall                  *SubscriptionManagementRequestCallType                  `json:"subscriptionManagementRequestCall,omitempty" eebus:"fct:subscriptionManagementRequestCall"`
	SupplyConditionDescriptionListData                 *SupplyConditionDescriptionListDataType                 `json:"supplyConditionDescriptionListData,omitempty" eebus:"fct:supplyConditionDescriptionListData"`
	SupplyConditionListData                            *SupplyConditionListDataType                            `json:"supplyConditionListData,omitempty" eebus:"fct:supplyConditionListData"`
	SupplyConditionThresholdRelationListData           *SupplyConditionThresholdRelationListDataType           `json:"supplyConditionThresholdRelationListData,omitempty" eebus:"fct:supplyConditionThresholdRelationListData"`
	TariffBoundaryRelationListData                     *TariffBoundaryRelationListDataType                     `json:"tariffBoundaryRelationListData,omitempty" eebus:"fct:tariffBoundaryRelationListData"`
	TariffDescriptionListData                          *TariffDescriptionListDataType                          `json:"tariffDescriptionListData,omitempty" eebus:"fct:tariffDescriptionListData"`
	TariffListData                                     *TariffListDataType                                     `json:"tariffListData,omitempty" eebus:"fct:tariffListData"`
	TariffOverallConstraintsData                       *TariffOverallConstraintsDataType                       `json:"tariffOverallConstraintsData,omitempty" eebus:"fct:tariffOverallConstraintsData"`
	TariffTierRelationListData                         *TariffTierRelationListDataType                         `json:"tariffTierRelationListData,omitempty" eebus:"fct:tariffTierRelationListData"`
	TaskManagementJobDescriptionListData               *TaskManagementJobDescriptionListDataType               `json:"taskManagementJobDescriptionListData,omitempty" eebus:"fct:taskManagementJobDescriptionListData"`
	TaskManagementJobListData                          *TaskManagementJobListDataType                          `json:"taskManagementJobListData,omitempty" eebus:"fct:taskManagementJobListData"`
	TaskManagementJobRelationListData                  *TaskManagementJobRelationListDataType                  `json:"taskManagementJobRelationListData,omitempty" eebus:"fct:taskManagementJobRelationListData"`
	TaskManagementOverviewData                         *TaskManagementOverviewDataType                         `json:"taskManagementOverviewData,omitempty" eebus:"fct:taskManagementOverviewData"`
	ThresholdConstraintsListData                       *ThresholdConstraintsListDataType                       `json:"thresholdConstraintsListData,omitempty" eebus:"fct:thresholdConstraintsListData"`
	ThresholdDescriptionListData                       *ThresholdDescriptionListDataType                       `json:"thresholdDescriptionListData,omitempty" eebus:"fct:thresholdDescriptionListData"`
	ThresholdListData                                  *ThresholdListDataType                                  `json:"thresholdListData,omitempty" eebus:"fct:thresholdListData"`
	TierBoundaryDescriptionListData                    *TierBoundaryDescriptionListDataType                    `json:"tierBoundaryDescriptionListData,omitempty" eebus:"fct:tierBoundaryDescriptionListData"`
	TierBoundaryListData                               *TierBoundaryListDataType                               `json:"tierBoundaryListData,omitempty" eebus:"fct:tierBoundaryListData"`
	TierDescriptionListData                            *TierDescriptionListDataType                            `json:"tierDescriptionListData,omitempty" eebus:"fct:tierDescriptionListData"`
	TierIncentiveRelationListData                      *TierIncentiveRelationListDataType                      `json:"tierIncentiveRelationListData,omitempty" eebus:"fct:tierIncentiveRelationListData"`
	TierListData                                       *TierListDataType                                       `json:"tierListData,omitempty" eebus:"fct:tierListData"`
	TimeDistributorData                                *TimeDistributorDataType                                `json:"timeDistributorData,omitempty" eebus:"fct:timeDistributorData"`
	TimeDistributorEnquiryCall                         *TimeDistributorEnquiryCallType                         `json:"timeDistributorEnquiryCall,omitempty" eebus:"fct:timeDistributorEnquiryCall"`
	TimeInformationData                                *TimeInformationDataType                                `json:"timeInformationData,omitempty" eebus:"fct:timeInformationData"`
	TimePrecisionData                                  *TimePrecisionDataType                                  `json:"timePrecisionData,omitempty" eebus:"fct:timePrecisionData"`
	TimeSeriesConstraintsListData                      *TimeSeriesConstraintsListDataType                      `json:"timeSeriesConstraintsListData,omitempty" eebus:"fct:timeSeriesConstraintsListData"`
	TimeSeriesDescriptionListData                      *TimeSeriesDescriptionListDataType                      `json:"timeSeriesDescriptionListData,omitempty" eebus:"fct:timeSeriesDescriptionListData"`
	TimeSeriesListData                                 *TimeSeriesListDataType                                 `json:"timeSeriesListData,omitempty" eebus:"fct:timeSeriesListData"`
	TimeTableConstraintsListData                       *TimeTableConstraintsListDataType                       `json:"timeTableConstraintsListData,omitempty" eebus:"fct:timeTableConstraintsListData"`
	TimeTableDescriptionListData                       *TimeTableDescriptionListDataType                       `json:"timeTableDescriptionListData,omitempty" eebus:"fct:timeTableDescriptionListData"`
	TimeTableListData                                  *TimeTableListDataType                                  `json:"timeTableListData,omitempty" eebus:"fct:timeTableListData"`
	UseCaseInformationListData                         *UseCaseInformationListDataType                         `json:"useCaseInformationListData,omitempty" eebus:"fct:useCaseInformationListData"`

	// DataExtendGroup
	ManufacturerSpecificExtension *string                     `json:"manufacturerSpecificExtension,omitempty"`
	LastUpdateAt                  *AbsoluteOrRelativeTimeType `json:"lastUpdateAt,omitempty"`
}
