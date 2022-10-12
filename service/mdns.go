package service

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/godbus/dbus/v5"
	"github.com/holoplot/go-avahi"
	"github.com/libp2p/zeroconf/v2"
)

type MdnsEntry struct {
	Name       string
	Identifier string   // mandatory
	Path       string   // mandatory
	Register   bool     // mandatory
	Brand      string   // optional
	Type       string   // optional
	Model      string   // optional
	Host       string   // mandatory
	Port       int      // mandatory
	Addresses  []net.IP // mandatory
}

type MdnsSearch interface {
	ReportMdnsEntries(entries map[string]MdnsEntry)
}

type mdns struct {
	serviceDescription *ServiceDescription
	ski                string

	isAnnounced         bool
	isSearchingServices bool

	cancelChan chan bool

	// the currently available mDNS entries with the SKI as the key in the map
	entries map[string]MdnsEntry

	// the registered callbacks
	searchDelegates []MdnsSearch

	// The zeroconf service for mDNS related tasks
	zc *zeroconf.Server

	// The alternative avahi mDNS service
	av           *avahi.Server
	avEntryGroup *avahi.EntryGroup

	mux sync.Mutex
}

func newMDNS(ski string, serviceDescription *ServiceDescription) (*mdns, error) {
	m := &mdns{
		ski:                ski,
		serviceDescription: serviceDescription,
		entries:            make(map[string]MdnsEntry),
		cancelChan:         make(chan bool),
	}

	if av, err := m.setupAvahi(); err == nil {
		m.av = av
	}

	// on startup always start mDNS announcement
	if err := m.Announce(); err != nil {
		return nil, err
	}

	// catch signals
	go func() {
		signalC := make(chan os.Signal, 1)
		signal.Notify(signalC, os.Interrupt, syscall.SIGTERM)

		<-signalC // wait for signal

		m.Unannounce()
	}()

	return m, nil
}

// setup avahi for mDNS
func (m *mdns) setupAvahi() (*avahi.Server, error) {
	dbusConn, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}

	avahiServer, err := avahi.ServerNew(dbusConn)
	if err != nil {
		return nil, err
	}

	if _, err := avahiServer.GetAPIVersion(); err != nil {
		return nil, err
	}

	return avahiServer, nil
}

// Return allowed interfaces for mDNS
func (m *mdns) interfaces() ([]net.Interface, []int32, error) {
	var ifaces []net.Interface
	var ifaceIndexes []int32

	if len(m.serviceDescription.Interfaces) > 0 {
		ifaces = make([]net.Interface, len(m.serviceDescription.Interfaces))
		ifaceIndexes = make([]int32, len(m.serviceDescription.Interfaces))
		for i, ifaceName := range m.serviceDescription.Interfaces {
			iface, err := net.InterfaceByName(ifaceName)
			if err != nil {
				return nil, nil, err
			}
			ifaces[i] = *iface
			ifaceIndexes[i] = int32(iface.Index)
		}
	}

	if len(ifaces) == 0 {
		ifaces = nil
		ifaceIndexes = []int32{avahi.InterfaceUnspec}
	}

	return ifaces, ifaceIndexes, nil
}

// Announces the service to the network via mDNS
// A CEM service should always invoke this on startup
// Any other service should only invoke this whenever it is not connected to a CEM service
func (m *mdns) Announce() error {
	if m.isAnnounced {
		return nil
	}

	ifaces, ifaceIndexes, err := m.interfaces()
	if err != nil {
		return err
	}

	serviceIdentifier := fmt.Sprintf("%s-%s-%s", m.serviceDescription.Brand, m.serviceDescription.Model, m.serviceDescription.SerialNumber)
	if len(m.serviceDescription.Identifier) > 0 {
		serviceIdentifier = m.serviceDescription.Identifier
	}

	txt := []string{ // SHIP 7.3.2
		"txtvers=1",
		"path=" + shipWebsocketPath,
		"id=" + serviceIdentifier,
		"ski=" + m.ski,
		"brand=" + m.serviceDescription.Brand,
		"model=" + m.serviceDescription.Model,
		"type=" + string(m.serviceDescription.DeviceType),
		"register=" + fmt.Sprintf("%v", m.serviceDescription.RegisterAutoAccept),
	}

	fmt.Println("mDNS: Announce")

	if m.av == nil {
		// use Zeroconf library if avahi is not available
		if mDNSServer, err := zeroconf.Register(serviceIdentifier, shipZeroConfServiceType, shipZeroConfDomain, m.serviceDescription.Port, txt, ifaces); err == nil {
			m.zc = mDNSServer

			m.isAnnounced = true
			return nil
		}

		return err
	}

	// avahi
	entryGroup, err := m.av.EntryGroupNew()
	if err != nil {
		return err
	}

	var btxt [][]byte
	for _, t := range txt {
		btxt = append(btxt, []byte(t))
	}

	for _, iface := range ifaceIndexes {
		err = entryGroup.AddService(iface, avahi.ProtoUnspec, 0, serviceIdentifier, shipZeroConfServiceType, shipZeroConfDomain, "", uint16(m.serviceDescription.Port), btxt)
		if err != nil {
			return err
		}
	}

	err = entryGroup.Commit()
	if err != nil {
		return err
	}

	m.avEntryGroup = entryGroup
	m.isAnnounced = true

	return nil
}

