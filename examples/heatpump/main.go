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

	"github.com/enbility/eebus-go/api"
	features "github.com/enbility/eebus-go/features/client"
	"github.com/enbility/eebus-go/service"
	shipapi "github.com/enbility/ship-go/api"
	"github.com/enbility/ship-go/cert"
	"github.com/enbility/ship-go/logging"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

var remoteSki string

type heatpump struct {
	myService *service.Service
}

func (h *heatpump) run() {
	var err error
	var certificate tls.Certificate

	if len(os.Args) == 5 {
		remoteSki = os.Args[2]

		certificate, err = tls.LoadX509KeyPair(os.Args[3], os.Args[4])
		if err != nil {
			usage()
			log.Fatal(err)
		}
	} else {
		certificate, err = cert.CreateCertificate("Demo", "Demo", "DE", "Demo-Unit-02")
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

	configuration, err := api.NewConfiguration(
		"Demo", "Demo", "HeatPump", "234567890",
		[]shipapi.DeviceCategoryType{shipapi.DeviceCategoryTypeHVAC},
		model.DeviceTypeTypeGeneric,
		[]model.EntityTypeType{model.EntityTypeTypeHeatPumpAppliance},
		port, certificate, time.Second*4)
	if err != nil {
		log.Fatal(err)
	}
	configuration.SetAlternateIdentifier("Demo-EVSE-234567890")

	h.myService = service.NewService(configuration, h)
	h.myService.SetLogging(h)

	if err = h.myService.Setup(); err != nil {
		fmt.Println(err)
		return
	}

	if len(remoteSki) == 0 {
		os.Exit(0)
	}

	h.AddFeatures()
	_ = spine.Events.Subscribe(h)

	h.myService.RegisterRemoteSKI(remoteSki, "")

	h.myService.Start()
	// defer h.myService.Shutdown()
}

// EEBUSServiceHandler

func (h *heatpump) RemoteSKIConnected(service api.ServiceInterface, ski string) {}

func (h *heatpump) RemoteSKIDisconnected(service api.ServiceInterface, ski string) {}

func (h *heatpump) VisibleRemoteServicesUpdated(service api.ServiceInterface, entries []shipapi.RemoteService) {
}

func (h *heatpump) ServiceShipIDUpdate(ski string, shipdID string) {}

func (h *heatpump) ServicePairingDetailUpdate(ski string, detail *shipapi.ConnectionStateDetail) {
	if ski == remoteSki && detail.State() == shipapi.ConnectionStateRemoteDeniedTrust {
		fmt.Println("The remote service denied trust. Exiting.")
		h.myService.CancelPairingWithSKI(ski)
		h.myService.UnregisterRemoteSKI(ski)
		h.myService.Shutdown()
		os.Exit(0)
	}
}

func (h *heatpump) AllowWaitingForTrust(ski string) bool {
	return ski == remoteSki
}

func (h *heatpump) AddFeatures() {
	entityAddress := []model.AddressEntityType{1}
	localEntity := spine.NewEntityLocal(h.myService.LocalDevice(), model.EntityTypeTypeCEM, entityAddress, time.Second*4)
	h.myService.LocalDevice().AddEntity(localEntity)

	_ = localEntity.GetOrAddFeature(model.FeatureTypeTypeGeneric, model.RoleTypeServer)
	_ = localEntity.GetOrAddFeature(model.FeatureTypeTypeDeviceClassification, model.RoleTypeClient)
}

func (h *heatpump) HandleResult(errorMsg spineapi.ResponseMessage) {}

func (h *heatpump) HandleEvent(payload spineapi.EventPayload) {
	switch payload.EventType {
	case spineapi.EventTypeEntityChange:
		entityType := payload.Entity.EntityType()

		switch payload.ChangeType {
		case spineapi.ElementChangeAdd:
			switch entityType {
			case model.EntityTypeTypeGeneric:
				h.deviceInformationConnected(payload.Entity)
				h.hemsConnected(payload.Entity)
			}
		}
	}
}

func (h *heatpump) deviceInformationConnected(remoteEntity spineapi.EntityRemoteInterface) {
	localDevice := h.myService.LocalDevice()
	localEntity := localDevice.Entities()[1]

	nodeEntity := remoteEntity.Device().Entities()[0]

	deviceClassification, _ := features.NewDeviceClassification(localEntity, nodeEntity)

	if deviceClassification != nil {
		if _, err := deviceClassification.RequestManufacturerDetails(); err != nil {
			logging.Log().Debug(err)
		}
	}
}

func (h *heatpump) hemsConnected(entity spineapi.EntityRemoteInterface) {}

// main app
func usage() {
	fmt.Println("First Run:")
	fmt.Println("  go run /examples/heatpump/main.go <serverport>")
	fmt.Println()
	fmt.Println("General Usage:")
	fmt.Println("  go run /examples/heatpump/main.go <serverport> <hems-ski> <crtfile> <keyfile>")
}

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}

	h := heatpump{}
	h.run()

	// Clean exit to make sure mdns shutdown is invoked
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	// User exit
}

// Logging interface

func (h *heatpump) Trace(args ...interface{}) {
	h.print("TRACE", args...)
}

func (h *heatpump) Tracef(format string, args ...interface{}) {
	h.printFormat("TRACE", format, args...)
}

func (h *heatpump) Debug(args ...interface{}) {
	h.print("DEBUG", args...)
}

func (h *heatpump) Debugf(format string, args ...interface{}) {
	h.printFormat("DEBUG", format, args...)
}

func (h *heatpump) Info(args ...interface{}) {
	h.print("INFO ", args...)
}

func (h *heatpump) Infof(format string, args ...interface{}) {
	h.printFormat("INFO ", format, args...)
}

func (h *heatpump) Error(args ...interface{}) {
	h.print("ERROR", args...)
}

func (h *heatpump) Errorf(format string, args ...interface{}) {
	h.printFormat("ERROR", format, args...)
}

func (h *heatpump) currentTimestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func (h *heatpump) print(msgType string, args ...interface{}) {
	value := fmt.Sprintln(args...)
	fmt.Printf("%s %s %s", h.currentTimestamp(), msgType, value)
}

func (h *heatpump) printFormat(msgType, format string, args ...interface{}) {
	value := fmt.Sprintf(format, args...)
	fmt.Println(h.currentTimestamp(), msgType, value)
}
