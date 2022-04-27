package spine

import (
	"errors"
	"fmt"

	"github.com/DerAndereAndi/eebus-go/spine/model"
)

func (r *NodeManagementImpl) RequestDetailedDiscovery(remoteDeviceAddress *model.AddressDeviceType, sender Sender) (*model.MsgCounterType, error) {
	cmd := model.CmdType{
		NodeManagementDetailedDiscoveryData: &model.NodeManagementDetailedDiscoveryDataType{},
	}

	rfAdress := featureAddressType(NodeManagementFeatureId, EntityAddressType(remoteDeviceAddress, DeviceInformationAddressEntity))

	return sender.Request(model.CmdClassifierTypeRead, rfAdress, false, []model.CmdType{cmd})
}

func (r *NodeManagementImpl) readDetailedDiscoveryData(featureRemote *FeatureRemoteImpl, requestHeader *model.HeaderType) error {

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

	return featureRemote.Sender().Reply(requestHeader, r.Address(), cmd)
}

func (r *NodeManagementImpl) replyDetailedDiscoveryData(message *Message, data *model.NodeManagementDetailedDiscoveryDataType) error {
	remoteDevice := message.featureRemote.Device()

	deviceDescription := data.DeviceInformation.Description
	if deviceDescription == nil {
		return errors.New("nodemanagement.replyDetailedDiscoveryData: invalid DeviceInformation.Description")
	}

	remoteDevice.UpdateDevice(deviceDescription)
	if err := remoteDevice.AddEntityAndFeatures(data); err != nil {
		return err
	}

	// payload := events.DeviceStateChangedPayload{
	// 	Device:      remoteDevice,
	// 	StateChange: model.NetworkManagementStateChangeTypeAdded,
	// }
	// events.DeviceStateChanged.Publish(payload)

	return nil
}

