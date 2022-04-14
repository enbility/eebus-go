package util

import (
	"bytes"
	"encoding/json"
	"strings"

	ordered "gitlab.com/c0b/go-ordered-json"
)

// process incoming json strings and transform it to match the model structure
func JsonFromEEBUSJson(json []byte) []byte {
	var result = bytes.ReplaceAll(json, []byte("[{"), []byte("{"))
	result = bytes.ReplaceAll(result, []byte("},{"), []byte(","))
	result = bytes.ReplaceAll(result, []byte("}]"), []byte("}"))
	result = bytes.ReplaceAll(result, []byte("[]"), []byte("{}"))

	return result
}

// convert objects in json to be arrays with each field being an array alement as eebus expects it
func process_eebus_json_hierarchie_level(data interface{}) interface{} {
	switch data.(type) {
	case *ordered.OrderedMap:
		var new_array []interface{} = make([]interface{}, 0)

		orderedData := data.(*ordered.OrderedMap)
		iter := orderedData.EntriesIter()
		for {
			pair, ok := iter()
			if !ok {
				break
			}
			var new_value = process_eebus_json_hierarchie_level(pair.Value)
			var new_object = map[string]interface{}{pair.Key: new_value}
			new_array = append(new_array, new_object)
		}
		return new_array

	case []interface{}:
		var new_array []interface{} = make([]interface{}, 0)
		for _, value := range data.([]interface{}) {
			var new_value = process_eebus_json_hierarchie_level(value)
			new_array = append(new_array, new_value)
		}
		return new_array
	default:
		return data
	}
}

func JsonIntoEEBUSJson(data []byte) (string, error) {
	// EEBUS defines the items to be ordered in the array,
	// so we can't use map[string]interface{} as that would
	// cause a random order when Unmarshalling
	var temp *ordered.OrderedMap = ordered.NewOrderedMap()

	if err := json.Unmarshal(data, &temp); err != nil {
		return "", err
	}

	var result = process_eebus_json_hierarchie_level(temp)

	var b, err = json.Marshal(result)
	if err != nil {
		return "", err
	}

	var json = string(b)

	// we are lazy: fix the first item being put into an array
	json = strings.TrimPrefix(json, "[")
	json = strings.TrimSuffix(json, "]")

	return json, nil
}
