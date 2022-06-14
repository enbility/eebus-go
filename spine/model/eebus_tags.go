package model

import (
	"fmt"
	"reflect"
	"strings"
)

type EEBusTag string

const EEBusTagFunction EEBusTag = "fct"

const EEBusTagName string = "eebus"

func EEBusTags(field reflect.StructField) map[EEBusTag]string {
	result := make(map[EEBusTag]string)
	tags := field.Tag.Get(EEBusTagName)
	for _, tag := range strings.Split(tags, ";") {
		pair := strings.Split(tag, ":")
		if len(pair) != 2 {
			fmt.Printf("Malformatted eebus tag: '%s'", tags)
		} else {
			result[EEBusTag(pair[0])] = pair[1]
		}
	}

	return result
}
