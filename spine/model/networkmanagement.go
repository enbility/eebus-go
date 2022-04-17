package model

type NetworkManagementNativeSetupType string

type NetworkManagementScanSetupType string

type NetworkManagementSetupType string

type NetworkManagementCandidateSetupType string

type NetworkManagementTechnologyAddressType string

type NetworkManagementCommunicationsTechnologyInformationType string

type NetworkManagementMinimumTrustLevelType string

type NetworkManagementProcessTimeoutType string

type NetworkManagementFeatureSetType string

const (
	NetworkManagementFeatureSetTypeGateway NetworkManagementFeatureSetType = "gateway"
	NetworkManagementFeatureSetTypeRouter  NetworkManagementFeatureSetType = "router"
	NetworkManagementFeatureSetTypeSmart   NetworkManagementFeatureSetType = "smart"
	NetworkManagementFeatureSetTypeSimple  NetworkManagementFeatureSetType = "simple"
)

type NetworkManagementProcessStateStateType string

const (
	NetworkManagementProcessStateStateTypeSucceeded NetworkManagementProcessStateStateType = "succeeded"
	NetworkManagementProcessStateStateTypeFailed    NetworkManagementProcessStateStateType = "failed"
	NetworkManagementProcessStateStateTypeAborted   NetworkManagementProcessStateStateType = "aborted"
)

type NetworkManagementStateChangeType string

const (
	NetworkManagementStateChangeTypeAdded    NetworkManagementStateChangeType = "added"
	NetworkManagementStateChangeTypeRemoved  NetworkManagementStateChangeType = "removed"
	NetworkManagementStateChangeTypeModified NetworkManagementStateChangeType = "modified"
)

type NetworkManagementAddNodeCallType struct {
	NodeAddress *FeatureAddressType                  `json:"nodeAddress,omitempty"`
	NativeSetup *NetworkManagementNativeSetupType    `json:"nativeSetup,omitempty"`
	Timeout     *NetworkManagementProcessTimeoutType `json:"timeout,omitempty"`
	Label       *LabelType                           `json:"label,omitempty"`
	Description *DescriptionType                     `json:"description,omitempty"`
}

type NetworkManagementAddNodeCallElementsType struct {
	NodeAddress *FeatureAddressElementsType `json:"nodeAddress,omitempty"`
	NativeSetup *ElementTagType             `json:"nativeSetup,omitempty"`
	Timeout     *ElementTagType             `json:"timeout,omitempty"`
	Label       *ElementTagType             `json:"label,omitempty"`
	Description *ElementTagType             `json:"description,omitempty"`
}

type NetworkManagementRemoveNodeCallType struct {
	NodeAddress *FeatureAddressType                  `json:"nodeAddress,omitempty"`
	Timeout     *NetworkManagementProcessTimeoutType `json:"timeout,omitempty"`
}

type NetworkManagementRemoveNodeCallElementsType struct {
	NodeAddress *FeatureAddressElementsType `json:"nodeAddress,omitempty"`
	Timeout     *ElementTagType             `json:"timeout,omitempty"`
}

type NetworkManagementModifyNodeCallType struct {
	NodeAddress *FeatureAddressType                  `json:"nodeAddress,omitempty"`
	NativeSetup *NetworkManagementNativeSetupType    `json:"nativeSetup,omitempty"`
	Timeout     *NetworkManagementProcessTimeoutType `json:"timeout,omitempty"`
	Label       *LabelType                           `json:"label,omitempty"`
	Description *DescriptionType                     `json:"description,omitempty"`
}

type NetworkManagementModifyNodeCallElementsType struct {
	NodeAddress *FeatureAddressElementsType `json:"nodeAddress,omitempty"`
	NativeSetup *ElementTagType             `json:"nativeSetup,omitempty"`
	Timeout     *ElementTagType             `json:"timeout,omitempty"`
	Label       *ElementTagType             `json:"label,omitempty"`
	Description *ElementTagType             `json:"description,omitempty"`
}

type NetworkManagementScanNetworkCallType struct {
	ScanSetup *NetworkManagementScanSetupType      `json:"scanSetup,omitempty"`
	Timeout   *NetworkManagementProcessTimeoutType `json:"timeout,omitempty"`
}

type NetworkManagementScanNetworkCallElementsType struct {
	ScanSetup *ElementTagType `json:"scanSetup,omitempty"`
	Timeout   *ElementTagType `json:"timeout,omitempty"`
}

