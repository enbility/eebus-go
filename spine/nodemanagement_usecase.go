package spine

import (
	"errors"
	"fmt"

	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
)

func (r *NodeManagementImpl) RequestUseCaseData(remoteDeviceSki string, remoteDeviceAddress *model.AddressDeviceType, sender Sender) (*model.MsgCounterType, *model.ErrorType) {
	rfAdress := featureAddressType(NodeManagementFeatureId, EntityAddressType(remoteDeviceAddress, DeviceInformationAddressEntity))
	cmd := model.CmdType{
		NodeManagementUseCaseData: &model.NodeManagementUseCaseDataType{},
	}
	return r.RequestDataBySenderAddress(cmd, sender, remoteDeviceSki, rfAdress, defaultMaxResponseDelay)
}

func (r *NodeManagementImpl) NotifyUseCaseData(remoteDevice *DeviceRemoteImpl) (*model.MsgCounterType, error) {
	rfAdress := featureAddressType(NodeManagementFeatureId, EntityAddressType(remoteDevice.address, DeviceInformationAddressEntity))
	rEntity := remoteDevice.Entity([]model.AddressEntityType{model.AddressEntityType(DeviceInformationEntityId)})

	featureRemote := remoteDevice.FeatureByEntityTypeAndRole(rEntity, model.FeatureTypeTypeNodeManagement, model.RoleTypeSpecial)

	fd := r.functionData(model.FunctionTypeNodeManagementUseCaseData)
	cmd := fd.NotifyCmdType(nil, nil, false, nil)

	return featureRemote.Sender().Notify(r.Address(), rfAdress, cmd)
}

func (r *NodeManagementImpl) processReadUseCaseData(featureRemote FeatureRemote, requestHeader *model.HeaderType) error {
	cmd := r.functionData(model.FunctionTypeNodeManagementUseCaseData).ReplyCmdType(false)

	return featureRemote.Sender().Reply(requestHeader, r.Address(), cmd)
}

func (r *NodeManagementImpl) processReplyUseCaseData(message *Message, data *model.NodeManagementUseCaseDataType) error {
	message.FeatureRemote.UpdateData(model.FunctionTypeNodeManagementUseCaseData, data, nil, nil)

	// the data was updated, so send an event, other event handlers may watch out for this as well
	payload := EventPayload{
		Ski:           message.FeatureRemote.Device().ski,
		EventType:     EventTypeDataChange,
		ChangeType:    ElementChangeUpdate,
		Feature:       message.FeatureRemote,
		Device:        message.FeatureRemote.Device(),
		Entity:        message.FeatureRemote.Entity(),
		CmdClassifier: util.Ptr(message.CmdClassifier),
		Data:          data,
	}
	Events.Publish(payload)

	return nil
}

func (r *NodeManagementImpl) handleMsgUseCaseData(message *Message, data *model.NodeManagementUseCaseDataType) error {
	switch message.CmdClassifier {
	case model.CmdClassifierTypeRead:
		return r.processReadUseCaseData(message.FeatureRemote, message.RequestHeader)

	case model.CmdClassifierTypeReply:
		if err := r.pendingRequests.Remove(message.DeviceRemote.ski, *message.RequestHeader.MsgCounterReference); err != nil {
			return errors.New(err.String())
		}
		return r.processReplyUseCaseData(message, data)

	case model.CmdClassifierTypeNotify:
		return r.processReplyUseCaseData(message, data)

	default:
		return fmt.Errorf("nodemanagement.handleUseCaseData: NodeManagementUseCaseData CmdClassifierType not implemented: %s", message.CmdClassifier)
	}
}
