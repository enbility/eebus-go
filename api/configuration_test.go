package api

import (
	"crypto/tls"
	"testing"
	"time"

	"github.com/enbility/ship-go/cert"
	"github.com/enbility/ship-go/mdns"
	spinemodel "github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestConfigurationSuite(t *testing.T) {
	suite.Run(t, new(ConfigurationSuite))
}

type ConfigurationSuite struct {
	suite.Suite
}

func (s *ConfigurationSuite) Test_Configuration() {
	certificate, _ := cert.CreateCertificate("unit", "org", "DE", "CN")
	vendor := "vendor"
	brand := "brand"
	model := "model"
	serial := "serial"
	port := 4567
	volt := 230.0
	heartbeatTimeout := time.Second * 4
	entityTypes := []spinemodel.EntityTypeType{spinemodel.EntityTypeTypeCEM}

	config, err := NewConfiguration("", brand, model, serial, spinemodel.DeviceTypeTypeEnergyManagementSystem,
		entityTypes, 0, certificate, volt, heartbeatTimeout)

	assert.Nil(s.T(), config)
	assert.NotNil(s.T(), err)

	config, err = NewConfiguration("", brand, model, serial, spinemodel.DeviceTypeTypeEnergyManagementSystem,
		entityTypes, port, certificate, volt, heartbeatTimeout)

	assert.Nil(s.T(), config)
	assert.NotNil(s.T(), err)

	config, err = NewConfiguration(vendor, "", model, serial, spinemodel.DeviceTypeTypeEnergyManagementSystem,
		entityTypes, port, certificate, 230, heartbeatTimeout)

	assert.Nil(s.T(), config)
	assert.NotNil(s.T(), err)

	config, err = NewConfiguration(vendor, brand, "", serial, spinemodel.DeviceTypeTypeEnergyManagementSystem,
		entityTypes, port, certificate, 230, heartbeatTimeout)

	assert.Nil(s.T(), config)
	assert.NotNil(s.T(), err)

	config, err = NewConfiguration(vendor, brand, model, "", spinemodel.DeviceTypeTypeEnergyManagementSystem,
		entityTypes, port, certificate, 230, heartbeatTimeout)

	assert.Nil(s.T(), config)
	assert.NotNil(s.T(), err)

	config, err = NewConfiguration(vendor, brand, model, serial, "",
		entityTypes, port, certificate, 230, heartbeatTimeout)

	assert.Nil(s.T(), config)
	assert.NotNil(s.T(), err)

	config, err = NewConfiguration(vendor, brand, model, serial, spinemodel.DeviceTypeTypeEnergyManagementSystem,
		[]spinemodel.EntityTypeType{}, port, certificate, 230, heartbeatTimeout)

	assert.Nil(s.T(), config)
	assert.NotNil(s.T(), err)

	config, err = NewConfiguration(vendor, brand, model, serial, spinemodel.DeviceTypeTypeEnergyManagementSystem,
		entityTypes, port, certificate, 230, heartbeatTimeout)

	assert.NotNil(s.T(), config)
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), mdns.MdnsProviderSelectionAll, config.MdnsProviderSelection())

	config.SetMdnsProviderSelection(mdns.MdnsProviderSelectionAvahiOnly)
	assert.Equal(s.T(), mdns.MdnsProviderSelectionAvahiOnly, config.MdnsProviderSelection())

	ifaces := []string{"lo", "eth0"}
	config.SetInterfaces(ifaces)
	assert.Equal(s.T(), 2, len(config.interfaces))

	ifacesValue := config.Interfaces()
	assert.Equal(s.T(), ifaces, ifacesValue)

	config.SetRegisterAutoAccept(true)
	assert.Equal(s.T(), true, config.registerAutoAccept)
	registerValue := config.RegisterAutoAccept()
	assert.Equal(s.T(), true, registerValue)

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

	portValue := config.Port()
	assert.Equal(s.T(), port, portValue)

	heartbeatValue := config.HeartbeatTimeout()
	assert.Equal(s.T(), heartbeatTimeout, heartbeatValue)

	vendorValue := config.VendorCode()
	assert.Equal(s.T(), vendor, vendorValue)

	deviceValue := config.DeviceBrand()
	assert.Equal(s.T(), brand, deviceValue)

	modelValue := config.DeviceModel()
	assert.Equal(s.T(), model, modelValue)

	serialValue := config.DeviceSerialNumber()
	assert.Equal(s.T(), serial, serialValue)

	deviceTypeValue := config.DeviceType()
	assert.Equal(s.T(), spinemodel.DeviceTypeTypeEnergyManagementSystem, deviceTypeValue)

	entityValues := config.EntityTypes()
	assert.Equal(s.T(), entityTypes, entityValues)
	featuresetValue := config.FeatureSet()
	assert.Equal(s.T(), spinemodel.NetworkManagementFeatureSetTypeSmart, featuresetValue)

	testCert := tls.Certificate{}
	config.SetCertificate(testCert)
	certValue := config.Certificate()
	assert.Equal(s.T(), testCert, certValue)
}
