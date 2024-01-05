package mdns

import "net"

//go:generate mockery --name=MdnsProvider

type MdnsProvider interface {
	CheckAvailability() bool
	Shutdown()
	Announce(serviceName string, port int, txt []string) error
	Unannounce()
	ResolveEntries(cancelChan chan bool, callback func(elements map[string]string, name, host string, addresses []net.IP, port int, remove bool))
}
