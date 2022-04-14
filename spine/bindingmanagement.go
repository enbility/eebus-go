package spine

type BindingIdType uint

type BindingManagementEntryDataType struct {
	BindingId     *BindingIdType      `json:"bindingId,omitempty"`
	ClientAddress *FeatureAddressType `json:"clientAddress,omitempty"`
	ServerAddress *FeatureAddressType `json:"serverAddress,omitempty"`
	Label         *LabelType          `json:"label,omitempty"`
	Description   *DescriptionType    `json:"description,omitempty"`
}

type BindingManagementEntryDataElementsType struct {
	BindingId     *ElementTagType `json:"bindingId,omitempty"`
	ClientAddress *ElementTagType `json:"clientAddress,omitempty"`
	ServerAddress *ElementTagType `json:"serverAddress,omitempty"`
	Label         *ElementTagType `json:"label,omitempty"`
	Description   *ElementTagType `json:"description,omitempty"`
}

type BindingManagementEntryListDataType struct {
	BindingManagementEntryData []BindingManagementEntryDataType `json:"bindingManagementEntryData,omitempty"`
}

type BindingManagementEntryListDataSelectorsType struct {
	BindingId     *BindingIdType      `json:"bindingId,omitempty"`
	ClientAddress *FeatureAddressType `json:"clientAddress,omitempty"`
	ServerAddress *FeatureAddressType `json:"serverAddress,omitempty"`
}

type BindingManagementRequestCallType struct {
	ClientAddress     *FeatureAddressType `json:"clientAddress,omitempty"`
	ServerAddress     *FeatureAddressType `json:"serverAddress,omitempty"`
	ServerFeatureType *FeatureTypeType    `json:"serverFeatureType,omitempty"`
}

type BindingManagementRequestCallElementsType struct {
	ClientAddress     *ElementTagType `json:"clientAddress,omitempty"`
	ServerAddress     *ElementTagType `json:"serverAddress,omitempty"`
	ServerFeatureType *ElementTagType `json:"serverFeatureType,omitempty"`
}

type BindingManagementDeleteCallType struct {
	BindingId     *BindingIdType      `json:"bindingId,omitempty"`
	ClientAddress *FeatureAddressType `json:"clientAddress,omitempty"`
	ServerAddress *FeatureAddressType `json:"serverAddress,omitempty"`
}

type BindingManagementDeleteCallElementsType struct {
	BindingId     *ElementTagType `json:"bindingId,omitempty"`
	ClientAddress *ElementTagType `json:"clientAddress,omitempty"`
	ServerAddress *ElementTagType `json:"serverAddress,omitempty"`
}
