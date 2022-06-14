package model

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/DerAndereAndi/eebus-go/util"
)

func (r *MsgCounterType) String() string {
	if r == nil {
		return ""
	}
	return fmt.Sprintf("%d", *r)
}

type CmdData struct {
	FieldName string
	Function  *FunctionType
	Value     any
}

// Get the data and some meta data of the current value
func (cmd *CmdType) Data() (*CmdData, error) {
	t := reflect.ValueOf(*cmd)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Kind() == reflect.Ptr {
			if !f.IsNil() {
				sf := t.Type().Field(i)
				function := sf.Tag.Get("function")
				var ft *FunctionType = nil
				if len(function) > 0 {
					ft = util.Ptr(FunctionType(function))
				}
				return &CmdData{
					FieldName: sf.Name,
					Function:  ft,
					Value:     f.Interface(),
				}, nil
			}
		}
	}
	return nil, errors.New("Data not found in Cmd")
}

// Get the non empty field name of the data type
func (cmd *CmdType) DataName() string {
	data, err := cmd.Data()
	if err != nil {
		return "unknown"
	}
	return data.FieldName
}
