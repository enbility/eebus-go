package spine

import "github.com/DerAndereAndi/eebus-go/spine/model"

type Message struct {
	RequestHeader *model.HeaderType
	CmdClassifier model.CmdClassifierType
	Cmd           model.CmdType
	IsPartial     bool
	featureRemote *FeatureRemoteImpl
	deviceRemote  *DeviceRemoteImpl
}
