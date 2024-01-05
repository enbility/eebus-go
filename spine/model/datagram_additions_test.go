package model

import (
	"testing"

	"github.com/enbility/eebus-go/util"
)

func TestPrintMessageOverview_Read_Send(t *testing.T) {
	datagram := &DatagramType{
		Header: HeaderType{
			AddressSource: &FeatureAddressType{
				Device: util.Ptr(AddressDeviceType("localdevice")),
			},
			AddressDestination: &FeatureAddressType{
				Device: util.Ptr(AddressDeviceType("localdevice")),
			},
			MsgCounter:    util.Ptr(MsgCounterType(1)),
			CmdClassifier: util.Ptr(CmdClassifierTypeRead),
		},
		Payload: PayloadType{
			Cmd: []CmdType{
				{},
			},
		},
	}

	datagram.PrintMessageOverview(false, "", "")
}

func TestPrintMessageOverview_Read_Recv(t *testing.T) {
	datagram := &DatagramType{
		Header: HeaderType{
			AddressSource: &FeatureAddressType{},
			AddressDestination: &FeatureAddressType{
				Device: util.Ptr(AddressDeviceType("localdevice")),
			},
			MsgCounter:    util.Ptr(MsgCounterType(1)),
			CmdClassifier: util.Ptr(CmdClassifierTypeRead),
		},
		Payload: PayloadType{
			Cmd: []CmdType{
				{},
			},
		},
	}

	datagram.PrintMessageOverview(true, "", "")
}

func TestPrintMessageOverview_Reply_Recv(t *testing.T) {
	datagram := &DatagramType{
		Header: HeaderType{
			AddressSource: &FeatureAddressType{},
			AddressDestination: &FeatureAddressType{
				Device: util.Ptr(AddressDeviceType("localdevice")),
			},
			MsgCounter:          util.Ptr(MsgCounterType(1)),
			MsgCounterReference: util.Ptr(MsgCounterType(1)),
			CmdClassifier:       util.Ptr(CmdClassifierTypeReply),
		},
		Payload: PayloadType{
			Cmd: []CmdType{
				{},
			},
		},
	}

	datagram.PrintMessageOverview(true, "", "")
}

func TestPrintMessageOverview_Result_Recv(t *testing.T) {
	datagram := &DatagramType{
		Header: HeaderType{
			AddressSource: &FeatureAddressType{},
			AddressDestination: &FeatureAddressType{
				Device: util.Ptr(AddressDeviceType("localdevice")),
			},
			MsgCounter:          util.Ptr(MsgCounterType(1)),
			MsgCounterReference: util.Ptr(MsgCounterType(1)),
			CmdClassifier:       util.Ptr(CmdClassifierTypeResult),
		},
		Payload: PayloadType{
			Cmd: []CmdType{
				{
					ResultData: &ResultDataType{
						ErrorNumber: util.Ptr(ErrorNumberType(1)),
					},
				},
			},
		},
	}

	datagram.PrintMessageOverview(true, "", "")
}

func TestPrintMessageOverview_Write_Recv(t *testing.T) {
	datagram := &DatagramType{
		Header: HeaderType{
			AddressSource: &FeatureAddressType{},
			AddressDestination: &FeatureAddressType{
				Device: util.Ptr(AddressDeviceType("localdevice")),
			},
			MsgCounter:          util.Ptr(MsgCounterType(1)),
			MsgCounterReference: util.Ptr(MsgCounterType(1)),
			CmdClassifier:       util.Ptr(CmdClassifierTypeWrite),
		},
		Payload: PayloadType{
			Cmd: []CmdType{
				{},
			},
		},
	}

	datagram.PrintMessageOverview(true, "", "")
}
