package mdns

import (
	"net"

	"github.com/enbility/ship-go/logging"
	"github.com/godbus/dbus/v5"
	"github.com/holoplot/go-avahi"
)

type AvahiProvider struct {
	ifaceIndexes []int32

	avServer     *avahi.Server
	avEntryGroup *avahi.EntryGroup
}

func NewAvahiProvider(ifaceIndexes []int32) *AvahiProvider {
	return &AvahiProvider{
		ifaceIndexes: ifaceIndexes,
	}
}

var _ MdnsProvider = (*AvahiProvider)(nil)

func (a *AvahiProvider) CheckAvailability() bool {
	dbusConn, err := dbus.SystemBus()
	if err != nil {
		return false
	}

	a.avServer, err = avahi.ServerNew(dbusConn)
	if err != nil {
		return false
	}

	if _, err := a.avServer.GetAPIVersion(); err != nil {
		return false
	}

	avBrowser, err := a.avServer.ServiceBrowserNew(avahi.InterfaceUnspec, avahi.ProtoUnspec, shipZeroConfServiceType, shipZeroConfDomain, 0)
	if err != nil {
		return false
	}

	if avBrowser != nil {
		a.avServer.ServiceBrowserFree(avBrowser)
		return true
	}

	return false
}

func (a *AvahiProvider) Shutdown() {
	if a.avServer == nil {
		return
	}

	a.avServer.Close()
	a.avServer = nil
	a.avEntryGroup = nil
}

func (a *AvahiProvider) Announce(serviceName string, port int, txt []string) error {
	logging.Log().Debug("mdns: using avahi")

	entryGroup, err := a.avServer.EntryGroupNew()
	if err != nil {
		return err
	}

	var btxt [][]byte
	for _, t := range txt {
		btxt = append(btxt, []byte(t))
	}

	for _, iface := range a.ifaceIndexes {
		err = entryGroup.AddService(iface, avahi.ProtoUnspec, 0, serviceName, shipZeroConfServiceType, shipZeroConfDomain, "", uint16(port), btxt)
		if err != nil {
			return err
		}
	}

	err = entryGroup.Commit()
	if err != nil {
		return err
	}

	a.avEntryGroup = entryGroup

	return nil
}

func (a *AvahiProvider) Unannounce() {
	if a.avEntryGroup == nil {
		return
	}

	a.avServer.EntryGroupFree(a.avEntryGroup)
	a.avEntryGroup = nil
}

func (a *AvahiProvider) ResolveEntries(cancelChan chan bool, callback func(elements map[string]string, name, host string, addresses []net.IP, port int, remove bool)) {
	var err error
	var end bool

	var avBrowser *avahi.ServiceBrowser

	// instead of limiting search on specific allowed interfaces, we allow all and filter the results
	if avBrowser, err = a.avServer.ServiceBrowserNew(avahi.InterfaceUnspec, avahi.ProtoUnspec, shipZeroConfServiceType, shipZeroConfDomain, 0); err != nil {
		logging.Log().Debug("mdns: error setting up avahi browser:", err)
		return
	}

	if avBrowser == nil {
		logging.Log().Debug("mdns: avahi browser is not available")
		return
	}

	for !end {
		select {
		case <-cancelChan:
			end = true
			break
		case service := <-avBrowser.AddChannel:
			a.processService(service, false, callback)
		case service := <-avBrowser.RemoveChannel:
			a.processService(service, true, callback)
		}
	}

	a.avServer.ServiceBrowserFree(avBrowser)
}

// process an avahi mDNS service
// as avahi returns a service per interface, we need to combine them
func (a *AvahiProvider) processService(service avahi.Service, remove bool, callback func(elements map[string]string, name, host string, addresses []net.IP, port int, remove bool)) {
	// check if the service is within the allowed list
	allow := false
	if len(a.ifaceIndexes) == 1 && a.ifaceIndexes[0] == avahi.InterfaceUnspec {
		allow = true
	} else {
		for _, iface := range a.ifaceIndexes {
			if service.Interface == iface {
				allow = true
				break
			}
		}
	}

	if !allow {
		logging.Log().Debug("avahi - ignoring service as its interface is not in the allowed list:", service.Name)
		return
	}

	resolved, err := a.avServer.ResolveService(service.Interface, service.Protocol, service.Name, service.Type, service.Domain, avahi.ProtoUnspec, 0)
	if err != nil {
		logging.Log().Debug("avahi - error resolving service:", service, "error:", err)
		return
	}

	// convert [][]byte to []string manually
	var txt []string
	for _, element := range resolved.Txt {
		txt = append(txt, string(element))
	}
	elements := parseTxt(txt)

	// convert address to net.IP
	address := net.ParseIP(resolved.Address)
	// if the address can not be used, ignore the entry
	if address == nil || address.IsUnspecified() {
		logging.Log().Debug("avahi - service provides unusable address:", service.Name)
		return
	}

	// Ignore IPv6 addresses for now
	if address.To4() == nil {
		return
	}

	callback(elements, resolved.Name, resolved.Host, []net.IP{address}, int(resolved.Port), remove)
}
