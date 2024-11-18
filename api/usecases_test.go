package api

import (
	"encoding/json"
	"testing"

	"github.com/enbility/spine-go/mocks"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"

	"github.com/stretchr/testify/assert"
)

func Test_RemoteEntityScenarios_Marshal(t *testing.T) {
	item := RemoteEntityScenarios{
		Entity:    nil,
		Scenarios: nil,
	}
	value, err := json.Marshal(item)
	assert.Nil(t, err)
	assert.NotNil(t, value)
	assert.Equal(t, `{"Entity":null,"Scenarios":null}`, string(value))

	item = RemoteEntityScenarios{
		Entity:    nil,
		Scenarios: []uint{1, 2, 3},
	}
	value, err = json.Marshal(item)
	assert.Nil(t, err)
	assert.NotNil(t, value)
	assert.Equal(t, `{"Entity":null,"Scenarios":[1,2,3]}`, string(value))

	device := mocks.NewDeviceRemoteInterface(t)
	deviceAddress := model.AddressDeviceType("test")
	device.EXPECT().Address().Return(&deviceAddress).Times(1)
	entity := spine.NewEntityRemote(device, model.EntityTypeTypeCEM, []model.AddressEntityType{1, 1})

	item = RemoteEntityScenarios{
		Entity:    entity,
		Scenarios: []uint{1, 2, 3},
	}
	value, err = json.Marshal(item)
	assert.Nil(t, err)
	assert.NotNil(t, value)
	assert.Equal(t, `{"Entity":{"Device":"test","Entity":[1,1]},"Scenarios":[1,2,3]}`, string(value))
}
