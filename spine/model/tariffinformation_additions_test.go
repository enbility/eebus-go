package model_test

import (
	"testing"

	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func TestTariffListDataType_Update(t *testing.T) {
	sut := model.TariffListDataType{
		TariffData: []model.TariffDataType{
			{
				TariffId:     util.Ptr(model.TariffIdType(0)),
				ActiveTierId: []model.TierIdType{0},
			},
			{
				TariffId:     util.Ptr(model.TariffIdType(1)),
				ActiveTierId: []model.TierIdType{0},
			},
		},
	}

	newData := model.TariffListDataType{
		TariffData: []model.TariffDataType{
			{
				TariffId:     util.Ptr(model.TariffIdType(1)),
				ActiveTierId: []model.TierIdType{1},
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.TariffData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.TariffId))
	assert.Equal(t, 0, int(item1.ActiveTierId[0]))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.TariffId))
	assert.Equal(t, 1, int(item2.ActiveTierId[0]))
}

func TestTariffTierRelationListDataType_Update(t *testing.T) {
	sut := model.TariffTierRelationListDataType{
		TariffTierRelationData: []model.TariffTierRelationDataType{
			{
				TariffId: util.Ptr(model.TariffIdType(0)),
				TierId:   []model.TierIdType{0},
			},
			{
				TariffId: util.Ptr(model.TariffIdType(1)),
				TierId:   []model.TierIdType{0},
			},
		},
	}

	newData := model.TariffTierRelationListDataType{
		TariffTierRelationData: []model.TariffTierRelationDataType{
			{
				TariffId: util.Ptr(model.TariffIdType(1)),
				TierId:   []model.TierIdType{1},
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.TariffTierRelationData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.TariffId))
	assert.Equal(t, 0, int(item1.TierId[0]))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.TariffId))
	assert.Equal(t, 1, int(item2.TierId[0]))
}

func TestTariffBoundaryRelationListDataType_Update(t *testing.T) {
	sut := model.TariffBoundaryRelationListDataType{
		TariffBoundaryRelationData: []model.TariffBoundaryRelationDataType{
			{
				TariffId:   util.Ptr(model.TariffIdType(0)),
				BoundaryId: []model.TierBoundaryIdType{0},
			},
			{
				TariffId:   util.Ptr(model.TariffIdType(1)),
				BoundaryId: []model.TierBoundaryIdType{0},
			},
		},
	}

	newData := model.TariffBoundaryRelationListDataType{
		TariffBoundaryRelationData: []model.TariffBoundaryRelationDataType{
			{
				TariffId:   util.Ptr(model.TariffIdType(1)),
				BoundaryId: []model.TierBoundaryIdType{1},
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.TariffBoundaryRelationData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.TariffId))
	assert.Equal(t, 0, int(item1.BoundaryId[0]))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.TariffId))
	assert.Equal(t, 1, int(item2.BoundaryId[0]))
}

func TestTariffDescriptionListDataType_Update(t *testing.T) {
	sut := model.TariffDescriptionListDataType{
		TariffDescriptionData: []model.TariffDescriptionDataType{
			{
				TariffId:      util.Ptr(model.TariffIdType(0)),
				CommodityId:   util.Ptr(model.CommodityIdType(0)),
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Description:   util.Ptr(model.DescriptionType("old")),
			},
			{
				TariffId:      util.Ptr(model.TariffIdType(1)),
				CommodityId:   util.Ptr(model.CommodityIdType(0)),
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Description:   util.Ptr(model.DescriptionType("old")),
			},
		},
	}

	newData := model.TariffDescriptionListDataType{
		TariffDescriptionData: []model.TariffDescriptionDataType{
			{
				TariffId:      util.Ptr(model.TariffIdType(1)),
				CommodityId:   util.Ptr(model.CommodityIdType(0)),
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				Description:   util.Ptr(model.DescriptionType("new")),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.TariffDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.TariffId))
	assert.Equal(t, "old", string(*item1.Description))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.TariffId))
	assert.Equal(t, "new", string(*item2.Description))
}

