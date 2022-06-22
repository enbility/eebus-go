package spine

import "github.com/DerAndereAndi/eebus-go/spine/model"

type Message struct {
	RequestHeader *model.HeaderType
	CmdClassifier model.CmdClassifierType
	Cmd           model.CmdType
	FilterPartial *model.FilterType
	FilterDelete  *model.FilterType
	FeatureRemote *FeatureRemoteImpl
	DeviceRemote  *DeviceRemoteImpl
}
