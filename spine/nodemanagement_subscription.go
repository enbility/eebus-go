package spine

import (
	"fmt"

	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/DerAndereAndi/eebus-go/util"
	"github.com/ahmetb/go-linq/v3"
)

func NewNodeManagementSubscriptionRequestCallType(clientAddress *model.FeatureAddressType, serverAddress *model.FeatureAddressType, featureType model.FeatureTypeType) *model.NodeManagementSubscriptionRequestCallType {
	return &model.NodeManagementSubscriptionRequestCallType{
		SubscriptionRequest: &model.SubscriptionManagementRequestCallType{
			ClientAddress:     clientAddress,
			ServerAddress:     serverAddress,
			ServerFeatureType: &featureType,
		},
	}
}

// route subscription request calls to the appropriate feature implementation and add the subscription to the current list
func (r *NodeManagementImpl) readSubscriptionData(message *Message) error {

	var remoteDeviceSubscriptions []model.SubscriptionManagementEntryDataType
	remoteDeviceSubscriptionEntries := r.Device().SubscriptionManager().Subscriptions(message.FeatureRemote.Device())
	linq.From(remoteDeviceSubscriptionEntries).SelectT(func(s *SubscriptionEntry) model.SubscriptionManagementEntryDataType {
		return model.SubscriptionManagementEntryDataType{
			SubscriptionId: util.Ptr(model.SubscriptionIdType(s.id)),
			ServerAddress:  s.serverFeature.Address(),
			ClientAddress:  s.clientFeature.Address(),
		}
	}).ToSlice(&remoteDeviceSubscriptions)

	cmd := model.CmdType{
		NodeManagementSubscriptionData: &model.NodeManagementSubscriptionDataType{
			SubscriptionEntry: remoteDeviceSubscriptions,
		},
	}

	return message.FeatureRemote.Sender().Reply(message.RequestHeader, r.Address(), cmd)
}

func (r *NodeManagementImpl) handleMsgSubscriptionData(message *Message) error {
	switch message.CmdClassifier {
	case model.CmdClassifierTypeCall:
		return r.readSubscriptionData(message)

	default:
		return fmt.Errorf("nodemanagement.handleSubscriptionDeleteCall: NodeManagementSubscriptionRequestCall CmdClassifierType not implemented: %s", message.CmdClassifier)
	}
}

func (r *NodeManagementImpl) handleMsgSubscriptionRequestCall(message *Message, data *model.NodeManagementSubscriptionRequestCallType) error {
	switch message.CmdClassifier {
	case model.CmdClassifierTypeCall:
		return r.Device().SubscriptionManager().AddSubscription(r.Device(), message.FeatureRemote.Device(), *data.SubscriptionRequest)

	default:
		return fmt.Errorf("nodemanagement.handleSubscriptionRequestCall: NodeManagementSubscriptionRequestCall CmdClassifierType not implemented: %s", message.CmdClassifier)
	}
}

func (r *NodeManagementImpl) handleMsgSubscriptionDeleteCall(message *Message, data *model.NodeManagementSubscriptionDeleteCallType) error {
	switch message.CmdClassifier {
	case model.CmdClassifierTypeCall:
		return r.Device().SubscriptionManager().RemoveSubscription(*data.SubscriptionDelete, message.FeatureRemote.Device())

	default:
		return fmt.Errorf("nodemanagement.handleSubscriptionDeleteCall: NodeManagementSubscriptionRequestCall CmdClassifierType not implemented: %s", message.CmdClassifier)
	}
}
