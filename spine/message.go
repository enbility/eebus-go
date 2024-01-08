package spine

import "github.com/enbility/eebus-go/spine/model"

type Message struct {
	RequestHeader *model.HeaderType
	CmdClassifier model.CmdClassifierType
	Cmd           model.CmdType
	FilterPartial *model.FilterType
	FilterDelete  *model.FilterType
	FeatureRemote FeatureRemote
	EntityRemote  EntityRemote
	DeviceRemote  DeviceRemote
}

type ResultMessage struct {
	MsgCounterReference model.MsgCounterType  // required
	Result              *model.ResultDataType // required, may not be nil
	FeatureLocal        FeatureLocal          // required, may not be nil
	FeatureRemote       FeatureRemote         // required, may not be nil
	EntityRemote        EntityRemote          // required, may not be nil
	DeviceRemote        DeviceRemote          // required, may not be nil
}
