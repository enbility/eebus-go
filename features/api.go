package features

import (
	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type Feature interface {
	Subscribe() error
	Bind() error
	AddResultCallback(msgCounterReference model.MsgCounterType, function func(msg api.ResultMessage))
}
