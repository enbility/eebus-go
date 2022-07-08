package util

import "reflect"

func Type[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

// checks if type T implements interface I
func Implements[T any, I any]() bool {
	_, supported := any((*T)(nil)).(I)
	return supported
}
