package spine

import (
	"fmt"

	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/DerAndereAndi/eebus-go/util"
)

const NodeManagementFeatureId uint = 0

func NodeManagementAddress(deviceAdress *model.AddressDeviceType) *model.FeatureAddressType {
	return &model.FeatureAddressType{
		Entity:  []model.AddressEntityType{0},
		Feature: util.Ptr(model.AddressFeatureType(NodeManagementFeatureId)),
		Device:  deviceAdress,
	}
}

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

	f.AddFunctionType(model.FunctionTypeNodeManagementDetailedDiscoveryData, true, false)
	f.AddFunctionType(model.FunctionTypeNodeManagementUseCaseData, true, false)
	f.AddFunctionType(model.FunctionTypeNodeManagementSubscriptionData, true, false)
	f.AddFunctionType(model.FunctionTypeNodeManagementSubscriptionRequestCall, false, false)
	f.AddFunctionType(model.FunctionTypeNodeManagementSubscriptionDeleteCall, false, false)
	f.AddFunctionType(model.FunctionTypeNodeManagementBindingData, true, false)
	f.AddFunctionType(model.FunctionTypeNodeManagementBindingRequestCall, false, false)
	f.AddFunctionType(model.FunctionTypeNodeManagementBindingDeleteCall, false, false)
	if f.Device().FeatureSet() != nil && *f.Device().FeatureSet() != model.NetworkManagementFeatureSetTypeSimple {
		f.AddFunctionType(model.FunctionTypeNodeManagementDestinationListData, true, false)
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

	case message.Cmd.NodeManagementBindingRequestCall != nil:
		if err := r.handleMsgBindingRequestCall(message, message.Cmd.NodeManagementBindingRequestCall); err != nil {
			return NewErrorType(model.ErrorNumberTypeGeneralError, err.Error())
		}

	case message.Cmd.NodeManagementBindingDeleteCall != nil:
		if err := r.handleMsgBindingDeleteCall(message, message.Cmd.NodeManagementBindingDeleteCall); err != nil {
			return NewErrorType(model.ErrorNumberTypeGeneralError, err.Error())
		}

	case message.Cmd.NodeManagementBindingData != nil:
		if err := r.handleMsgBindingData(message); err != nil {
			return NewErrorType(model.ErrorNumberTypeGeneralError, err.Error())
		}

	case message.Cmd.NodeManagementUseCaseData != nil:
		if err := r.handleMsgUseCaseData(message, message.Cmd.NodeManagementUseCaseData); err != nil {
			return NewErrorType(model.ErrorNumberTypeGeneralError, err.Error())
		}

	case message.Cmd.NodeManagementDestinationListData != nil:
		if err := r.handleMsgDestinationListData(message, message.Cmd.NodeManagementDestinationListData); err != nil {
			return NewErrorType(model.ErrorNumberTypeGeneralError, err.Error())
		}

	default:
		return NewErrorType(model.ErrorNumberTypeCommandNotSupported, fmt.Sprintf("nodemanagement.Handle: Cmd data not implemented: %s", message.Cmd.DataName()))
	}

	return nil
}
