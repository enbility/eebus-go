package model

import (
	"reflect"
	"strings"

	"github.com/enbility/eebus-go/logging"
)

type EEBusTag string

const (
	EEBusTagFunction EEBusTag = "fct"
	EEBusTagType     EEBusTag = "typ"
	EEBusTagKey      EEBusTag = "key"
)

type EEBusTagTypeType string

const (
	EEBusTagTypeTypeSelector EEBusTagTypeType = "selector"
	EEbusTagTypeTypeElements EEBusTagTypeType = "elements"
)

const EEBusTagName string = "eebus"

func EEBusTags(field reflect.StructField) map[EEBusTag]string {
	result := make(map[EEBusTag]string)
	tags := field.Tag.Get(EEBusTagName)
	if len(tags) == 0 {
		return result
	}
	for _, tag := range strings.Split(tags, ",") {
		pair := strings.Split(tag, ":")
		if len(pair) == 1 {
			// boolean tags like "key"
			result[EEBusTag(pair[0])] = "true"
		} else if len(pair) == 2 {
			result[EEBusTag(pair[0])] = pair[1]
		} else {
			logging.Log.Errorf("error: malformatted eebus tag: '%s'", tags)
		}
	}

	return result
}
