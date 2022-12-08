package model_test

import (
	"testing"

	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func TestMeasurementListDataType_Update(t *testing.T) {
	sut := model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(1),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				Value:         model.NewScaledNumberType(1),
			},
		},
	}

	newData := model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				Value:         model.NewScaledNumberType(10),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.MeasurementData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.MeasurementId))
	assert.Equal(t, 1.0, item1.Value.GetValue())
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.MeasurementId))
	assert.Equal(t, 10.0, item2.Value.GetValue())
}

func TestMeasurementConstraintsListDataType_Update(t *testing.T) {
	sut := model.MeasurementConstraintsListDataType{
		MeasurementConstraintsData: []model.MeasurementConstraintsDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				ValueStepSize: model.NewScaledNumberType(1),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				ValueStepSize: model.NewScaledNumberType(1),
			},
		},
	}

	newData := model.MeasurementConstraintsListDataType{
		MeasurementConstraintsData: []model.MeasurementConstraintsDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				ValueStepSize: model.NewScaledNumberType(10),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.MeasurementConstraintsData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.MeasurementId))
	assert.Equal(t, 1.0, item1.ValueStepSize.GetValue())
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.MeasurementId))
	assert.Equal(t, 10.0, item2.ValueStepSize.GetValue())
}

func TestMeasurementDescriptionListDataType_Update(t *testing.T) {
	sut := model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				ScopeType:     util.Ptr(model.ScopeTypeTypeACCurrent),
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				ScopeType:     util.Ptr(model.ScopeTypeTypeACCurrent),
			},
		},
	}

	newData := model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				ScopeType:     util.Ptr(model.ScopeTypeTypeACPower),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.MeasurementDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.MeasurementId))
	assert.Equal(t, model.ScopeTypeTypeACCurrent, *item1.ScopeType)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.MeasurementId))
	assert.Equal(t, model.ScopeTypeTypeACPower, *item2.ScopeType)
}

func TestMeasurementThresholdRelationListDataType_Update(t *testing.T) {
	sut := model.MeasurementThresholdRelationListDataType{
		MeasurementThresholdRelationData: []model.MeasurementThresholdRelationDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				ThresholdId:   []model.ThresholdIdType{model.ThresholdIdType(0)},
			},
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				ThresholdId:   []model.ThresholdIdType{model.ThresholdIdType(0)},
			},
		},
	}

	newData := model.MeasurementThresholdRelationListDataType{
		MeasurementThresholdRelationData: []model.MeasurementThresholdRelationDataType{
			{
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				ThresholdId:   []model.ThresholdIdType{model.ThresholdIdType(1)},
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.MeasurementThresholdRelationData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.MeasurementId))
	assert.Equal(t, 0, int(item1.ThresholdId[0]))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.MeasurementId))
	assert.Equal(t, 1, int(item2.ThresholdId[0]))
}
