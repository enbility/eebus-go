package service

import (
	"crypto/tls"
	"errors"
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
		[]model.EntityTypeType{model.EntityTypeTypeCEM}, 4729, certificate, time.Second*4, nil, nil)
	assert.Nil(s.T(), nil, err)

	s.sut = NewService(s.config, s.serviceReader)
}

func (s *ServiceSuite) Test_AddUseCase() {
	ucMock := mocks.NewUseCaseInterface(s.T())
	ucMock.EXPECT().AddFeatures().Return(nil).Once()
	ucMock.EXPECT().AddUseCase().Return().Once()

	s.sut.AddUseCase(ucMock)
}

func (s *ServiceSuite) Test_AddUseCase_Error() {
	ucMock := mocks.NewUseCaseInterface(s.T())
	ucMock.EXPECT().AddFeatures().Return(assert.AnError).Once()

	err := s.sut.AddUseCase(ucMock)
	assert.Equal(s.T(), assert.AnError, err)
}

func (s *ServiceSuite) Test_EEBUSHandler() {
	testSki := "test"

	s.sut.spineLocalDevice = s.localDevice

	testIdentity := shipapi.NewServiceIdentity(testSki, "", "")

	entry := shipapi.RemoteMdnsService{
		Ski: testSki,
	}

	entries := []shipapi.RemoteMdnsService{entry}
	s.serviceReader.EXPECT().VisibleRemoteMdnsServicesUpdated(mock.Anything, mock.Anything).Return()
	s.sut.VisibleRemoteMdnsServicesUpdated(entries)

	s.serviceReader.EXPECT().RemoteServiceConnected(mock.Anything, mock.Anything).Return()
	s.sut.RemoteServiceConnected(testIdentity)

	s.serviceReader.EXPECT().RemoteServiceDisconnected(mock.Anything, mock.Anything).Return()
	s.localDevice.EXPECT().RemoveRemoteDeviceConnection(testSki).Return()
	s.sut.RemoteServiceDisconnected(testIdentity)

	s.serviceReader.EXPECT().ServiceUpdated(mock.Anything).Return()
	s.sut.ServiceUpdated(testIdentity)

	s.serviceReader.EXPECT().ServicePairingDetailUpdate(mock.Anything, mock.Anything).Return()
	detail := &shipapi.ConnectionStateDetail{}
	s.sut.ServicePairingDetailUpdate(testIdentity, detail)

	s.sut.UserIsAbleToApproveOrCancelPairingRequests(true)
	result := s.sut.AllowWaitingForTrust(testIdentity)
	assert.Equal(s.T(), true, result)

	// Test AllowWaitingForTrust returns false when pairing not possible
	s.sut.UserIsAbleToApproveOrCancelPairingRequests(false)
	result = s.sut.AllowWaitingForTrust(testIdentity)
	assert.Equal(s.T(), false, result)

	// Test ServiceAutoTrusted
	s.serviceReader.EXPECT().ServiceAutoTrusted(mock.Anything, mock.Anything).Return()
	s.sut.ServiceAutoTrusted(testIdentity)

	// Test ServiceAutoTrustFailed
	testErr := errors.New("pairing failed")
	s.serviceReader.EXPECT().ServiceAutoTrustFailed(mock.Anything, mock.Anything, mock.Anything).Return()
	s.sut.ServiceAutoTrustFailed(testIdentity, testErr)

	// Test ServiceAutoTrustRemoved
	s.serviceReader.EXPECT().ServiceAutoTrustRemoved(mock.Anything, mock.Anything, mock.Anything).Return()
	s.sut.ServiceAutoTrustRemoved(testIdentity, "device replaced")

	conf := s.sut.Configuration()
	assert.Equal(s.T(), s.sut.configuration, conf)

	lService := s.sut.LocalService()
	assert.Equal(s.T(), s.sut.localService, lService)
}

