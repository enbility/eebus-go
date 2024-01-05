package model

import (
	"testing"

	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func TestHvacSystemFunctionListDataType_Update(t *testing.T) {
	sut := HvacSystemFunctionListDataType{
		HvacSystemFunctionData: []HvacSystemFunctionDataType{
			{
				SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(0)),
				IsOverrunActive:  util.Ptr(false),
			},
			{
				SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(1)),
				IsOverrunActive:  util.Ptr(false),
			},
		},
	}

	newData := HvacSystemFunctionListDataType{
		HvacSystemFunctionData: []HvacSystemFunctionDataType{
			{
				SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(1)),
				IsOverrunActive:  util.Ptr(true),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

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

func TestHvacSystemFunctionOperationModeRelationListDataType_Update(t *testing.T) {
	sut := HvacSystemFunctionOperationModeRelationListDataType{
		HvacSystemFunctionOperationModeRelationData: []HvacSystemFunctionOperationModeRelationDataType{
			{
				SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(0)),
				OperationModeId:  util.Ptr(HvacOperationModeIdType(0)),
			},
			{
				SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(1)),
				OperationModeId:  util.Ptr(HvacOperationModeIdType(0)),
			},
		},
	}

	newData := HvacSystemFunctionOperationModeRelationListDataType{
		HvacSystemFunctionOperationModeRelationData: []HvacSystemFunctionOperationModeRelationDataType{
			{
				SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(1)),
				OperationModeId:  util.Ptr(HvacOperationModeIdType(1)),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

	data := sut.HvacSystemFunctionOperationModeRelationData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SystemFunctionId))
	assert.Equal(t, 0, int(*item1.OperationModeId))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SystemFunctionId))
	assert.Equal(t, 1, int(*item2.OperationModeId))
}

func TestHvacSystemFunctionSetpointRelationListDataType_Update(t *testing.T) {
	sut := HvacSystemFunctionSetpointRelationListDataType{
		HvacSystemFunctionSetpointRelationData: []HvacSystemFunctionSetpointRelationDataType{
			{
				SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(0)),
				OperationModeId:  util.Ptr(HvacOperationModeIdType(0)),
			},
			{
				SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(1)),
				OperationModeId:  util.Ptr(HvacOperationModeIdType(0)),
			},
		},
	}

	newData := HvacSystemFunctionSetpointRelationListDataType{
		HvacSystemFunctionSetpointRelationData: []HvacSystemFunctionSetpointRelationDataType{
			{
				SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(1)),
				OperationModeId:  util.Ptr(HvacOperationModeIdType(1)),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

	data := sut.HvacSystemFunctionSetpointRelationData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SystemFunctionId))
	assert.Equal(t, 0, int(*item1.OperationModeId))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SystemFunctionId))
	assert.Equal(t, 1, int(*item2.OperationModeId))
}

func TestHvacSystemFunctionPowerSequenceRelationListDataType_Update(t *testing.T) {
	sut := HvacSystemFunctionPowerSequenceRelationListDataType{
		HvacSystemFunctionPowerSequenceRelationData: []HvacSystemFunctionPowerSequenceRelationDataType{
			{
				SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(0)),
				SequenceId:       []PowerSequenceIdType{0},
			},
			{
				SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(1)),
				SequenceId:       []PowerSequenceIdType{0},
			},
		},
	}

	newData := HvacSystemFunctionPowerSequenceRelationListDataType{
		HvacSystemFunctionPowerSequenceRelationData: []HvacSystemFunctionPowerSequenceRelationDataType{
			{
				SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(1)),
				SequenceId:       []PowerSequenceIdType{1},
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

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
	sut := HvacSystemFunctionDescriptionListDataType{
		HvacSystemFunctionDescriptionData: []HvacSystemFunctionDescriptionDataType{
			{
				SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(0)),
				Description:      util.Ptr(DescriptionType("old")),
			},
			{
				SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(1)),
				Description:      util.Ptr(DescriptionType("old")),
			},
		},
	}

	newData := HvacSystemFunctionDescriptionListDataType{
		HvacSystemFunctionDescriptionData: []HvacSystemFunctionDescriptionDataType{
			{
				SystemFunctionId: util.Ptr(HvacSystemFunctionIdType(1)),
				Description:      util.Ptr(DescriptionType("new")),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

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
	sut := HvacOperationModeDescriptionListDataType{
		HvacOperationModeDescriptionData: []HvacOperationModeDescriptionDataType{
			{
				OperationModeId: util.Ptr(HvacOperationModeIdType(0)),
				Description:     util.Ptr(DescriptionType("old")),
			},
			{
				OperationModeId: util.Ptr(HvacOperationModeIdType(1)),
				Description:     util.Ptr(DescriptionType("old")),
			},
		},
	}

	newData := HvacOperationModeDescriptionListDataType{
		HvacOperationModeDescriptionData: []HvacOperationModeDescriptionDataType{
			{
				OperationModeId: util.Ptr(HvacOperationModeIdType(1)),
				Description:     util.Ptr(DescriptionType("new")),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

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
	sut := HvacOverrunListDataType{
		HvacOverrunData: []HvacOverrunDataType{
			{
				OverrunId:                 util.Ptr(HvacOverrunIdType(0)),
				IsOverrunStatusChangeable: util.Ptr(false),
			},
			{
				OverrunId:                 util.Ptr(HvacOverrunIdType(1)),
				IsOverrunStatusChangeable: util.Ptr(false),
			},
		},
	}

	newData := HvacOverrunListDataType{
		HvacOverrunData: []HvacOverrunDataType{
			{
				OverrunId:                 util.Ptr(HvacOverrunIdType(1)),
				IsOverrunStatusChangeable: util.Ptr(true),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

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
	sut := HvacOverrunDescriptionListDataType{
		HvacOverrunDescriptionData: []HvacOverrunDescriptionDataType{
			{
				OverrunId:   util.Ptr(HvacOverrunIdType(0)),
				Description: util.Ptr(DescriptionType("old")),
			},
			{
				OverrunId:   util.Ptr(HvacOverrunIdType(1)),
				Description: util.Ptr(DescriptionType("old")),
			},
		},
	}

	newData := HvacOverrunDescriptionListDataType{
		HvacOverrunDescriptionData: []HvacOverrunDescriptionDataType{
			{
				OverrunId:   util.Ptr(HvacOverrunIdType(1)),
				Description: util.Ptr(DescriptionType("new")),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

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
