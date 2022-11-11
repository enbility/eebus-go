package model_test

import (
	"encoding/json"
	"testing"

	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/stretchr/testify/assert"
)

// verifies that a subset of existing items will be updated with identified new values
func TestElectricalConnectionPermittedValueSetListDataType_Update_Modify(t *testing.T) {
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

	// Act
	sut.NewUpdater(&newData, model.NewFilterTypePartial(), nil).DoUpdate()

	// check the non changing items
	assert.Equal(t, 4, len(sut.ElectricalConnectionPermittedValueSetData))
	item1 := sut.ElectricalConnectionPermittedValueSetData[0]
	assert.Equal(t, 0, int(*item1.ElectricalConnectionId))
	assert.Equal(t, 0, int(*item1.ParameterId))
	assert.Equal(t, 1, len(item1.PermittedValueSet))
	// check properties of updated item
	item2 := sut.ElectricalConnectionPermittedValueSetData[1]
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
func TestElectricalConnectionPermittedValueSetListDataType_Update_DeleteModify(t *testing.T) {
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

	electricalConnecctionId := model.ElectricalConnectionIdType(0)
	parameterId := model.ElectricalConnectionParameterIdType(0)

	delete := &model.FilterType{
		CmdControl: &model.CmdControlType{
			Delete: &model.ElementTagType{},
		},
		ElectricalConnectionPermittedValueSetListDataSelectors: &model.ElectricalConnectionPermittedValueSetListDataSelectorsType{
			ElectricalConnectionId: &electricalConnecctionId,
			ParameterId:            &parameterId,
		},
	}

	// Act
	sut.NewUpdater(&newData, model.NewFilterTypePartial(), delete).DoUpdate()

	// check the deleted item is gone
	assert.Equal(t, 3, len(sut.ElectricalConnectionPermittedValueSetData))
	// check properties of updated item
	item1 := sut.ElectricalConnectionPermittedValueSetData[0]
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

	electricalConnecctionId := model.ElectricalConnectionIdType(0)
	parameterId := model.ElectricalConnectionParameterIdType(0)

	delete := &model.FilterType{
		CmdControl: &model.CmdControlType{
			Delete: &model.ElementTagType{},
		},
		ElectricalConnectionPermittedValueSetListDataSelectors: &model.ElectricalConnectionPermittedValueSetListDataSelectorsType{
			ElectricalConnectionId: &electricalConnecctionId,
			ParameterId:            &parameterId,
		},
	}

	// Act
	sut.NewUpdater(nil, nil, delete).DoUpdate()

	// check the deleted item is added again
	assert.Equal(t, 3, len(sut.ElectricalConnectionPermittedValueSetData))
	// check properties of remaining item
	item1 := sut.ElectricalConnectionPermittedValueSetData[0]
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
func TestElectricalConnectionPermittedValueSetListDataType_Update_DeleteAdd(t *testing.T) {
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

	electricalConnecctionId := model.ElectricalConnectionIdType(0)
	parameterId := model.ElectricalConnectionParameterIdType(0)

	delete := &model.FilterType{
		CmdControl: &model.CmdControlType{
			Delete: &model.ElementTagType{},
		},
		ElectricalConnectionPermittedValueSetListDataSelectors: &model.ElectricalConnectionPermittedValueSetListDataSelectorsType{
			ElectricalConnectionId: &electricalConnecctionId,
			ParameterId:            &parameterId,
		},
	}

	// Act
	sut.NewUpdater(&newData, model.NewFilterTypePartial(), delete).DoUpdate()

	// check the deleted item is added again
	assert.Equal(t, 4, len(sut.ElectricalConnectionPermittedValueSetData))
	item1 := sut.ElectricalConnectionPermittedValueSetData[0]
	assert.Equal(t, 0, int(*item1.ElectricalConnectionId))
	assert.Equal(t, 0, int(*item1.ParameterId))
	assert.Equal(t, 1, len(item1.PermittedValueSet))
	// check properties of updated item
	item2 := sut.ElectricalConnectionPermittedValueSetData[1]
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
	sut.NewUpdater(&newData, model.NewFilterTypePartial(), nil).DoUpdate()

	// new item should be added
	if assert.Equal(t, 2, len(sut.ElectricalConnectionPermittedValueSetData)) {
		item1 := sut.ElectricalConnectionPermittedValueSetData[0]
		assert.Equal(t, 1, int(*item1.ElectricalConnectionId))
		assert.Equal(t, 1, int(*item1.ParameterId))
		assert.Equal(t, 1, len(item1.PermittedValueSet))
		// check properties of added item
		item2 := sut.ElectricalConnectionPermittedValueSetData[1]
		assert.Equal(t, 1, int(*item2.ElectricalConnectionId))
		assert.Equal(t, 2, int(*item2.ParameterId))
		assert.Equal(t, 2, len(item2.PermittedValueSet))
	}
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
	sut.NewUpdater(&newData, model.NewFilterTypePartial(), nil).DoUpdate()

	// the new item should not be added
	if assert.Equal(t, 2, len(sut.ElectricalConnectionPermittedValueSetData)) {
		item1 := sut.ElectricalConnectionPermittedValueSetData[0]
		assert.Equal(t, 1, int(*item1.ElectricalConnectionId))
		assert.Equal(t, 1, int(*item1.ParameterId))
		if assert.Equal(t, 1, len(item1.PermittedValueSet)) {
			valueSet := item1.PermittedValueSet[0]
			if assert.Equal(t, 1, len(valueSet.Range)) {
				// the values of the item in the payload should be copied to the first item
				assert.Equal(t, 30, int(*valueSet.Range[0].Min.Number))
				assert.Equal(t, 0, int(*valueSet.Range[0].Min.Scale))
				assert.Equal(t, 36, int(*valueSet.Range[0].Max.Number))
				assert.Equal(t, 0, int(*valueSet.Range[0].Max.Scale))
			}
		}

		item2 := sut.ElectricalConnectionPermittedValueSetData[1]
		assert.Equal(t, 1, int(*item2.ElectricalConnectionId))
		assert.Equal(t, 2, int(*item2.ParameterId))
		if assert.Equal(t, 1, len(item2.PermittedValueSet)) {
			valueSet := item2.PermittedValueSet[0]
			if assert.Equal(t, 1, len(valueSet.Range)) {
				// the values of the item in the payload should be also copied to the second item
				assert.Equal(t, 30, int(*valueSet.Range[0].Min.Number))
				assert.Equal(t, 0, int(*valueSet.Range[0].Min.Scale))
				assert.Equal(t, 36, int(*valueSet.Range[0].Max.Number))
				assert.Equal(t, 0, int(*valueSet.Range[0].Max.Scale))
			}
		}
	}
}
