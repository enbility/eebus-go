package spine

import (
	"errors"
	"fmt"

	"github.com/enbility/eebus-go/logging"
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

	cmd := model.CmdType{
		NodeManagementUseCaseData: &model.NodeManagementUseCaseDataType{
			UseCaseInformation: r.entity.Device().UseCaseManager().UseCaseInformation(),
		},
	}

	return featureRemote.Sender().Notify(r.Address(), rfAdress, cmd)
}

func (r *NodeManagementImpl) processReadUseCaseData(featureRemote *FeatureRemoteImpl, requestHeader *model.HeaderType) error {

	cmd := model.CmdType{
		NodeManagementUseCaseData: &model.NodeManagementUseCaseDataType{
			UseCaseInformation: r.entity.Device().UseCaseManager().UseCaseInformation(),
		},
	}

	return featureRemote.Sender().Reply(requestHeader, r.Address(), cmd)
}

func (r *NodeManagementImpl) processReplyUseCaseData(message *Message, data model.NodeManagementUseCaseDataType) error {
	useCaseInformation := data.UseCaseInformation
	if useCaseInformation == nil {
		return errors.New("nodemanagement.replyUseCaseData: invalid UseCaseInformation")
	}

	remoteUseCaseManager := message.FeatureRemote.Device().UseCaseManager()
	remoteUseCaseManager.RemoveAll()

	for _, useCaseInfo := range useCaseInformation {
		// this is mandatory
		var actor model.UseCaseActorType
		if useCaseInfo.Actor != nil {
			actor = model.UseCaseActorType(*useCaseInfo.Actor)
		} else {
			logging.Log().Debug("actor is missing in useCaseInformation")
			break
		}

		for _, useCaseSupport := range useCaseInfo.UseCaseSupport {

			// this is mandatory
			var useCaseName model.UseCaseNameType
			if useCaseSupport.UseCaseName != nil {
				useCaseName = model.UseCaseNameType(*useCaseSupport.UseCaseName)
			} else {
				logging.Log().Debug("useCaseName is missing in useCaseSupport")
				continue
			}

			// this is optional
			useCaseAvailable := true
			if useCaseSupport.UseCaseAvailable != nil {
				useCaseAvailable = *useCaseSupport.UseCaseAvailable
			}

			var useCaseVersion model.SpecificationVersionType
			if useCaseSupport.UseCaseVersion != nil {
				useCaseVersion = model.SpecificationVersionType(*useCaseSupport.UseCaseVersion)
			}

			var useCaseDocumemtSubRevision string
			if useCaseSupport.UseCaseDocumentSubRevision != nil {
				useCaseDocumemtSubRevision = *useCaseSupport.UseCaseDocumentSubRevision
			}

			if useCaseSupport.ScenarioSupport == nil {
				logging.Log().Errorf("scenarioSupport is missing in useCaseSupport %s", useCaseName)
				continue
			}

			remoteUseCaseManager.Add(
				actor,
				useCaseName,
				useCaseVersion,
				useCaseDocumemtSubRevision,
				useCaseAvailable,
				useCaseSupport.ScenarioSupport)
		}
	}

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
		return r.processReplyUseCaseData(message, *data)

	case model.CmdClassifierTypeNotify:
		return r.processReplyUseCaseData(message, *data)

	default:
		return fmt.Errorf("nodemanagement.handleUseCaseData: NodeManagementUseCaseData CmdClassifierType not implemented: %s", message.CmdClassifier)
	}
}
