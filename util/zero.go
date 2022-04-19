package util

import "reflect"

func Zero[T any]() (ret T) {
	return
}

func IsZero[T comparable](v T) bool {
	return v == *new(T)
}

func IsNil[T any](v T) bool {
	return !reflect.ValueOf(v).IsValid()
}
