package util_test

import (
	"fmt"
	"testing"

	"github.com/DerAndereAndi/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	id   int
	data string
}

func calcHash(s *testStruct) string {
	return fmt.Sprintf("%d", s.id)
}

func TestUnion_NewData(t *testing.T) {
	existingData := []testStruct{
		{id: 1, data: "data1"},
	}

	newData := []testStruct{
		{id: 2, data: "data2"},
	}

	// Act
	result := util.Union(existingData, newData, calcHash)

	if assert.Equal(t, 2, len(result)) {
		assert.Equal(t, 1, result[0].id)
		assert.Equal(t, "data1", result[0].data)
		assert.Equal(t, 2, result[1].id)
		assert.Equal(t, "data2", result[1].data)
	}
}

func TestUnion_NewAndUpdateData(t *testing.T) {
	existingData := []testStruct{
		{id: 1, data: "data1"},
		{id: 2, data: "data2"},
	}

	newData := []testStruct{
		{id: 2, data: "data22"},
		{id: 3, data: "data33"},
	}

	// Act
	result := util.Union(existingData, newData, calcHash)

	if assert.Equal(t, 3, len(result)) {
		assert.Equal(t, 1, result[0].id)
		assert.Equal(t, "data1", result[0].data)
		assert.Equal(t, 2, result[1].id)
		assert.Equal(t, "data22", result[1].data)
		assert.Equal(t, 3, result[2].id)
		assert.Equal(t, "data33", result[2].data)
	}
}
