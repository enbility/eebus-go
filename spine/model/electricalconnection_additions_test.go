package model_test

import (
	"encoding/json"
	"testing"

	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/stretchr/testify/assert"
)

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
	sut.Update(&newData, model.NewFilterTypePartial(), nil)

	if assert.Equal(t, 2, len(sut.ElectricalConnectionPermittedValueSetData)) {
		item1 := sut.ElectricalConnectionPermittedValueSetData[0]
		assert.Equal(t, 1, int(*item1.ElectricalConnectionId))
		assert.Equal(t, 1, int(*item1.ParameterId))
		assert.Equal(t, 1, len(item1.PermittedValueSet))
		item2 := sut.ElectricalConnectionPermittedValueSetData[1]
		assert.Equal(t, 1, int(*item2.ElectricalConnectionId))
		assert.Equal(t, 2, int(*item2.ParameterId))
		assert.Equal(t, 2, len(item2.PermittedValueSet))
	}
}
