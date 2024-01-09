package model

import (
	"fmt"
)

func (d *DatagramType) PrintMessageOverview(send bool, localFeature, remoteFeature string) string {
	var result string

	transmission := "Send"
	device := ""
	if d.Header.AddressDestination != nil && d.Header.AddressDestination.Device != nil {
		device = string(*d.Header.AddressDestination.Device)
	}
	if !send {
		transmission = "Recv"
		if d.Header.AddressSource.Device != nil {
			device = string(*d.Header.AddressSource.Device)
		}
		device = fmt.Sprintf("%s:%s to %s", device, remoteFeature, localFeature)
	}

	cmdClassifier := *d.Header.CmdClassifier
	msgCounter := *d.Header.MsgCounter
	cmd := d.Payload.Cmd[0]

	switch cmdClassifier {
	case CmdClassifierTypeRead:
		result = fmt.Sprintf("%s: %s %s %d %s", transmission, device, cmdClassifier, msgCounter, cmd.DataName())
	case CmdClassifierTypeReply:
		msgCounterRef := *d.Header.MsgCounterReference
		result = fmt.Sprintf("%s: %s %s %d %d %s", transmission, device, cmdClassifier, msgCounter, msgCounterRef, cmd.DataName())
	case CmdClassifierTypeResult:
		msgCounterRef := *d.Header.MsgCounterReference
		errorNumber := *d.Payload.Cmd[0].ResultData.ErrorNumber
		result = fmt.Sprintf("%s: %s %s %d %d %s %d", transmission, device, cmdClassifier, msgCounter, msgCounterRef, cmd.DataName(), errorNumber)
	default:
		result = fmt.Sprintf("%s: %s %s %d %s", transmission, device, cmdClassifier, msgCounter, cmd.DataName())
	}

	return result
}
