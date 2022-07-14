package spine

import (
	"errors"
	"fmt"

	"github.com/DerAndereAndi/eebus-go/spine/model"
)

// request detailed discovery data from a remote device
func (r *NodeManagementImpl) RequestDetailedDiscovery(remoteDeviceAddress *model.AddressDeviceType, sender Sender) (*model.MsgCounterType, *ErrorType) {
	rfAdress := featureAddressType(NodeManagementFeatureId, EntityAddressType(remoteDeviceAddress, DeviceInformationAddressEntity))
	cmd := model.CmdType{
		NodeManagementDetailedDiscoveryData: &model.NodeManagementDetailedDiscoveryDataType{},
	}
	return r.RequestDataBySenderAddress(cmd, sender, rfAdress, defaultMaxResponseDelay)
}

// handle incoming detailed discovery read call
func (r *NodeManagementImpl) processReadDetailedDiscoveryData(deviceRemote *DeviceRemoteImpl, requestHeader *model.HeaderType) error {
	if deviceRemote == nil {
		return errors.New("nodemanagement.readDetailedDiscoveryData: invalid deviceRemote")
	}

	var entityInformation []model.NodeManagementDetailedDiscoveryEntityInformationType
	var featureInformation []model.NodeManagementDetailedDiscoveryFeatureInformationType

	for _, e := range r.Device().Entities() {
		entityInformation = append(entityInformation, *e.Information())

		for _, f := range e.Features() {
			featureInformation = append(featureInformation, *f.Information())
		}
	}

	cmd := model.CmdType{
		NodeManagementDetailedDiscoveryData: &model.NodeManagementDetailedDiscoveryDataType{
			SpecificationVersionList: &model.NodeManagementSpecificationVersionListType{
				SpecificationVersion: []model.SpecificationVersionDataType{model.SpecificationVersionDataType(SpecificationVersion)},
			},
			DeviceInformation:  r.Device().Information(),
			EntityInformation:  entityInformation,
			FeatureInformation: featureInformation,
		},
	}

	return deviceRemote.Sender().Reply(requestHeader, r.Address(), cmd)
}

// handle incoming detailed discovery reply data
func (r *NodeManagementImpl) processReplyDetailedDiscoveryData(message *Message, data *model.NodeManagementDetailedDiscoveryDataType) error {
	remoteDevice := message.DeviceRemote

	deviceDescription := data.DeviceInformation.Description
	if deviceDescription == nil {
		return errors.New("nodemanagement.replyDetailedDiscoveryData: invalid DeviceInformation.Description")
	}

	remoteDevice.UpdateDevice(deviceDescription)
	entities, err := remoteDevice.AddEntityAndFeatures(true, data)
	if err != nil {
		return err
	}

	// publish event for remote device added
	payload := EventPayload{
		Ski:        remoteDevice.ski,
		EventType:  EventTypeDeviceChange,
		ChangeType: ElementChangeAdd,
		Device:     remoteDevice,
		Feature:    message.FeatureRemote,
		Data:       data,
	}
	Events.Publish(payload)

	// publish event for each added remote entity
	for _, entity := range entities {
		payload := EventPayload{
			Ski:        remoteDevice.ski,
			EventType:  EventTypeEntityChange,
			ChangeType: ElementChangeAdd,
			Entity:     entity,
			Data:       data,
		}
		Events.Publish(payload)
	}

	return nil
}

