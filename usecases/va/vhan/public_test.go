package vhan

import (
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func (s *VaVHANSuite) Test_Name() {
	// incompatible entity type
	name, err := s.sut.Name(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), "", name)

	// no user data available yet
	name, err = s.sut.Name(s.heatingAreaEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), "", name)

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.heatingAreaEntity, model.FeatureTypeTypeDeviceClassification, model.RoleTypeServer)

	// user data without a label
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeDeviceClassificationUserData, &model.DeviceClassificationUserDataType{}, nil, nil)
	assert.Nil(s.T(), fErr)

	name, err = s.sut.Name(s.heatingAreaEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), "", name)

	// user data with a label
	userData := &model.DeviceClassificationUserDataType{
		UserLabel: util.Ptr(model.LabelType("Living Room")),
	}
	_, fErr = rFeature.UpdateData(true, model.FunctionTypeDeviceClassificationUserData, userData, nil, nil)
	assert.Nil(s.T(), fErr)

	name, err = s.sut.Name(s.heatingAreaEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), "Living Room", name)
}
