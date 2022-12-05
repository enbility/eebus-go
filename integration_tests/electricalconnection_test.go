package integrationtests

import (
	"testing"

	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	ec_permittedvaluesetlistdata_recv_notify_partial_file_path = "./testdata/ec_permittedvaluesetlistdata_recv_notify_partial.json"
	ec_descriptionlistdata_recv_reply_file_path                = "./testdata/ec_descriptionListData_recv_reply.json"
	ec_parameterdescriptionlistdata_recv_reply_file_path       = "./testdata/ec_parameterDescriptionListData_recv_reply.json"
	ec_subscriptionRequestCall_recv_result_file_path           = "./testdata/ec_subscriptionRequestCall_recv_result.json"
)

func TestElectricalConnectionSuite(t *testing.T) {
	suite.Run(t, new(ElectricalConnectionSuite))
}

type ElectricalConnectionSuite struct {
	suite.Suite
	sut *spine.DeviceLocalImpl

	remoteSki string

	readHandler  spine.SpineDataProcessing
	writeHandler *WriteMessageHandler
}

func (s *ElectricalConnectionSuite) SetupSuite() {
}

func (s *ElectricalConnectionSuite) BeforeTest(suiteName, testName string) {
	s.sut, s.remoteSki, s.readHandler, s.writeHandler = beforeTest(suiteName, testName, 1, model.FeatureTypeTypeElectricalConnection, model.RoleTypeClient)
	initialCommunication(s.T(), s.readHandler, s.writeHandler)
}

func (s *ElectricalConnectionSuite) AfterTest(suiteName, testName string) {
}

func (s *ElectricalConnectionSuite) TestDescriptionListData_RecvReply() {
	// Act
	msgCounter, _ := s.readHandler.HandleIncomingSpineMesssage(loadFileData(s.T(), ec_descriptionlistdata_recv_reply_file_path))
	waitForAck(s.T(), msgCounter, s.writeHandler)

	// Assert
	remoteDevice := s.sut.RemoteDeviceForSki(s.remoteSki)
	assert.NotNil(s.T(), remoteDevice)

	ecFeature := remoteDevice.FeatureByEntityTypeAndRole(
		remoteDevice.Entity(spine.NewAddressEntityType([]uint{1, 1})),
		model.FeatureTypeTypeElectricalConnection,
		model.RoleTypeServer)
	assert.NotNil(s.T(), ecFeature)

	fdata := ecFeature.Data(model.FunctionTypeElectricalConnectionDescriptionListData)
	if !assert.NotNil(s.T(), fdata) {
		return
	}
	data := fdata.(*model.ElectricalConnectionDescriptionListDataType)

	if !assert.Equal(s.T(), 1, len(data.ElectricalConnectionDescriptionData)) {
		return
	}

	item1 := data.ElectricalConnectionDescriptionData[0]
	assert.Equal(s.T(), 0, int(*item1.ElectricalConnectionId))
	assert.Equal(s.T(), string(model.ElectricalConnectionVoltageTypeTypeAc), string(*item1.PowerSupplyType))
	assert.Equal(s.T(), 1, int(*item1.AcConnectedPhases))
	assert.Equal(s.T(), string(model.EnergyDirectionTypeConsume), string(*item1.PositiveEnergyDirection))
}

