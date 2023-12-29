package util

import (
	"encoding/json"
	"strings"
)

// standardize the provided SKI strings
func NormalizeSKI(ski string) string {
	ski = strings.ReplaceAll(ski, " ", "")
	ski = strings.ReplaceAll(ski, "-", "")
	ski = strings.ToLower(ski)

	return ski
}

// quick way to a struct into another
func DeepCopy[A any](source, dest A) {
	byt, _ := json.Marshal(source)
	_ = json.Unmarshal(byt, dest)
}