func (s *ServiceSuite) Test_ConnectionsHub() {
	testSki := "test"
	testIdentity := shipapi.NewServiceIdentity(testSki, "", "")

	s.sut.connectionsHub = s.conHub
	s.sut.mdns = s.mdns
	s.sut.spineLocalDevice = s.localDevice
	s.sut.localService, _ = shipapi.NewServiceDetails(testSki, "", "")

	s.conHub.EXPECT().PairingDetailFor(mock.Anything).Return(nil)
	s.sut.PairingDetailFor(testIdentity)

	s.conHub.EXPECT().ServiceFor(mock.Anything).Return(nil)
	details := s.sut.RemoteServiceFor(testIdentity)
	assert.Nil(s.T(), details)

	s.localDevice.EXPECT().SetupRemoteDevice(mock.Anything, s).Return(nil)
	s.sut.SetupRemoteService(testIdentity, s)

	s.conHub.EXPECT().SetAutoAccept(mock.Anything).Return()
	s.sut.SetAutoAccept(true)
	assert.True(s.T(), s.sut.IsAutoAcceptEnabled())

	s.conHub.EXPECT().GeneratePairingQR().Return("text", nil)
	qrCode, err := s.sut.QRCodeText()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), "text", qrCode)
	s.conHub.EXPECT().RegisterRemoteService(mock.Anything).Return()
	s.sut.RegisterRemoteService(testIdentity)

	s.conHub.EXPECT().UnregisterRemoteService(mock.Anything).Return()
	s.sut.UnregisterRemoteService(testIdentity)

	s.conHub.EXPECT().CancelPairing(mock.Anything).Return()
	s.sut.CancelPairing(testIdentity)

	s.conHub.EXPECT().DisconnectService(mock.Anything, mock.Anything).Return()
	s.sut.DisconnectService(testIdentity, "reason")

	// Test GetLocalCertificateFingerprint
	s.conHub.EXPECT().GetLocalCertificateFingerprint().Return("fingerprint123", nil)
	fp, err := s.sut.GetLocalCertificateFingerprint()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), "fingerprint123", fp)

	// Test StartAnnouncementTo
	target := shipapi.PairingTarget{
		SKI:    "targetSki",
		ShipID: "targetShipId",
	}
	s.conHub.EXPECT().StartAnnouncementTo(mock.Anything).Return(nil)
	err = s.sut.StartAnnouncementTo(target)
	assert.Nil(s.T(), err)

	// Test StopAnnouncementTo
	s.conHub.EXPECT().StopAnnouncementTo("targetShipId").Return(nil)
	err = s.sut.StopAnnouncementTo("targetShipId")
	assert.Nil(s.T(), err)

	// Test IsAnnouncingTo
	s.conHub.EXPECT().IsAnnouncingTo("targetShipId").Return(true)
	isAnnouncing := s.sut.IsAnnouncingTo("targetShipId")
	assert.True(s.T(), isAnnouncing)

	// Test GetActiveAnnouncements
	s.conHub.EXPECT().GetActiveAnnouncements().Return([]string{"ship1", "ship2"})
	announcements := s.sut.GetActiveAnnouncements()
	assert.Equal(s.T(), []string{"ship1", "ship2"}, announcements)

	// Test GetTrustedAddCuDevice
	service, _ := shipapi.NewServiceDetails("", "fpValue", "shipIdValue")
	s.conHub.EXPECT().GetTrustedAddCuDevice().Return(service)
	svc := s.sut.GetTrustedAddCuDevice()
	assert.Equal(s.T(), "fpValue", svc.Fingerprint())
	assert.Equal(s.T(), "shipIdValue", svc.ShipID())
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
	s.conHub.EXPECT().Start().Return(nil).Once()
	_ = s.sut.Start()

	time.Sleep(time.Millisecond * 200)

	isRunning := s.sut.IsRunning()
	assert.True(s.T(), isRunning)

	// nothing should happen
	_ = s.sut.Start()

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
		[]model.EntityTypeType{model.EntityTypeTypeCEM}, 4729, certificate, time.Second*4, nil, nil)
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
	s.conHub.EXPECT().Start().Return(nil)
	_ = s.sut.Start()

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
		[]model.EntityTypeType{model.EntityTypeTypeCEM}, 4729, certificate, time.Second*4, nil, nil)
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