func (s *ElectricalConnectionSuite) TestParameterDescriptionListData_RecvReply() {
	// Act
	msgCounter, _ := s.readHandler.HandleIncomingSpineMesssage(loadFileData(s.T(), ec_parameterdescriptionlistdata_recv_reply_file_path))
	waitForAck(s.T(), msgCounter, s.writeHandler)

	// Assert
	remoteDevice := s.sut.RemoteDeviceForSki(s.remoteSki)
	assert.NotNil(s.T(), remoteDevice)

	ecFeature := remoteDevice.FeatureByEntityTypeAndRole(
		remoteDevice.Entity(spine.NewAddressEntityType([]uint{1, 1})),
		model.FeatureTypeTypeElectricalConnection,
		model.RoleTypeServer)
	assert.NotNil(s.T(), ecFeature)

	fdata := ecFeature.Data(model.FunctionTypeElectricalConnectionParameterDescriptionListData)
	if !assert.NotNil(s.T(), fdata) {
		return
	}
	data := fdata.(*model.ElectricalConnectionParameterDescriptionListDataType)

	if !assert.Equal(s.T(), 4, len(data.ElectricalConnectionParameterDescriptionData)) {
		return
	}

	item1 := data.ElectricalConnectionParameterDescriptionData[0]
	assert.Equal(s.T(), 0, int(*item1.ElectricalConnectionId))
	assert.Equal(s.T(), 1, int(*item1.ParameterId))
	assert.Equal(s.T(), 1, int(*item1.MeasurementId))
	assert.Equal(s.T(), string(model.ElectricalConnectionVoltageTypeTypeAc), string(*item1.VoltageType))
	assert.Equal(s.T(), string(model.ElectricalConnectionPhaseNameTypeA), string(*item1.AcMeasuredPhases))
	assert.Equal(s.T(), string(model.ElectricalConnectionPhaseNameTypeNeutral), string(*item1.AcMeasuredInReferenceTo))
	assert.Equal(s.T(), string(model.ElectricalConnectionAcMeasurementTypeTypeReal), string(*item1.AcMeasurementType))
	assert.Equal(s.T(), string(model.ElectricalConnectionMeasurandVariantTypeRms), string(*item1.AcMeasurementVariant))
}

func (s *ElectricalConnectionSuite) TestPermittedValueSetListData_RecvNotifyPartial() {
	// Act
	msgCounter, _ := s.readHandler.HandleIncomingSpineMesssage(loadFileData(s.T(), ec_permittedvaluesetlistdata_recv_notify_partial_file_path))
	waitForAck(s.T(), msgCounter, s.writeHandler)

	// Assert
	remoteDevice := s.sut.RemoteDeviceForSki(s.remoteSki)
	assert.NotNil(s.T(), remoteDevice)

	ecFeature := remoteDevice.FeatureByEntityTypeAndRole(
		remoteDevice.Entity(spine.NewAddressEntityType([]uint{1, 1})),
		model.FeatureTypeTypeElectricalConnection,
		model.RoleTypeServer)
	assert.NotNil(s.T(), ecFeature)

	fdata := ecFeature.Data(model.FunctionTypeElectricalConnectionPermittedValueSetListData)
	if !assert.NotNil(s.T(), fdata) {
		return
	}
	data := fdata.(*model.ElectricalConnectionPermittedValueSetListDataType)

	if !assert.Equal(s.T(), 3, len(data.ElectricalConnectionPermittedValueSetData)) {
		return
	}

	item1 := data.ElectricalConnectionPermittedValueSetData[0]
	assert.Equal(s.T(), 0, int(*item1.ElectricalConnectionId))
	assert.Equal(s.T(), 1, int(*item1.ParameterId))
	assert.Equal(s.T(), 1, len(item1.PermittedValueSet))
	assert.Equal(s.T(), 1, len(item1.PermittedValueSet[0].Range))
	assert.NotNil(s.T(), item1.PermittedValueSet[0].Range)
	assert.Equal(s.T(), 6, int(*item1.PermittedValueSet[0].Range[0].Min.Number))
	assert.Equal(s.T(), 0, int(*item1.PermittedValueSet[0].Range[0].Min.Scale))
	assert.Equal(s.T(), 16, int(*item1.PermittedValueSet[0].Range[0].Max.Number))
	assert.Equal(s.T(), 0, int(*item1.PermittedValueSet[0].Range[0].Max.Scale))
	assert.Nil(s.T(), item1.PermittedValueSet[0].Value)
}