// handle incoming detailed discovery notify data
func (r *NodeManagementImpl) processNotifyDetailedDiscoveryData(message *Message, data *model.NodeManagementDetailedDiscoveryDataType) error {
	// is this a partial request?
	if message.FilterPartial == nil {
		return errors.New("the received NodeManagementDetailedDiscovery.notify dataset should be partial")
	}

	if data.EntityInformation == nil || len(data.EntityInformation) == 0 || data.EntityInformation[0].Description == nil || data.EntityInformation[0].Description.LastStateChange == nil {
		return errors.New("the received NodeManagementDetailedDiscovery.notify dataset is incomplete")
	}

	lastStateChange := *data.EntityInformation[0].Description.LastStateChange
	remoteDevice := message.FeatureRemote.Device()

	// is this addition?
	if lastStateChange == model.NetworkManagementStateChangeTypeAdded {
		entities, err := remoteDevice.AddEntityAndFeatures(false, data)
		if err != nil {
			return err
		}

		// publish event for each remote entity added
		for _, entity := range entities {
			payload := EventPayload{
				Ski:        remoteDevice.ski,
				EventType:  EventTypeEntityChange,
				ChangeType: ElementChangeAdd,
				Entity:     entity,
				Feature:    message.FeatureRemote,
				Data:       data,
			}
			Events.Publish(payload)
		}
	}

	// removal example:
	// {"data":[{"header":[{"protocolId":"ee1.0"}]},{"payload":{"datagram":[{"header":[{"specificationVersion":"1.1.1"},{"addressSource":[{"device":"d:_i:19667_PorscheEVSE-00016544"},{"entity":[0]},{"feature":0}]},{"addressDestination":[{"device":"EVCC_HEMS"},{"entity":[0]},{"feature":0}]},{"msgCounter":4835},{"cmdClassifier":"notify"}]},{"payload":[{"cmd":[[{"function":"nodeManagementDetailedDiscoveryData"},{"filter":[[{"cmdControl":[{"partial":[]}]}]]},{"nodeManagementDetailedDiscoveryData":[{"deviceInformation":[{"description":[{"deviceAddress":[{"device":"d:_i:19667_PorscheEVSE-00016544"}]}]}]},{"entityInformation":[[{"description":[{"entityAddress":[{"entity":[1,1]}]},{"lastStateChange":"removed"}]}]]}]}]]}]}]}}]}
	// {
	// 	"cmd": [[
	// 			{"function": "nodeManagementDetailedDiscoveryData"},
	// 			{"filter": [[{"cmdControl": [{"partial": []}]}]]},
	// 			{"nodeManagementDetailedDiscoveryData": [
	// 					{"deviceInformation": [{"description": [{"deviceAddress": [{"device": "d:_i:19667_PorscheEVSE-00016544"}]}]}]},
	// 					{"entityInformation": [[
	// 							{
	// 								"description": [
	// 									{"entityAddress": [{"entity": [1,1]}]},
	// 									{"lastStateChange": "removed"}
	// 								]
	// 							}
	// 						]]
	// 					}
	// 				]
	// 			}
	// 	]]
	// }

	// is this removal?
	if lastStateChange == model.NetworkManagementStateChangeTypeRemoved {
		for _, ei := range data.EntityInformation {
			if err := remoteDevice.CheckEntityInformation(false, ei); err != nil {
				return err
			}

			entityAddress := ei.Description.EntityAddress.Entity
			remoteDevice.RemoveByAddress(entityAddress)

			// TODO: How to identify that an entity was removed? payload is missing option for this
			// payload := EventPayload{
			// 	EventType:  EventTypeEntityChange,
			// 	ChangeType: ElementChangeRemove,
			// 	Entity:     remoteDevice,
			// }
			// Events.Publish(payload)
		}
	}

	return nil
}

// func (f *NodeManagement) announceFeatureDiscovery(e spine.Entity) error {
// 	entity := f.Entity()
// 	if entity == nil {
// 		return errors.New("announceFeatureDiscovery: entity not found")
// 	}
// 	device := entity.Device()
// 	if device == nil {
// 		return errors.New("announceFeatureDiscovery: device not found")
// 	}
// 	entities := device.Entities()
// 	if entities == nil {
// 		return errors.New("announceFeatureDiscovery: entities not found")
// 	}

// 	for _, le := range entities {
// 		for _, lf := range le.Features() {

// 			// connect client to server features
// 			for _, rf := range e.Features() {
// 				lr := lf.Role()
// 				rr := rf.Role()
// 				rolesValid := (lr == model.RoleTypeSpecial && rr == model.RoleTypeSpecial) || (lr == model.RoleTypeClient && rr == model.RoleTypeServer)
// 				if lf.Type() == rf.Type() && rolesValid {
// 					if cf, ok := lf.(spine.ClientFeature); ok {
// 						if err := cf.ServerFound(rf); err != nil {
// 							return err
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}

// 	return nil
// }

func (r *NodeManagementImpl) handleMsgDetailedDiscoveryData(message *Message, data *model.NodeManagementDetailedDiscoveryDataType) error {
	switch message.CmdClassifier {
	case model.CmdClassifierTypeRead:
		return r.processReadDetailedDiscoveryData(message.DeviceRemote, message.RequestHeader)

	case model.CmdClassifierTypeReply:
		if err := r.pendingRequests.Remove(*message.RequestHeader.MsgCounterReference); err != nil {
			return errors.New(err.String())
		}
		return r.processReplyDetailedDiscoveryData(message, data)

	case model.CmdClassifierTypeNotify:
		return r.processNotifyDetailedDiscoveryData(message, data)

	default:
		return fmt.Errorf("nodemanagement.handleDetailedDiscoveryData: NodeManagementDetailedDiscoveryData CmdClassifierType not implemented: %s", message.CmdClassifier)
	}
}
