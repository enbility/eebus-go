package spine

import (
	"errors"
	"fmt"

	"github.com/DerAndereAndi/eebus-go/spine/model"
)

func (r *NodeManagementImpl) readUseCaseData(featureRemote *FeatureRemoteImpl, requestHeader *model.HeaderType) error {

	cmd := model.CmdType{
		NodeManagementUseCaseData: &model.NodeManagementUseCaseDataType{
			UseCaseInformation: r.entity.Device().UseCaseManager().UseCaseInformation(),
		},
	}

	return featureRemote.Sender().Reply(requestHeader, r.Address(), cmd)
}

func (r *NodeManagementImpl) replyUseCaseData(message *Message, data model.NodeManagementUseCaseDataType) error {
	useCaseInformation := data.UseCaseInformation
	if useCaseInformation == nil {
		return errors.New("nodemanagement.replyUseCaseData: invalid UseCaseInformation")
	}

	remoteUseCaseManager := message.featureRemote.Device().UseCaseManager()
	for _, useCaseInfo := range useCaseInformation {
		for _, useCaseSupport := range useCaseInfo.UseCaseSupport {
			remoteUseCaseManager.Add(
				model.UseCaseActorType(*useCaseInfo.Actor),
				model.UseCaseNameType(*useCaseSupport.UseCaseName),
				useCaseSupport.ScenarioSupport)
		}
	}

	return nil
}

func (r *NodeManagementImpl) handleMsgUseCaseData(message *Message, data *model.NodeManagementUseCaseDataType) error {
	switch message.CmdClassifier {
	case model.CmdClassifierTypeRead:
		return r.readUseCaseData(message.featureRemote, message.RequestHeader)

	case model.CmdClassifierTypeReply:
		return r.replyUseCaseData(message, *data)

	case model.CmdClassifierTypeNotify:
		return r.replyUseCaseData(message, *data)

	default:
		return fmt.Errorf("nodemanagement.handleUseCaseData: NodeManagementUseCaseData CmdClassifierType not implemented: %s", message.CmdClassifier)
	}
}
