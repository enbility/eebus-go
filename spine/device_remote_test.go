package spine

import (
	"testing"
	"time"

	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	nm_usecaseinformationlistdata_recv_reply_file_path = "../spine/testdata/nm_usecaseinformationlistdata_recv_reply.json"
)

func TestDeviceRemoteSuite(t *testing.T) {
	suite.Run(t, new(DeviceRemoteSuite))
}

type DeviceRemoteSuite struct {
	suite.Suite

	localDevice  *DeviceLocalImpl
	remoteDevice *DeviceRemoteImpl
	remoteEntity EntityRemote
}

func (s *DeviceRemoteSuite) WriteSpineMessage([]byte) {}

func (s *DeviceRemoteSuite) SetupSuite() {}

func (s *DeviceRemoteSuite) BeforeTest(suiteName, testName string) {
	s.localDevice = NewDeviceLocalImpl("brand", "model", "serial", "code", "address", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart, time.Second*4)

	ski := "test"
	sender := NewSender(s)
	s.remoteDevice = NewDeviceRemoteImpl(s.localDevice, ski, sender)
	s.remoteDevice.address = util.Ptr(model.AddressDeviceType("test"))
	s.localDevice.AddRemoteDevice(ski, s)

	s.remoteEntity = NewEntityRemoteImpl(s.remoteDevice, model.EntityTypeTypeEVSE, []model.AddressEntityType{1})

	feature := NewFeatureRemoteImpl(0, s.remoteEntity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	s.remoteEntity.AddFeature(feature)

	s.remoteDevice.AddEntity(s.remoteEntity)
}

func (s *DeviceRemoteSuite) Test_RemoveByAddress() {
	assert.Equal(s.T(), 2, len(s.remoteDevice.Entities()))

	s.remoteDevice.RemoveByAddress([]model.AddressEntityType{2})
	assert.Equal(s.T(), 2, len(s.remoteDevice.Entities()))

	s.remoteDevice.RemoveByAddress([]model.AddressEntityType{1})
	assert.Equal(s.T(), 1, len(s.remoteDevice.Entities()))
}

func (s *DeviceRemoteSuite) Test_FeatureByEntityTypeAndRole() {
	entity := s.remoteDevice.Entity([]model.AddressEntityType{1})
	assert.NotNil(s.T(), entity)

	assert.Equal(s.T(), 1, len(entity.Features()))

	feature := s.remoteDevice.FeatureByEntityTypeAndRole(entity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeClient)
	assert.Nil(s.T(), feature)

	feature = s.remoteDevice.FeatureByEntityTypeAndRole(entity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	assert.NotNil(s.T(), feature)

	s.remoteDevice.RemoveByAddress([]model.AddressEntityType{1})
	assert.Equal(s.T(), 1, len(s.remoteDevice.Entities()))

	_ = s.remoteDevice.Entity([]model.AddressEntityType{0})
	s.remoteDevice.RemoveByAddress([]model.AddressEntityType{0})
	assert.Equal(s.T(), 0, len(s.remoteDevice.Entities()))

	feature = s.remoteDevice.FeatureByEntityTypeAndRole(entity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	assert.Nil(s.T(), feature)
}

func (s *DeviceRemoteSuite) Test_VerifyUseCaseScenariosAndFeaturesSupport_ElliJSON() {
	_, _ = s.remoteDevice.HandleIncomingSpineMesssage(loadFileData(s.T(), nm_usecaseinformationlistdata_recv_reply_file_path))

	result := s.remoteDevice.VerifyUseCaseScenariosAndFeaturesSupport(
		model.UseCaseActorTypeBatterySystem,
		model.UseCaseNameTypeControlOfBattery,
		[]model.UseCaseScenarioSupportType{},
		[]model.FeatureTypeType{},
	)
	assert.Equal(s.T(), false, result)

	result = s.remoteDevice.VerifyUseCaseScenariosAndFeaturesSupport(
		model.UseCaseActorTypeEVSE,
		model.UseCaseNameTypeEVSECommissioningAndConfiguration,
		[]model.UseCaseScenarioSupportType{},
		[]model.FeatureTypeType{},
	)
	assert.Equal(s.T(), true, result)
}

func (s *DeviceRemoteSuite) Test_VerifyUseCaseScenariosAndFeaturesSupport() {
	result := s.remoteDevice.VerifyUseCaseScenariosAndFeaturesSupport(
		model.UseCaseActorTypeEVSE,
		model.UseCaseNameTypeEVSECommissioningAndConfiguration,
		[]model.UseCaseScenarioSupportType{},
		[]model.FeatureTypeType{},
	)
	assert.Equal(s.T(), false, result)

	nodeMgmtEntity := s.remoteDevice.Entity(DeviceInformationAddressEntity)
	nodeMgmt := nodeMgmtEntity.Feature(util.Ptr(model.AddressFeatureType(NodeManagementFeatureId)))

	// initialize with empty data
	newData := &model.NodeManagementUseCaseDataType{
		UseCaseInformation: []model.UseCaseInformationDataType{},
	}
	nodeMgmt.UpdateData(model.FunctionTypeNodeManagementUseCaseData, newData, nil, nil)

	data := nodeMgmt.DataCopy(model.FunctionTypeNodeManagementUseCaseData).(*model.NodeManagementUseCaseDataType)

	address := model.FeatureAddressType{
		Device: s.remoteDevice.Address(),
		Entity: s.remoteEntity.Address().Entity,
	}

	data.AddUseCaseSupport(
		address,
		model.UseCaseActorTypeBatterySystem,
		model.UseCaseNameTypeControlOfBattery,
		model.SpecificationVersionType("1.0.0"),
		"",
		true,
		[]model.UseCaseScenarioSupportType{1},
	)
	nodeMgmt.SetData(model.FunctionTypeNodeManagementUseCaseData, data)
	data = nodeMgmt.DataCopy(model.FunctionTypeNodeManagementUseCaseData).(*model.NodeManagementUseCaseDataType)

	result = s.remoteDevice.VerifyUseCaseScenariosAndFeaturesSupport(
		model.UseCaseActorTypeEVSE,
		model.UseCaseNameTypeEVSECommissioningAndConfiguration,
		nil,
		nil,
	)
	assert.Equal(s.T(), false, result)

	data.AddUseCaseSupport(
		address,
		model.UseCaseActorTypeEVSE,
		model.UseCaseNameTypeEVCommissioningAndConfiguration,
		model.SpecificationVersionType("1.0.0"),
		"",
		true,
		[]model.UseCaseScenarioSupportType{1},
	)
	nodeMgmt.SetData(model.FunctionTypeNodeManagementUseCaseData, data)
	data = nodeMgmt.DataCopy(model.FunctionTypeNodeManagementUseCaseData).(*model.NodeManagementUseCaseDataType)

	result = s.remoteDevice.VerifyUseCaseScenariosAndFeaturesSupport(
		model.UseCaseActorTypeEVSE,
		model.UseCaseNameTypeEVSECommissioningAndConfiguration,
		nil,
		nil,
	)
	assert.Equal(s.T(), false, result)

	data.AddUseCaseSupport(
		address,
		model.UseCaseActorTypeEVSE,
		model.UseCaseNameTypeEVSECommissioningAndConfiguration,
		model.SpecificationVersionType("1.0.0"),
		"",
		false,
		[]model.UseCaseScenarioSupportType{1},
	)
	nodeMgmt.SetData(model.FunctionTypeNodeManagementUseCaseData, data)
	data = nodeMgmt.DataCopy(model.FunctionTypeNodeManagementUseCaseData).(*model.NodeManagementUseCaseDataType)

	result = s.remoteDevice.VerifyUseCaseScenariosAndFeaturesSupport(
		model.UseCaseActorTypeEVSE,
		model.UseCaseNameTypeEVSECommissioningAndConfiguration,
		nil,
		nil,
	)
	assert.Equal(s.T(), true, result)

	data.AddUseCaseSupport(
		address,
		model.UseCaseActorTypeEVSE,
		model.UseCaseNameTypeEVSECommissioningAndConfiguration,
		model.SpecificationVersionType("1.0.0"),
		"",
		true,
		[]model.UseCaseScenarioSupportType{1},
	)
	nodeMgmt.SetData(model.FunctionTypeNodeManagementUseCaseData, data)

	result = s.remoteDevice.VerifyUseCaseScenariosAndFeaturesSupport(
		model.UseCaseActorTypeEVSE,
		model.UseCaseNameTypeEVSECommissioningAndConfiguration,
		nil,
		nil,
	)
	assert.Equal(s.T(), true, result)

	result = s.remoteDevice.VerifyUseCaseScenariosAndFeaturesSupport(
		model.UseCaseActorTypeEVSE,
		model.UseCaseNameTypeEVSECommissioningAndConfiguration,
		[]model.UseCaseScenarioSupportType{2},
		nil,
	)
	assert.Equal(s.T(), false, result)

	result = s.remoteDevice.VerifyUseCaseScenariosAndFeaturesSupport(
		model.UseCaseActorTypeEVSE,
		model.UseCaseNameTypeEVSECommissioningAndConfiguration,
		[]model.UseCaseScenarioSupportType{1},
		nil,
	)
	assert.Equal(s.T(), true, result)

	result = s.remoteDevice.VerifyUseCaseScenariosAndFeaturesSupport(
		model.UseCaseActorTypeEVSE,
		model.UseCaseNameTypeEVSECommissioningAndConfiguration,
		[]model.UseCaseScenarioSupportType{1},
		[]model.FeatureTypeType{model.FeatureTypeTypeElectricalConnection},
	)
	assert.Equal(s.T(), false, result)

	entity := s.remoteDevice.Entity([]model.AddressEntityType{1})
	assert.NotNil(s.T(), entity)

	feature := NewFeatureRemoteImpl(0, entity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeClient)
	entity.AddFeature(feature)

	result = s.remoteDevice.VerifyUseCaseScenariosAndFeaturesSupport(
		model.UseCaseActorTypeEVSE,
		model.UseCaseNameTypeEVSECommissioningAndConfiguration,
		[]model.UseCaseScenarioSupportType{1},
		[]model.FeatureTypeType{model.FeatureTypeTypeElectricalConnection},
	)
	assert.Equal(s.T(), false, result)

	feature = NewFeatureRemoteImpl(0, entity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	entity.AddFeature(feature)

	result = s.remoteDevice.VerifyUseCaseScenariosAndFeaturesSupport(
		model.UseCaseActorTypeEVSE,
		model.UseCaseNameTypeEVSECommissioningAndConfiguration,
		[]model.UseCaseScenarioSupportType{1},
		[]model.FeatureTypeType{model.FeatureTypeTypeElectricalConnection},
	)
	assert.Equal(s.T(), true, result)

	s.remoteDevice.RemoveByAddress(feature.Address().Entity)

	result = s.remoteDevice.VerifyUseCaseScenariosAndFeaturesSupport(
		model.UseCaseActorTypeEVSE,
		model.UseCaseNameTypeEVSECommissioningAndConfiguration,
		[]model.UseCaseScenarioSupportType{1},
		[]model.FeatureTypeType{model.FeatureTypeTypeElectricalConnection},
	)
	assert.Equal(s.T(), false, result)
}
