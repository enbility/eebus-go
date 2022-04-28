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

type hems struct {
	myService *service.EEBUSService
}

func (h *hems) run() {

	serviceDescription := &service.ServiceDescription{
		DeviceBrand:        "Demo",
		DeviceModel:        "HEMS",
		DeviceSerialNumber: "123456789",
		DeviceIdentifier:   "Demo-HEMS-123456789",
		DeviceType:         model.DeviceTypeTypeEnergyManagementSystem,
		RemoteDeviceTypes: []model.DeviceTypeType{
			model.DeviceTypeTypeChargingStation,
		},
	}

	h.myService = service.NewEEBUSService(serviceDescription, h)

	var err error
	var certificate tls.Certificate

	serviceDescription.Port, err = strconv.Atoi(os.Args[1])
	if err != nil {
		usage()
		log.Fatal(err)
	}

	remoteSki := os.Args[2]

	if len(os.Args) == 5 {
		certificate, err = tls.LoadX509KeyPair(os.Args[3], os.Args[4])
		if err != nil {
			usage()
			log.Fatal(err)
		}
	} else {
		certificate, err = service.CreateCertificate("Demo", "Demo", "DE", "Demo-Unit-01")
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

	h.myService.Start()
	// defer h.myService.Shutdown()

	remoteService := service.ServiceDetails{
		SKI: remoteSki,
	}
	h.myService.RegisterRemoteService(remoteService)
}

var _ service.EEBUSServiceDelegate = (*hems)(nil)

// handle a request to trust a remote service
func (h *hems) RemoteServiceTrustRequested(ski string) {
	// we directly trust it in this example
	h.myService.UpdateRemoteServiceTrust(ski, true)
}

// report the Ship ID of a newly trusted connection
func (h *hems) RemoteServiceShipIDReported(ski string, shipID string) {
	// we should associated the Ship ID with the SKI and store it
	// so the next connection can start trusted
	fmt.Println("SKI", ski, "has Ship ID:", shipID)
}

func usage() {
	fmt.Println("Usage: go run /cmd/hems/main.go <serverport> <evse-ski> <crtfile> <keyfile>")
}

func main() {
	if len(os.Args) < 4 {
		usage()
		return
	}

	h := hems{}
	h.run()

	// Clean exit to make sure mdns shutdown is invoked
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	select {
	case <-sig:
		// User exit
	}
}
