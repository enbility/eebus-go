package spine

import (
	"errors"
	"fmt"

	"github.com/DerAndereAndi/eebus-go/spine/model"
)

func (r *NodeManagementImpl) RequestDestinationListData(remoteDeviceAddress *model.AddressDeviceType, sender Sender) (*model.MsgCounterType, *ErrorType) {
	return nil, NewErrorTypeFromString("Not implemented")
}

func (r *NodeManagementImpl) readDestinationListData(featureRemote *FeatureRemoteImpl, requestHeader *model.HeaderType) error {
	data := []model.NodeManagementDestinationDataType{
		r.Device().DestinationData(),
	}
	// add other remote devices here

	cmd := model.CmdType{
		NodeManagementDestinationListData: &model.NodeManagementDestinationListDataType{
			NodeManagementDestinationData: data,
		},
	}

	return featureRemote.Sender().Reply(requestHeader, r.Address(), cmd)
}

func (r *NodeManagementImpl) replyDestinationListData(message *Message, data model.NodeManagementDestinationListDataType) error {
	return errors.New("Not implemented")
}

func (r *NodeManagementImpl) handleMsgDestinationListData(message *Message, data *model.NodeManagementDestinationListDataType) error {
	switch message.CmdClassifier {
	case model.CmdClassifierTypeRead:
		return r.readDestinationListData(message.FeatureRemote, message.RequestHeader)

	case model.CmdClassifierTypeReply:
		if err := r.pendingRequests.Remove(*message.RequestHeader.MsgCounterReference); err != nil {
			return errors.New(err.String())
		}
		return r.replyDestinationListData(message, *data)

	case model.CmdClassifierTypeNotify:
		return r.replyDestinationListData(message, *data)

	default:
		return fmt.Errorf("nodemanagement.handleMsgDestinationListData: NodeManagementDestinationListDataType CmdClassifierType not implemented: %s", message.CmdClassifier)
	}
}
