package model_test

import (
	"encoding/json"
	"testing"

	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func TestElectricalConnectionStateListDataType_Update(t *testing.T) {
	sut := model.ElectricalConnectionStateListDataType{
		ElectricalConnectionStateData: []model.ElectricalConnectionStateDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				CurrentEnergyMode:      util.Ptr(model.EnergyModeTypeProduce),
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(1)),
				CurrentEnergyMode:      util.Ptr(model.EnergyModeTypeProduce),
			},
		},
	}

	newData := model.ElectricalConnectionStateListDataType{
		ElectricalConnectionStateData: []model.ElectricalConnectionStateDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(1)),
				CurrentEnergyMode:      util.Ptr(model.EnergyModeTypeConsume),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.ElectricalConnectionStateData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.ElectricalConnectionId))
	assert.Equal(t, model.EnergyModeTypeProduce, *item1.CurrentEnergyMode)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.ElectricalConnectionId))
	assert.Equal(t, model.EnergyModeTypeConsume, *item2.CurrentEnergyMode)
}

// verifies that a subset of existing items will be updated with identified new values
func TestElectricalConnectionPermittedValueSetListDataType_Update_Modify(t *testing.T) {
	sut := model.ElectricalConnectionPermittedValueSetListDataType{
		ElectricalConnectionPermittedValueSetData: []model.ElectricalConnectionPermittedValueSetDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(0)),
				PermittedValueSet: []model.ScaledNumberSetType{
					{
						Range: []model.ScaledNumberRangeType{
							{
								Min: model.NewScaledNumberType(1),
							},
						},
					},
				},
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(1)),
				PermittedValueSet: []model.ScaledNumberSetType{
					{
						Range: []model.ScaledNumberRangeType{
							{
								Min: model.NewScaledNumberType(6),
								Max: model.NewScaledNumberType(16),
							},
						},
					},
				},
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(2)),
				PermittedValueSet: []model.ScaledNumberSetType{
					{
						Range: []model.ScaledNumberRangeType{
							{
								Min: model.NewScaledNumberType(6),
								Max: model.NewScaledNumberType(16),
							},
						},
					},
				},
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(3)),
				PermittedValueSet: []model.ScaledNumberSetType{
					{
						Range: []model.ScaledNumberRangeType{
							{
								Min: model.NewScaledNumberType(6),
								Max: model.NewScaledNumberType(16),
							},
						},
					},
				},
			},
		},
	}

	newData := model.ElectricalConnectionPermittedValueSetListDataType{
		ElectricalConnectionPermittedValueSetData: []model.ElectricalConnectionPermittedValueSetDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(1)),
				PermittedValueSet: []model.ScaledNumberSetType{
					{
						Range: []model.ScaledNumberRangeType{
							{
								Min: model.NewScaledNumberType(2),
								Max: model.NewScaledNumberType(16),
							},
						},
					},
				},
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(2)),
				PermittedValueSet: []model.ScaledNumberSetType{
					{
						Range: []model.ScaledNumberRangeType{
							{
								Min: model.NewScaledNumberType(2),
								Max: model.NewScaledNumberType(16),
							},
						},
					},
				},
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(3)),
				PermittedValueSet: []model.ScaledNumberSetType{
					{
						Range: []model.ScaledNumberRangeType{
							{
								Min: model.NewScaledNumberType(2),
								Max: model.NewScaledNumberType(16),
							},
						},
					},
				},
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.ElectricalConnectionPermittedValueSetData
	// check the non changing items
	assert.Equal(t, 4, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.ElectricalConnectionId))
	assert.Equal(t, 0, int(*item1.ParameterId))
	assert.Equal(t, 1, len(item1.PermittedValueSet))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 0, int(*item2.ElectricalConnectionId))
	assert.Equal(t, 1, int(*item2.ParameterId))
	assert.Equal(t, 1, len(item2.PermittedValueSet))
	valueSet := item2.PermittedValueSet[0]
	assert.Equal(t, 1, len(valueSet.Range))
	rangeSet := valueSet.Range[0]
	assert.Equal(t, 2.0, rangeSet.Min.GetValue())
	assert.Equal(t, 16.0, rangeSet.Max.GetValue())
}

