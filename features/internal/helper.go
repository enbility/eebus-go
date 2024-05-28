package internal

import (
	"reflect"

	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

func featureDataCopyOfType[T any](
	featureLocal spineapi.FeatureLocalInterface,
	featureRemote spineapi.FeatureRemoteInterface,
	function model.FunctionType) (*T, error) {
	var data *T
	var err error

	if featureLocal != nil {
		data, err = spine.LocalFeatureDataCopyOfType[*T](featureLocal, function)
	} else {
		data, err = spine.RemoteFeatureDataCopyOfType[*T](featureRemote, function)
	}

	return data, err
}

func searchFilterInItem[T any](item T, filter T) bool {
	v := reflect.ValueOf(item)

	match := true
	for i := 0; i < v.NumField(); i++ {
		filterField := reflect.ValueOf(filter).Field(i)
		itemField := v.Field(i)

		if filterField.Kind() != reflect.Ptr || itemField.Kind() != reflect.Ptr {
			continue
		}

		if (!filterField.IsNil() && !itemField.IsNil() && filterField.Elem().Interface() != itemField.Elem().Interface()) ||
			(!filterField.IsNil() && itemField.IsNil()) {
			match = false
			break
		}
	}

	return match
}

func searchFilterInList[T any](list []T, filter T) []T {
	var result []T

	for _, item := range list {
		match := searchFilterInItem[T](item, filter)

		if match {
			result = append(result, item)
		}
	}

	return result
}
