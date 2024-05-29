package internal

import (
	"testing"

	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func Test_IsHeartbeat(t *testing.T) {
	payload := spineapi.EventPayload{}
	result := IsHeartbeat(payload)
	assert.False(t, result)

	payload.Data = &model.DeviceDiagnosisHeartbeatDataType{}
	result = IsHeartbeat(payload)
	assert.False(t, result)

	payload.Function = model.FunctionTypeDeviceDiagnosisHeartbeatData
	result = IsHeartbeat(payload)
	assert.False(t, result)

	payload.EventType = spineapi.EventTypeDataChange
	payload.ChangeType = spineapi.ElementChangeUpdate
	payload.CmdClassifier = util.Ptr(model.CmdClassifierTypeNotify)
	result = IsHeartbeat(payload)
	assert.True(t, result)
}