// verifies that a subset of existing items will be updated with identified new values
func TestElectricalConnectionPermittedValueSetListDataType_Update_Modify_Selector(t *testing.T) {
	existingDataJson := `{
		"electricalConnectionPermittedValueSetData":[
			{
				"electricalConnectionId":0,
				"parameterId":0,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":1,"scale":0}
							}
						]
					}
				]
			},
			{
				"electricalConnectionId":0,
				"parameterId":1,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":6,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			},
			{
				"electricalConnectionId":0,
				"parameterId":2,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":6,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			},
			{
				"electricalConnectionId":0,
				"parameterId":3,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":6,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			}
		]
	}`

	var sut model.ElectricalConnectionPermittedValueSetListDataType
	err := json.Unmarshal([]byte(existingDataJson), &sut)
	if assert.Nil(t, err) == false {
		return
	}

	newDataJson := `{
		"electricalConnectionPermittedValueSetData":[
			{
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":2,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			}
		]
	}`

	var newData model.ElectricalConnectionPermittedValueSetListDataType
	err = json.Unmarshal([]byte(newDataJson), &newData)
	if assert.Nil(t, err) == false {
		return
	}

	partial := &model.FilterType{
		CmdControl: &model.CmdControlType{
			Partial: &model.ElementTagType{},
		},
		ElectricalConnectionPermittedValueSetListDataSelectors: &model.ElectricalConnectionPermittedValueSetListDataSelectorsType{
			ElectricalConnectionId: util.Ptr[model.ElectricalConnectionIdType](0),
			ParameterId:            util.Ptr[model.ElectricalConnectionParameterIdType](1),
		},
	}

	// Act
	sut.UpdateList(&newData, partial, nil)

	data := sut.ElectricalConnectionPermittedValueSetData
	// check the non changing items
	assert.Equal(t, 4, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.ElectricalConnectionId))
	assert.Equal(t, 0, int(*item1.ParameterId))
	assert.Equal(t, 1, len(item1.PermittedValueSet))
	item3 := data[2]
	assert.Equal(t, 0, int(*item3.ElectricalConnectionId))
	assert.Equal(t, 2, int(*item3.ParameterId))
	assert.Equal(t, 1, len(item3.PermittedValueSet))
	valueSet := item3.PermittedValueSet[0]
	assert.Equal(t, 1, len(valueSet.Range))
	rangeSet := valueSet.Range[0]
	assert.Equal(t, 6.0, rangeSet.Min.GetValue())
	assert.Equal(t, 16.0, rangeSet.Max.GetValue())

	// check properties of updated item
	item2 := sut.ElectricalConnectionPermittedValueSetData[1]
	assert.Equal(t, 0, int(*item2.ElectricalConnectionId))
	assert.Equal(t, 1, int(*item2.ParameterId))
	assert.Equal(t, 1, len(item2.PermittedValueSet))
	valueSet = item2.PermittedValueSet[0]
	assert.Equal(t, 1, len(valueSet.Range))
	rangeSet = valueSet.Range[0]
	assert.Equal(t, 2.0, rangeSet.Min.GetValue())
	assert.Equal(t, 16.0, rangeSet.Max.GetValue())
}

// verifies that a subset of existing items will be updated with identified new values
func TestElectricalConnectionPermittedValueSetListDataType_Update_Delete_Modify(t *testing.T) {
	existingDataJson := `{
		"electricalConnectionPermittedValueSetData":[
			{
				"electricalConnectionId":0,
				"parameterId":0,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":1,"scale":0}
							}
						]
					}
				]
			},
			{
				"electricalConnectionId":0,
				"parameterId":1,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":6,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			},
			{
				"electricalConnectionId":0,
				"parameterId":2,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":6,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			},
			{
				"electricalConnectionId":0,
				"parameterId":3,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":6,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			}
		]
	}`

	var sut model.ElectricalConnectionPermittedValueSetListDataType
	err := json.Unmarshal([]byte(existingDataJson), &sut)
	if assert.Nil(t, err) == false {
		return
	}

	newDataJson := `{
		"electricalConnectionPermittedValueSetData":[
			{
				"electricalConnectionId":0,
				"parameterId":1,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":2,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			},
			{
				"electricalConnectionId":0,
				"parameterId":2,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":2,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			},
			{
				"electricalConnectionId":0,
				"parameterId":3,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":2,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			}
		]
	}`

	var newData model.ElectricalConnectionPermittedValueSetListDataType
	err = json.Unmarshal([]byte(newDataJson), &newData)
	if assert.Nil(t, err) == false {
		return
	}

	delete := &model.FilterType{
		CmdControl: &model.CmdControlType{
			Delete: &model.ElementTagType{},
		},
		ElectricalConnectionPermittedValueSetListDataSelectors: &model.ElectricalConnectionPermittedValueSetListDataSelectorsType{
			ElectricalConnectionId: util.Ptr[model.ElectricalConnectionIdType](0),
			ParameterId:            util.Ptr[model.ElectricalConnectionParameterIdType](0),
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), delete)

	data := sut.ElectricalConnectionPermittedValueSetData
	// check the deleted item is gone
	assert.Equal(t, 3, len(data))
	// check properties of updated item
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.ElectricalConnectionId))
	assert.Equal(t, 1, int(*item1.ParameterId))
	assert.Equal(t, 1, len(item1.PermittedValueSet))
	valueSet := item1.PermittedValueSet[0]
	assert.Equal(t, 1, len(valueSet.Range))
	rangeSet := valueSet.Range[0]
	assert.Equal(t, 2.0, rangeSet.Min.GetValue())
	assert.Equal(t, 16.0, rangeSet.Max.GetValue())
}

