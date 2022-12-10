package model

import (
	"fmt"
	"reflect"
)

// creates an hash key by using fields that have eebus tag "key"
func hashKey(data any) string {
	result := ""

	keys := keyFieldNames(data)

	if len(keys) == 0 {
		return result
	}

	v := reflect.ValueOf(data)

	for _, fieldName := range keys {
		f := v.FieldByName(fieldName)

		if f.IsNil() || !f.IsValid() {
			return result
		}

		switch f.Elem().Kind() {
		case reflect.String:
			value := f.Elem().String()

			if len(result) > 0 {
				result = fmt.Sprintf("%s|", result)
			}
			result = fmt.Sprintf("%s%s", result, value)

		case reflect.Uint:
			value := f.Elem().Uint()

			if len(result) > 0 {
				result = fmt.Sprintf("%s|", result)
			}
			result = fmt.Sprintf("%s%d", result, value)

		default:
			return result
		}
	}

	return result
}

// Merges two slices into one. The item in the first slice will be replaced by the one in the second slice
// if the hash key is the same. Items in the second slice which are not in the first will be added.
func Merge[T any](s1 []T, s2 []T) []T {
	result := []T{}

	m2 := ToMap(s2)

	// go through the first slice
	m1 := make(map[string]T, len(s1))
	for _, s1Item := range s1 {
		s1ItemHash := hashKey(s1Item)
		// s1ItemHash := s1Item.HashKey()
		s2Item, exist := m2[s1ItemHash]
		if exist {
			// the item in the first slice will be replaces by the one of the second slice
			result = append(result, s2Item)
		} else {
			result = append(result, s1Item)
		}

		m1[s1ItemHash] = s1Item
	}

	// append items which were not in the first slice
	for _, s2Item := range s2 {
		s2ItemHash := hashKey(s2Item)
		// s2ItemHash := s2Item.HashKey()
		_, exist := m1[s2ItemHash]
		if !exist {
			result = append(result, s2Item)
		}
	}

	return result
}

func ToMap[T any](s []T) map[string]T {
	result := make(map[string]T, len(s))
	for _, item := range s {
		result[hashKey(item)] = item
	}
	return result
}

/*
func FindFirst[T any](s []T, predicate func(i T) bool) *T {
	for _, item := range s {
		if predicate(item) {
			return &item
		}
	}
	return nil
}

func Values[K comparable, V any](m map[K]V) []V {
	ret := make([]V, 0, len(m))
	for _, v := range m {
		ret = append(ret, v)
	}
	return ret
}

// casts all elements in slice s to type D
func CastElements[S any, D any](s []S) []D {
	result := make([]D, len(s))
	for i, item := range s {
		result[i] = any(item).(D)
	}
	return result
}
*/
