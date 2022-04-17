package model

type PurposeIdType string

type ChannelIdType uint

type DataTunnelingHeaderType struct {
	PurposeId  *PurposeIdType `json:"purposeId,omitempty"`
	ChannelId  *ChannelIdType `json:"channelId,omitempty"`
	SequenceId *uint          `json:"sequenceId,omitempty"`
}

type DataTunnelingHeaderElementsType struct {
	PurposeId  *ElementTagType `json:"purposeId,omitempty"`
	ChannelId  *ElementTagType `json:"channelId,omitempty"`
	SequenceId *ElementTagType `json:"sequenceId,omitempty"`
}

type DataTunnelingCallType struct {
	Header  *DataTunnelingHeaderType `json:"header,omitempty"`
	Payload *string                  `json:"payload,omitempty"`
}

type DataTunnelingCallElementsType struct {
	Header  *DataTunnelingHeaderElementsType `json:"header,omitempty"`
	Payload *ElementTagType                  `json:"payload,omitempty"`
}