func (s *ServiceSuite) Test_QRCodeText_NoSetup() {
	// QRCodeText should return error when connectionsHub is nil (not set up)
	qr, err := s.sut.QRCodeText()
	assert.NotNil(s.T(), err)
	assert.Empty(s.T(), qr)
}

func (s *ServiceSuite) Test_QRCodeText_Error() {
	// QRCodeText should propagate errors from connectionsHub
	s.sut.connectionsHub = s.conHub
	expectedError := errors.New("qr generation failed")
	s.conHub.EXPECT().GeneratePairingQR().Return("", expectedError)
	qr, err := s.sut.QRCodeText()
	assert.NotNil(s.T(), err)
	assert.Empty(s.T(), qr)
	assert.Equal(s.T(), err.Error(), expectedError.Error())
}

func (s *ServiceSuite) Test_Start_Error() {
	// Start should propagate errors from connectionsHub.Start()
	s.sut.connectionsHub = s.conHub
	expectedError := errors.New("start failed")
	s.conHub.EXPECT().Start().Return(expectedError).Once()
	err := s.sut.Start()
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), err.Error(), expectedError.Error())
	assert.False(s.T(), s.sut.IsRunning())
}

func (s *ServiceSuite) Test_RemoteServiceDisconnected_NilLocalDevice() {
	// RemoteServiceDisconnected should not panic when spineLocalDevice is nil
	testIdentity := shipapi.NewServiceIdentity("test", "", "")
	s.sut.spineLocalDevice = nil
	s.serviceReader.EXPECT().RemoteServiceDisconnected(mock.Anything, mock.Anything).Return()
	s.sut.RemoteServiceDisconnected(testIdentity)
}

func (s *ServiceSuite) Test_IsAnnouncingTo_False() {
	s.sut.connectionsHub = s.conHub
	s.conHub.EXPECT().IsAnnouncingTo("nonexistent").Return(false)
	result := s.sut.IsAnnouncingTo("nonexistent")
	assert.False(s.T(), result)
}

func (s *ServiceSuite) Test_GetActiveAnnouncements_Empty() {
	s.sut.connectionsHub = s.conHub
	s.conHub.EXPECT().GetActiveAnnouncements().Return([]string{})
	announcements := s.sut.GetActiveAnnouncements()
	assert.Empty(s.T(), announcements)
}

func (s *ServiceSuite) Test_GetTrustedAddCuDevice_Empty() {
	s.sut.connectionsHub = s.conHub
	s.conHub.EXPECT().GetTrustedAddCuDevice().Return(nil)
	svc := s.sut.GetTrustedAddCuDevice()
	assert.Nil(s.T(), svc)
}

func (s *ServiceSuite) Test_StartAnnouncementTo_Error() {
	s.sut.connectionsHub = s.conHub
	target := shipapi.PairingTarget{
		SKI:    "ski",
		ShipID: "shipId",
	}
	expectedError := errors.New("announcement failed")
	s.conHub.EXPECT().StartAnnouncementTo(mock.Anything).Return(expectedError)
	err := s.sut.StartAnnouncementTo(target)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), err.Error(), expectedError.Error())
}

func (s *ServiceSuite) Test_StopAnnouncementTo_Error() {
	s.sut.connectionsHub = s.conHub
	expectedError := errors.New("stop failed")
	s.conHub.EXPECT().StopAnnouncementTo("ship1").Return(expectedError)
	err := s.sut.StopAnnouncementTo("ship1")
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), err.Error(), expectedError.Error())
}

func (s *ServiceSuite) Test_GetLocalCertificateFingerprint_Error() {
	s.sut.connectionsHub = s.conHub
	expectedError := errors.New("fingerprint error")
	s.conHub.EXPECT().GetLocalCertificateFingerprint().Return("", expectedError)
	fp, err := s.sut.GetLocalCertificateFingerprint()
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), err.Error(), expectedError.Error())
	assert.Empty(s.T(), fp)
}
