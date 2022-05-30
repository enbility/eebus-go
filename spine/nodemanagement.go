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

func (r *NodeManagementImpl) HandleMessage(message *Message) *ErrorType {
	switch {
	case message.Cmd.ResultData != nil:
		if err := r.processResult(message); err != nil {
			_ = r.pendingRequests.Remove(*message.RequestHeader.MsgCounterReference)
			return err
		}

	case message.Cmd.NodeManagementDetailedDiscoveryData != nil:
		if err := r.handleMsgDetailedDiscoveryData(message, message.Cmd.NodeManagementDetailedDiscoveryData); err != nil {
			return NewErrorType(model.ErrorNumberTypeGeneralError, err.Error())
		}

	case message.Cmd.NodeManagementSubscriptionRequestCall != nil:
		if err := r.handleMsgSubscriptionRequestCall(message, message.Cmd.NodeManagementSubscriptionRequestCall); err != nil {
			return NewErrorType(model.ErrorNumberTypeGeneralError, err.Error())
		}

	case message.Cmd.NodeManagementSubscriptionDeleteCall != nil:
		if err := r.handleMsgSubscriptionDeleteCall(message, message.Cmd.NodeManagementSubscriptionDeleteCall); err != nil {
			return NewErrorType(model.ErrorNumberTypeGeneralError, err.Error())
		}

	case message.Cmd.NodeManagementSubscriptionData != nil:
		if err := r.handleMsgSubscriptionData(message); err != nil {
			return NewErrorType(model.ErrorNumberTypeGeneralError, err.Error())
		}

	case message.Cmd.NodeManagementUseCaseData != nil:
		if err := r.handleMsgUseCaseData(message, message.Cmd.NodeManagementUseCaseData); err != nil {
			return NewErrorType(model.ErrorNumberTypeGeneralError, err.Error())
		}

	default:
		return NewErrorType(model.ErrorNumberTypeCommandNotSupported, fmt.Sprintf("nodemanagement.Handle: Cmd data not implemented: %s", message.Cmd.DataName()))
	}

	return nil
}
