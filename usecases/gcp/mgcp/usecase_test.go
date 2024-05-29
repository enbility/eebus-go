package mgcp

import (
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func (s *GcpMGCPSuite) Test_UpdateUseCaseAvailability() {
	s.sut.UpdateUseCaseAvailability(true)
}

func (s *GcpMGCPSuite) Test_IsUseCaseSupported() {
	data, err := s.sut.IsUseCaseSupported(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), false, data)

	data, err = s.sut.IsUseCaseSupported(s.smgwEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), false, data)

	ucData := &model.NodeManagementUseCaseDataType{
		UseCaseInformation: []model.UseCaseInformationDataType{
			{
				Actor: util.Ptr(model.UseCaseActorTypeGridConnectionPoint),
				UseCaseSupport: []model.UseCaseSupportType{
					{
						UseCaseName:      util.Ptr(model.UseCaseNameTypeMonitoringOfGridConnectionPoint),
						UseCaseAvailable: util.Ptr(true),
						ScenarioSupport:  []model.UseCaseScenarioSupportType{2, 3, 4},
					},
				},
			},
		},
	}

	nodemgmtEntity := s.remoteDevice.Entity([]model.AddressEntityType{0})
	nodeFeature := s.remoteDevice.FeatureByEntityTypeAndRole(nodemgmtEntity, model.FeatureTypeTypeNodeManagement, model.RoleTypeSpecial)
	fErr := nodeFeature.UpdateData(model.FunctionTypeNodeManagementUseCaseData, ucData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.IsUseCaseSupported(s.smgwEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), false, data)

	descData := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				ScopeType:     util.Ptr(model.ScopeTypeTypeACPower),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				ScopeType:     util.Ptr(model.ScopeTypeTypeGridFeedIn),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(2)),
				ScopeType:     util.Ptr(model.ScopeTypeTypeGridConsumption),
			},
		},
	}

	measurementFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.smgwEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	fErr = measurementFeature.UpdateData(model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.IsUseCaseSupported(s.smgwEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), false, data)

	elData := &model.ElectricalConnectionDescriptionListDataType{
		ElectricalConnectionDescriptionData: []model.ElectricalConnectionDescriptionDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
			},
		},
	}

	elFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.smgwEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	fErr = elFeature.UpdateData(model.FunctionTypeElectricalConnectionDescriptionListData, elData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.IsUseCaseSupported(s.smgwEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), true, data)
}
