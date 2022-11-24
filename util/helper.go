package util

import "strings"

// standardize the provided SKI strings
func NormalizeSKI(ski string) string {
	ski = strings.ReplaceAll(ski, " ", "")
	ski = strings.ReplaceAll(ski, "-", "")
	ski = strings.ToLower(ski)

	return ski
}
