package model_test

import (
	"testing"

	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
)

func TestPrintMessageOverview_Read_Send(t *testing.T) {
	datagram := &model.DatagramType{
		Header: model.HeaderType{
			AddressSource: &model.FeatureAddressType{
				Device: util.Ptr(model.AddressDeviceType("localdevice")),
			},
			AddressDestination: &model.FeatureAddressType{
				Device: util.Ptr(model.AddressDeviceType("localdevice")),
			},
			MsgCounter:    util.Ptr(model.MsgCounterType(1)),
			CmdClassifier: util.Ptr(model.CmdClassifierTypeRead),
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{
				{},
			},
		},
	}

	datagram.PrintMessageOverview(false, "", "")
}

func TestPrintMessageOverview_Read_Recv(t *testing.T) {
	datagram := &model.DatagramType{
		Header: model.HeaderType{
			AddressSource: &model.FeatureAddressType{},
			AddressDestination: &model.FeatureAddressType{
				Device: util.Ptr(model.AddressDeviceType("localdevice")),
			},
			MsgCounter:    util.Ptr(model.MsgCounterType(1)),
			CmdClassifier: util.Ptr(model.CmdClassifierTypeRead),
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{
				{},
			},
		},
	}

	datagram.PrintMessageOverview(true, "", "")
}

func TestPrintMessageOverview_Reply_Recv(t *testing.T) {
	datagram := &model.DatagramType{
		Header: model.HeaderType{
			AddressSource: &model.FeatureAddressType{},
			AddressDestination: &model.FeatureAddressType{
				Device: util.Ptr(model.AddressDeviceType("localdevice")),
			},
			MsgCounter:          util.Ptr(model.MsgCounterType(1)),
			MsgCounterReference: util.Ptr(model.MsgCounterType(1)),
			CmdClassifier:       util.Ptr(model.CmdClassifierTypeReply),
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{
				{},
			},
		},
	}

	datagram.PrintMessageOverview(true, "", "")
}

func TestPrintMessageOverview_Result_Recv(t *testing.T) {
	datagram := &model.DatagramType{
		Header: model.HeaderType{
			AddressSource: &model.FeatureAddressType{},
			AddressDestination: &model.FeatureAddressType{
				Device: util.Ptr(model.AddressDeviceType("localdevice")),
			},
			MsgCounter:          util.Ptr(model.MsgCounterType(1)),
			MsgCounterReference: util.Ptr(model.MsgCounterType(1)),
			CmdClassifier:       util.Ptr(model.CmdClassifierTypeResult),
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{
				{
					ResultData: &model.ResultDataType{
						ErrorNumber: util.Ptr(model.ErrorNumberType(1)),
					},
				},
			},
		},
	}

	datagram.PrintMessageOverview(true, "", "")
}

func TestPrintMessageOverview_Write_Recv(t *testing.T) {
	datagram := &model.DatagramType{
		Header: model.HeaderType{
			AddressSource: &model.FeatureAddressType{},
			AddressDestination: &model.FeatureAddressType{
				Device: util.Ptr(model.AddressDeviceType("localdevice")),
			},
			MsgCounter:          util.Ptr(model.MsgCounterType(1)),
			MsgCounterReference: util.Ptr(model.MsgCounterType(1)),
			CmdClassifier:       util.Ptr(model.CmdClassifierTypeWrite),
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{
				{},
			},
		},
	}

	datagram.PrintMessageOverview(true, "", "")
}