// verifies that a subset of existing items will be updated with identified new values
func TestElectricalConnectionPermittedValueSetListDataType_Update_Delete(t *testing.T) {
	existingDataJson := `{
		"electricalConnectionPermittedValueSetData":[
			{
				"electricalConnectionId":0,
				"parameterId":0,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":1,"scale":0}
							}
						]
					}
				]
			},
			{
				"electricalConnectionId":0,
				"parameterId":1,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":6,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			},
			{
				"electricalConnectionId":0,
				"parameterId":2,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":6,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			},
			{
				"electricalConnectionId":0,
				"parameterId":3,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":6,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			}
		]
	}`

	var sut model.ElectricalConnectionPermittedValueSetListDataType
	err := json.Unmarshal([]byte(existingDataJson), &sut)
	if assert.Nil(t, err) == false {
		return
	}

	delete := &model.FilterType{
		CmdControl: &model.CmdControlType{
			Delete: &model.ElementTagType{},
		},
		ElectricalConnectionPermittedValueSetListDataSelectors: &model.ElectricalConnectionPermittedValueSetListDataSelectorsType{
			ElectricalConnectionId: util.Ptr[model.ElectricalConnectionIdType](0),
			ParameterId:            util.Ptr[model.ElectricalConnectionParameterIdType](0),
		},
	}

	// Act
	sut.UpdateList(nil, nil, delete)

	data := sut.ElectricalConnectionPermittedValueSetData
	// check the deleted item is added again
	assert.Equal(t, 3, len(data))
	// check properties of remaining item
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.ElectricalConnectionId))
	assert.Equal(t, 1, int(*item1.ParameterId))
	assert.Equal(t, 1, len(item1.PermittedValueSet))
	valueSet := item1.PermittedValueSet[0]
	assert.Equal(t, 1, len(valueSet.Range))
	rangeSet := valueSet.Range[0]
	assert.Equal(t, 6.0, rangeSet.Min.GetValue())
	assert.Equal(t, 16.0, rangeSet.Max.GetValue())
}

// verifies that a subset of existing items will be updated with identified new values
func TestElectricalConnectionPermittedValueSetListDataType_Update_Delete_Element(t *testing.T) {
	existingDataJson := `{
		"electricalConnectionPermittedValueSetData":[
			{
				"electricalConnectionId":0,
				"parameterId":0,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":1,"scale":0}
							}
						]
					}
				]
			},
			{
				"electricalConnectionId":0,
				"parameterId":1,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":6,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			},
			{
				"electricalConnectionId":0,
				"parameterId":2,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":6,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			},
			{
				"electricalConnectionId":0,
				"parameterId":3,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":6,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			}
		]
	}`

	var sut model.ElectricalConnectionPermittedValueSetListDataType
	err := json.Unmarshal([]byte(existingDataJson), &sut)
	if assert.Nil(t, err) == false {
		return
	}

	delete := &model.FilterType{
		CmdControl: &model.CmdControlType{
			Delete: &model.ElementTagType{},
		},
		ElectricalConnectionPermittedValueSetDataElements: &model.ElectricalConnectionPermittedValueSetDataElementsType{
			PermittedValueSet: &model.ElementTagType{},
		},
		ElectricalConnectionPermittedValueSetListDataSelectors: &model.ElectricalConnectionPermittedValueSetListDataSelectorsType{
			ElectricalConnectionId: util.Ptr[model.ElectricalConnectionIdType](0),
			ParameterId:            util.Ptr[model.ElectricalConnectionParameterIdType](0),
		},
	}

	// Act
	sut.UpdateList(nil, nil, delete)

	data := sut.ElectricalConnectionPermittedValueSetData
	// check no items are deleted
	assert.Equal(t, 4, len(data))
	// check permitted value is removed from item with ID 0
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.ElectricalConnectionId))
	assert.Equal(t, 0, int(*item1.ParameterId))
	var nilValue []model.ScaledNumberSetType
	assert.Equal(t, nilValue, item1.PermittedValueSet)

	// check properties of remaining item
	item2 := data[1]
	assert.Equal(t, 0, int(*item2.ElectricalConnectionId))
	assert.Equal(t, 1, int(*item2.ParameterId))
	assert.Equal(t, 1, len(item2.PermittedValueSet))
	valueSet := item2.PermittedValueSet[0]
	assert.Equal(t, 1, len(valueSet.Range))
	rangeSet := valueSet.Range[0]
	assert.Equal(t, 6.0, rangeSet.Min.GetValue())
	assert.Equal(t, 16.0, rangeSet.Max.GetValue())
}

