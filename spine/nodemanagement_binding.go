package spine

import (
	"fmt"

	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/DerAndereAndi/eebus-go/util"
	"github.com/ahmetb/go-linq/v3"
)

func NewNodeManagementBindingRequestCallType(clientAddress *model.FeatureAddressType, serverAddress *model.FeatureAddressType, featureType model.FeatureTypeType) *model.NodeManagementBindingRequestCallType {
	return &model.NodeManagementBindingRequestCallType{
		BindingRequest: &model.BindingManagementRequestCallType{
			ClientAddress:     clientAddress,
			ServerAddress:     serverAddress,
			ServerFeatureType: &featureType,
		},
	}
}

// route bindings request calls to the appropriate feature implementation and add the bindings to the current list
func (r *NodeManagementImpl) processReadBindingData(message *Message) error {

	var remoteDeviceBindings []model.BindingManagementEntryDataType
	remoteDeviceBindingEntries := r.Device().BindingManager().Bindings(message.FeatureRemote.Device())
	linq.From(remoteDeviceBindingEntries).SelectT(func(s *BindingEntry) model.BindingManagementEntryDataType {
		return model.BindingManagementEntryDataType{
			BindingId:     util.Ptr(model.BindingIdType(s.id)),
			ServerAddress: s.serverFeature.Address(),
			ClientAddress: s.clientFeature.Address(),
		}
	}).ToSlice(&remoteDeviceBindings)

	cmd := model.CmdType{
		NodeManagementBindingData: &model.NodeManagementBindingDataType{
			BindingEntry: remoteDeviceBindings,
		},
	}

	return message.FeatureRemote.Sender().Reply(message.RequestHeader, r.Address(), cmd)
}

func (r *NodeManagementImpl) handleMsgBindingData(message *Message) error {
	switch message.CmdClassifier {
	case model.CmdClassifierTypeCall:
		return r.processReadBindingData(message)

	default:
		return fmt.Errorf("nodemanagement.handleBindingDeleteCall: NodeManagementBindingRequestCall CmdClassifierType not implemented: %s", message.CmdClassifier)
	}
}

func (r *NodeManagementImpl) handleMsgBindingRequestCall(message *Message, data *model.NodeManagementBindingRequestCallType) error {
	switch message.CmdClassifier {
	case model.CmdClassifierTypeCall:
		return r.Device().BindingManager().AddBinding(r.Device(), message.FeatureRemote.Device(), *data.BindingRequest)

	default:
		return fmt.Errorf("nodemanagement.handleBindingRequestCall: NodeManagementBindingRequestCall CmdClassifierType not implemented: %s", message.CmdClassifier)
	}
}

func (r *NodeManagementImpl) handleMsgBindingDeleteCall(message *Message, data *model.NodeManagementBindingDeleteCallType) error {
	switch message.CmdClassifier {
	case model.CmdClassifierTypeCall:
		return r.Device().BindingManager().RemoveBinding(*data.BindingDelete, message.FeatureRemote.Device())

	default:
		return fmt.Errorf("nodemanagement.handleBindingDeleteCall: NodeManagementBindingRequestCall CmdClassifierType not implemented: %s", message.CmdClassifier)
	}
}
