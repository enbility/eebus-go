package features

import (
	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
)

type Feature interface {
	SubscribeForEntity() error
	AddResultCallback(msgCounterReference model.MsgCounterType, function func(msg spine.ResultMessage))
}
