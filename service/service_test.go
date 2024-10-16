package service

import (
	"crypto/tls"
	"testing"
	"time"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/mocks"
	shipapi "github.com/enbility/ship-go/api"
	"github.com/enbility/ship-go/cert"
	"github.com/enbility/ship-go/logging"
	shipmocks "github.com/enbility/ship-go/mocks"
	spinemocks "github.com/enbility/spine-go/mocks"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

type ServiceSuite struct {
	suite.Suite

	config *api.Configuration

	sut *Service

	serviceReader *mocks.ServiceReaderInterface
	conHub        *shipmocks.HubInterface
	mdns          *shipmocks.MdnsInterface
	logging       *shipmocks.LoggingInterface
	localDevice   *spinemocks.DeviceLocalInterface
}

func (s *ServiceSuite) WriteShipMessageWithPayload(message []byte) {}

func (s *ServiceSuite) BeforeTest(suiteName, testName string) {
	s.serviceReader = mocks.NewServiceReaderInterface(s.T())

	s.conHub = shipmocks.NewHubInterface(s.T())

	s.mdns = shipmocks.NewMdnsInterface(s.T())

	s.logging = shipmocks.NewLoggingInterface(s.T())

	s.localDevice = spinemocks.NewDeviceLocalInterface(s.T())

	certificate := tls.Certificate{}
	var err error
	s.config, err = api.NewConfiguration(
		"vendor", "brand", "model", "serial",
		[]shipapi.DeviceCategoryType{shipapi.DeviceCategoryTypeEnergyManagementSystem},
		model.DeviceTypeTypeEnergyManagementSystem,
		[]model.EntityTypeType{model.EntityTypeTypeCEM}, 4729, certificate, time.Second*4)
	assert.Nil(s.T(), nil, err)

	s.sut = NewService(s.config, s.serviceReader)
}

func (s *ServiceSuite) Test_AddUseCase() {
	ucMock := mocks.NewUseCaseInterface(s.T())
	ucMock.EXPECT().AddFeatures().Return().Once()
	ucMock.EXPECT().AddUseCase().Return().Once()

	s.sut.AddUseCase(ucMock)
}

func (s *ServiceSuite) Test_EEBUSHandler() {
	testSki := "test"

	s.sut.spineLocalDevice = s.localDevice

	entry := shipapi.RemoteService{
		Ski: testSki,
	}

	entries := []shipapi.RemoteService{entry}
	s.serviceReader.EXPECT().VisibleRemoteServicesUpdated(mock.Anything, mock.Anything).Return()
	s.sut.VisibleRemoteServicesUpdated(entries)

	s.serviceReader.EXPECT().RemoteSKIConnected(mock.Anything, mock.Anything).Return()
	s.sut.RemoteSKIConnected(testSki)

	s.serviceReader.EXPECT().RemoteSKIDisconnected(mock.Anything, mock.Anything).Return()
	s.localDevice.EXPECT().RemoveRemoteDeviceConnection(testSki).Return()
	s.sut.RemoteSKIDisconnected(testSki)

	s.serviceReader.EXPECT().ServiceShipIDUpdate(mock.Anything, mock.Anything).Return()
	s.sut.ServiceShipIDUpdate(testSki, "shipid")

	s.serviceReader.EXPECT().ServicePairingDetailUpdate(mock.Anything, mock.Anything).Return()
	detail := &shipapi.ConnectionStateDetail{}
	s.sut.ServicePairingDetailUpdate(testSki, detail)

	s.sut.UserIsAbleToApproveOrCancelPairingRequests(true)
	result := s.sut.AllowWaitingForTrust(testSki)
	assert.Equal(s.T(), true, result)

	conf := s.sut.Configuration()
	assert.Equal(s.T(), s.sut.configuration, conf)

	lService := s.sut.LocalService()
	assert.Equal(s.T(), s.sut.localService, lService)
}

func (s *ServiceSuite) Test_ConnectionsHub() {
	testSki := "test"

	s.sut.connectionsHub = s.conHub
	s.sut.mdns = s.mdns
	s.sut.spineLocalDevice = s.localDevice
	s.sut.localService = shipapi.NewServiceDetails(testSki)

	s.conHub.EXPECT().PairingDetailForSki(mock.Anything).Return(nil)
	s.sut.PairingDetailForSki(testSki)

	s.conHub.EXPECT().ServiceForSKI(mock.Anything).Return(nil)
	details := s.sut.RemoteServiceForSKI(testSki)
	assert.Nil(s.T(), details)

	s.localDevice.EXPECT().SetupRemoteDevice(mock.Anything, s).Return(nil)
	s.sut.SetupRemoteDevice(testSki, s)

	s.conHub.EXPECT().SetAutoAccept(mock.Anything).Return()
	s.sut.SetAutoAccept(true)
	assert.True(s.T(), s.sut.IsAutoAcceptEnabled())

	s.mdns.EXPECT().QRCodeText().Return("text")
	assert.Equal(s.T(), "text", s.sut.QRCodeText())

	s.conHub.EXPECT().RegisterRemoteSKI(mock.Anything).Return()
	s.sut.RegisterRemoteSKI(testSki)

	s.conHub.EXPECT().UnregisterRemoteSKI(mock.Anything).Return()
	s.sut.UnregisterRemoteSKI(testSki)

	s.conHub.EXPECT().CancelPairingWithSKI(mock.Anything).Return()
	s.sut.CancelPairingWithSKI(testSki)

	s.conHub.EXPECT().DisconnectSKI(mock.Anything, mock.Anything).Return()
	s.sut.DisconnectSKI(testSki, "reason")
}

func (s *ServiceSuite) Test_SetLogging() {
	s.sut.SetLogging(nil)
	assert.Equal(s.T(), &logging.NoLogging{}, logging.Log())

	s.sut.SetLogging(s.logging)
	assert.Equal(s.T(), s.logging, logging.Log())

	s.sut.SetLogging(&logging.NoLogging{})
	assert.Equal(s.T(), &logging.NoLogging{}, logging.Log())
}

func (s *ServiceSuite) Test_Setup() {
	err := s.sut.Setup()
	assert.NotNil(s.T(), err)

	certificate, err := cert.CreateCertificate("unit", "org", "de", "cn")
	assert.Nil(s.T(), err)
	s.config.SetCertificate(certificate)

	err = s.sut.Setup()
	assert.Nil(s.T(), err)

	address := s.sut.LocalDevice().Address()
	assert.Equal(s.T(), "d:_n:vendor_model-serial", string(*address))

	s.sut.connectionsHub = s.conHub
	s.conHub.EXPECT().Start().Once()
	s.sut.Start()

	time.Sleep(time.Millisecond * 200)

	isRunning := s.sut.IsRunning()
	assert.True(s.T(), isRunning)

	// nothing should happen
	s.sut.Start()

	s.conHub.EXPECT().Shutdown().Once()
	s.sut.Shutdown()

	// nothing should happen
	s.sut.Shutdown()

	device := s.sut.LocalDevice()
	assert.NotNil(s.T(), device)
}

func (s *ServiceSuite) Test_Setup_IANA() {
	var err error
	certificate := tls.Certificate{}
	s.config, err = api.NewConfiguration(
		"12345", "brand", "model", "serial",
		[]shipapi.DeviceCategoryType{shipapi.DeviceCategoryTypeEnergyManagementSystem},
		model.DeviceTypeTypeEnergyManagementSystem,
		[]model.EntityTypeType{model.EntityTypeTypeCEM}, 4729, certificate, time.Second*4)
	assert.Nil(s.T(), nil, err)

	s.sut = NewService(s.config, s.serviceReader)

	err = s.sut.Setup()
	assert.NotNil(s.T(), err)

	certificate, err = cert.CreateCertificate("unit", "org", "de", "cn")
	assert.Nil(s.T(), err)
	s.config.SetCertificate(certificate)

	err = s.sut.Setup()
	assert.Nil(s.T(), err)

	address := s.sut.LocalDevice().Address()
	assert.Equal(s.T(), "d:_i:12345_model-serial", string(*address))

	s.sut.connectionsHub = s.conHub
	s.conHub.EXPECT().Start()
	s.sut.Start()

	time.Sleep(time.Millisecond * 200)

	s.conHub.EXPECT().Shutdown()
	s.sut.Shutdown()

	device := s.sut.LocalDevice()
	assert.NotNil(s.T(), device)
}

func (s *ServiceSuite) Test_Setup_Error_DeviceName() {
	var err error
	certificate := tls.Certificate{}
	s.config, err = api.NewConfiguration(
		"1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890",
		"brand",
		"modelmodelmodelmodelmodelmodelmodelmodelmodelmodelmodelmodelmodelmodelmodelmodelmodelmodelmodelmodel",
		"serialserialserialserialserialserialserialserialserialserialserialserialserialserialserialserialserial",
		[]shipapi.DeviceCategoryType{shipapi.DeviceCategoryTypeEnergyManagementSystem},
		model.DeviceTypeTypeEnergyManagementSystem,
		[]model.EntityTypeType{model.EntityTypeTypeCEM}, 4729, certificate, time.Second*4)
	assert.Nil(s.T(), nil, err)

	s.sut = NewService(s.config, s.serviceReader)

	err = s.sut.Setup()
	assert.NotNil(s.T(), err)

	certificate, err = cert.CreateCertificate("unit", "org", "de", "cn")
	assert.Nil(s.T(), err)
	s.config.SetCertificate(certificate)

	err = s.sut.Setup()
	assert.NotNil(s.T(), err)
}
