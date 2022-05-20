package model

import (
	"fmt"
	"reflect"
)

func (r *MsgCounterType) String() string {
	if r == nil {
		return ""
	}
	return fmt.Sprintf("%d", *r)
}

// Get the non empty field name of the data type
func (cmd CmdType) DataName() string {
	t := reflect.ValueOf(cmd)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Kind() == reflect.Ptr {
			if !f.IsNil() {
				return t.Type().Field(i).Name
			}
		}
	}

	return "unknown"
}
