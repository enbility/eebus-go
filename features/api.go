package features

import (
	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type FeatureInterface interface {
	// check if there is a subscription to the remote feature
	HasSubscription() bool

	// subscribe to the feature of the entity
	Subscribe() (*model.MsgCounterType, error)

	// check if there is a binding to the remote feature
	HasBinding() bool

	// bind to the feature of the entity
	Bind() (*model.MsgCounterType, error)

	// add a callback function to be invoked once a result to a msgCounter came in
	AddResultCallback(msgCounterReference model.MsgCounterType, function func(msg api.ResultMessage))
}
