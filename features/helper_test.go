package features_test

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/enbility/eebus-go/util"
	shipapi "github.com/enbility/ship-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
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

var _ shipapi.SpineDataConnection = (*WriteMessageHandler)(nil)

func (t *WriteMessageHandler) WriteSpineMessage(message []byte) {
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

func setupFeatures(t assert.TestingT, dataCon shipapi.SpineDataConnection, featureFunctions []featureFunctions) (spineapi.EntityLocal, spineapi.EntityRemote) {
	localDevice := spine.NewDeviceLocalImpl("TestBrandName", "TestDeviceModel", "TestSerialNumber", "TestDeviceCode",
		"TestDeviceAddress", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart, time.Second*4)
	localEntity := spine.NewEntityLocalImpl(localDevice, model.EntityTypeTypeCEM, spine.NewAddressEntityType([]uint{1}))
	localDevice.AddEntity(localEntity)

	for i, item := range featureFunctions {
		f := spine.NewFeatureLocalImpl(uint(i+1), localEntity, item.featureType, model.RoleTypeServer)
		localEntity.AddFeature(f)
	}

	remoteDeviceName := "remoteDevice"
	sender := spine.NewSender(dataCon)
	remoteDevice := spine.NewDeviceRemoteImpl(localDevice, "test", sender)
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
				Role:        util.Ptr(model.RoleTypeClient),
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

	return localEntity, remoteEntities[0]
}
