package spine

import (
	"fmt"

	"github.com/ahmetb/go-linq/v3"
	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
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

func NewNodeManagementSubscriptionDeleteCallType(clientAddress *model.FeatureAddressType, serverAddress *model.FeatureAddressType, featureType model.FeatureTypeType) *model.NodeManagementSubscriptionDeleteCallType {
	return &model.NodeManagementSubscriptionDeleteCallType{
		SubscriptionDelete: &model.SubscriptionManagementDeleteCallType{
			ClientAddress: clientAddress,
			ServerAddress: serverAddress,
		},
	}
}

// route subscription request calls to the appropriate feature implementation and add the subscription to the current list
func (r *NodeManagementImpl) processReadSubscriptionData(message *Message) error {

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
		return r.processReadSubscriptionData(message)

	default:
		return fmt.Errorf("nodemanagement.handleSubscriptionDeleteCall: NodeManagementSubscriptionRequestCall CmdClassifierType not implemented: %s", message.CmdClassifier)
	}
}

func (r *NodeManagementImpl) handleMsgSubscriptionRequestCall(message *Message, data *model.NodeManagementSubscriptionRequestCallType) error {
	switch message.CmdClassifier {
	case model.CmdClassifierTypeCall:
		subscriptionMgr := r.Device().SubscriptionManager()

		err := subscriptionMgr.AddSubscription(message.FeatureRemote.Device(), *data.SubscriptionRequest)
		if err == nil {
			r.Device().HeartbeatManager().UpdateHeartbeatOnSubscriptions()
		}

		return err

	default:
		return fmt.Errorf("nodemanagement.handleSubscriptionRequestCall: NodeManagementSubscriptionRequestCall CmdClassifierType not implemented: %s", message.CmdClassifier)
	}
}

func (r *NodeManagementImpl) handleMsgSubscriptionDeleteCall(message *Message, data *model.NodeManagementSubscriptionDeleteCallType) error {
	switch message.CmdClassifier {
	case model.CmdClassifierTypeCall:
		subscriptionMgr := r.Device().SubscriptionManager()

		err := subscriptionMgr.RemoveSubscription(*data.SubscriptionDelete, message.FeatureRemote.Device())
		if err == nil {
			r.Device().HeartbeatManager().UpdateHeartbeatOnSubscriptions()
		}

		return err
	default:
		return fmt.Errorf("nodemanagement.handleSubscriptionDeleteCall: NodeManagementSubscriptionRequestCall CmdClassifierType not implemented: %s", message.CmdClassifier)
	}
}
