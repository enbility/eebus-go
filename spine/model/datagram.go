package model

type Datagram struct {
	Datagram DatagramType `json:"datagram"`
}

type DatagramType struct {
	Header  HeaderType  `json:"header"`
	Payload PayloadType `json:"payload"`
}

type HeaderType struct {
	SpecificationVersion *SpecificationVersionType   `json:"specificationVersion,omitempty"`
	AddressSource        *FeatureAddressType         `json:"addressSource,omitempty"`
	AddressDestination   *FeatureAddressType         `json:"addressDestination,omitempty"`
	AddressOriginator    *FeatureAddressType         `json:"addressOriginator,omitempty"`
	MsgCounter           *MsgCounterType             `json:"msgCounter,omitempty"`
	MsgCounterReference  *MsgCounterType             `json:"msgCounterReference,omitempty"`
	CmdClassifier        *CmdClassifierType          `json:"cmdClassifier,omitempty"`
	AckRequest           *bool                       `json:"ackRequest,omitempty"`
	Timestamp            *AbsoluteOrRelativeTimeType `json:"timestamp,omitempty"`
}

type PayloadType struct {
	Cmd []CmdType `json:"cmd"`
}
