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

// FilterData stores the function field name and
// selector field name for a function
type FilterData struct {
	Elements any
	Selector any
	Function *FunctionType
}

func (f *FilterData) SelectorMatch(item any) bool {
	if f.Selector == nil {
		return false
	}

	v := reflect.ValueOf(f.Selector).Elem()
	t := reflect.TypeOf(f.Selector).Elem()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() != reflect.Ptr {
			continue
		}

		if field.IsNil() {
			continue
		}

		fieldname := t.Field(i).Name
		value := field.Elem().Interface()

		itemV := reflect.ValueOf(item).Elem()
		itemF := itemV.FieldByName(fieldname)
		if !itemF.IsValid() {
			continue
		}

		itemValue := itemF.Elem().Interface()
		if itemValue != value {
			return false
		}
	}

	return true
}

// Get the data and some meta data for the current value
func (f *FilterType) Data() (*FilterData, error) {
	var elements any = nil
	var selector any = nil
	var function string

	v := reflect.ValueOf(*f)
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if f.Kind() != reflect.Ptr {
			continue
		}

		if f.IsNil() {
			continue
		}

		sf := v.Type().Field(i)
		// Exclude the generic fields
		if sf.Name == "CmdControl" || sf.Name == "FilterId" {
			continue
		}

		eebusTags := EEBusTags(sf)
		fname, exists := eebusTags[EEBusTagFunction]
		if !exists || len(fname) == 0 {
			continue
		}
		typ, exists := eebusTags[EEBusTagType]
		if !exists || len(typ) == 0 {
			continue
		}

		function = fname

		switch typ {
		case string(EEBusTagTypeTypeSelector):
			selector = f.Interface()
		case string(EEbusTagTypeTypeElements):
			elements = f.Interface()
		}
	}

	if len(function) == 0 {
		return nil, errors.New("Data not found in Filter")
	}

	ft := util.Ptr(FunctionType(function))

	return &FilterData{
		Elements: elements,
		Selector: selector,
		Function: ft,
	}, nil
}

// CmdData stores the function field name for a cmd field
type CmdData struct {
	FieldName string
	Function  *FunctionType
	Value     any
}

// Get the data and some meta data of the current value
func (cmd *CmdType) Data() (*CmdData, error) {
	v := reflect.ValueOf(*cmd)
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if f.Kind() != reflect.Ptr {
			continue
		}

		if f.IsNil() {
			continue
		}

		sf := v.Type().Field(i)
		// Exclude the CmdOptionGroup fields
		if sf.Name == "Function" || sf.Name == "Filter" {
			continue
		}

		eebusTags := EEBusTags(sf)
		function, exists := eebusTags[EEBusTagFunction]
		if !exists {
			continue
		}

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

func (cmd *CmdType) ExtractFilter() (filterPartial *FilterType, filterDelete *FilterType) {
	if cmd != nil && cmd.Filter != nil && len(cmd.Filter) > 0 {
		for i := range cmd.Filter {
			if cmd.Filter[i].CmdControl.Partial != nil {
				filterPartial = &cmd.Filter[i]
			} else if cmd.Filter[i].CmdControl.Delete != nil {
				filterDelete = &cmd.Filter[i]
			}
		}
	}
	return
}

func NewFilterTypePartial() *FilterType {
	return &FilterType{CmdControl: &CmdControlType{Partial: &ElementTagType{}}}
}