// verifies that a subset of existing items will be updated with identified new values
func TestElectricalConnectionPermittedValueSetListDataType_Update_Delete_OnlyElement(t *testing.T) {
	existingDataJson := `{
		"electricalConnectionPermittedValueSetData":[
			{
				"electricalConnectionId":0,
				"parameterId":0,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":1,"scale":0}
							}
						]
					}
				]
			},
			{
				"electricalConnectionId":0,
				"parameterId":1,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":6,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			},
			{
				"electricalConnectionId":0,
				"parameterId":2,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":6,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			},
			{
				"electricalConnectionId":0,
				"parameterId":3,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":6,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			}
		]
	}`

	var sut model.ElectricalConnectionPermittedValueSetListDataType
	err := json.Unmarshal([]byte(existingDataJson), &sut)
	if assert.Nil(t, err) == false {
		return
	}

	delete := &model.FilterType{
		CmdControl: &model.CmdControlType{
			Delete: &model.ElementTagType{},
		},
		ElectricalConnectionPermittedValueSetDataElements: &model.ElectricalConnectionPermittedValueSetDataElementsType{
			PermittedValueSet: &model.ElementTagType{},
		},
	}

	// Act
	sut.UpdateList(nil, nil, delete)

	data := sut.ElectricalConnectionPermittedValueSetData
	// check no items are deleted
	assert.Equal(t, 4, len(data))
	// check permitted value is removed from item with ID 0
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.ElectricalConnectionId))
	assert.Equal(t, 0, int(*item1.ParameterId))
	var nilValue []model.ScaledNumberSetType
	assert.Equal(t, nilValue, item1.PermittedValueSet)

	// check properties
	item2 := data[1]
	assert.Equal(t, 0, int(*item2.ElectricalConnectionId))
	assert.Equal(t, 1, int(*item2.ParameterId))
	assert.Equal(t, nilValue, item2.PermittedValueSet)

	item3 := data[2]
	assert.Equal(t, 0, int(*item3.ElectricalConnectionId))
	assert.Equal(t, 2, int(*item3.ParameterId))
	assert.Equal(t, nilValue, item3.PermittedValueSet)

	item4 := data[3]
	assert.Equal(t, 0, int(*item4.ElectricalConnectionId))
	assert.Equal(t, 3, int(*item4.ParameterId))
	assert.Equal(t, nilValue, item4.PermittedValueSet)
}

