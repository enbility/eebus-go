package spine

type SubscriptionIdType uint

type SubscriptionManagementEntryDataType struct {
	SubscriptionId *SubscriptionIdType `json:"subscriptionId,omitempty"`
	ClientAddress  *FeatureAddressType `json:"clientAddress,omitempty"`
	ServerAddress  *FeatureAddressType `json:"serverAddress,omitempty"`
	Label          *LabelType          `json:"label,omitempty"`
	Description    *DescriptionType    `json:"description,omitempty"`
}

type SubscriptionManagementEntryDataElementsType struct {
	SubscriptionId *ElementTagType             `json:"subscriptionId,omitempty"`
	ClientAddress  *FeatureAddressElementsType `json:"clientAddress,omitempty"`
	ServerAddress  *FeatureAddressElementsType `json:"serverAddress,omitempty"`
	Label          *ElementTagType             `json:"label,omitempty"`
	Description    *ElementTagType             `json:"description,omitempty"`
}

type SubscriptionManagementEntryListDataType struct {
	SubscriptionManagementEntryData []SubscriptionManagementEntryDataType `json:"subscriptionManagementEntryData,omitempty"`
}

type SubscriptionManagementEntryListDataSelectorsType struct {
	SubscriptionId *SubscriptionIdType `json:"subscriptionId,omitempty"`
	ClientAddress  *FeatureAddressType `json:"clientAddress,omitempty"`
	ServerAddress  *FeatureAddressType `json:"serverAddress,omitempty"`
}

type SubscriptionManagementRequestCallType struct {
	ClientAddress     *FeatureAddressType `json:"clientAddress,omitempty"`
	ServerAddress     *FeatureAddressType `json:"serverAddress,omitempty"`
	ServerFeatureType *FeatureTypeType    `json:"serverFeatureType,omitempty"`
}

type SubscriptionManagementRequestCallElementsType struct {
	ClientAddress     *FeatureAddressElementsType `json:"clientAddress,omitempty"`
	ServerAddress     *FeatureAddressElementsType `json:"serverAddress,omitempty"`
	ServerFeatureType *ElementTagType             `json:"serverFeatureType,omitempty"`
}

type SubscriptionManagementDeleteCallType struct {
	SubscriptionId *SubscriptionIdType `json:"subscriptionId,omitempty"`
	ClientAddress  *FeatureAddressType `json:"clientAddress,omitempty"`
	ServerAddress  *FeatureAddressType `json:"serverAddress,omitempty"`
}

type SubscriptionManagementDeleteCallElementsType struct {
	SubscriptionId *ElementTagType             `json:"subscriptionId,omitempty"`
	ClientAddress  *FeatureAddressElementsType `json:"clientAddress,omitempty"`
	ServerAddress  *FeatureAddressElementsType `json:"serverAddress,omitempty"`
}
