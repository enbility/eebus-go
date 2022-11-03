package spine

import "github.com/DerAndereAndi/eebus-go/logging"

var log logging.Logging

// Sets a custom logging implementation
// By default NoLogging is used, so no logs are printed
// This is used by service.SetLogging()
func SetLogging(logger logging.Logging) {
	if logger == nil {
		return
	}
	log = logger
}
