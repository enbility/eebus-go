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
	"time"

	"github.com/DerAndereAndi/eebus-go/service"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

type hems struct {
	myService *service.EEBUSService
}

func (h *hems) run() {
	var err error
	var certificate tls.Certificate
	var remoteSki string

	if len(os.Args) == 5 {
		remoteSki = os.Args[2]

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

	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		usage()
		log.Fatal(err)
	}

	serviceDescription, err := service.NewServiceDescription(
		"Demo", "Demo", "HEMS", "123456789",
		model.DeviceTypeTypeEnergyManagementSystem, port, certificate)
	if err != nil {
		log.Fatal(err)
	}
	serviceDescription.SetAlternateIdentifier("Demo-HEMS-123456789")

	h.myService = service.NewEEBUSService(serviceDescription, h)
	h.myService.SetLogging(h)

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
	h.myService.PairRemoteService(remoteService)
}

// EEBUSServiceDelegate

// // handle a request to trust a remote service
// func (h *hems) RemoteServiceTrustRequested(service *service.EEBUSService, ski string) {
// 	// we directly trust it in this example
// 	h.myService.UpdateRemoteServiceTrust(ski, true)
// }

// report the Ship ID of a newly trusted connection
func (h *hems) RemoteServiceShipIDReported(service *service.EEBUSService, ski string, shipID string) {
	// we should associated the Ship ID with the SKI and store it
	// so the next connection can start trusted
	fmt.Println("SKI", ski, "has Ship ID:", shipID)
}

func (h *hems) RemoteSKIConnected(service *service.EEBUSService, ski string) {}

func (h *hems) RemoteSKIDisconnected(service *service.EEBUSService, ski string) {}

// UCEvseCommisioningConfigurationCemDelegate

// handle device state updates from the remote EVSE device
func (h *hems) HandleEVSEDeviceState(ski string, failure bool, errorCode string) {
	fmt.Println("EVSE Error State:", failure, errorCode)
}

// main app
func usage() {
	fmt.Println("First Run:")
	fmt.Println("  go run /cmd/hems/main.go <serverport>")
	fmt.Println()
	fmt.Println("General Usage:")
	fmt.Println("  go run /cmd/hems/main.go <serverport> <evse-ski> <crtfile> <keyfile>")
}

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}

	h := hems{}
	h.run()

	// Clean exit to make sure mdns shutdown is invoked
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	// User exit
}

// Logging interface

func (h *hems) Trace(args ...interface{}) {
	h.print("TRACE", args...)
}

func (h *hems) Tracef(format string, args ...interface{}) {
	h.printFormat("TRACE", format, args...)
}

func (h *hems) Debug(args ...interface{}) {
	h.print("DEBUG", args...)
}

func (h *hems) Debugf(format string, args ...interface{}) {
	h.printFormat("DEBUG", format, args...)
}

func (h *hems) Info(args ...interface{}) {
	h.print("INFO ", args...)
}

func (h *hems) Infof(format string, args ...interface{}) {
	h.printFormat("INFO ", format, args...)
}

func (h *hems) Error(args ...interface{}) {
	h.print("ERROR", args...)
}

func (h *hems) Errorf(format string, args ...interface{}) {
	h.printFormat("ERROR", format, args...)
}

func (h *hems) currentTimestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func (h *hems) print(msgType string, args ...interface{}) {
	value := fmt.Sprintln(args...)
	fmt.Printf("%s %s %s", h.currentTimestamp(), msgType, value)
}

func (h *hems) printFormat(msgType, format string, args ...interface{}) {
	value := fmt.Sprintf(format, args...)
	fmt.Println(h.currentTimestamp(), msgType, value)
}