type NetworkManagementDiscoverCallType struct {
	DiscoverAddress *FeatureAddressType `json:"discoverAddress,omitempty"`
}

type NetworkManagementDiscoverCallElementsType struct {
	DiscoverAddress *FeatureAddressElementsType `json:"discoverAddress,omitempty"`
}

type NetworkManagementAbortCallType struct{}

type NetworkManagementAbortCallElementsType struct{}

type NetworkManagementProcessStateDataType struct {
	State       *NetworkManagementProcessStateStateType `json:"state,omitempty"`
	Description *DescriptionType                        `json:"description,omitempty"`
}

type NetworkManagementProcessStateDataElementsType struct {
	State       *ElementTagType `json:"state,omitempty"`
	Description *ElementTagType `json:"description,omitempty"`
}

type NetworkManagementJoiningModeDataType struct {
	Setup *NetworkManagementSetupType `json:"setup,omitempty"`
}

type NetworkManagementJoiningModeDataElementsType struct {
	Setup *ElementTagType `json:"setup,omitempty"`
}

type NetworkManagementReportCandidateDataType struct {
	CandidateSetup    *NetworkManagementCandidateSetupType `json:"candidateSetup,omitempty"`
	SetupUsableForAdd *bool                                `json:"setupUsableForAdd,omitempty"`
	Label             *LabelType                           `json:"label,omitempty"`
	Description       *DescriptionType                     `json:"description,omitempty"`
}

type NetworkManagementReportCandidateDataElementsType struct {
	CandidateSetup    *ElementTagType `json:"candidateSetup,omitempty"`
	SetupUsableForAdd *ElementTagType `json:"setupUsableForAdd,omitempty"`
	Label             *ElementTagType `json:"label,omitempty"`
	Description       *ElementTagType `json:"description,omitempty"`
}

type NetworkManagementDeviceDescriptionDataType struct {
	DeviceAddress                       *DeviceAddressType                                        `json:"deviceAddress,omitempty"`
	DeviceType                          *DeviceTypeType                                           `json:"deviceType,omitempty"`
	NetworkManagementResponsibleAddress *FeatureAddressType                                       `json:"networkManagementResponsibleAddress,omitempty"`
	NativeSetup                         *NetworkManagementNativeSetupType                         `json:"nativeSetup,omitempty"`
	TechnologyAddress                   *NetworkManagementTechnologyAddressType                   `json:"technologyAddress,omitempty"`
	CommunicationsTechnologyInformation *NetworkManagementCommunicationsTechnologyInformationType `json:"communicationsTechnologyInformation,omitempty"`
	NetworkFeatureSet                   *NetworkManagementFeatureSetType                          `json:"networkFeatureSet,omitempty"`
	LastStateChange                     *NetworkManagementStateChangeType                         `json:"lastStateChange,omitempty"`
	MinimumTrustLevel                   *NetworkManagementMinimumTrustLevelType                   `json:"minimumTrustLevel,omitempty"`
	Label                               *LabelType                                                `json:"label,omitempty"`
	Description                         *DescriptionType                                          `json:"description,omitempty"`
}

type NetworkManagementDeviceDescriptionDataElementsType struct {
	DeviceAddress                       *ElementTagType `json:"deviceAddress,omitempty"`
	DeviceType                          *ElementTagType `json:"deviceType,omitempty"`
	NetworkManagementResponsibleAddress *ElementTagType `json:"networkManagementResponsibleAddress,omitempty"`
	NativeSetup                         *ElementTagType `json:"nativeSetup,omitempty"`
	TechnologyAddress                   *ElementTagType `json:"technologyAddress,omitempty"`
	CommunicationsTechnologyInformation *ElementTagType `json:"communicationsTechnologyInformation,omitempty"`
	NetworkFeatureSet                   *ElementTagType `json:"networkFeatureSet,omitempty"`
	LastStateChange                     *ElementTagType `json:"lastStateChange,omitempty"`
	MinimumTrustLevel                   *ElementTagType `json:"minimumTrustLevel,omitempty"`
	Label                               *ElementTagType `json:"label,omitempty"`
	Description                         *ElementTagType `json:"description,omitempty"`
}

type NetworkManagementDeviceDescriptionListDataType struct {
	NetworkManagementDeviceDescriptionData []NetworkManagementDeviceDescriptionDataType `json:"networkManagementDeviceDescriptionData,omitempty"`
}

