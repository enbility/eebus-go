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

	"github.com/enbility/eebus-go/features"
	"github.com/enbility/eebus-go/logging"
	"github.com/enbility/eebus-go/service"
	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
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

	configuration, err := service.NewConfiguration(
		"Demo", "Demo", "HEMS", "123456789",
		model.DeviceTypeTypeEnergyManagementSystem, port, certificate, 230)
	if err != nil {
		log.Fatal(err)
	}
	configuration.SetAlternateIdentifier("Demo-HEMS-123456789")

	h.myService = service.NewEEBUSService(configuration, h)
	h.myService.SetLogging(h)

	if err = h.myService.Setup(); err != nil {
		fmt.Println(err)
		return
	}

	if len(remoteSki) == 0 {
		os.Exit(0)
	}

	spine.Events.Subscribe(h)

	h.myService.Start()
	// defer h.myService.Shutdown()

	remoteService := service.NewServiceDetails(remoteSki)
	h.myService.PairRemoteService(remoteService)
}

// EEBUSServiceHandler

func (h *hems) RemoteSKIConnected(service *service.EEBUSService, ski string) {}

func (h *hems) RemoteSKIDisconnected(service *service.EEBUSService, ski string) {}

func (h *hems) ReportServiceShipID(ski string, shipdID string) {}

// UCEvseCommisioningConfigurationCemDelegate

// handle device state updates from the remote EVSE device
func (h *hems) HandleEVSEDeviceState(ski string, failure bool, errorCode string) {
	fmt.Println("EVSE Error State:", failure, errorCode)
}

func (h *hems) HandleEvent(payload spine.EventPayload) {
	if payload.Entity != nil {
		entityType := payload.Entity.EntityType()
		if entityType != model.EntityTypeTypeGeneric {
			return
		}
	}

	switch payload.EventType {
	case spine.EventTypeEntityChange:
		switch payload.ChangeType {
		case spine.ElementChangeAdd:
			h.heatpumpConnected(payload.Entity)
		}
	}
}

func (h *hems) heatpumpConnected(entity *spine.EntityRemoteImpl) {
	localDevice := h.myService.LocalDevice()

	deviceClassification, _ := features.NewDeviceClassification(model.RoleTypeClient, model.RoleTypeServer, localDevice, entity.Device().Entity([]model.AddressEntityType{0}))
	electricalConnection, _ := features.NewElectricalConnection(model.RoleTypeClient, model.RoleTypeClient, localDevice, entity)
	measurement, _ := features.NewMeasurement(model.RoleTypeClient, model.RoleTypeClient, localDevice, entity)
	hvac, _ := features.NewHVAC(model.RoleTypeClient, model.RoleTypeClient, localDevice, entity)

	if deviceClassification != nil {
		if _, err := deviceClassification.RequestManufacturerDetails(); err != nil {
			logging.Log.Debug(err)
		}
	}

	if electricalConnection != nil {
		if err := electricalConnection.SubscribeForEntity(); err != nil {
			logging.Log.Error(err)
		}

		if err := electricalConnection.RequestDescriptions(); err != nil {
			logging.Log.Error(err)
		}

		if err := electricalConnection.RequestParameterDescriptions(); err != nil {
			logging.Log.Error(err)
		}
	}

	if measurement != nil {
		if err := measurement.SubscribeForEntity(); err != nil {
			logging.Log.Error(err)
		}

		if err := measurement.RequestDescriptions(); err != nil {
			logging.Log.Error(err)
		}

		if err := measurement.RequestConstraints(); err != nil {
			logging.Log.Error(err)
		}

		if _, err := measurement.RequestValues(); err != nil {
			logging.Log.Error(err)
		}
	}

	if hvac != nil {
		if err := hvac.SubscribeForEntity(); err != nil {
			logging.Log.Error(err)
		}

		if _, err := hvac.RequestOverrunDescriptions(); err != nil {
			logging.Log.Error(err)
		}

		if _, err := hvac.RequestOverrunValues(); err != nil {
			logging.Log.Error(err)
		}

		if _, err := hvac.RequestSystemFunctionDescriptions(); err != nil {
			logging.Log.Error(err)
		}
	}
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
