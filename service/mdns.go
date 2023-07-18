package service

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/enbility/eebus-go/logging"
	"github.com/enbility/eebus-go/service/mdns"
	"github.com/enbility/eebus-go/util"
	"github.com/holoplot/go-avahi"
)

type MdnsEntry struct {
	Name       string
	Ski        string
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

//go:generate mockgen -destination=mock_mdns.go -package=service github.com/enbility/eebus-go/service MdnsSearch,MdnsService

// implemented by hubConnection, used by mdns
type MdnsSearch interface {
	ReportMdnsEntries(entries map[string]MdnsEntry)
}

// implemented by mdns, used by hubConnection
type MdnsService interface {
	SetupMdnsService() error
	ShutdownMdnsService()
	AnnounceMdnsEntry() error
	UnannounceMdnsEntry()
	RegisterMdnsSearch(cb MdnsSearch)
	UnregisterMdnsSearch(cb MdnsSearch)
}

type mdnsManager struct {
	configuration *Configuration
	ski           string

	isAnnounced         bool
	isSearchingServices bool

	cancelChan chan bool

	// the currently available mDNS entries with the SKI as the key in the map
	entries map[string]MdnsEntry

	// the registered callback, only connectionsHub is using this
	searchDelegate MdnsSearch

	mdnsProvider mdns.MdnsProvider

	mux sync.Mutex
}

func newMDNS(ski string, configuration *Configuration) *mdnsManager {
	m := &mdnsManager{
		ski:           ski,
		configuration: configuration,
		entries:       make(map[string]MdnsEntry),
		cancelChan:    make(chan bool),
	}

	return m
}

// Return allowed interfaces for mDNS
func (m *mdnsManager) interfaces() ([]net.Interface, []int32, error) {
	var ifaces []net.Interface
	var ifaceIndexes []int32

	if len(m.configuration.interfaces) > 0 {
		ifaces = make([]net.Interface, len(m.configuration.interfaces))
		ifaceIndexes = make([]int32, len(m.configuration.interfaces))
		for i, ifaceName := range m.configuration.interfaces {
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

var _ MdnsService = (*mdnsManager)(nil)

func (m *mdnsManager) SetupMdnsService() error {
	ifaces, ifaceIndexes, err := m.interfaces()
	if err != nil {
		return err
	}

	m.mdnsProvider = mdns.NewAvahiProvider(ifaceIndexes)
	if !m.mdnsProvider.CheckAvailability() {
		m.mdnsProvider.Shutdown()

		// Avahi is not availble, use Zeroconf
		m.mdnsProvider = mdns.NewZeroconfProvider(ifaces)
		if !m.mdnsProvider.CheckAvailability() {
			return errors.New("No mDNS provider available")
		}
	}

	// on startup always start mDNS announcement
	if err := m.AnnounceMdnsEntry(); err != nil {
		return err
	}

	// catch signals
	go func() {
		signalC := make(chan os.Signal, 1)
		signal.Notify(signalC, os.Interrupt, syscall.SIGTERM)

		<-signalC // wait for signal

		m.ShutdownMdnsService()
	}()

	return nil
}

// Announces the service to the network via mDNS
// A CEM service should always invoke this on startup
// Any other service should only invoke this whenever it is not connected to a CEM service
func (m *mdnsManager) AnnounceMdnsEntry() error {
	if m.isAnnounced {
		return nil
	}

	serviceIdentifier := m.configuration.Identifier()

	txt := []string{ // SHIP 7.3.2
		"txtvers=1",
		"path=" + shipWebsocketPath,
		"id=" + serviceIdentifier,
		"ski=" + m.ski,
		"brand=" + m.configuration.deviceBrand,
		"model=" + m.configuration.deviceModel,
		"type=" + string(m.configuration.deviceType),
		"register=" + fmt.Sprintf("%v", m.configuration.registerAutoAccept),
	}

	logging.Log.Debug("mdns: announce")

	serviceName := m.configuration.MdnsServiceName()

	if err := m.mdnsProvider.Announce(serviceName, m.configuration.port, txt); err != nil {
		logging.Log.Debug("mdns: failure announcing service", err)
		return err
	}

	m.isAnnounced = true

	return nil
}

// Stop the mDNS announcement on the network
func (m *mdnsManager) UnannounceMdnsEntry() {
	if !m.isAnnounced {
		return
	}

	m.mdnsProvider.Unannounce()
	logging.Log.Debug("mdns: stop announcement")

	m.isAnnounced = false
}

// Shutdown all of mDNS
func (m *mdnsManager) ShutdownMdnsService() {
	m.UnannounceMdnsEntry()
	m.stopResolvingEntries()

	m.mdnsProvider.Shutdown()
	m.mdnsProvider = nil
}

// Register a callback to be invoked for found mDNS entries
func (m *mdnsManager) RegisterMdnsSearch(cb MdnsSearch) {
	if m.searchDelegate != cb {
		m.searchDelegate = cb
	}

	m.mux.Lock()

	if !m.isSearchingServices {
		m.isSearchingServices = true
		m.mux.Unlock()
		logging.Log.Debug("mdns: start search")
		m.resolveEntries()
		return
	}

	defer m.mux.Unlock()

	// do we already know some entries?
	if len(m.entries) == 0 {
		return
	}

	// maybe entries are already found
	mdnsEntries := m.entries

	go m.searchDelegate.ReportMdnsEntries(mdnsEntries)
}

// Remove a callback for found mDNS entries and stop searching if no callbacks are left
func (m *mdnsManager) UnregisterMdnsSearch(cb MdnsSearch) {
	m.mux.Lock()
	defer m.mux.Unlock()

	m.searchDelegate = nil

	m.stopResolvingEntries()
}

// search for mDNS entries and report them
func (m *mdnsManager) resolveEntries() {
	if m.mdnsProvider == nil {
		m.isSearchingServices = false
		return
	}
	go func() {
		m.mdnsProvider.ResolveEntries(m.cancelChan, m.processMdnsEntry)

		m.mux.Lock()
		m.isSearchingServices = false
		m.mux.Unlock()
	}()
}

// stop searching for mDNS entries
func (m *mdnsManager) stopResolvingEntries() {
	if m.cancelChan == nil {
		return
	}

	if util.IsChannelClosed(m.cancelChan) {
		return
	}

	logging.Log.Debug("mdns: stop search")

	m.cancelChan <- true
}

// process an mDNS entry and manage mDNS entries map
func (m *mdnsManager) processMdnsEntry(elements map[string]string, name, host string, addresses []net.IP, port int, remove bool) {
	// check for mandatory text elements
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

	m.mux.Lock()
	defer m.mux.Unlock()

	updated := true

	_, exists := m.entries[ski]

	if remove && exists {
		// remove
		// there will be a remove for each address with avahi, but we'll delete it right away
		delete(m.entries, ski)
	} else if exists {
		// update
		updated = false

		// avahi sends an item for each network address, merge them
		entry := m.entries[ski]

		// we assume only network addresses are added
		for _, address := range addresses {
			// only add if it is not added yet
			isNewElement := true

			for _, item := range entry.Addresses {
				if item.String() == address.String() {
					isNewElement = false
					break
				}
			}

			if isNewElement {
				entry.Addresses = append(entry.Addresses, address)
				updated = true
			}
		}

		m.entries[ski] = entry
	} else if !exists && !remove {
		// new
		newEntry := MdnsEntry{
			Name:       name,
			Ski:        ski,
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

		logging.Log.Debug("ski:", ski, "name:", name, "brand:", brand, "model:", model, "typ:", deviceType, "identifier:", identifier, "register:", register, "host:", host, "port:", port, "addresses:", addresses)
	} else {
		return
	}

	if m.searchDelegate != nil && updated {
		mdnsEntries := m.entries
		go m.searchDelegate.ReportMdnsEntries(mdnsEntries)
	}
}
