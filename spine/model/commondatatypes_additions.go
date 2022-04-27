package model

import (
	"fmt"
)

func (r *FeatureAddressType) String() string {
	if r == nil {
		return ""
	}

	var result string = ""
	if r.Device != nil {
		result += string(*r.Device)
	}
	result += ":["
	for _, id := range r.Entity {
		result += fmt.Sprintf("%d,", id)
	}
	result += "]:"
	if r.Feature != nil {
		result += fmt.Sprintf("%d", *r.Feature)
	}
	return result
}
