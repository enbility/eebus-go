package service

import (
	"crypto/tls"
	"errors"
	"testing"
	"time"

	spineModel "github.com/enbility/eebus-go/spine/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestTypesSuite(t *testing.T) {
	suite.Run(t, new(TypesSuite))
}

type TypesSuite struct {
	suite.Suite
}

func (s *TypesSuite) SetupSuite()   {}
func (s *TypesSuite) TearDownTest() {}

func (s *TypesSuite) BeforeTest(suiteName, testName string) {}

func (s *TypesSuite) Test_ConnectionState() {
	conState := NewConnectionStateDetail(ConnectionStateNone, nil)
	assert.Equal(s.T(), ConnectionStateNone, conState.State())
	assert.Nil(s.T(), conState.Error())

	conState.SetState(ConnectionStateError)
	assert.Equal(s.T(), ConnectionStateError, conState.State())

	conState.SetError(errors.New("test"))
	assert.NotNil(s.T(), conState.Error())
}

func (s *TypesSuite) Test_ServiceDetails() {
	testSki := "test"

	details := NewServiceDetails(testSki)
	assert.NotNil(s.T(), details)
}

func (s *TypesSuite) Test_Configuration() {
	certificate := tls.Certificate{}
	vendor := "vendor"
	brand := "brand"
	model := "model"
	serial := "serial"
	port := 4567
	volt := 230.0

	config, err := NewConfiguration("", brand, model, serial, spineModel.DeviceTypeTypeEnergyManagementSystem,
		[]spineModel.EntityTypeType{spineModel.EntityTypeTypeCEM}, port, certificate, volt, time.Second*4)

	assert.Nil(s.T(), config)
	assert.NotNil(s.T(), err)

	config, err = NewConfiguration(vendor, "", model, serial, spineModel.DeviceTypeTypeEnergyManagementSystem,
		[]spineModel.EntityTypeType{spineModel.EntityTypeTypeCEM}, port, certificate, 230, time.Second*4)

	assert.Nil(s.T(), config)
	assert.NotNil(s.T(), err)

	config, err = NewConfiguration(vendor, brand, "", serial, spineModel.DeviceTypeTypeEnergyManagementSystem,
		[]spineModel.EntityTypeType{spineModel.EntityTypeTypeCEM}, port, certificate, 230, time.Second*4)

	assert.Nil(s.T(), config)
	assert.NotNil(s.T(), err)

	config, err = NewConfiguration(vendor, brand, model, "", spineModel.DeviceTypeTypeEnergyManagementSystem,
		[]spineModel.EntityTypeType{spineModel.EntityTypeTypeCEM}, port, certificate, 230, time.Second*4)

	assert.Nil(s.T(), config)
	assert.NotNil(s.T(), err)

	config, err = NewConfiguration(vendor, brand, model, serial, "",
		[]spineModel.EntityTypeType{spineModel.EntityTypeTypeCEM}, port, certificate, 230, time.Second*4)

	assert.Nil(s.T(), config)
	assert.NotNil(s.T(), err)

	config, err = NewConfiguration(vendor, brand, model, serial, spineModel.DeviceTypeTypeEnergyManagementSystem,
		[]spineModel.EntityTypeType{}, port, certificate, 230, time.Second*4)

	assert.Nil(s.T(), config)
	assert.NotNil(s.T(), err)

	config, err = NewConfiguration(vendor, brand, model, serial, spineModel.DeviceTypeTypeEnergyManagementSystem,
		[]spineModel.EntityTypeType{spineModel.EntityTypeTypeCEM}, port, certificate, 230, time.Second*4)

	assert.NotNil(s.T(), config)
	assert.Nil(s.T(), err)

	ifaces := []string{"lo", "eth0"}
	config.SetInterfaces(ifaces)
	assert.Equal(s.T(), 2, len(config.interfaces))

	config.SetRegisterAutoAccept(true)
	assert.Equal(s.T(), true, config.registerAutoAccept)

	id := config.generateIdentifier()
	assert.NotEqual(s.T(), "", id)

	id = config.Identifier()
	assert.NotEqual(s.T(), "", id)

	id = config.MdnsServiceName()
	assert.NotEqual(s.T(), "", id)

	alternate := "alternate"

	config.SetAlternateIdentifier(alternate)
	id = config.Identifier()
	assert.Equal(s.T(), alternate, id)

	config.SetAlternateMdnsServiceName(alternate)
	id = config.MdnsServiceName()
	assert.Equal(s.T(), alternate, id)

	voltage := config.Voltage()
	assert.Equal(s.T(), volt, voltage)
}
