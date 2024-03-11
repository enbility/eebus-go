package features

import (
	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// Feature interface were the local feature role is client and the remote feature role is server
type FeatureInterface interface {
	// check if there is a subscription to the remote feature
	HasSubscription() bool

	// subscribe to the feature of the entity
	Subscribe() (*model.MsgCounterType, error)

	// check if there is a binding to the remote feature
	HasBinding() bool

	// bind to the feature of the entity
	Bind() (*model.MsgCounterType, error)

	// add a callback function to be invoked once a result or reply message for a msgCounter came in
	AddResponseCallback(msgCounterReference model.MsgCounterType, function func(msg api.ResponseMessage)) error
}
