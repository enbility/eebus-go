package mdns

import (
	"strings"
)

// parse mDNS text fields
func parseTxt(txt []string) map[string]string {
	result := make(map[string]string)

	for _, item := range txt {
		s := strings.Split(item, "=")
		if len(s) != 2 {
			continue
		}
		result[s[0]] = s[1]
	}

	return result
}
