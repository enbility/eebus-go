package spine

import "github.com/enbility/eebus-go/spine/model"

type Message struct {
	RequestHeader *model.HeaderType
	CmdClassifier model.CmdClassifierType
	Cmd           model.CmdType
	FilterPartial *model.FilterType
	FilterDelete  *model.FilterType
	FeatureRemote *FeatureRemoteImpl
	EntityRemote  *EntityRemoteImpl
	DeviceRemote  *DeviceRemoteImpl
}

type ResultMessage struct {
	MsgCounterReference model.MsgCounterType  // required
	Result              *model.ResultDataType // required, may not be nil
	FeatureLocal        *FeatureLocalImpl     // required, may not be nil
	FeatureRemote       *FeatureRemoteImpl    // required, may not be nil
	EntityRemote        *EntityRemoteImpl     // required, may not be nil
	DeviceRemote        *DeviceRemoteImpl     // required, may not be nil
}
