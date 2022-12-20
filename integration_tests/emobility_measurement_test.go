package integrationtests

import (
	"testing"

	"github.com/enbility/eebus-go/features"
	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestEmobilityMeasurementSuite(t *testing.T) {
	suite.Run(t, new(EmobilityMeasurementSuite))
}

type EmobilityMeasurementSuite struct {
	suite.Suite
	spine.SpineDataConnection

	sut *spine.DeviceLocalImpl

	measurement          *features.Measurement
	electricalconnection *features.ElectricalConnection

	remoteSki string

	readHandler  spine.SpineDataProcessing
	writeHandler *WriteMessageHandler
}

func (s *EmobilityMeasurementSuite) SetupSuite() {
}

func (s *EmobilityMeasurementSuite) BeforeTest(suiteName, testName string) {
	s.sut = spine.NewDeviceLocalImpl("TestBrandName", "TestDeviceModel", "TestSerialNumber", "TestDeviceCode",
		"TestDeviceAddress", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart)
	localEntity := spine.NewEntityLocalImpl(s.sut, model.EntityTypeTypeCEM, spine.NewAddressEntityType([]uint{1}))
	s.sut.AddEntity(localEntity)

	f := spine.NewFeatureLocalImpl(1, localEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeClient)
	localEntity.AddFeature(f)
	f = spine.NewFeatureLocalImpl(2, localEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeClient)
	localEntity.AddFeature(f)

	s.remoteSki = "TestRemoteSki"

	s.writeHandler = &WriteMessageHandler{}
	s.readHandler = s.sut.AddRemoteDevice(s.remoteSki, s.writeHandler)

	initialCommunication(s.T(), s.readHandler, s.writeHandler)
}

func (s *EmobilityMeasurementSuite) AfterTest(suiteName, testName string) {
}

func (s *EmobilityMeasurementSuite) TestGetValuesPerPhaseForScope() {
	remoteEntity := s.sut.RemoteDeviceForSki(s.remoteSki).Entity([]model.AddressEntityType{1, 1})
	assert.NotNil(s.T(), remoteEntity)

	var err error
	s.measurement, err = features.NewMeasurement(model.RoleTypeClient, model.RoleTypeServer, s.sut, remoteEntity)
	assert.Nil(s.T(), err)

	s.electricalconnection, err = features.NewElectricalConnection(model.RoleTypeClient, model.RoleTypeServer, s.sut, remoteEntity)
	assert.Nil(s.T(), err)

	// Act
	msgCounter, _ := s.readHandler.HandleIncomingSpineMesssage(loadFileData(s.T(), ec_parameterdescriptionlistdata_recv_reply_file_path))
	waitForAck(s.T(), msgCounter, s.writeHandler)

	msgCounter, _ = s.readHandler.HandleIncomingSpineMesssage(loadFileData(s.T(), m_descriptionListData_recv_reply_file_path))
	waitForAck(s.T(), msgCounter, s.writeHandler)

	msgCounter, _ = s.readHandler.HandleIncomingSpineMesssage(loadFileData(s.T(), m_measurementListData_recv_notify_file_path))
	waitForAck(s.T(), msgCounter, s.writeHandler)

	measurement := model.MeasurementTypeTypeCurrent
	commodity := model.CommodityTypeTypeElectricity
	scope := model.ScopeTypeTypeACCurrent
	data, err := s.measurement.GetDataForTypeCommodityScope(measurement, commodity, scope)

	// Assert
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
	assert.Equal(s.T(), 1, len(data))
	assert.Equal(s.T(), 5.0, data[0].Value.GetValue())

	measurement = model.MeasurementTypeTypePower
	scope = model.ScopeTypeTypeACPower
	data, err = s.measurement.GetDataForTypeCommodityScope(measurement, commodity, scope)

	// Assert
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
	assert.Equal(s.T(), 1, len(data))
	assert.Equal(s.T(), 1185.0, data[0].Value.GetValue())

	measurement = model.MeasurementTypeTypeEnergy
	scope = model.ScopeTypeTypeCharge
	data, err = s.measurement.GetDataForTypeCommodityScope(measurement, commodity, scope)

	// Assert
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
	assert.Equal(s.T(), 1, len(data))
	assert.Equal(s.T(), 1825.0, data[0].Value.GetValue())
}
