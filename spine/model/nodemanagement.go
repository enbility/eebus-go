package model

type NodeManagementSpecificationVersionListType struct {
	SpecificationVersion []SpecificationVersionDataType `json:"specificationVersion,omitempty"`
}
type NodeManagementSpecificationVersionListElementsType struct {
	SpecificationVersion *SpecificationVersionDataElementsType `json:"specificationVersion,omitempty"`
}

type NodeManagementDetailedDiscoveryDeviceInformationType struct {
	Description *NetworkManagementDeviceDescriptionDataType `json:"description,omitempty"`
}

type NodeManagementDetailedDiscoveryDeviceInformationElementsType struct {
	Description *NetworkManagementDeviceDescriptionDataElementsType `json:"description,omitempty"`
}

type NodeManagementDetailedDiscoveryEntityInformationType struct {
	Description *NetworkManagementEntityDescriptionDataType `json:"description,omitempty"`
}

type NodeManagementDetailedDiscoveryEntityInformationElementsType struct {
	Description *NetworkManagementEntityDescriptionDataElementsType `json:"description,omitempty"`
}

type NodeManagementDetailedDiscoveryFeatureInformationType struct {
	Description *NetworkManagementFeatureDescriptionDataType `json:"description,omitempty"`
}

type NodeManagementDetailedDiscoveryFeatureInformationElementsType struct {
	Description *NetworkManagementFeatureDescriptionDataElementsType `json:"description,omitempty"`
}

type NodeManagementDetailedDiscoveryDataType struct {
	SpecificationVersionList *NodeManagementSpecificationVersionListType             `json:"specificationVersionList,omitempty"`
	DeviceInformation        *NodeManagementDetailedDiscoveryDeviceInformationType   `json:"deviceInformation,omitempty"`
	EntityInformation        []NodeManagementDetailedDiscoveryEntityInformationType  `json:"entityInformation,omitempty"`
	FeatureInformation       []NodeManagementDetailedDiscoveryFeatureInformationType `json:"featureInformation,omitempty"`
}

type NodeManagementDetailedDiscoveryDataElementsType struct {
	SpecificationVersionList *NodeManagementSpecificationVersionListElementsType            `json:"specificationVersionList,omitempty"`
	DeviceInformation        *NodeManagementDetailedDiscoveryDeviceInformationElementsType  `json:"deviceInformation,omitempty"`
	EntityInformation        *NodeManagementDetailedDiscoveryEntityInformationElementsType  `json:"entityInformation,omitempty"`
	FeatureInformation       *NodeManagementDetailedDiscoveryFeatureInformationElementsType `json:"featureInformation,omitempty"`
}

type NodeManagementDetailedDiscoveryDataSelectorsType struct {
	DeviceInformation  *NetworkManagementDeviceDescriptionListDataSelectorsType  `json:"deviceInformation,omitempty"`
	EntityInformation  *NetworkManagementEntityDescriptionListDataSelectorsType  `json:"entityInformation,omitempty"`
	FeatureInformation *NetworkManagementFeatureDescriptionListDataSelectorsType `json:"featureInformation,omitempty"`
}

type NodeManagementBindingDataType struct {
	BindingEntry []BindingManagementEntryDataType `json:"bindingEntry,omitempty"`
}

type NodeManagementBindingDataElementsType struct {
	BindingEntry *BindingManagementEntryDataElementsType `json:"bindingEntry,omitempty"`
}

type NodeManagementBindingDataSelectorsType struct {
	BindingEntry *BindingManagementEntryListDataSelectorsType `json:"bindingEntry,omitempty"`
}

type NodeManagementBindingRequestCallType struct {
	BindingRequest *BindingManagementRequestCallType `json:"bindingRequest,omitempty"`
}

type NodeManagementBindingRequestCallElementsType struct {
	BindingRequest *BindingManagementRequestCallElementsType `json:"bindingRequest,omitempty"`
}

type NodeManagementBindingDeleteCallType struct {
	BindingDelete *BindingManagementDeleteCallType `json:"bindingDelete,omitempty"`
}

type NodeManagementBindingDeleteCallElementsType struct {
	BindingDelete *BindingManagementDeleteCallElementsType `json:"bindingDelete,omitempty"`
}

type NodeManagementSubscriptionDataType struct {
	SubscriptionEntry []SubscriptionManagementEntryDataType `json:"subscriptionEntry,omitempty"`
}

type NodeManagementSubscriptionDataElementsType struct {
	SubscriptionEntry *SubscriptionManagementEntryDataElementsType `json:"subscriptionEntry,omitempty"`
}

type NodeManagementSubscriptionDataSelectorsType struct {
	SubscriptionEntry *SubscriptionManagementEntryListDataSelectorsType `json:"subscriptionEntry,omitempty"`
}

type NodeManagementSubscriptionRequestCallType struct {
	SubscriptionRequest *SubscriptionManagementRequestCallType `json:"subscriptionRequest,omitempty"`
}

type NodeManagementSubscriptionRequestCallElementsType struct {
	SubscriptionRequest *SubscriptionManagementRequestCallElementsType `json:"subscriptionRequest,omitempty"`
}

type NodeManagementSubscriptionDeleteCallType struct {
	SubscriptionDelete *SubscriptionManagementDeleteCallType `json:"subscriptionDelete,omitempty"`
}

type NodeManagementSubscriptionDeleteCallElementsType struct {
	SubscriptionDelete *SubscriptionManagementDeleteCallElementsType `json:"subscriptionDelete,omitempty"`
}

type NodeManagementDestinationDataType struct {
	DeviceDescription *NetworkManagementDeviceDescriptionDataType `json:"deviceDescription,omitempty"`
}

type NodeManagementDestinationDataElementsType struct {
	DeviceDescription *NetworkManagementDeviceDescriptionDataElementsType `json:"deviceDescription,omitempty"`
}

type NodeManagementDestinationListDataType struct {
	NodeManagementDestinationData []NodeManagementDestinationDataType `json:"nodeManagementDestinationData,omitempty"`
}

type NodeManagementDestinationListDataSelectorsType struct {
	DeviceDescription *NetworkManagementDeviceDescriptionListDataSelectorsType `json:"deviceDescription,omitempty"`
}

type NodeManagementUseCaseDataType struct {
	UseCaseInformation []UseCaseInformationDataType `json:"useCaseInformation,omitempty"`
}

type NodeManagementUseCaseDataElementsType struct {
	UseCaseInformation *UseCaseInformationDataElementsType `json:"useCaseInformation,omitempty"`
}

type NodeManagementUseCaseDataSelectorsType struct {
	UseCaseInformation *UseCaseInformationListDataSelectorsType `json:"useCaseInformation,omitempty"`
}
