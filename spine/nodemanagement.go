package spine

import (
	"fmt"

	"github.com/DerAndereAndi/eebus-go/spine/model"
)

const NodeManagementFeatureId uint = 0

type NodeManagementImpl struct {
	*FeatureLocalImpl
	entity *EntityLocalImpl
}

func NewNodeManagementImpl(id uint, entity *EntityLocalImpl) *NodeManagementImpl {
	f := &NodeManagementImpl{
		FeatureLocalImpl: NewFeatureLocalImpl(
			id, entity,
			model.FeatureTypeTypeNodeManagement,
			model.RoleTypeSpecial),
		entity: entity,
	}

	return f
}

func (r *NodeManagementImpl) Device() *DeviceLocalImpl {
	return r.entity.Device()
}

func (r *NodeManagementImpl) HandleMessage(message *Message) error {
	if message.Cmd.ResultData != nil {
		return r.processResult(message.CmdClassifier)
	}

	switch {
	case message.Cmd.NodeManagementDetailedDiscoveryData != nil:
		return r.handleMsgDetailedDiscoveryData(message, message.Cmd.NodeManagementDetailedDiscoveryData)

	case message.Cmd.NodeManagementSubscriptionRequestCall != nil:
		return r.handleMsgSubscriptionRequestCall(message, message.Cmd.NodeManagementSubscriptionRequestCall)

	case message.Cmd.NodeManagementSubscriptionDeleteCall != nil:
		return r.handleMsgSubscriptionDeleteCall(message, message.Cmd.NodeManagementSubscriptionDeleteCall)

	case message.Cmd.NodeManagementSubscriptionData != nil:
		return r.handleMsgSubscriptionData(message)

	case message.Cmd.NodeManagementUseCaseData != nil:
		return r.handleMsgUseCaseData(message, message.Cmd.NodeManagementUseCaseData)

	default:
		return fmt.Errorf("nodemanagement.Handle: Cmd data not implemented: %s", message.Cmd.DataName())
	}
}

func (r *NodeManagementImpl) processResult(cmdClassifier model.CmdClassifierType) error {
	switch cmdClassifier {
	case model.CmdClassifierTypeResult:
		// TODO process the return result data for the message sent with the ID in msgCounterReference
		// error numbers explained in Resource Spec 3.11
		return nil

	default:
		return fmt.Errorf("ResultData CmdClassifierType %s not implemented", cmdClassifier)
	}
}
