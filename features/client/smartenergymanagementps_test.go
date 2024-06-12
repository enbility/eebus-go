package client

import (
	"testing"
	"time"

	shipapi "github.com/enbility/ship-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestSmartEnergyManagementPsSuite(t *testing.T) {
	suite.Run(t, new(SmartEnergyManagementPsSuite))
}

type SmartEnergyManagementPsSuite struct {
	suite.Suite

	localEntity  spineapi.EntityLocalInterface
	remoteEntity spineapi.EntityRemoteInterface

	smartenergymgmtps *SmartEnergyManagementPs
	sentMessage       []byte
}

var _ shipapi.ShipConnectionDataWriterInterface = (*SmartEnergyManagementPsSuite)(nil)

func (s *SmartEnergyManagementPsSuite) WriteShipMessageWithPayload(message []byte) {
	s.sentMessage = message
}

func (s *SmartEnergyManagementPsSuite) BeforeTest(suiteName, testName string) {
	s.localEntity, s.remoteEntity = setupFeatures(
		s.T(),
		s,
		[]featureFunctions{
			{
				featureType: model.FeatureTypeTypeSmartEnergyManagementPs,
				functions: []model.FunctionType{
					model.FunctionTypeSmartEnergyManagementPsData,
				},
			},
		},
	)

	var err error
	s.smartenergymgmtps, err = NewSmartEnergyManagementPs(s.localEntity, nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), s.smartenergymgmtps)

	s.smartenergymgmtps, err = NewSmartEnergyManagementPs(s.localEntity, s.remoteEntity)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), s.smartenergymgmtps)
}

func (s *SmartEnergyManagementPsSuite) Test_RequestData() {
	counter, err := s.smartenergymgmtps.RequestData()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *SmartEnergyManagementPsSuite) Test_WriteData() {
	counter, err := s.smartenergymgmtps.WriteData(nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), counter)

	data := &model.SmartEnergyManagementPsDataType{
		Alternatives: []model.SmartEnergyManagementPsAlternativesType{
			{
				Relation: &model.SmartEnergyManagementPsAlternativesRelationType{
					AlternativesId: util.Ptr(model.AlternativesIdType(1)),
				},
				PowerSequence: []model.SmartEnergyManagementPsPowerSequenceType{
					{
						Description: &model.PowerSequenceDescriptionDataType{
							SequenceId: util.Ptr(model.PowerSequenceIdType(1)),
						},
						State: &model.PowerSequenceStateDataType{
							State: util.Ptr(model.PowerSequenceStateTypeInactive),
						},
						Schedule: &model.PowerSequenceScheduleDataType{
							StartTime: model.NewAbsoluteOrRelativeTimeTypeFromDuration(time.Minute * 28),
							EndTime:   model.NewAbsoluteOrRelativeTimeTypeFromDuration(time.Hour*2 + time.Minute*28),
						},
						PowerTimeSlot: []model.SmartEnergyManagementPsPowerTimeSlotType{
							{
								ValueList: &model.SmartEnergyManagementPsPowerTimeSlotValueListType{
									Value: []model.PowerTimeSlotValueDataType{
										{
											ValueType: util.Ptr(model.PowerTimeSlotValueTypeTypePower),
											Value:     model.NewScaledNumberType(4444),
										},
										{
											ValueType: util.Ptr(model.PowerTimeSlotValueTypeTypePowerMax),
											Value:     model.NewScaledNumberType(9999),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	counter, err = s.smartenergymgmtps.WriteData(data)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}