func TestTierBoundaryListDataType_Update(t *testing.T) {
	sut := model.TierBoundaryListDataType{
		TierBoundaryData: []model.TierBoundaryDataType{
			{
				BoundaryId:         util.Ptr(model.TierBoundaryIdType(0)),
				TimeTableId:        util.Ptr(model.TimeTableIdType(0)),
				LowerBoundaryValue: model.NewScaledNumberType(1),
			},
			{
				BoundaryId:         util.Ptr(model.TierBoundaryIdType(1)),
				TimeTableId:        util.Ptr(model.TimeTableIdType(0)),
				LowerBoundaryValue: model.NewScaledNumberType(1),
			},
		},
	}

	newData := model.TierBoundaryListDataType{
		TierBoundaryData: []model.TierBoundaryDataType{
			{
				BoundaryId:         util.Ptr(model.TierBoundaryIdType(1)),
				TimeTableId:        util.Ptr(model.TimeTableIdType(0)),
				LowerBoundaryValue: model.NewScaledNumberType(10),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.TierBoundaryData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.BoundaryId))
	assert.Equal(t, 1.0, item1.LowerBoundaryValue.GetValue())
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.BoundaryId))
	assert.Equal(t, 10.0, item2.LowerBoundaryValue.GetValue())
}

func TestTierBoundaryDescriptionListDataType_Update(t *testing.T) {
	sut := model.TierBoundaryDescriptionListDataType{
		TierBoundaryDescriptionData: []model.TierBoundaryDescriptionDataType{
			{
				BoundaryId:               util.Ptr(model.TierBoundaryIdType(0)),
				ValidForTierId:           util.Ptr(model.TierIdType(0)),
				SwitchToTierIdWhenLower:  util.Ptr(model.TierIdType(0)),
				SwitchToTierIdWhenHigher: util.Ptr(model.TierIdType(0)),
				Description:              util.Ptr(model.DescriptionType("old")),
			},
			{
				BoundaryId:               util.Ptr(model.TierBoundaryIdType(1)),
				ValidForTierId:           util.Ptr(model.TierIdType(0)),
				SwitchToTierIdWhenLower:  util.Ptr(model.TierIdType(0)),
				SwitchToTierIdWhenHigher: util.Ptr(model.TierIdType(0)),
				Description:              util.Ptr(model.DescriptionType("old")),
			},
		},
	}

	newData := model.TierBoundaryDescriptionListDataType{
		TierBoundaryDescriptionData: []model.TierBoundaryDescriptionDataType{
			{
				BoundaryId:               util.Ptr(model.TierBoundaryIdType(1)),
				ValidForTierId:           util.Ptr(model.TierIdType(0)),
				SwitchToTierIdWhenLower:  util.Ptr(model.TierIdType(0)),
				SwitchToTierIdWhenHigher: util.Ptr(model.TierIdType(0)),
				Description:              util.Ptr(model.DescriptionType("new")),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.TierBoundaryDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.BoundaryId))
	assert.Equal(t, "old", string(*item1.Description))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.BoundaryId))
	assert.Equal(t, "new", string(*item2.Description))
}

func TestCommodityListDataType_Update(t *testing.T) {
	sut := model.CommodityListDataType{
		CommodityData: []model.CommodityDataType{
			{
				CommodityId: util.Ptr(model.CommodityIdType(0)),
				Description: util.Ptr(model.DescriptionType("old")),
			},
			{
				CommodityId: util.Ptr(model.CommodityIdType(1)),
				Description: util.Ptr(model.DescriptionType("old")),
			},
		},
	}

	newData := model.CommodityListDataType{
		CommodityData: []model.CommodityDataType{
			{
				CommodityId: util.Ptr(model.CommodityIdType(1)),
				Description: util.Ptr(model.DescriptionType("new")),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.CommodityData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.CommodityId))
	assert.Equal(t, "old", string(*item1.Description))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.CommodityId))
	assert.Equal(t, "new", string(*item2.Description))
}

func TestTierListDataType_Update(t *testing.T) {
	sut := model.TierListDataType{
		TierData: []model.TierDataType{
			{
				TierId:            util.Ptr(model.TierIdType(0)),
				TimeTableId:       util.Ptr(model.TimeTableIdType(0)),
				ActiveIncentiveId: []model.IncentiveIdType{0},
			},
			{
				TierId:            util.Ptr(model.TierIdType(1)),
				TimeTableId:       util.Ptr(model.TimeTableIdType(0)),
				ActiveIncentiveId: []model.IncentiveIdType{0},
			},
		},
	}

	newData := model.TierListDataType{
		TierData: []model.TierDataType{
			{
				TierId:            util.Ptr(model.TierIdType(1)),
				TimeTableId:       util.Ptr(model.TimeTableIdType(0)),
				ActiveIncentiveId: []model.IncentiveIdType{1},
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.TierData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.TierId))
	assert.Equal(t, 0, int(item1.ActiveIncentiveId[0]))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.TierId))
	assert.Equal(t, 1, int(item2.ActiveIncentiveId[0]))
}