// verifies that a subset of existing items will be updated with identified new values
func TestElectricalConnectionPermittedValueSetListDataType_Update_Delete_Add(t *testing.T) {
	existingDataJson := `{
		"electricalConnectionPermittedValueSetData":[
			{
				"electricalConnectionId":0,
				"parameterId":0,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":1,"scale":0}
							}
						]
					}
				]
			},
			{
				"electricalConnectionId":0,
				"parameterId":1,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":6,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			},
			{
				"electricalConnectionId":0,
				"parameterId":2,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":6,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			},
			{
				"electricalConnectionId":0,
				"parameterId":3,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":6,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			}
		]
	}`

	var sut model.ElectricalConnectionPermittedValueSetListDataType
	err := json.Unmarshal([]byte(existingDataJson), &sut)
	if assert.Nil(t, err) == false {
		return
	}

	newDataJson := `{
		"electricalConnectionPermittedValueSetData":[
			{
				"electricalConnectionId":0,
				"parameterId":0,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":1,"scale":0}
							}
						]
					}
				]
			},
			{
				"electricalConnectionId":0,
				"parameterId":1,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":2,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			},
			{
				"electricalConnectionId":0,
				"parameterId":2,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":2,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			},
			{
				"electricalConnectionId":0,
				"parameterId":3,
				"permittedValueSet":[
					{
						"range":[
							{
								"min":{"number":2,"scale":0},
								"max":{"number":16,"scale":0}
							}
						]
					}
				]
			}
		]
	}`

	var newData model.ElectricalConnectionPermittedValueSetListDataType
	err = json.Unmarshal([]byte(newDataJson), &newData)
	if assert.Nil(t, err) == false {
		return
	}

	delete := &model.FilterType{
		CmdControl: &model.CmdControlType{
			Delete: &model.ElementTagType{},
		},
		ElectricalConnectionPermittedValueSetListDataSelectors: &model.ElectricalConnectionPermittedValueSetListDataSelectorsType{
			ElectricalConnectionId: util.Ptr[model.ElectricalConnectionIdType](0),
			ParameterId:            util.Ptr[model.ElectricalConnectionParameterIdType](0),
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), delete)

	data := sut.ElectricalConnectionPermittedValueSetData
	// check the deleted item is added again
	assert.Equal(t, 4, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.ElectricalConnectionId))
	assert.Equal(t, 0, int(*item1.ParameterId))
	assert.Equal(t, 1, len(item1.PermittedValueSet))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 0, int(*item2.ElectricalConnectionId))
	assert.Equal(t, 1, int(*item2.ParameterId))
	assert.Equal(t, 1, len(item2.PermittedValueSet))
	valueSet := item2.PermittedValueSet[0]
	assert.Equal(t, 1, len(valueSet.Range))
	rangeSet := valueSet.Range[0]
	assert.Equal(t, 2.0, rangeSet.Min.GetValue())
	assert.Equal(t, 16.0, rangeSet.Max.GetValue())
}

// verifies that an item in the payload which is not in the existing data will be added
func TestElectricalConnectionPermittedValueSetListDataType_Update_NewItem(t *testing.T) {
	existingDataJson := `{
		"electricalConnectionPermittedValueSetData": [
		  {
			"electricalConnectionId": 1,
			"parameterId": 1,
			"permittedValueSet": [
			  {
				"range": [
				  {
					"min": { "number": 3, "scale": 0 },
					"max": { "number": 6, "scale": 0 }
				  }
				]
			  }
			]
		  }
		]
	}`

	var sut model.ElectricalConnectionPermittedValueSetListDataType
	err := json.Unmarshal([]byte(existingDataJson), &sut)
	if assert.Nil(t, err) == false {
		return
	}

	newDataJson := `{
		"electricalConnectionPermittedValueSetData": [
		  {
			"electricalConnectionId": 1,
			"parameterId": 2,
			"permittedValueSet": [
			  {
				"range": [
				  {
					"min": { "number": 9, "scale": 0 },
					"max": { "number": 19, "scale": 0 }
				  }
				]
			  },
			  {
				"range": [
				  {
					"min": { "number": 30, "scale": 0 },
					"max": { "number": 36, "scale": 0 }
				  }
				]
			  }
			]
		  }
		]
	}`

	var newData model.ElectricalConnectionPermittedValueSetListDataType
	err = json.Unmarshal([]byte(newDataJson), &newData)
	if assert.Nil(t, err) == false {
		return
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.ElectricalConnectionPermittedValueSetData
	// new item should be added
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 1, int(*item1.ElectricalConnectionId))
	assert.Equal(t, 1, int(*item1.ParameterId))
	assert.Equal(t, 1, len(item1.PermittedValueSet))
	// check properties of added item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.ElectricalConnectionId))
	assert.Equal(t, 2, int(*item2.ParameterId))
	assert.Equal(t, 2, len(item2.PermittedValueSet))
}

