package service

import (
	"net"
	"testing"
	"time"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/cert"
	"github.com/enbility/eebus-go/mocks"
	"github.com/enbility/eebus-go/util"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestMdnsSuite(t *testing.T) {
	suite.Run(t, new(MdnsSuite))
}

type MdnsSuite struct {
	suite.Suite

	sut *mdnsManager

	config *api.Configuration

	mdnsService  *mocks.MdnsService
	mdnsSearch   *mocks.MdnsSearch
	mdnsProvider *mocks.MdnsProvider
}

func (s *MdnsSuite) SetupSuite()   {}
func (s *MdnsSuite) TearDownTest() {}

func (s *MdnsSuite) BeforeTest(suiteName, testName string) {
	s.mdnsService = mocks.NewMdnsService(s.T())

	s.mdnsSearch = mocks.NewMdnsSearch(s.T())
	s.mdnsSearch.EXPECT().ReportMdnsEntries(mock.Anything).Maybe()

	s.mdnsProvider = mocks.NewMdnsProvider(s.T())
	s.mdnsProvider.On("ResolveEntries", mock.Anything, mock.Anything).Maybe().Return()
	s.mdnsProvider.On("Shutdown").Maybe().Return()

	certificate, _ := cert.CreateCertificate("unit", "org", "DE", "CN")

	s.config, _ = api.NewConfiguration(
		"vendor", "brand", "model", "serial", model.DeviceTypeTypeEnergyManagementSystem,
		[]model.EntityTypeType{model.EntityTypeTypeCEM}, 4729, certificate, 230.0, time.Second*4)

	s.sut = newMDNS("test", s.config)
	s.sut.mdnsProvider = s.mdnsProvider
}

func (s *MdnsSuite) Test_SetupMdnsService() {
	err := s.sut.SetupMdnsService()
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), true, s.sut.isAnnounced)

	s.sut.UnannounceMdnsEntry()
	assert.Equal(s.T(), false, s.sut.isAnnounced)

	s.sut.UnannounceMdnsEntry()
	assert.Equal(s.T(), false, s.sut.isAnnounced)

	ifaces, err := net.Interfaces()
	assert.NotEqual(s.T(), 0, len(ifaces))
	assert.Nil(s.T(), err)

	// we don't have access to iface names on CI
	if !util.IsRunningOnCI() {
		s.config.SetInterfaces([]string{ifaces[0].Name})
		err = s.sut.SetupMdnsService()
		assert.Nil(s.T(), err)
	}

	s.config.SetInterfaces([]string{"noifacename"})
	err = s.sut.SetupMdnsService()
	assert.NotNil(s.T(), err)

	assert.Equal(s.T(), false, s.sut.isSearchingServices)
	s.sut.setIsSearchingServices(true)
	assert.Equal(s.T(), true, s.sut.isSearchingServices)

}

func (s *MdnsSuite) Test_ShutdownMdnsService() {
	s.sut.ShutdownMdnsService()
	assert.Nil(s.T(), s.sut.mdnsProvider)
}

func (s *MdnsSuite) Test_MdnsEntry() {
	testSki := "test"

	entries := s.sut.mdnsEntries()
	assert.Equal(s.T(), 0, len(entries))

	entry := &api.MdnsEntry{
		Ski: testSki,
	}

	s.sut.setMdnsEntry(testSki, entry)
	entries = s.sut.mdnsEntries()
	assert.Equal(s.T(), 1, len(entries))

	theEntry, ok := s.sut.mdnsEntry(testSki)
	assert.Equal(s.T(), true, ok)
	assert.NotNil(s.T(), theEntry)

	copyEntries := s.sut.copyMdnsEntries()
	assert.Equal(s.T(), 1, len(copyEntries))

	s.sut.removeMdnsEntry(testSki)
	entries = s.sut.mdnsEntries()
	assert.Equal(s.T(), 0, len(entries))
	assert.Equal(s.T(), 1, len(copyEntries))
}

func (s *MdnsSuite) Test_MdnsSearch() {
	assert.Equal(s.T(), false, s.sut.isSearchingServices)
	s.sut.RegisterMdnsSearch(s.mdnsSearch)
	assert.Equal(s.T(), true, s.sut.isSearchingServices)

	s.sut.setIsSearchingServices(true)
	assert.Equal(s.T(), true, s.sut.isSearchingServices)

	s.sut.RegisterMdnsSearch(s.mdnsSearch)

	testSki := "test"

	entry := &api.MdnsEntry{
		Ski: testSki,
	}
	s.sut.setMdnsEntry(testSki, entry)
	entries := s.sut.mdnsEntries()
	assert.Equal(s.T(), 1, len(entries))

	s.sut.setIsSearchingServices(false)

	s.sut.RegisterMdnsSearch(s.mdnsSearch)

	// wait a bit as ResolveEntries is called in a goroutine
	time.Sleep(time.Millisecond * 200)

	s.sut.UnregisterMdnsSearch(s.mdnsSearch)
}

func (s *MdnsSuite) Test_ProcessMdnsEntry() {
	elements := make(map[string]string, 1)

	name := "name"
	host := "host"
	ips := []net.IP{}
	port := 4567

	s.sut.processMdnsEntry(elements, name, host, ips, port, false)
	assert.Equal(s.T(), 0, len(s.sut.mdnsEntries()))

	elements["txtvers"] = "2"
	elements["id"] = "id"
	elements["path"] = "/ship"
	elements["ski"] = "testski"
	elements["register"] = "falsee"

	s.sut.processMdnsEntry(elements, name, host, ips, port, false)
	assert.Equal(s.T(), 0, len(s.sut.mdnsEntries()))

	elements["txtvers"] = "1"
	s.sut.processMdnsEntry(elements, name, host, ips, port, false)
	assert.Equal(s.T(), 0, len(s.sut.mdnsEntries()))

	elements["ski"] = s.sut.ski
	s.sut.processMdnsEntry(elements, name, host, ips, port, false)
	assert.Equal(s.T(), 0, len(s.sut.mdnsEntries()))

	elements["ski"] = "testski"
	s.sut.processMdnsEntry(elements, name, host, ips, port, false)
	assert.Equal(s.T(), 0, len(s.sut.mdnsEntries()))

	elements["register"] = "false"
	s.sut.processMdnsEntry(elements, name, host, ips, port, false)
	assert.Equal(s.T(), 1, len(s.sut.mdnsEntries()))

	elements["brand"] = "brand"
	elements["type"] = "type"
	elements["model"] = "model"
	s.sut.processMdnsEntry(elements, name, host, ips, port, false)
	assert.Equal(s.T(), 1, len(s.sut.mdnsEntries()))

	ips = []net.IP{[]byte("127.0.0.1")}
	s.sut.processMdnsEntry(elements, name, host, ips, port, false)
	assert.Equal(s.T(), 1, len(s.sut.mdnsEntries()))

	s.sut.searchDelegate = s.mdnsSearch
	s.sut.processMdnsEntry(elements, name, host, ips, port, false)
	assert.Equal(s.T(), 1, len(s.sut.mdnsEntries()))

	s.sut.processMdnsEntry(elements, name, host, ips, port, true)
	assert.Equal(s.T(), 0, len(s.sut.mdnsEntries()))
}
