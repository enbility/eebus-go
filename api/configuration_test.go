package api

import (
	"crypto/tls"
	"testing"
	"time"

	shipapi "github.com/enbility/ship-go/api"
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
	categories := []shipapi.DeviceCategoryType{shipapi.DeviceCategoryTypeEnergyManagementSystem}
	port := 4567
	heartbeatTimeout := time.Second * 4
	entityTypes := []spinemodel.EntityTypeType{spinemodel.EntityTypeTypeCEM}

	config, err := NewConfiguration("", brand, model, serial,
		categories,
		spinemodel.DeviceTypeTypeEnergyManagementSystem,
		entityTypes, 0, certificate, heartbeatTimeout)

	assert.Nil(s.T(), config)
	assert.NotNil(s.T(), err)

	config, err = NewConfiguration("", brand, model, serial,
		categories,
		spinemodel.DeviceTypeTypeEnergyManagementSystem,
		entityTypes, port, certificate, heartbeatTimeout)

	assert.Nil(s.T(), config)
	assert.NotNil(s.T(), err)

	config, err = NewConfiguration(vendor, "", model, serial,
		categories,
		spinemodel.DeviceTypeTypeEnergyManagementSystem,
		entityTypes, port, certificate, heartbeatTimeout)

	assert.Nil(s.T(), config)
	assert.NotNil(s.T(), err)

	config, err = NewConfiguration(vendor, brand, "", serial,
		categories,
		spinemodel.DeviceTypeTypeEnergyManagementSystem,
		entityTypes, port, certificate, heartbeatTimeout)

	assert.Nil(s.T(), config)
	assert.NotNil(s.T(), err)

	config, err = NewConfiguration(vendor, brand, model, "",
		categories,
		spinemodel.DeviceTypeTypeEnergyManagementSystem,
		entityTypes, port, certificate, heartbeatTimeout)

	assert.Nil(s.T(), config)
	assert.NotNil(s.T(), err)

	config, err = NewConfiguration(vendor, brand, model, serial,
		nil,
		spinemodel.DeviceTypeTypeEnergyManagementSystem,
		entityTypes, port, certificate, heartbeatTimeout)

	assert.Nil(s.T(), config)
	assert.NotNil(s.T(), err)

	config, err = NewConfiguration(vendor, brand, model, serial,
		categories,
		"",
		entityTypes, port, certificate, heartbeatTimeout)

	assert.Nil(s.T(), config)
	assert.NotNil(s.T(), err)

	config, err = NewConfiguration(vendor, brand, model, serial,
		categories,
		spinemodel.DeviceTypeTypeEnergyManagementSystem,
		[]spinemodel.EntityTypeType{}, port, certificate, heartbeatTimeout)

	assert.Nil(s.T(), config)
	assert.NotNil(s.T(), err)

	config, err = NewConfiguration(vendor, brand, model, serial,
		categories,
		spinemodel.DeviceTypeTypeEnergyManagementSystem,
		entityTypes, port, certificate, heartbeatTimeout)

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

	categoryValue := config.DeviceCategories()
	assert.Equal(s.T(), categories, categoryValue)

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
