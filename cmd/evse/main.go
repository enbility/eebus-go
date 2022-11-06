package main

import (
	"crypto/ecdsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/DerAndereAndi/eebus-go/service"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

type evse struct {
	myService *service.EEBUSService
}

func (h *evse) run() {

	serviceDescription := &service.ServiceDescription{
		Brand:        "Demo",
		Model:        "EVSE",
		SerialNumber: "234567890",
		Identifier:   "Demo-EVSE-234567890",
		DeviceType:   model.DeviceTypeTypeChargingStation,
	}

	h.myService = service.NewEEBUSService(serviceDescription, h)
	h.myService.SetLogging(h)

	var err error
	var certificate tls.Certificate
	var remoteSki string

	serviceDescription.Port, err = strconv.Atoi(os.Args[1])
	if err != nil {
		usage()
		log.Fatal(err)
	}

	fmt.Println(os.Args)
	if len(os.Args) == 5 {
		remoteSki = os.Args[2]

		certificate, err = tls.LoadX509KeyPair(os.Args[3], os.Args[4])
		if err != nil {
			usage()
			log.Fatal(err)
		}
	} else {
		certificate, err = service.CreateCertificate("Demo", "Demo", "DE", "Demo-Unit-02")
		if err != nil {
			log.Fatal(err)
		}

		pemdata := pem.EncodeToMemory(&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: certificate.Certificate[0],
		})
		fmt.Println(string(pemdata))

		b, err := x509.MarshalECPrivateKey(certificate.PrivateKey.(*ecdsa.PrivateKey))
		if err != nil {
			log.Fatal(err)
		}
		pemdata = pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: b})
		fmt.Println(string(pemdata))
	}

	serviceDescription.Certificate = certificate

	if err = h.myService.Setup(); err != nil {
		fmt.Println(err)
		return
	}

	if len(remoteSki) == 0 {
		os.Exit(0)
	}

	h.myService.Start()
	// defer h.myService.Shutdown()

	remoteService := service.ServiceDetails{
		SKI: remoteSki,
	}
	h.myService.RegisterRemoteService(remoteService)
}

// handle a request to trust a remote service
func (h *evse) RemoteServiceTrustRequested(ski string) {
	// we directly trust it in this example
	h.myService.UpdateRemoteServiceTrust(ski, true)
}

// report the Ship ID of a newly trusted connection
func (h *evse) RemoteServiceShipIDReported(ski string, shipID string) {
	// we should associated the Ship ID with the SKI and store it
	// so the next connection can start trusted
	fmt.Println("SKI", ski, "has Ship ID:", shipID)
}

func (h *evse) RemoteSKIConnected(ski string) {}

func (h *evse) RemoteSKIDisconnected(ski string) {}

// main app
func usage() {
	fmt.Println("First Run:")
	fmt.Println("  go run /cmd/evse/main.go <serverport>")
	fmt.Println()
	fmt.Println("General Usage:")
	fmt.Println("  go run /cmd/evse/main.go <serverport> <hems-ski> <crtfile> <keyfile>")
}

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}

	h := evse{}
	h.run()

	// Clean exit to make sure mdns shutdown is invoked
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	// User exit
}

// Logging interface

func (h *evse) Trace(args ...interface{}) {
	fmt.Println(args...)
}

func (h *evse) Tracef(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (h *evse) Debug(args ...interface{}) {
	fmt.Println(args...)
}

func (h *evse) Debugf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (h *evse) Info(args ...interface{}) {
	fmt.Println(args...)
}

func (h *evse) Infof(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (h *evse) Error(args ...interface{}) {
	fmt.Println(args...)
}

func (h *evse) Errorf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}
