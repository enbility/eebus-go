package internal_test

import (
	"encoding/json"
	"sync"
	"time"

	shipapi "github.com/enbility/ship-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

type featureFunctions struct {
	featureType model.FeatureTypeType
	functions   []model.FunctionType
}

type WriteMessageHandler struct {
	sentMessages [][]byte

	mux sync.Mutex
}

var _ shipapi.ShipConnectionDataWriterInterface = (*WriteMessageHandler)(nil)

func (t *WriteMessageHandler) WriteShipMessageWithPayload(message []byte) {
	t.mux.Lock()
	defer t.mux.Unlock()

	t.sentMessages = append(t.sentMessages, message)
}

func (t *WriteMessageHandler) LastMessage() []byte {
	t.mux.Lock()
	defer t.mux.Unlock()

	if len(t.sentMessages) == 0 {
		return nil
	}

	return t.sentMessages[len(t.sentMessages)-1]
}

func (t *WriteMessageHandler) MessageWithReference(msgCounterReference *model.MsgCounterType) []byte {
	t.mux.Lock()
	defer t.mux.Unlock()

	var datagram model.Datagram

	for _, msg := range t.sentMessages {
		if err := json.Unmarshal(msg, &datagram); err != nil {
			return nil
		}
		if datagram.Datagram.Header.MsgCounterReference == nil {
			continue
		}
		if uint(*datagram.Datagram.Header.MsgCounterReference) != uint(*msgCounterReference) {
			continue
		}
		if datagram.Datagram.Payload.Cmd[0].ResultData != nil {
			continue
		}

		return msg
	}

	return nil
}

func (t *WriteMessageHandler) ResultWithReference(msgCounterReference *model.MsgCounterType) []byte {
	t.mux.Lock()
	defer t.mux.Unlock()

	var datagram model.Datagram

	for _, msg := range t.sentMessages {
		if err := json.Unmarshal(msg, &datagram); err != nil {
			return nil
		}
		if datagram.Datagram.Header.MsgCounterReference == nil {
			continue
		}
		if uint(*datagram.Datagram.Header.MsgCounterReference) != uint(*msgCounterReference) {
			continue
		}
		if datagram.Datagram.Payload.Cmd[0].ResultData == nil {
			continue
		}

		return msg
	}

	return nil
}

func setupFeatures(
	t assert.TestingT,
	dataCon shipapi.ShipConnectionDataWriterInterface,
	featureFunctions []featureFunctions,
) (spineapi.EntityLocalInterface, spineapi.EntityRemoteInterface) {
	localDevice := spine.NewDeviceLocal("TestBrandName", "TestDeviceModel", "TestSerialNumber", "TestDeviceCode",
		"TestDeviceAddress", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart)
	localEntity := spine.NewEntityLocal(localDevice, model.EntityTypeTypeCEM, spine.NewAddressEntityType([]uint{1}), time.Second*4)

	for i, item := range featureFunctions {
		f := spine.NewFeatureLocal(uint(i+1), localEntity, item.featureType, model.RoleTypeClient)
		localEntity.AddFeature(f)

		f = spine.NewFeatureLocal(uint(i+1), localEntity, item.featureType, model.RoleTypeServer)
		for _, function := range item.functions {
			f.AddFunctionType(function, true, true)
		}
		localEntity.AddFeature(f)
	}

	localDevice.AddEntity(localEntity)

	remoteDeviceName := "remoteDevice"
	sender := spine.NewSender(dataCon)
	remoteDevice := spine.NewDeviceRemote(localDevice, "test", sender)
	data := &model.NodeManagementDetailedDiscoveryDataType{
		DeviceInformation: &model.NodeManagementDetailedDiscoveryDeviceInformationType{
			Description: &model.NetworkManagementDeviceDescriptionDataType{
				DeviceAddress: &model.DeviceAddressType{
					Device: util.Ptr(model.AddressDeviceType(remoteDeviceName)),
				},
			},
		},
		EntityInformation: []model.NodeManagementDetailedDiscoveryEntityInformationType{
			{
				Description: &model.NetworkManagementEntityDescriptionDataType{
					EntityAddress: &model.EntityAddressType{
						Device: util.Ptr(model.AddressDeviceType(remoteDeviceName)),
						Entity: []model.AddressEntityType{1},
					},
					EntityType: util.Ptr(model.EntityTypeTypeEVSE),
				},
			},
		},
	}

	var features []model.NodeManagementDetailedDiscoveryFeatureInformationType
	for i, item := range featureFunctions {
		featureI := model.NodeManagementDetailedDiscoveryFeatureInformationType{
			Description: &model.NetworkManagementFeatureDescriptionDataType{
				FeatureAddress: &model.FeatureAddressType{
					Device:  util.Ptr(model.AddressDeviceType(remoteDeviceName)),
					Entity:  []model.AddressEntityType{1},
					Feature: util.Ptr(model.AddressFeatureType(i + 1)),
				},
				FeatureType: util.Ptr(item.featureType),
				Role:        util.Ptr(model.RoleTypeServer),
			},
		}
		var supportedFcts []model.FunctionPropertyType
		for _, function := range item.functions {
			supportedFct := model.FunctionPropertyType{
				Function: util.Ptr(function),
				PossibleOperations: &model.PossibleOperationsType{
					Read:  &model.PossibleOperationsReadType{},
					Write: &model.PossibleOperationsWriteType{},
				},
			}

			supportedFcts = append(supportedFcts, supportedFct)
		}
		featureI.Description.SupportedFunction = supportedFcts
		features = append(features, featureI)
	}
	data.FeatureInformation = features

	remoteEntities, err := remoteDevice.AddEntityAndFeatures(true, data)
	assert.Nil(t, err)
	assert.NotNil(t, remoteEntities)
	assert.NotEqual(t, 0, len(remoteEntities))
	remoteDevice.UpdateDevice(data.DeviceInformation.Description)

	localDevice.AddRemoteDeviceForSki("test", remoteDevice)

	return localEntity, remoteEntities[0]
}
