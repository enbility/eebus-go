package integrationtests

import (
	"testing"
	"time"

	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	m_subscriptionRequestCall_recv_result_file_path = "./testdata/m_subscriptionRequestCall_recv_result.json"
	m_descriptionListData_recv_reply_file_path      = "./testdata/m_descriptionListData_recv_reply.json"
	m_measurementListData_recv_notify_file_path     = "./testdata/m_measurementListData_recv_notify.json"
)

func TestMeasurementSuite(t *testing.T) {
	suite.Run(t, new(MeasurementSuite))
}

type MeasurementSuite struct {
	suite.Suite
	sut *spine.DeviceLocalImpl

	remoteSki string

	readHandler  spine.SpineDataProcessing
	writeHandler *WriteMessageHandler
}

func (s *MeasurementSuite) SetupSuite() {
}

func (s *MeasurementSuite) BeforeTest(suiteName, testName string) {
	s.sut, s.remoteSki, s.readHandler, s.writeHandler = beforeTest(suiteName, testName, 2, model.FeatureTypeTypeMeasurement, model.RoleTypeClient)
	initialCommunication(s.T(), s.readHandler, s.writeHandler)
}

func (s *MeasurementSuite) AfterTest(suiteName, testName string) {
}

func (s *MeasurementSuite) TestDescriptionList_Recv() {
	// Act
	msgCounter, _ := s.readHandler.HandleIncomingSpineMesssage(loadFileData(s.T(), m_descriptionListData_recv_reply_file_path))
	waitForAck(s.T(), msgCounter, s.writeHandler)

	// Assert
	remoteDevice := s.sut.RemoteDeviceForSki(s.remoteSki)
	assert.NotNil(s.T(), remoteDevice)

	mFeature := remoteDevice.FeatureByEntityTypeAndRole(
		remoteDevice.Entity(spine.NewAddressEntityType([]uint{1, 1})),
		model.FeatureTypeTypeMeasurement,
		model.RoleTypeServer)
	assert.NotNil(s.T(), mFeature)

	fdata := mFeature.DataCopy(model.FunctionTypeMeasurementDescriptionListData)
	if !assert.NotNil(s.T(), fdata) {
		return
	}
	data := fdata.(*model.MeasurementDescriptionListDataType)

	if !assert.Equal(s.T(), 3, len(data.MeasurementDescriptionData)) {
		return
	}

	item1 := data.MeasurementDescriptionData[0]
	assert.Equal(s.T(), 1, int(*item1.MeasurementId))
	assert.Equal(s.T(), string(model.MeasurementTypeTypeCurrent), string(*item1.MeasurementType))
	assert.Equal(s.T(), string(model.CommodityTypeTypeElectricity), string(*item1.CommodityType))
	assert.Equal(s.T(), string(model.UnitOfMeasurementTypeA), string(*item1.Unit))
	assert.Equal(s.T(), string(model.ScopeTypeTypeACCurrent), string(*item1.ScopeType))
}

func (s *MeasurementSuite) TestMeasurementList_Recv() {
	// Act
	msgCounter, _ := s.readHandler.HandleIncomingSpineMesssage(loadFileData(s.T(), m_measurementListData_recv_notify_file_path))
	waitForAck(s.T(), msgCounter, s.writeHandler)

	// Assert
	remoteDevice := s.sut.RemoteDeviceForSki(s.remoteSki)
	assert.NotNil(s.T(), remoteDevice)

	mFeature := remoteDevice.FeatureByEntityTypeAndRole(
		remoteDevice.Entity(spine.NewAddressEntityType([]uint{1, 1})),
		model.FeatureTypeTypeMeasurement,
		model.RoleTypeServer)
	assert.NotNil(s.T(), mFeature)

	fdata := mFeature.DataCopy(model.FunctionTypeMeasurementListData)
	if !assert.NotNil(s.T(), fdata) {
		return
	}
	data := fdata.(*model.MeasurementListDataType)

	if !assert.Equal(s.T(), 3, len(data.MeasurementData)) {
		return
	}

	item1 := data.MeasurementData[0]
	assert.Equal(s.T(), 1, int(*item1.MeasurementId))
	assert.Equal(s.T(), string(model.MeasurementValueTypeTypeValue), string(*item1.ValueType))
	assert.Equal(s.T(), 5.0, item1.Value.GetValue())
	timestamp, err := item1.Timestamp.GetDateTimeType().GetTime()
	assert.Nil(s.T(), err)
	compareTimestamp := time.Date(
		2022, 11, 19, 15, 21, 50, 3000000, time.UTC)
	assert.Equal(s.T(), compareTimestamp, timestamp)
	assert.Equal(s.T(), string(model.MeasurementValueSourceTypeMeasuredValue), string(*item1.ValueSource))
}

func (s *MeasurementSuite) TestMeasurementByScope_Recv() {
	// Act
	msgCounter, _ := s.readHandler.HandleIncomingSpineMesssage(loadFileData(s.T(), m_descriptionListData_recv_reply_file_path))
	waitForAck(s.T(), msgCounter, s.writeHandler)

	// Act
	msgCounter, _ = s.readHandler.HandleIncomingSpineMesssage(loadFileData(s.T(), m_measurementListData_recv_notify_file_path))
	waitForAck(s.T(), msgCounter, s.writeHandler)

	// Assert
	remoteDevice := s.sut.RemoteDeviceForSki(s.remoteSki)
	assert.NotNil(s.T(), remoteDevice)

	mFeature := remoteDevice.FeatureByEntityTypeAndRole(
		remoteDevice.Entity(spine.NewAddressEntityType([]uint{1, 1})),
		model.FeatureTypeTypeMeasurement,
		model.RoleTypeServer)
	assert.NotNil(s.T(), mFeature)

	fdata := mFeature.DataCopy(model.FunctionTypeMeasurementDescriptionListData)
	if !assert.NotNil(s.T(), fdata) {
		return
	}
	descData := fdata.(*model.MeasurementDescriptionListDataType)

	if !assert.Equal(s.T(), 3, len(descData.MeasurementDescriptionData)) {
		return
	}

	fdata = mFeature.DataCopy(model.FunctionTypeMeasurementListData)
	if !assert.NotNil(s.T(), fdata) {
		return
	}
	mData := fdata.(*model.MeasurementListDataType)

	if !assert.Equal(s.T(), 3, len(mData.MeasurementData)) {
		return
	}

	item1 := mData.MeasurementData[0]
	assert.Equal(s.T(), 1, int(*item1.MeasurementId))
	assert.Equal(s.T(), string(model.MeasurementValueTypeTypeValue), string(*item1.ValueType))
	assert.Equal(s.T(), 5.0, item1.Value.GetValue())
	timestamp, err := item1.Timestamp.GetDateTimeType().GetTime()
	assert.Nil(s.T(), err)
	compareTimestamp := time.Date(
		2022, 11, 19, 15, 21, 50, 3000000, time.UTC)
	assert.Equal(s.T(), compareTimestamp, timestamp)
	assert.Equal(s.T(), string(model.MeasurementValueSourceTypeMeasuredValue), string(*item1.ValueSource))
}
