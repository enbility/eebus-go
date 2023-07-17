package mdns

import "net"

const shipZeroConfServiceType = "_ship._tcp"
const shipZeroConfDomain = "local."

type MdnsProvider interface {
	CheckAvailability() bool
	Shutdown()
	Announce(serviceName string, port int, txt []string) error
	Unannounce()
	ResolveEntries(cancelChan chan bool)
}

type MdnsManager interface {
	ProcessMdnsEntry(elements map[string]string, name, host string, addresses []net.IP, port int, remove bool)
}