func (r *NodeManagementImpl) notifyDetailedDiscoveryData(message *Message, data *model.NodeManagementDetailedDiscoveryDataType) error {
	// is this a partial request?
	if !message.IsPartial {
		return errors.New("the received NodeManagementDetailedDiscovery.notify dataset should be partial")
	}

	if data.EntityInformation == nil || len(data.EntityInformation) == 0 || data.EntityInformation[0].Description == nil || data.EntityInformation[0].Description.LastStateChange == nil {
		return errors.New("the received NodeManagementDetailedDiscovery.notify dataset is incomplete")
	}

	lastStateChange := *data.EntityInformation[0].Description.LastStateChange
	remoteDevice := message.featureRemote.Device()

	// addition exmaple:
	// {"data":[{"header":[{"protocolId":"ee1.0"}]},{"payload":{"datagram":[{"header":[{"specificationVersion":"1.1.1"},{"addressSource":[{"device":"d:_i:19667_PorscheEVSE-00016544"},{"entity":[0]},{"feature":0}]},{"addressDestination":[{"device":"EVCC_HEMS"},{"entity":[0]},{"feature":0}]},{"msgCounter":6023},{"cmdClassifier":"notify"}]},{"payload":[{"cmd":[[{"function":"nodeManagementDetailedDiscoveryData"},{"filter":[[{"cmdControl":[{"partial":[]}]}]]},{"nodeManagementDetailedDiscoveryData":[{"deviceInformation":[{"description":[{"deviceAddress":[{"device":"d:_i:19667_PorscheEVSE-00016544"}]}]}]},{"entityInformation":[[{"description":[{"entityAddress":[{"entity":[1,1]}]},{"entityType":"EV"},{"lastStateChange":"added"},{"description":"Electric Vehicle"}]}]]},{"featureInformation":[[{"description":[{"featureAddress":[{"entity":[1,1]},{"feature":1}]},{"featureType":"LoadControl"},{"role":"server"},{"supportedFunction":[[{"function":"loadControlLimitDescriptionListData"},{"possibleOperations":[{"read":[]}]}],[{"function":"loadControlLimitListData"},{"possibleOperations":[{"read":[]},{"write":[]}]}]]},{"description":"Load Control"}]}],[{"description":[{"featureAddress":[{"entity":[1,1]},{"feature":2}]},{"featureType":"ElectricalConnection"},{"role":"server"},{"supportedFunction":[[{"function":"electricalConnectionParameterDescriptionListData"},{"possibleOperations":[{"read":[]}]}],[{"function":"electricalConnectionDescriptionListData"},{"possibleOperations":[{"read":[]}]}],[{"function":"electricalConnectionPermittedValueSetListData"},{"possibleOperations":[{"read":[]}]}]]},{"description":"Electrical Connection"}]}],[{"description":[{"featureAddress":[{"entity":[1,1]},{"feature":3}]},{"featureType":"Measurement"},{"specificUsage":["Electrical"]},{"role":"server"},{"supportedFunction":[[{"function":"measurementListData"},{"possibleOperations":[{"read":[]}]}],[{"function":"measurementDescriptionListData"},{"possibleOperations":[{"read":[]}]}]]},{"description":"Measurements"}]}],[{"description":[{"featureAddress":[{"entity":[1,1]},{"feature":5}]},{"featureType":"DeviceConfiguration"},{"role":"server"},{"supportedFunction":[[{"function":"deviceConfigurationKeyValueDescriptionListData"},{"possibleOperations":[{"read":[]}]}],[{"function":"deviceConfigurationKeyValueListData"},{"possibleOperations":[{"read":[]}]}]]},{"description":"Device Configuration EV"}]}],[{"description":[{"featureAddress":[{"entity":[1,1]},{"feature":6}]},{"featureType":"DeviceClassification"},{"role":"server"},{"supportedFunction":[[{"function":"deviceClassificationManufacturerData"},{"possibleOperations":[{"read":[]}]}]]},{"description":"Device Classification for EV"}]}],[{"description":[{"featureAddress":[{"entity":[1,1]},{"feature":7}]},{"featureType":"TimeSeries"},{"role":"server"},{"supportedFunction":[[{"function":"timeSeriesConstraintsListData"},{"possibleOperations":[{"read":[]}]}],[{"function":"timeSeriesDescriptionListData"},{"possibleOperations":[{"read":[]}]}],[{"function":"timeSeriesListData"},{"possibleOperations":[{"read":[]},{"write":[]}]}]]},{"description":"Time Series"}]}],[{"description":[{"featureAddress":[{"entity":[1,1]},{"feature":8}]},{"featureType":"IncentiveTable"},{"role":"server"},{"supportedFunction":[[{"function":"incentiveTableConstraintsData"},{"possibleOperations":[{"read":[]}]}],[{"function":"incentiveTableData"},{"possibleOperations":[{"read":[]},{"write":[]}]}],[{"function":"incentiveTableDescriptionData"},{"possibleOperations":[{"read":[]},{"write":[]}]}]]},{"description":"Incentive Table"}]}],[{"description":[{"featureAddress":[{"entity":[1,1]},{"feature":9}]},{"featureType":"DeviceDiagnosis"},{"role":"server"},{"supportedFunction":[[{"function":"deviceDiagnosisStateData"},{"possibleOperations":[{"read":[]}]}]]},{"description":"Device Diagnosis EV"}]}],[{"description":[{"featureAddress":[{"entity":[1,1]},{"feature":10}]},{"featureType":"Identification"},{"role":"server"},{"supportedFunction":[[{"function":"identificationListData"},{"possibleOperations":[{"read":[]}]}]]},{"description":"Identification for EV"}]}]]}]}]]}]}]}}]}
	// {"cmd":[[
	// 	{"function":"nodeManagementDetailedDiscoveryData"},
	// 	{"filter":[[{"cmdControl":[{"partial":[]}]}]]},
	// 	{"nodeManagementDetailedDiscoveryData":[
	// 		{"deviceInformation":[{"description":[{"deviceAddress":[{"device":"d:_i:19667_PorscheEVSE-00016544"}]}]}]},
	// 		{"entityInformation":[[
	// 			{"description":[
	// 				{"entityAddress":[{"entity":[1,1]}]},{"entityType":"EV"},
	// 				{"lastStateChange":"added"},
	// 				{"description":"Electric Vehicle"}]}
	// 		]]},
	// 		{"featureInformation":[
	// 			[{"description":[
	// 				{"featureAddress":[{"entity":[1,1]},{"feature":1}]},{"featureType":"LoadControl"},{"role":"server"},
	// 				{"supportedFunction":[
	// 					[{"function":"loadControlLimitDescriptionListData"},{"possibleOperations":[{"read":[]}]}],
	// 					[{"function":"loadControlLimitListData"},{"possibleOperations":[{"read":[]},{"write":[]}]}]
	// 				]},
	// 				{"description":"Load Control"}
	// 			]}],
	// 			[{"description":[
	// 				{"featureAddress":[{"entity":[1,1]},{"feature":2}]},{"featureType":"ElectricalConnection"},{"role":"server"},
	// 				{"supportedFunction":[
	// 					[{"function":"electricalConnectionParameterDescriptionListData"},{"possibleOperations":[{"read":[]}]}],
	// 					[{"function":"electricalConnectionDescriptionListData"},{"possibleOperations":[{"read":[]}]}],[{"function":"electricalConnectionPermittedValueSetListData"},{"possibleOperations":[{"read":[]}]}]
	// 				]},
	// 				{"description":"Electrical Connection"}
	// 			]}],
	// 			[{"description":[
	// 				{"featureAddress":[{"entity":[1,1]},{"feature":3}]},{"featureType":"Measurement"},{"specificUsage":["Electrical"]},{"role":"server"},
	// 				{"supportedFunction":[
	// 					[{"function":"measurementListData"},{"possibleOperations":[{"read":[]}]}],
	// 					[{"function":"measurementDescriptionListData"},{"possibleOperations":[{"read":[]}]}]
	// 				]},
	// 				{"description":"Measurements"}
	// 			]}],
	// 			[{"description":[
	// 				{"featureAddress":[{"entity":[1,1]},{"feature":5}]},{"featureType":"DeviceConfiguration"},{"role":"server"},
	// 				{"supportedFunction":[
	// 					[{"function":"deviceConfigurationKeyValueDescriptionListData"},{"possibleOperations":[{"read":[]}]}],
	// 					[{"function":"deviceConfigurationKeyValueListData"},{"possibleOperations":[{"read":[]}]}]
	// 				]},
	// 				{"description":"Device Configuration EV"}
	// 			]}],
	// 			[{"description":[
	// 				{"featureAddress":[{"entity":[1,1]},{"feature":6}]},{"featureType":"DeviceClassification"},{"role":"server"},
	// 				{"supportedFunction":[
	// 					[{"function":"deviceClassificationManufacturerData"},{"possibleOperations":[{"read":[]}]}]
	// 				]},
	// 				{"description":"Device Classification for EV"}]
	// 			}],
	// 			[{"description":[
	// 				{"featureAddress":[{"entity":[1,1]},{"feature":7}]},{"featureType":"TimeSeries"},{"role":"server"},
	// 				{"supportedFunction":[
	// 					[{"function":"timeSeriesConstraintsListData"},{"possibleOperations":[{"read":[]}]}],
	// 					[{"function":"timeSeriesDescriptionListData"},{"possibleOperations":[{"read":[]}]}],
	// 					[{"function":"timeSeriesListData"},{"possibleOperations":[{"read":[]},{"write":[]}]}]
	// 				]},
	// 				{"description":"Time Series"}]
	// 			}],
	// 			[{"description":[
	// 				{"featureAddress":[{"entity":[1,1]},{"feature":8}]},{"featureType":"IncentiveTable"},{"role":"server"},
	// 				{"supportedFunction":[
	// 					[{"function":"incentiveTableConstraintsData"},{"possibleOperations":[{"read":[]}]}],
	// 					[{"function":"incentiveTableData"},{"possibleOperations":[{"read":[]},{"write":[]}]}],
	// 					[{"function":"incentiveTableDescriptionData"},{"possibleOperations":[{"read":[]},{"write":[]}]}]
	// 				]},
	// 				{"description":"Incentive Table"}
	// 			]}],
	// 			[{"description":[
	// 				{"featureAddress":[{"entity":[1,1]},{"feature":9}]},{"featureType":"DeviceDiagnosis"},{"role":"server"},
	// 				{"supportedFunction":[
	// 					[{"function":"deviceDiagnosisStateData"},{"possibleOperations":[{"read":[]}]}]
	// 				]},
	// 				{"description":"Device Diagnosis EV"}
	// 			]}],
	// 			[{"description":[
	// 				{"featureAddress":[{"entity":[1,1]},{"feature":10}]},{"featureType":"Identification"},{"role":"server"},
	// 				{"supportedFunction":[
	// 					[{"function":"identificationListData"},{"possibleOperations":[{"read":[]}]}]
	// 				]},
	// 				{"description":"Identification for EV"}
	// 			]}]
	// 		]}
	// 	]}
	// ]]}

	// is this addition?
	if lastStateChange == model.NetworkManagementStateChangeTypeAdded {
		if err := remoteDevice.AddEntityAndFeatures(data); err != nil {
			return err
		}

		// payload := events.DeviceStateChangedPayload{
		// 	Device:      remoteDevice,
		// 	StateChange: lastStateChange,
		// }
		// events.DeviceStateChanged.Publish(payload)
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
			if err := remoteDevice.CheckEntityInformation(ei); err != nil {
				return err
			}

			entityAddress := ei.Description.EntityAddress.Entity
			remoteDevice.RemoveByAddress(entityAddress)
		}

		// payload := events.DeviceStateChangedPayload{
		// 	Device:      remoteDevice,
		// 	StateChange: lastStateChange,
		// }
		// events.DeviceStateChanged.Publish(payload)
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
		return r.readDetailedDiscoveryData(message.featureRemote, message.RequestHeader)

	case model.CmdClassifierTypeReply:
		return r.replyDetailedDiscoveryData(message, data)

	case model.CmdClassifierTypeNotify:
		return r.notifyDetailedDiscoveryData(message, data)

	default:
		return fmt.Errorf("nodemanagement.handleDetailedDiscoveryData: NodeManagementDetailedDiscoveryData CmdClassifierType not implemented: %s", message.CmdClassifier)
	}
}
