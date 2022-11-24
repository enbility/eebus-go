package util_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/DerAndereAndi/eebus-go/ship/model"
	"github.com/DerAndereAndi/eebus-go/ship/util"
)

func TestJsonFromEEBUSJson(t *testing.T) {
	jsonTest := `{"datagram":[{"header":[{"specificationVersion":"1.2.0"},{"addressSource":[{"device":"d:_i:3210_EVSE"},{"entity":[1,1]},{"feature":6}]},{"addressDestination":[{"device":"d:_i:3210_HEMS"},{"entity":[1]},{"feature":1}]},{"msgCounter":194},{"msgCounterReference":4890},{"cmdClassifier":"reply"}]},{"payload":[{"cmd":[[{"deviceClassificationManufacturerData":[{"deviceName":""},{"deviceCode":""},{"brandName":""},{"powerSource":"mains3Phase"}]}]]}]}]}`
	jsonExpected := `{"datagram":{"header":{"specificationVersion":"1.2.0","addressSource":{"device":"d:_i:3210_EVSE","entity":[1,1],"feature":6},"addressDestination":{"device":"d:_i:3210_HEMS","entity":[1],"feature":1},"msgCounter":194,"msgCounterReference":4890,"cmdClassifier":"reply"},"payload":{"cmd":[{"deviceClassificationManufacturerData":{"deviceName":"","deviceCode":"","brandName":"","powerSource":"mains3Phase"}}]}}}`

	var json = util.JsonFromEEBUSJson([]byte(jsonTest))

	if string(json) != jsonExpected {
		t.Errorf("\nExpected:\n  %s\ngot:\n  %s", jsonExpected, json)
	}
}

func TestJsonIntoEEBUSJson(t *testing.T) {
	jsonTest := `{"datagram":{"header":{"specificationVersion":"1.2.0","addressSource":{"device":"d:_i:3210_EVSE","entity":[1,1],"feature":6},"addressDestination":{"device":"d:_i:3210_HEMS","entity":[1],"feature":1},"msgCounter":194,"msgCounterReference":4890,"cmdClassifier":"reply"},"payload":{"cmd":[{"deviceClassificationManufacturerData":{"deviceName":"","deviceCode":"","brandName":"","powerSource":"mains3Phase"}}]}}}`
	jsonExpected := `{"datagram":[{"header":[{"specificationVersion":"1.2.0"},{"addressSource":[{"device":"d:_i:3210_EVSE"},{"entity":[1,1]},{"feature":6}]},{"addressDestination":[{"device":"d:_i:3210_HEMS"},{"entity":[1]},{"feature":1}]},{"msgCounter":194},{"msgCounterReference":4890},{"cmdClassifier":"reply"}]},{"payload":[{"cmd":[[{"deviceClassificationManufacturerData":[{"deviceName":""},{"deviceCode":""},{"brandName":""},{"powerSource":"mains3Phase"}]}]]}]}]}`

	var json, err = util.JsonIntoEEBUSJson([]byte(jsonTest))
	if err != nil {
		println(err.Error())
		t.Errorf("\nExpected:\n  %s\ngot:\n  %s", jsonExpected, json)
	}

	if json != jsonExpected {
		t.Errorf("\nExpected:\n  %s\ngot:\n  %s", jsonExpected, json)
	}
}

const payloadPlaceholder = `{"place":"holder"}`

func TestShipJsonIntoEEBUSJson(t *testing.T) {
	spineTest := `{"datagram":{"header":{"specificationVersion":"1.2.0","addressSource":{"device":"Demo-EVSE-234567890","entity":[0],"feature":0},"addressDestination":{"device":"Demo-HEMS-123456789","entity":[0],"feature":0},"msgCounter":1,"cmdClassifier":"read"},"payload":{"cmd":[{"nodeManagementDetailedDiscoveryData":{}}]}}}`
	jsonExpected := `{"data":[{"header":[{"protocolId":"ee1.0"}]},{"payload":{"datagram":[{"header":[{"specificationVersion":"1.2.0"},{"addressSource":[{"device":"Demo-EVSE-234567890"},{"entity":[0]},{"feature":0}]},{"addressDestination":[{"device":"Demo-HEMS-123456789"},{"entity":[0]},{"feature":0}]},{"msgCounter":1},{"cmdClassifier":"read"}]},{"payload":[{"cmd":[[{"nodeManagementDetailedDiscoveryData":[]}]]}]}]}}]}`

	// TODO: move this test into connection_test using "transformSpineDataIntoShipJson()"
	spineMsg, err := util.JsonIntoEEBUSJson([]byte(spineTest))
	if err != nil {
		t.Errorf(err.Error())
	}
	payload := json.RawMessage([]byte(spineMsg))

	shipMessage := model.ShipData{
		Data: model.DataType{
			Header: model.HeaderType{
				ProtocolId: model.ShipProtocolId,
			},
			Payload: json.RawMessage([]byte(payloadPlaceholder)),
		},
	}

	msg, err := json.Marshal(shipMessage)
	if err != nil {
		t.Errorf(err.Error())
	}

	json, err := util.JsonIntoEEBUSJson(msg)
	if err != nil {
		println(err.Error())
		t.Errorf("\nExpected:\n  %s\ngot:\n  %s", jsonExpected, json)
	}

	json = strings.ReplaceAll(json, `[`+payloadPlaceholder+`]`, string(payload))

	if json != jsonExpected {
		t.Errorf("\nExpected:\n  %s\ngot:\n  %s", jsonExpected, json)
	}
}
