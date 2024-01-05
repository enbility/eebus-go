package mdns

import (
	"context"
	"net"

	"github.com/DerAndereAndi/zeroconf/v2"
	"github.com/enbility/eebus-go/logging"
)

type ZeroconfProvider struct {
	ifaces []net.Interface

	zc *zeroconf.Server
}

func NewZeroconfProvider(ifaces []net.Interface) *ZeroconfProvider {
	return &ZeroconfProvider{
		ifaces: ifaces,
	}
}

var _ MdnsProvider = (*ZeroconfProvider)(nil)

func (z *ZeroconfProvider) CheckAvailability() bool {
	return true
}

func (z *ZeroconfProvider) Shutdown() {}

func (z *ZeroconfProvider) Announce(serviceName string, port int, txt []string) error {
	logging.Log().Debug("mdns: using zeroconf")

	// use Zeroconf library if avahi is not available
	// Set TTL to 2 minutes as defined in SHIP chapter 7
	mDNSServer, err := zeroconf.Register(serviceName, shipZeroConfServiceType, shipZeroConfDomain, port, txt, z.ifaces, zeroconf.TTL(120))
	if err != nil {
		return err
	}

	z.zc = mDNSServer

	return nil
}

func (z *ZeroconfProvider) Unannounce() {
	if z.zc == nil {
		return
	}

	z.zc.Shutdown()
	z.zc = nil
}

func (z *ZeroconfProvider) ResolveEntries(cancelChan chan bool, callback func(elements map[string]string, name, host string, addresses []net.IP, port int, remove bool)) {
	var end bool

	zcEntries := make(chan *zeroconf.ServiceEntry)
	zcRemoved := make(chan *zeroconf.ServiceEntry)
	defer close(zcEntries)
	defer close(zcRemoved)

	// for Zeroconf we need a context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		_ = zeroconf.Browse(ctx, shipZeroConfServiceType, shipZeroConfDomain, zcEntries, zcRemoved)
	}()

	for !end {
		select {
		case <-ctx.Done():
			end = true
			break
		case <-cancelChan:
			ctx.Done()
		case service := <-zcRemoved:
			// Zeroconf has issues with merging mDNS data and sometimes reports incomplete records
			if len(service.Text) == 0 {
				continue
			}

			elements := parseTxt(service.Text)

			addresses := service.AddrIPv4
			callback(elements, service.Instance, service.HostName, addresses, service.Port, true)

		case service := <-zcEntries:
			// Zeroconf has issues with merging mDNS data and sometimes reports incomplete records
			if len(service.Text) == 0 {
				continue
			}

			elements := parseTxt(service.Text)

			addresses := service.AddrIPv4
			// Only use IPv4 for now
			// addresses = append(addresses, service.AddrIPv6...)
			callback(elements, service.Instance, service.HostName, addresses, service.Port, false)
		}
	}
}