// verifies that an item in the payload which has no identifiers will be copied to all existing data
// (see EEBus_SPINE_TS_ProtocolSpecification.pdf, Table 7: Considered cmdOptions combinations for classifier "notify")
func TestElectricalConnectionPermittedValueSetListDataType_UpdateWithoutIdenifiers(t *testing.T) {
	existingDataJson := `{
		"electricalConnectionPermittedValueSetData": [
		  {
			"electricalConnectionId": 1,
			"parameterId": 1,
			"permittedValueSet": [
			  {
				"range": [
				  {
					"min": { "number": 3, "scale": 0 },
					"max": { "number": 6, "scale": 0 }
				  }
				]
			  }
			]
		  },
		  {
			"electricalConnectionId": 1,
			"parameterId": 2,
			"permittedValueSet": [
			  {
				"range": [
				  {
					"min": { "number": 6, "scale": 0 },
					"max": { "number": 12, "scale": 0 }
				  }
				]
			  }
			]
		  }		]
	}`

	var sut model.ElectricalConnectionPermittedValueSetListDataType
	err := json.Unmarshal([]byte(existingDataJson), &sut)
	if assert.Nil(t, err) == false {
		return
	}

	// item with no identifiers
	newDataJson := `{
		"electricalConnectionPermittedValueSetData": [
		  {
			"permittedValueSet": [
			  {
				"range": [
				  {
					"min": { "number": 30, "scale": 0 },
					"max": { "number": 36, "scale": 0 }
				  }
				]
			  }
			]
		  }
		]
	}`

	var newData model.ElectricalConnectionPermittedValueSetListDataType
	err = json.Unmarshal([]byte(newDataJson), &newData)
	if assert.Nil(t, err) == false {
		return
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.ElectricalConnectionPermittedValueSetData
	// the new item should not be added
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 1, int(*item1.ElectricalConnectionId))
	assert.Equal(t, 1, int(*item1.ParameterId))
	assert.Equal(t, 1, len(item1.PermittedValueSet))
	valueSet := item1.PermittedValueSet[0]
	assert.Equal(t, 1, len(valueSet.Range))
	// the values of the item in the payload should be copied to the first item
	assert.Equal(t, 30, int(*valueSet.Range[0].Min.Number))
	assert.Equal(t, 0, int(*valueSet.Range[0].Min.Scale))
	assert.Equal(t, 36, int(*valueSet.Range[0].Max.Number))
	assert.Equal(t, 0, int(*valueSet.Range[0].Max.Scale))

	item2 := data[1]
	assert.Equal(t, 1, int(*item2.ElectricalConnectionId))
	assert.Equal(t, 2, int(*item2.ParameterId))
	assert.Equal(t, 1, len(item2.PermittedValueSet))
	valueSet = item2.PermittedValueSet[0]
	assert.Equal(t, 1, len(valueSet.Range))
	// the values of the item in the payload should be also copied to the second item
	assert.Equal(t, 30, int(*valueSet.Range[0].Min.Number))
	assert.Equal(t, 0, int(*valueSet.Range[0].Min.Scale))
	assert.Equal(t, 36, int(*valueSet.Range[0].Max.Number))
	assert.Equal(t, 0, int(*valueSet.Range[0].Max.Scale))
}

func TestElectricalConnectionDescriptionListDataType_Update(t *testing.T) {
	sut := model.ElectricalConnectionDescriptionListDataType{
		ElectricalConnectionDescriptionData: []model.ElectricalConnectionDescriptionDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				PowerSupplyType:        util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(1)),
				PowerSupplyType:        util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
			},
		},
	}

	newData := model.ElectricalConnectionDescriptionListDataType{
		ElectricalConnectionDescriptionData: []model.ElectricalConnectionDescriptionDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(1)),
				PowerSupplyType:        util.Ptr(model.ElectricalConnectionVoltageTypeTypeDc),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.ElectricalConnectionDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.ElectricalConnectionId))
	assert.Equal(t, model.ElectricalConnectionVoltageTypeTypeAc, *item1.PowerSupplyType)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.ElectricalConnectionId))
	assert.Equal(t, model.ElectricalConnectionVoltageTypeTypeDc, *item2.PowerSupplyType)
}

func TestElectricalConnectionParameterDescriptionListDataType_Update(t *testing.T) {
	sut := model.ElectricalConnectionParameterDescriptionListDataType{
		ElectricalConnectionParameterDescriptionData: []model.ElectricalConnectionParameterDescriptionDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(0)),
				VoltageType:            util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(1)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(0)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(0)),
				VoltageType:            util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
			},
		},
	}

	newData := model.ElectricalConnectionParameterDescriptionListDataType{
		ElectricalConnectionParameterDescriptionData: []model.ElectricalConnectionParameterDescriptionDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(1)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(0)),
				VoltageType:            util.Ptr(model.ElectricalConnectionVoltageTypeTypeDc),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, model.NewFilterTypePartial(), nil)

	data := sut.ElectricalConnectionParameterDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.ElectricalConnectionId))
	assert.Equal(t, model.ElectricalConnectionVoltageTypeTypeAc, *item1.VoltageType)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.ElectricalConnectionId))
	assert.Equal(t, model.ElectricalConnectionVoltageTypeTypeDc, *item2.VoltageType)
}
