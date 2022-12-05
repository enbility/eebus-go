package spine

import (
	"errors"
	"fmt"

	"github.com/enbility/eebus-go/logging"
	"github.com/enbility/eebus-go/spine/model"
)

func (r *NodeManagementImpl) RequestUseCaseData(remoteDeviceSki string, remoteDeviceAddress *model.AddressDeviceType, sender Sender) (*model.MsgCounterType, *ErrorType) {
	rfAdress := featureAddressType(NodeManagementFeatureId, EntityAddressType(remoteDeviceAddress, DeviceInformationAddressEntity))
	cmd := model.CmdType{
		NodeManagementUseCaseData: &model.NodeManagementUseCaseDataType{},
	}
	return r.RequestDataBySenderAddress(cmd, sender, remoteDeviceSki, rfAdress, defaultMaxResponseDelay)
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
	for _, useCaseInfo := range useCaseInformation {
		// this is mandatory
		var actor model.UseCaseActorType
		if useCaseInfo.Actor != nil {
			actor = model.UseCaseActorType(*useCaseInfo.Actor)
		} else {
			logging.Log.Error("actor is missing in useCaseInformation")
			break
		}

		for _, useCaseSupport := range useCaseInfo.UseCaseSupport {

			// this is mandatory
			var useCaseName model.UseCaseNameType
			if useCaseSupport.UseCaseName != nil {
				useCaseName = model.UseCaseNameType(*useCaseSupport.UseCaseName)
			} else {
				logging.Log.Error("useCaseName is missing in useCaseSupport")
				continue
			}

			// this is optional
			var useCaseVersion model.SpecificationVersionType
			if useCaseSupport.UseCaseVersion != nil {
				useCaseVersion = model.SpecificationVersionType(*useCaseSupport.UseCaseVersion)
			}

			if useCaseSupport.ScenarioSupport == nil {
				logging.Log.Errorf("scenarioSupport is missing in useCaseSupport %s", useCaseName)
				continue
			}

			remoteUseCaseManager.Add(
				actor,
				useCaseName,
				useCaseVersion,
				useCaseSupport.ScenarioSupport)
		}
	}

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
