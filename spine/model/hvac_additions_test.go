package model_test

import (
	"testing"

	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func TestHvacSystemFunctionListDataType_Update(t *testing.T) {
	sut := model.HvacSystemFunctionListDataType{
		HvacSystemFunctionData: []model.HvacSystemFunctionDataType{
			{
				SystemFunctionId:       util.Ptr(model.HvacSystemFunctionIdType(0)),
				CurrentOperationModeId: util.Ptr(model.HvacOperationModeIdType(0)),
				CurrentSetpointId:      util.Ptr(model.SetpointIdType(0)),
				IsOverrunActive:        util.Ptr(false),
			},
			{
				SystemFunctionId:       util.Ptr(model.HvacSystemFunctionIdType(1)),
				CurrentOperationModeId: util.Ptr(model.HvacOperationModeIdType(0)),
				CurrentSetpointId:      util.Ptr(model.SetpointIdType(0)),
				IsOverrunActive:        util.Ptr(false),
			},
		},
	}

	newData := model.HvacSystemFunctionListDataType{
		HvacSystemFunctionData: []model.HvacSystemFunctionDataType{
			{
				SystemFunctionId:       util.Ptr(model.HvacSystemFunctionIdType(1)),
				CurrentOperationModeId: util.Ptr(model.HvacOperationModeIdType(0)),
				CurrentSetpointId:      util.Ptr(model.SetpointIdType(0)),
				IsOverrunActive:        util.Ptr(true),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.HvacSystemFunctionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SystemFunctionId))
	assert.Equal(t, false, *item1.IsOverrunActive)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SystemFunctionId))
	assert.Equal(t, true, *item2.IsOverrunActive)
}

func TestHvacSystemFunctionPowerSequenceRelationListDataType_Update(t *testing.T) {
	sut := model.HvacSystemFunctionPowerSequenceRelationListDataType{
		HvacSystemFunctionPowerSequenceRelationData: []model.HvacSystemFunctionPowerSequenceRelationDataType{
			{
				SystemFunctionId: util.Ptr(model.HvacSystemFunctionIdType(0)),
				SequenceId:       []model.PowerSequenceIdType{0},
			},
			{
				SystemFunctionId: util.Ptr(model.HvacSystemFunctionIdType(1)),
				SequenceId:       []model.PowerSequenceIdType{0},
			},
		},
	}

	newData := model.HvacSystemFunctionPowerSequenceRelationListDataType{
		HvacSystemFunctionPowerSequenceRelationData: []model.HvacSystemFunctionPowerSequenceRelationDataType{
			{
				SystemFunctionId: util.Ptr(model.HvacSystemFunctionIdType(1)),
				SequenceId:       []model.PowerSequenceIdType{1},
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.HvacSystemFunctionPowerSequenceRelationData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SystemFunctionId))
	assert.Equal(t, 0, int(item1.SequenceId[0]))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SystemFunctionId))
	assert.Equal(t, 1, int(item2.SequenceId[0]))
}

func TestHvacSystemFunctionDescriptionListDataType_Update(t *testing.T) {
	sut := model.HvacSystemFunctionDescriptionListDataType{
		HvacSystemFunctionDescriptionData: []model.HvacSystemFunctionDescriptionDataType{
			{
				SystemFunctionId: util.Ptr(model.HvacSystemFunctionIdType(0)),
				Description:      util.Ptr(model.DescriptionType("old")),
			},
			{
				SystemFunctionId: util.Ptr(model.HvacSystemFunctionIdType(1)),
				Description:      util.Ptr(model.DescriptionType("old")),
			},
		},
	}

	newData := model.HvacSystemFunctionDescriptionListDataType{
		HvacSystemFunctionDescriptionData: []model.HvacSystemFunctionDescriptionDataType{
			{
				SystemFunctionId: util.Ptr(model.HvacSystemFunctionIdType(1)),
				Description:      util.Ptr(model.DescriptionType("new")),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.HvacSystemFunctionDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SystemFunctionId))
	assert.Equal(t, "old", string(*item1.Description))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SystemFunctionId))
	assert.Equal(t, "new", string(*item2.Description))
}

func TestHvacOperationModeDescriptionListDataType_Update(t *testing.T) {
	sut := model.HvacOperationModeDescriptionListDataType{
		HvacOperationModeDescriptionData: []model.HvacOperationModeDescriptionDataType{
			{
				OperationModeId: util.Ptr(model.HvacOperationModeIdType(0)),
				Description:     util.Ptr(model.DescriptionType("old")),
			},
			{
				OperationModeId: util.Ptr(model.HvacOperationModeIdType(1)),
				Description:     util.Ptr(model.DescriptionType("old")),
			},
		},
	}

	newData := model.HvacOperationModeDescriptionListDataType{
		HvacOperationModeDescriptionData: []model.HvacOperationModeDescriptionDataType{
			{
				OperationModeId: util.Ptr(model.HvacOperationModeIdType(1)),
				Description:     util.Ptr(model.DescriptionType("new")),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.HvacOperationModeDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.OperationModeId))
	assert.Equal(t, "old", string(*item1.Description))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.OperationModeId))
	assert.Equal(t, "new", string(*item2.Description))
}

func TestHvacOverrunListDataType_Update(t *testing.T) {
	sut := model.HvacOverrunListDataType{
		HvacOverrunData: []model.HvacOverrunDataType{
			{
				OverrunId:                 util.Ptr(model.HvacOverrunIdType(0)),
				TimeTableId:               util.Ptr(model.TimeTableIdType(0)),
				IsOverrunStatusChangeable: util.Ptr(false),
			},
			{
				OverrunId:                 util.Ptr(model.HvacOverrunIdType(1)),
				TimeTableId:               util.Ptr(model.TimeTableIdType(0)),
				IsOverrunStatusChangeable: util.Ptr(false),
			},
		},
	}

	newData := model.HvacOverrunListDataType{
		HvacOverrunData: []model.HvacOverrunDataType{
			{
				OverrunId:                 util.Ptr(model.HvacOverrunIdType(1)),
				TimeTableId:               util.Ptr(model.TimeTableIdType(0)),
				IsOverrunStatusChangeable: util.Ptr(true),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.HvacOverrunData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.OverrunId))
	assert.Equal(t, false, *item1.IsOverrunStatusChangeable)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.OverrunId))
	assert.Equal(t, true, *item2.IsOverrunStatusChangeable)
}

func TestHvacOverrunDescriptionListDataType_Update(t *testing.T) {
	sut := model.HvacOverrunDescriptionListDataType{
		HvacOverrunDescriptionData: []model.HvacOverrunDescriptionDataType{
			{
				OverrunId:   util.Ptr(model.HvacOverrunIdType(0)),
				Description: util.Ptr(model.DescriptionType("old")),
			},
			{
				OverrunId:   util.Ptr(model.HvacOverrunIdType(1)),
				Description: util.Ptr(model.DescriptionType("old")),
			},
		},
	}

	newData := model.HvacOverrunDescriptionListDataType{
		HvacOverrunDescriptionData: []model.HvacOverrunDescriptionDataType{
			{
				OverrunId:   util.Ptr(model.HvacOverrunIdType(1)),
				Description: util.Ptr(model.DescriptionType("new")),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.HvacOverrunDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.OverrunId))
	assert.Equal(t, "old", string(*item1.Description))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.OverrunId))
	assert.Equal(t, "new", string(*item2.Description))
}
