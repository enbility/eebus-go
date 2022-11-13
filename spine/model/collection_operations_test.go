package model_test

import (
	"fmt"
	"testing"

	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/DerAndereAndi/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	id   *uint `eebus:"key"`
	data *string
}

func (r testStruct) HashKey() string {
	return fmt.Sprintf("%d", r.id)
}

func TestUnion_NewData(t *testing.T) {
	existingData := []testStruct{
		{id: util.Ptr(uint(1)), data: util.Ptr(string("data1"))},
	}

	newData := []testStruct{
		{id: util.Ptr(uint(2)), data: util.Ptr(string("data2"))},
	}

	// Act
	result := model.Merge(existingData, newData)

	if assert.Equal(t, 2, len(result)) {
		assert.Equal(t, 1, int(*result[0].id))
		assert.Equal(t, "data1", string(*result[0].data))
		assert.Equal(t, 2, int(*result[1].id))
		assert.Equal(t, "data2", string(*result[1].data))
	}
}

func TestUnion_NewAndUpdateData(t *testing.T) {
	existingData := []testStruct{
		{id: util.Ptr(uint(1)), data: util.Ptr(string("data1"))},
		{id: util.Ptr(uint(2)), data: util.Ptr(string("data2"))},
	}

	newData := []testStruct{
		{id: util.Ptr(uint(2)), data: util.Ptr(string("data22"))},
		{id: util.Ptr(uint(3)), data: util.Ptr(string("data33"))},
	}

	// Act
	result := model.Merge(existingData, newData)

	if assert.Equal(t, 3, len(result)) {
		assert.Equal(t, 1, int(*result[0].id))
		assert.Equal(t, "data1", string(*result[0].data))
		assert.Equal(t, 2, int(*result[1].id))
		assert.Equal(t, "data22", string(*result[1].data))
		assert.Equal(t, 3, int(*result[2].id))
		assert.Equal(t, "data33", string(*result[2].data))
	}
}