func TestTierIncentiveRelationListDataType_Update(t *testing.T) {
	sut := model.TierIncentiveRelationListDataType{
		TierIncentiveRelationData: []model.TierIncentiveRelationDataType{
			{
				TierId:      util.Ptr(model.TierIdType(0)),
				IncentiveId: []model.IncentiveIdType{0},
			},
			{
				TierId:      util.Ptr(model.TierIdType(1)),
				IncentiveId: []model.IncentiveIdType{0},
			},
		},
	}

	newData := model.TierIncentiveRelationListDataType{
		TierIncentiveRelationData: []model.TierIncentiveRelationDataType{
			{
				TierId:      util.Ptr(model.TierIdType(1)),
				IncentiveId: []model.IncentiveIdType{1},
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.TierIncentiveRelationData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.TierId))
	assert.Equal(t, 0, int(item1.IncentiveId[0]))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.TierId))
	assert.Equal(t, 1, int(item2.IncentiveId[0]))
}

func TestTierDescriptionListDataType_Update(t *testing.T) {
	sut := model.TierDescriptionListDataType{
		TierDescriptionData: []model.TierDescriptionDataType{
			{
				TierId:      util.Ptr(model.TierIdType(0)),
				Description: util.Ptr(model.DescriptionType("old")),
			},
			{
				TierId:      util.Ptr(model.TierIdType(1)),
				Description: util.Ptr(model.DescriptionType("old")),
			},
		},
	}

	newData := model.TierDescriptionListDataType{
		TierDescriptionData: []model.TierDescriptionDataType{
			{
				TierId:      util.Ptr(model.TierIdType(1)),
				Description: util.Ptr(model.DescriptionType("new")),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.TierDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.TierId))
	assert.Equal(t, "old", string(*item1.Description))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.TierId))
	assert.Equal(t, "new", string(*item2.Description))
}

func TestIncentiveListDataType_Update(t *testing.T) {
	sut := model.IncentiveListDataType{
		IncentiveData: []model.IncentiveDataType{
			{
				IncentiveId: util.Ptr(model.IncentiveIdType(0)),
				Value:       model.NewScaledNumberType(1),
			},
			{
				IncentiveId: util.Ptr(model.IncentiveIdType(1)),
				Value:       model.NewScaledNumberType(1),
			},
		},
	}

	newData := model.IncentiveListDataType{
		IncentiveData: []model.IncentiveDataType{
			{
				IncentiveId: util.Ptr(model.IncentiveIdType(1)),
				Value:       model.NewScaledNumberType(10),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.IncentiveData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.IncentiveId))
	assert.Equal(t, 1.0, item1.Value.GetValue())
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.IncentiveId))
	assert.Equal(t, 10.0, item2.Value.GetValue())
}

func TestIncentiveDescriptionListDataType_Update(t *testing.T) {
	sut := model.IncentiveDescriptionListDataType{
		IncentiveDescriptionData: []model.IncentiveDescriptionDataType{
			{
				IncentiveId: util.Ptr(model.IncentiveIdType(0)),
				Description: util.Ptr(model.DescriptionType("old")),
			},
			{
				IncentiveId: util.Ptr(model.IncentiveIdType(1)),
				Description: util.Ptr(model.DescriptionType("old")),
			},
		},
	}

	newData := model.IncentiveDescriptionListDataType{
		IncentiveDescriptionData: []model.IncentiveDescriptionDataType{
			{
				IncentiveId: util.Ptr(model.IncentiveIdType(1)),
				Description: util.Ptr(model.DescriptionType("new")),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.IncentiveDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.IncentiveId))
	assert.Equal(t, "old", string(*item1.Description))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.IncentiveId))
	assert.Equal(t, "new", string(*item2.Description))
}