// Stop the mDNS announcement on the network
func (m *mdns) Unannounce() {
	if !m.isAnnounced {
		return
	}

	if m.zc != nil {
		m.zc.Shutdown()
		m.zc = nil
	}
	if m.av != nil {
		m.av.EntryGroupFree(m.avEntryGroup)
		m.avEntryGroup = nil
	}
	fmt.Println("mDNS: Stop announcement")

	m.isAnnounced = false
}

// Shutdown all of mDNS
func (m *mdns) shutdown() {
	m.Unannounce()

	if m.av != nil {
		m.av.Close()
		m.av = nil
	}

	m.stopResolvingEntries()
}

// Register a callback to be invoked for found mDNS entries
func (m *mdns) RegisterMdnsSearch(cb MdnsSearch) {
	// check if callback is already registered
	registered := false
	for _, c := range m.searchDelegates {
		if c == cb {
			registered = true
		}
	}

	if !registered {
		m.searchDelegates = append(m.searchDelegates, cb)
	}

	m.mux.Lock()
	defer m.mux.Unlock()

	if !m.isSearchingServices {
		fmt.Println("mDNS: Start search")
		go m.resolveEntries()
	}
}

// Remove a callback for found mDNS entries and stop searching if no callbacks are left
func (m *mdns) UnregisterMdnsSearch(cb MdnsSearch) {
	var delegates []MdnsSearch

	for _, item := range m.searchDelegates {
		if item != cb {
			delegates = append(delegates, cb)
		}
	}

	m.searchDelegates = delegates

	if len(m.searchDelegates) == 0 {
		m.stopResolvingEntries()
	}
}

// search for mDNS entries and report them
// to be invoked in a background thread!
func (m *mdns) resolveEntries() {
	// for Zeroconf we need a context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var err error
	var avBrowser *avahi.ServiceBrowser

	zcEntries := make(chan *zeroconf.ServiceEntry)
	defer close(zcEntries)

	if m.av != nil {
		// instead of limiting search on specific allowed interfaces, we allow all and filter the results
		if avBrowser, err = m.av.ServiceBrowserNew(avahi.InterfaceUnspec, avahi.ProtoUnspec, shipZeroConfServiceType, shipZeroConfDomain, 0); err != nil {
			return
		}
	} else {
		go func() {
			_ = zeroconf.Browse(ctx, shipZeroConfServiceType, shipZeroConfDomain, zcEntries)
		}()
	}

	m.mux.Lock()
	m.isSearchingServices = true
	m.mux.Unlock()

	var end bool
	for !end {
		if m.av != nil {
			select {
			case <-m.cancelChan:
				end = true
				break
			case service := <-avBrowser.AddChannel:
				m.processAvahiService(service, false)
			case service := <-avBrowser.RemoveChannel:
				m.processAvahiService(service, true)
			}
		} else {
			select {
			case <-ctx.Done():
				end = true
				break
			case <-m.cancelChan:
				ctx.Done()
			case service := <-zcEntries:
				// Zeroconf has issues with merging mDNS data and sometimes reports non complety records
				if len(service.Text) == 0 {
					continue
				}

				elements := m.parseTxt(service.Text)

				addresses := service.AddrIPv4
				// Only use IPv4 for now
				// addresses = append(addresses, service.AddrIPv6...)
				m.processMdnsEntry(elements, service.Instance, service.HostName, addresses, service.Port, false)
			}
		}
	}

	if m.av != nil {
		m.av.ServiceBrowserFree(avBrowser)
	}

	m.mux.Lock()
	m.isSearchingServices = false
	m.mux.Unlock()
}

