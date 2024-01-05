package model

import (
	"testing"

	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func TestMeasurementListDataType_Update_Add(t *testing.T) {
	sut := MeasurementListDataType{
		MeasurementData: []MeasurementDataType{
			{
				MeasurementId: util.Ptr(MeasurementIdType(0)),
				ValueType:     util.Ptr(MeasurementValueTypeTypeAverageValue),
				Value:         NewScaledNumberType(1),
			},
			{
				MeasurementId: util.Ptr(MeasurementIdType(1)),
				ValueType:     util.Ptr(MeasurementValueTypeTypeAverageValue),
				Value:         NewScaledNumberType(1),
			},
		},
	}

	newData := MeasurementListDataType{
		MeasurementData: []MeasurementDataType{
			{
				MeasurementId: util.Ptr(MeasurementIdType(1)),
				ValueType:     util.Ptr(MeasurementValueTypeTypeValue),
				Value:         NewScaledNumberType(10),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

	data := sut.MeasurementData
	// check the non changing items
	assert.Equal(t, 3, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.MeasurementId))
	assert.Equal(t, MeasurementValueTypeTypeAverageValue, *item1.ValueType)
	assert.Equal(t, 1.0, item1.Value.GetValue())
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.MeasurementId))
	assert.Equal(t, MeasurementValueTypeTypeAverageValue, *item2.ValueType)
	assert.Equal(t, 1.0, item2.Value.GetValue())
}

func TestMeasurementListDataType_Update_Replace(t *testing.T) {
	sut := MeasurementListDataType{
		MeasurementData: []MeasurementDataType{
			{
				MeasurementId: util.Ptr(MeasurementIdType(0)),
				ValueType:     util.Ptr(MeasurementValueTypeTypeAverageValue),
				Value:         NewScaledNumberType(1),
			},
			{
				MeasurementId: util.Ptr(MeasurementIdType(1)),
				ValueType:     util.Ptr(MeasurementValueTypeTypeValue),
				Value:         NewScaledNumberType(1),
			},
		},
	}

	newData := MeasurementListDataType{
		MeasurementData: []MeasurementDataType{
			{
				MeasurementId: util.Ptr(MeasurementIdType(1)),
				ValueType:     util.Ptr(MeasurementValueTypeTypeValue),
				Value:         NewScaledNumberType(10),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

	data := sut.MeasurementData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.MeasurementId))
	assert.Equal(t, MeasurementValueTypeTypeAverageValue, *item1.ValueType)
	assert.Equal(t, 1.0, item1.Value.GetValue())
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.MeasurementId))
	assert.Equal(t, MeasurementValueTypeTypeValue, *item2.ValueType)
	assert.Equal(t, 10.0, item2.Value.GetValue())
}

func TestMeasurementConstraintsListDataType_Update(t *testing.T) {
	sut := MeasurementConstraintsListDataType{
		MeasurementConstraintsData: []MeasurementConstraintsDataType{
			{
				MeasurementId: util.Ptr(MeasurementIdType(0)),
				ValueStepSize: NewScaledNumberType(1),
			},
			{
				MeasurementId: util.Ptr(MeasurementIdType(1)),
				ValueStepSize: NewScaledNumberType(1),
			},
		},
	}

	newData := MeasurementConstraintsListDataType{
		MeasurementConstraintsData: []MeasurementConstraintsDataType{
			{
				MeasurementId: util.Ptr(MeasurementIdType(1)),
				ValueStepSize: NewScaledNumberType(10),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

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
	sut := MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []MeasurementDescriptionDataType{
			{
				MeasurementId: util.Ptr(MeasurementIdType(0)),
				ScopeType:     util.Ptr(ScopeTypeTypeACCurrent),
			},
			{
				MeasurementId: util.Ptr(MeasurementIdType(1)),
				ScopeType:     util.Ptr(ScopeTypeTypeACCurrent),
			},
		},
	}

	newData := MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []MeasurementDescriptionDataType{
			{
				MeasurementId: util.Ptr(MeasurementIdType(1)),
				ScopeType:     util.Ptr(ScopeTypeTypeACPower),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

	data := sut.MeasurementDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.MeasurementId))
	assert.Equal(t, ScopeTypeTypeACCurrent, *item1.ScopeType)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.MeasurementId))
	assert.Equal(t, ScopeTypeTypeACPower, *item2.ScopeType)
}

func TestMeasurementThresholdRelationListDataType_Update(t *testing.T) {
	sut := MeasurementThresholdRelationListDataType{
		MeasurementThresholdRelationData: []MeasurementThresholdRelationDataType{
			{
				MeasurementId: util.Ptr(MeasurementIdType(0)),
				ThresholdId:   []ThresholdIdType{ThresholdIdType(0)},
			},
			{
				MeasurementId: util.Ptr(MeasurementIdType(1)),
				ThresholdId:   []ThresholdIdType{ThresholdIdType(0)},
			},
		},
	}

	newData := MeasurementThresholdRelationListDataType{
		MeasurementThresholdRelationData: []MeasurementThresholdRelationDataType{
			{
				MeasurementId: util.Ptr(MeasurementIdType(1)),
				ThresholdId:   []ThresholdIdType{ThresholdIdType(1)},
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

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
