package model

import (
	"testing"

	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func TestTariffListDataType_Update(t *testing.T) {
	sut := TariffListDataType{
		TariffData: []TariffDataType{
			{
				TariffId:     util.Ptr(TariffIdType(0)),
				ActiveTierId: []TierIdType{0},
			},
			{
				TariffId:     util.Ptr(TariffIdType(1)),
				ActiveTierId: []TierIdType{0},
			},
		},
	}

	newData := TariffListDataType{
		TariffData: []TariffDataType{
			{
				TariffId:     util.Ptr(TariffIdType(1)),
				ActiveTierId: []TierIdType{1},
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

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
	sut := TariffTierRelationListDataType{
		TariffTierRelationData: []TariffTierRelationDataType{
			{
				TariffId: util.Ptr(TariffIdType(0)),
				TierId:   []TierIdType{0},
			},
			{
				TariffId: util.Ptr(TariffIdType(1)),
				TierId:   []TierIdType{0},
			},
		},
	}

	newData := TariffTierRelationListDataType{
		TariffTierRelationData: []TariffTierRelationDataType{
			{
				TariffId: util.Ptr(TariffIdType(1)),
				TierId:   []TierIdType{1},
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

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
	sut := TariffBoundaryRelationListDataType{
		TariffBoundaryRelationData: []TariffBoundaryRelationDataType{
			{
				TariffId:   util.Ptr(TariffIdType(0)),
				BoundaryId: []TierBoundaryIdType{0},
			},
			{
				TariffId:   util.Ptr(TariffIdType(1)),
				BoundaryId: []TierBoundaryIdType{0},
			},
		},
	}

	newData := TariffBoundaryRelationListDataType{
		TariffBoundaryRelationData: []TariffBoundaryRelationDataType{
			{
				TariffId:   util.Ptr(TariffIdType(1)),
				BoundaryId: []TierBoundaryIdType{1},
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

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
	sut := TariffDescriptionListDataType{
		TariffDescriptionData: []TariffDescriptionDataType{
			{
				TariffId:    util.Ptr(TariffIdType(0)),
				Description: util.Ptr(DescriptionType("old")),
			},
			{
				TariffId:    util.Ptr(TariffIdType(1)),
				Description: util.Ptr(DescriptionType("old")),
			},
		},
	}

	newData := TariffDescriptionListDataType{
		TariffDescriptionData: []TariffDescriptionDataType{
			{
				TariffId:    util.Ptr(TariffIdType(1)),
				Description: util.Ptr(DescriptionType("new")),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

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
	sut := TierBoundaryListDataType{
		TierBoundaryData: []TierBoundaryDataType{
			{
				BoundaryId:         util.Ptr(TierBoundaryIdType(0)),
				LowerBoundaryValue: NewScaledNumberType(1),
			},
			{
				BoundaryId:         util.Ptr(TierBoundaryIdType(1)),
				LowerBoundaryValue: NewScaledNumberType(1),
			},
		},
	}

	newData := TierBoundaryListDataType{
		TierBoundaryData: []TierBoundaryDataType{
			{
				BoundaryId:         util.Ptr(TierBoundaryIdType(1)),
				LowerBoundaryValue: NewScaledNumberType(10),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

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
	sut := TierBoundaryDescriptionListDataType{
		TierBoundaryDescriptionData: []TierBoundaryDescriptionDataType{
			{
				BoundaryId:  util.Ptr(TierBoundaryIdType(0)),
				Description: util.Ptr(DescriptionType("old")),
			},
			{
				BoundaryId:  util.Ptr(TierBoundaryIdType(1)),
				Description: util.Ptr(DescriptionType("old")),
			},
		},
	}

	newData := TierBoundaryDescriptionListDataType{
		TierBoundaryDescriptionData: []TierBoundaryDescriptionDataType{
			{
				BoundaryId:  util.Ptr(TierBoundaryIdType(1)),
				Description: util.Ptr(DescriptionType("new")),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

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
	sut := CommodityListDataType{
		CommodityData: []CommodityDataType{
			{
				CommodityId: util.Ptr(CommodityIdType(0)),
				Description: util.Ptr(DescriptionType("old")),
			},
			{
				CommodityId: util.Ptr(CommodityIdType(1)),
				Description: util.Ptr(DescriptionType("old")),
			},
		},
	}

	newData := CommodityListDataType{
		CommodityData: []CommodityDataType{
			{
				CommodityId: util.Ptr(CommodityIdType(1)),
				Description: util.Ptr(DescriptionType("new")),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

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
	sut := TierListDataType{
		TierData: []TierDataType{
			{
				TierId:            util.Ptr(TierIdType(0)),
				ActiveIncentiveId: []IncentiveIdType{0},
			},
			{
				TierId:            util.Ptr(TierIdType(1)),
				ActiveIncentiveId: []IncentiveIdType{0},
			},
		},
	}

	newData := TierListDataType{
		TierData: []TierDataType{
			{
				TierId:            util.Ptr(TierIdType(1)),
				ActiveIncentiveId: []IncentiveIdType{1},
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

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
	sut := TierIncentiveRelationListDataType{
		TierIncentiveRelationData: []TierIncentiveRelationDataType{
			{
				TierId:      util.Ptr(TierIdType(0)),
				IncentiveId: []IncentiveIdType{0},
			},
			{
				TierId:      util.Ptr(TierIdType(1)),
				IncentiveId: []IncentiveIdType{0},
			},
		},
	}

	newData := TierIncentiveRelationListDataType{
		TierIncentiveRelationData: []TierIncentiveRelationDataType{
			{
				TierId:      util.Ptr(TierIdType(1)),
				IncentiveId: []IncentiveIdType{1},
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

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
	sut := TierDescriptionListDataType{
		TierDescriptionData: []TierDescriptionDataType{
			{
				TierId:      util.Ptr(TierIdType(0)),
				Description: util.Ptr(DescriptionType("old")),
			},
			{
				TierId:      util.Ptr(TierIdType(1)),
				Description: util.Ptr(DescriptionType("old")),
			},
		},
	}

	newData := TierDescriptionListDataType{
		TierDescriptionData: []TierDescriptionDataType{
			{
				TierId:      util.Ptr(TierIdType(1)),
				Description: util.Ptr(DescriptionType("new")),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

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
	sut := IncentiveListDataType{
		IncentiveData: []IncentiveDataType{
			{
				IncentiveId: util.Ptr(IncentiveIdType(0)),
				Value:       NewScaledNumberType(1),
			},
			{
				IncentiveId: util.Ptr(IncentiveIdType(1)),
				Value:       NewScaledNumberType(1),
			},
		},
	}

	newData := IncentiveListDataType{
		IncentiveData: []IncentiveDataType{
			{
				IncentiveId: util.Ptr(IncentiveIdType(1)),
				Value:       NewScaledNumberType(10),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

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
	sut := IncentiveDescriptionListDataType{
		IncentiveDescriptionData: []IncentiveDescriptionDataType{
			{
				IncentiveId: util.Ptr(IncentiveIdType(0)),
				Description: util.Ptr(DescriptionType("old")),
			},
			{
				IncentiveId: util.Ptr(IncentiveIdType(1)),
				Description: util.Ptr(DescriptionType("old")),
			},
		},
	}

	newData := IncentiveDescriptionListDataType{
		IncentiveDescriptionData: []IncentiveDescriptionDataType{
			{
				IncentiveId: util.Ptr(IncentiveIdType(1)),
				Description: util.Ptr(DescriptionType("new")),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

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