// stop searching for mDNS entries
func (m *mdns) stopResolvingEntries() {
	if m.cancelChan != nil {
		fmt.Println("mDNS: stop search")

		m.cancelChan <- true
	}
}

// process an avahi mDNS service
// as avahi returns a service per interface, we need to combine them
func (m *mdns) processAvahiService(service avahi.Service, remove bool) {
	_, ifaceIndexes, err := m.interfaces()
	if err != nil {
		fmt.Println("Error getting interfaces:", err)
		return
	}

	// check if the service is within the allowed list
	allow := false
	if len(ifaceIndexes) == 1 && ifaceIndexes[0] == avahi.InterfaceUnspec {
		allow = true
	} else {
		for _, iface := range ifaceIndexes {
			if service.Interface == iface {
				allow = true
				break
			}
		}
	}

	if !allow {
		return
	}

	resolved, err := m.av.ResolveService(service.Interface, service.Protocol, service.Name, service.Type, service.Domain, avahi.ProtoUnspec, 0)
	if err != nil {
		return
	}

	// convert [][]byte to []string manually
	var txt []string
	for _, element := range resolved.Txt {
		txt = append(txt, string(element))
	}
	elements := m.parseTxt(txt)

	// convert address to net.IP
	address := net.ParseIP(resolved.Address)
	// if the address can not be used, ignore the entry
	if address == nil || address.IsUnspecified() {
		return
	}

	// Ignore IPv6 addresses for now
	if address.To4() == nil {
		return
	}

	m.processMdnsEntry(elements, resolved.Name, resolved.Host, []net.IP{address}, int(resolved.Port), remove)
}

// parse mDNS text fields
func (m *mdns) parseTxt(txt []string) map[string]string {
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

// process an mDNS entry and manage mDNS entries map
func (m *mdns) processMdnsEntry(elements map[string]string, name, host string, addresses []net.IP, port int, remove bool) {
	// check for mandatory text elements
	fmt.Println(fmt.Sprintf("mDNS: Processing entry: %s on %s",name, host))

	mapItems := []string{"txtvers", "id", "path", "ski", "register"}
	for _, item := range mapItems {
		if _, ok := elements[item]; !ok {
			return
		}
	}

	txtvers := elements["txtvers"]
	// value of mandatory txtvers has to be 1 or the response be ignored: SHIP 7.3.2
	if txtvers != "1" {
		return
	}

	identifier := elements["id"]
	path := elements["path"]
	ski := elements["ski"]

	// ignore own service
	if ski == m.ski {
		return
	}

	register := elements["register"]
	// register has to be a boolean
	if register != "true" && register != "false" {
		return
	}

	var deviceType, model, brand string

	if _, ok := elements["brand"]; ok {
		brand = elements["brand"]
	}
	if _, ok := elements["type"]; ok {
		deviceType = elements["type"]
	}
	if _, ok := elements["model"]; ok {
		model = elements["model"]
	}

	_, exists := m.entries[ski]

	if remove && exists {
		// remove
		// there will be a remove for each address with avahi, but we'll delete it right away
		delete(m.entries, ski)
	} else if exists {
		// update
		// avahi sends an item for each network address, merge them
		entry := m.entries[ski]

		// we assume only network addresses are added
		entry.Addresses = append(entry.Addresses, addresses...)

		m.entries[ski] = entry
	} else if !exists && !remove {
		// new
		newEntry := MdnsEntry{
			Name:       name,
			Identifier: identifier,
			Path:       path,
			Register:   register == "true",
			Brand:      brand,
			Type:       deviceType,
			Model:      model,
			Host:       host,
			Port:       port,
			Addresses:  addresses,
		}
		m.entries[ski] = newEntry

		fmt.Println("SKI:", ski, "Name:", name, "Brand:", brand, "Model:", model, "Typ:", deviceType, "Identifier:", identifier, "Register:", register)
	} else {
		return
	}

	for _, cb := range m.searchDelegates {
		cb.ReportMdnsEntries(m.entries)
	}
}