type NetworkManagementDeviceDescriptionListDataSelectorsType struct {
	DeviceAddress *DeviceAddressType `json:"deviceAddress,omitempty"`
	DeviceType    *DeviceTypeType    `json:"deviceType,omitempty"`
}

type NetworkManagementEntityDescriptionDataType struct {
	EntityAddress     *EntityAddressType                      `json:"entityAddress,omitempty"`
	EntityType        *EntityTypeType                         `json:"entityType,omitempty"`
	LastStateChange   *NetworkManagementStateChangeType       `json:"lastStateChange,omitempty"`
	MinimumTrustLevel *NetworkManagementMinimumTrustLevelType `json:"minimumTrustLevel,omitempty"`
	Label             *LabelType                              `json:"label,omitempty"`
	Description       *DescriptionType                        `json:"description,omitempty"`
}

type NetworkManagementEntityDescriptionDataElementsType struct {
	EntityAddress     *ElementTagType `json:"entityAddress,omitempty"`
	EntityType        *ElementTagType `json:"entityType,omitempty"`
	LastStateChange   *ElementTagType `json:"lastStateChange,omitempty"`
	MinimumTrustLevel *ElementTagType `json:"minimumTrustLevel,omitempty"`
	Label             *ElementTagType `json:"label,omitempty"`
	Description       *ElementTagType `json:"description,omitempty"`
}

type NetworkManagementEntityDescriptionListDataType struct {
	NetworkManagementEntityDescriptionData []NetworkManagementEntityDescriptionDataType `json:"networkManagementEntityDescriptionData,omitempty"`
}

type NetworkManagementEntityDescriptionListDataSelectorsType struct {
	EntityAddress *EntityAddressType `json:"entityAddress,omitempty"`
	EntityType    *EntityTypeType    `json:"entityType,omitempty"`
}

type NetworkManagementFeatureDescriptionDataType struct {
	FeatureAddress    *FeatureAddressType                     `json:"featureAddress,omitempty"`
	FeatureType       *FeatureTypeType                        `json:"featureType,omitempty"`
	SpecificUsage     []FeatureSpecificUsageType              `json:"specificUsage,omitempty"`
	FeatureGroup      *FeatureGroupType                       `json:"featureGroup,omitempty"`
	Role              *RoleType                               `json:"role,omitempty"`
	SupportedFunction []FunctionPropertyType                  `json:"supportedFunction,omitempty"`
	LastStateChange   *NetworkManagementStateChangeType       `json:"lastStateChange,omitempty"`
	MinimumTrustLevel *NetworkManagementMinimumTrustLevelType `json:"minimumTrustLevel,omitempty"`
	Label             *LabelType                              `json:"label,omitempty"`
	Description       *DescriptionType                        `json:"description,omitempty"`
	MaxResponseDelay  *MaxResponseDelayType                   `json:"maxResponseDelay,omitempty"`
}

type NetworkManagementFeatureDescriptionDataElementsType struct {
	FeatureAddress    *FeatureAddressElementsType   `json:"featureAddress,omitempty"`
	FeatureType       *ElementTagType               `json:"featureType,omitempty"`
	SpecificUsage     *ElementTagType               `json:"specificUsage,omitempty"`
	FeatureGroup      *ElementTagType               `json:"featureGroup,omitempty"`
	Role              *ElementTagType               `json:"role,omitempty"`
	SupportedFunction *FunctionPropertyElementsType `json:"supportedFunction,omitempty"`
	LastStateChange   *ElementTagType               `json:"lastStateChange,omitempty"`
	MinimumTrustLevel *ElementTagType               `json:"minimumTrustLevel,omitempty"`
	Label             *ElementTagType               `json:"label,omitempty"`
	Description       *ElementTagType               `json:"description,omitempty"`
	MaxResponseDelay  *ElementTagType               `json:"maxResponseDelay,omitempty"`
}

type NetworkManagementFeatureDescriptionListDataType struct {
	NetworkManagementFeatureDescriptionData []NetworkManagementFeatureDescriptionDataType `json:"networkManagementFeatureDescriptionData,omitempty"`
}

type NetworkManagementFeatureDescriptionListDataSelectorsType struct {
	FeatureAddress *FeatureAddressType `json:"featureAddress,omitempty"`
	FeatureType    *FeatureTypeType    `json:"featureType,omitempty"`
}
