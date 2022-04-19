//go:generate mockery --name=Sender
package spine

import "github.com/DerAndereAndi/eebus-go/spine/model"

type Sender interface {
	Request(cmdClassifier model.CmdClassifierType, senderAddress, destinationAddress *model.FeatureAddressType, ackRequest bool, cmd []model.CmdType) (*model.MsgCounterType, error)
	Reply(requestHeader *model.HeaderType, senderAddress *model.FeatureAddressType, cmd model.CmdType) error
	Subscribe(senderAddress, destinationAddress *model.FeatureAddressType, serverFeatureType model.FeatureTypeType) error
	Notify(senderAddress, destinationAddress *model.FeatureAddressType, cmd []model.CmdType) error
	SendAcknowledgementMessage(err error, featureSource *model.FeatureAddressType, featureDestination *model.FeatureAddressType, msgCounterReference *model.MsgCounterType) error
}
