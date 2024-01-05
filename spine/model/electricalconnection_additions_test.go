package model

import (
	"encoding/json"
	"testing"

	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func TestElectricalConnectionStateListDataType_Update(t *testing.T) {
	sut := ElectricalConnectionStateListDataType{
		ElectricalConnectionStateData: []ElectricalConnectionStateDataType{
			{
				ElectricalConnectionId: util.Ptr(ElectricalConnectionIdType(0)),
				CurrentEnergyMode:      util.Ptr(EnergyModeTypeProduce),
			},
			{
				ElectricalConnectionId: util.Ptr(ElectricalConnectionIdType(1)),
				CurrentEnergyMode:      util.Ptr(EnergyModeTypeProduce),
			},
		},
	}

	newData := ElectricalConnectionStateListDataType{
		ElectricalConnectionStateData: []ElectricalConnectionStateDataType{
			{
				ElectricalConnectionId: util.Ptr(ElectricalConnectionIdType(1)),
				CurrentEnergyMode:      util.Ptr(EnergyModeTypeConsume),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

	data := sut.ElectricalConnectionStateData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.ElectricalConnectionId))
	assert.Equal(t, EnergyModeTypeProduce, *item1.CurrentEnergyMode)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.ElectricalConnectionId))
	assert.Equal(t, EnergyModeTypeConsume, *item2.CurrentEnergyMode)
}

