package evsecc

import (
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func (s *EVSECCSuite) Test_Events() {
	payload := spineapi.EventPayload{
		Entity: s.mockRemoteEntity,
	}
	s.sut.HandleEvent(payload)

	payload.Entity = s.evseEntity
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeDeviceChange
	payload.ChangeType = spineapi.ElementChangeRemove
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeEntityChange
	payload.ChangeType = spineapi.ElementChangeAdd
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeEntityChange
	payload.ChangeType = spineapi.ElementChangeRemove
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeDataChange
	payload.ChangeType = spineapi.ElementChangeAdd
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeDataChange
	payload.ChangeType = spineapi.ElementChangeUpdate
	payload.Data = util.Ptr(model.DeviceClassificationManufacturerDataType{})
	s.sut.HandleEvent(payload)

	payload.Data = util.Ptr(model.DeviceDiagnosisStateDataType{})
	s.sut.HandleEvent(payload)
}

func (s *EVSECCSuite) Test_evseManufacturerDataUpdate() {
	payload := spineapi.EventPayload{
		Ski:    remoteSki,
		Device: s.remoteDevice,
		Entity: s.mockRemoteEntity,
	}
	s.sut.evseManufacturerDataUpdate(payload)

	payload.Entity = s.evseEntity
	s.sut.evseManufacturerDataUpdate(payload)

	data := &model.DeviceClassificationManufacturerDataType{
		BrandName: util.Ptr(model.DeviceClassificationStringType("test")),
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evseEntity, model.FeatureTypeTypeDeviceClassification, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeDeviceClassificationManufacturerData, data, nil, nil)
	assert.Nil(s.T(), fErr)

	s.sut.evseManufacturerDataUpdate(payload)
}

func (s *EVSECCSuite) Test_evseStateUpdate() {
	payload := spineapi.EventPayload{
		Ski:    remoteSki,
		Device: s.remoteDevice,
		Entity: s.mockRemoteEntity,
	}
	s.sut.evseStateUpdate(payload)

	payload.Entity = s.evseEntity
	s.sut.evseStateUpdate(payload)

	data := &model.DeviceDiagnosisStateDataType{
		OperatingState: util.Ptr(model.DeviceDiagnosisOperatingStateTypeNormalOperation),
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evseEntity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeDeviceDiagnosisStateData, data, nil, nil)
	assert.Nil(s.T(), fErr)

	s.sut.evseStateUpdate(payload)
}