// verifies that a subset of existing items will be updated with identified new values
func TestElectricalConnectionPermittedValueSetListDataType_Update_Modify(t *testing.T) {
	sut := ElectricalConnectionPermittedValueSetListDataType{
		ElectricalConnectionPermittedValueSetData: []ElectricalConnectionPermittedValueSetDataType{
			{
				ElectricalConnectionId: util.Ptr(ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(ElectricalConnectionParameterIdType(0)),
				PermittedValueSet: []ScaledNumberSetType{
					{
						Range: []ScaledNumberRangeType{
							{
								Min: NewScaledNumberType(1),
							},
						},
					},
				},
			},
			{
				ElectricalConnectionId: util.Ptr(ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(ElectricalConnectionParameterIdType(1)),
				PermittedValueSet: []ScaledNumberSetType{
					{
						Range: []ScaledNumberRangeType{
							{
								Min: NewScaledNumberType(6),
								Max: NewScaledNumberType(16),
							},
						},
					},
				},
			},
			{
				ElectricalConnectionId: util.Ptr(ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(ElectricalConnectionParameterIdType(2)),
				PermittedValueSet: []ScaledNumberSetType{
					{
						Range: []ScaledNumberRangeType{
							{
								Min: NewScaledNumberType(6),
								Max: NewScaledNumberType(16),
							},
						},
					},
				},
			},
			{
				ElectricalConnectionId: util.Ptr(ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(ElectricalConnectionParameterIdType(3)),
				PermittedValueSet: []ScaledNumberSetType{
					{
						Range: []ScaledNumberRangeType{
							{
								Min: NewScaledNumberType(6),
								Max: NewScaledNumberType(16),
							},
						},
					},
				},
			},
		},
	}

	newData := ElectricalConnectionPermittedValueSetListDataType{
		ElectricalConnectionPermittedValueSetData: []ElectricalConnectionPermittedValueSetDataType{
			{
				ElectricalConnectionId: util.Ptr(ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(ElectricalConnectionParameterIdType(1)),
				PermittedValueSet: []ScaledNumberSetType{
					{
						Range: []ScaledNumberRangeType{
							{
								Min: NewScaledNumberType(2),
								Max: NewScaledNumberType(16),
							},
						},
					},
				},
			},
			{
				ElectricalConnectionId: util.Ptr(ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(ElectricalConnectionParameterIdType(2)),
				PermittedValueSet: []ScaledNumberSetType{
					{
						Range: []ScaledNumberRangeType{
							{
								Min: NewScaledNumberType(2),
								Max: NewScaledNumberType(16),
							},
						},
					},
				},
			},
			{
				ElectricalConnectionId: util.Ptr(ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(ElectricalConnectionParameterIdType(3)),
				PermittedValueSet: []ScaledNumberSetType{
					{
						Range: []ScaledNumberRangeType{
							{
								Min: NewScaledNumberType(2),
								Max: NewScaledNumberType(16),
							},
						},
					},
				},
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

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

	var sut ElectricalConnectionPermittedValueSetListDataType
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

	var newData ElectricalConnectionPermittedValueSetListDataType
	err = json.Unmarshal([]byte(newDataJson), &newData)
	if assert.Nil(t, err) == false {
		return
	}

	partial := &FilterType{
		CmdControl: &CmdControlType{
			Partial: &ElementTagType{},
		},
		ElectricalConnectionPermittedValueSetListDataSelectors: &ElectricalConnectionPermittedValueSetListDataSelectorsType{
			ElectricalConnectionId: util.Ptr[ElectricalConnectionIdType](0),
			ParameterId:            util.Ptr[ElectricalConnectionParameterIdType](1),
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

	var sut ElectricalConnectionPermittedValueSetListDataType
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

	var newData ElectricalConnectionPermittedValueSetListDataType
	err = json.Unmarshal([]byte(newDataJson), &newData)
	if assert.Nil(t, err) == false {
		return
	}

	delete := &FilterType{
		CmdControl: &CmdControlType{
			Delete: &ElementTagType{},
		},
		ElectricalConnectionPermittedValueSetListDataSelectors: &ElectricalConnectionPermittedValueSetListDataSelectorsType{
			ElectricalConnectionId: util.Ptr[ElectricalConnectionIdType](0),
			ParameterId:            util.Ptr[ElectricalConnectionParameterIdType](0),
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), delete)

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

	var sut ElectricalConnectionPermittedValueSetListDataType
	err := json.Unmarshal([]byte(existingDataJson), &sut)
	if assert.Nil(t, err) == false {
		return
	}

	delete := &FilterType{
		CmdControl: &CmdControlType{
			Delete: &ElementTagType{},
		},
		ElectricalConnectionPermittedValueSetListDataSelectors: &ElectricalConnectionPermittedValueSetListDataSelectorsType{
			ElectricalConnectionId: util.Ptr[ElectricalConnectionIdType](0),
			ParameterId:            util.Ptr[ElectricalConnectionParameterIdType](0),
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

	var sut ElectricalConnectionPermittedValueSetListDataType
	err := json.Unmarshal([]byte(existingDataJson), &sut)
	if assert.Nil(t, err) == false {
		return
	}

	delete := &FilterType{
		CmdControl: &CmdControlType{
			Delete: &ElementTagType{},
		},
		ElectricalConnectionPermittedValueSetDataElements: &ElectricalConnectionPermittedValueSetDataElementsType{
			PermittedValueSet: &ElementTagType{},
		},
		ElectricalConnectionPermittedValueSetListDataSelectors: &ElectricalConnectionPermittedValueSetListDataSelectorsType{
			ElectricalConnectionId: util.Ptr[ElectricalConnectionIdType](0),
			ParameterId:            util.Ptr[ElectricalConnectionParameterIdType](0),
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
	var nilValue []ScaledNumberSetType
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

	var sut ElectricalConnectionPermittedValueSetListDataType
	err := json.Unmarshal([]byte(existingDataJson), &sut)
	if assert.Nil(t, err) == false {
		return
	}

	delete := &FilterType{
		CmdControl: &CmdControlType{
			Delete: &ElementTagType{},
		},
		ElectricalConnectionPermittedValueSetDataElements: &ElectricalConnectionPermittedValueSetDataElementsType{
			PermittedValueSet: &ElementTagType{},
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
	var nilValue []ScaledNumberSetType
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

	var sut ElectricalConnectionPermittedValueSetListDataType
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

	var newData ElectricalConnectionPermittedValueSetListDataType
	err = json.Unmarshal([]byte(newDataJson), &newData)
	if assert.Nil(t, err) == false {
		return
	}

	delete := &FilterType{
		CmdControl: &CmdControlType{
			Delete: &ElementTagType{},
		},
		ElectricalConnectionPermittedValueSetListDataSelectors: &ElectricalConnectionPermittedValueSetListDataSelectorsType{
			ElectricalConnectionId: util.Ptr[ElectricalConnectionIdType](0),
			ParameterId:            util.Ptr[ElectricalConnectionParameterIdType](0),
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), delete)

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

	var sut ElectricalConnectionPermittedValueSetListDataType
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

	var newData ElectricalConnectionPermittedValueSetListDataType
	err = json.Unmarshal([]byte(newDataJson), &newData)
	if assert.Nil(t, err) == false {
		return
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

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

	var sut ElectricalConnectionPermittedValueSetListDataType
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

	var newData ElectricalConnectionPermittedValueSetListDataType
	err = json.Unmarshal([]byte(newDataJson), &newData)
	if assert.Nil(t, err) == false {
		return
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

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
	sut := ElectricalConnectionDescriptionListDataType{
		ElectricalConnectionDescriptionData: []ElectricalConnectionDescriptionDataType{
			{
				ElectricalConnectionId: util.Ptr(ElectricalConnectionIdType(0)),
				PowerSupplyType:        util.Ptr(ElectricalConnectionVoltageTypeTypeAc),
			},
			{
				ElectricalConnectionId: util.Ptr(ElectricalConnectionIdType(1)),
				PowerSupplyType:        util.Ptr(ElectricalConnectionVoltageTypeTypeAc),
			},
		},
	}

	newData := ElectricalConnectionDescriptionListDataType{
		ElectricalConnectionDescriptionData: []ElectricalConnectionDescriptionDataType{
			{
				ElectricalConnectionId: util.Ptr(ElectricalConnectionIdType(1)),
				PowerSupplyType:        util.Ptr(ElectricalConnectionVoltageTypeTypeDc),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

	data := sut.ElectricalConnectionDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.ElectricalConnectionId))
	assert.Equal(t, ElectricalConnectionVoltageTypeTypeAc, *item1.PowerSupplyType)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.ElectricalConnectionId))
	assert.Equal(t, ElectricalConnectionVoltageTypeTypeDc, *item2.PowerSupplyType)
}

func TestElectricalConnectionParameterDescriptionListDataType_Update(t *testing.T) {
	sut := ElectricalConnectionParameterDescriptionListDataType{
		ElectricalConnectionParameterDescriptionData: []ElectricalConnectionParameterDescriptionDataType{
			{
				ElectricalConnectionId: util.Ptr(ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(ElectricalConnectionParameterIdType(0)),
				VoltageType:            util.Ptr(ElectricalConnectionVoltageTypeTypeAc),
			},
			{
				ElectricalConnectionId: util.Ptr(ElectricalConnectionIdType(1)),
				ParameterId:            util.Ptr(ElectricalConnectionParameterIdType(0)),
				MeasurementId:          util.Ptr(MeasurementIdType(0)),
				VoltageType:            util.Ptr(ElectricalConnectionVoltageTypeTypeAc),
			},
		},
	}

	newData := ElectricalConnectionParameterDescriptionListDataType{
		ElectricalConnectionParameterDescriptionData: []ElectricalConnectionParameterDescriptionDataType{
			{
				ElectricalConnectionId: util.Ptr(ElectricalConnectionIdType(1)),
				ParameterId:            util.Ptr(ElectricalConnectionParameterIdType(0)),
				VoltageType:            util.Ptr(ElectricalConnectionVoltageTypeTypeDc),
			},
		},
	}

	// Act
	sut.UpdateList(&newData, NewFilterTypePartial(), nil)

	data := sut.ElectricalConnectionParameterDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.ElectricalConnectionId))
	assert.Equal(t, ElectricalConnectionVoltageTypeTypeAc, *item1.VoltageType)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.ElectricalConnectionId))
	assert.Equal(t, ElectricalConnectionVoltageTypeTypeDc, *item2.VoltageType)
}
