package main

import (
	"crypto/ecdsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/service"
	ucapi "github.com/enbility/eebus-go/usecases/api"

	// "github.com/enbility/eebus-go/usecases/cem/vabd"
	// "github.com/enbility/eebus-go/usecases/cem/vapd"
	gcpmgcp "github.com/enbility/eebus-go/usecases/gcp/mgcp"
	mumpc "github.com/enbility/eebus-go/usecases/mu/mpc"
	shipapi "github.com/enbility/ship-go/api"
	"github.com/enbility/ship-go/cert"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

var remoteSki string

type smartmeter struct {
	myService *service.Service

	ucgcpmgcp ucapi.GcpMGCPInterface
	ucmumpc   ucapi.MuMPCInterface

	// uccemvabd ucapi.CemVABDInterface
	// uccemvapd ucapi.CemVAPDInterface

	gridPowerLimitFactor float64
	gridPower            float64
	gridPowerPerPhase    []float64
	gridConsumedEnergy   float64
	gridFeedInEnergy     float64
	gridCurrentPerPhase  []float64
	gridVoltagePerPhase  []float64
	gridFrequency        float64
}

func (h *smartmeter) run() {
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
		certificate, err = cert.CreateCertificate("Demo", "Demo", "DE", "Demo-Unit-01")
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
		"Demo", "Demo", "SmartMeter", "123456789",
		[]shipapi.DeviceCategoryType{shipapi.DeviceCategoryTypeEnergyManagementSystem},
		model.DeviceTypeTypeEnergyManagementSystem,
		[]model.EntityTypeType{model.EntityTypeTypeSubMeterElectricity},
		port, certificate, time.Second*4)
	if err != nil {
		log.Fatal(err)
	}
	configuration.SetAlternateIdentifier("Demo-SmartMeter-123456789")

	h.myService = service.NewService(configuration, h)
	h.myService.SetLogging(h)

	if err = h.myService.Setup(); err != nil {
		fmt.Println(err)
		return
	}

	localEntityCEM := h.myService.LocalDevice().EntityForType(model.EntityTypeTypeCEM)
	h.ucgcpmgcp = gcpmgcp.NewMGCP(localEntityCEM, h.OnMGCPEvent)
	h.myService.AddUseCase(h.ucgcpmgcp)
	// h.uccemvabd = vabd.NewVABD(localEntityCEM, h.OnVABDEvent)
	// h.myService.AddUseCase(h.uccemvabd)
	// h.uccemvapd = vapd.NewVAPD(localEntityCEM, h.OnVAPDEvent)
	// h.myService.AddUseCase(h.uccemvapd)

	localEntitySME := h.myService.LocalDevice().EntityForType(model.EntityTypeTypeSubMeterElectricity)
	h.ucmumpc = mumpc.NewMPC(localEntitySME, h.OnMPCEvent)
	h.myService.AddUseCase(h.ucmumpc)

	//_ = h.ucmamgcp.
	if len(remoteSki) == 0 {
		os.Exit(0)
	}

	h.gridPowerLimitFactor = 70
	h.gridPower = 3000
	h.gridPowerPerPhase = []float64{900, 1000, 1100}
	h.gridConsumedEnergy = 12345
	h.gridFeedInEnergy = -1000
	h.gridCurrentPerPhase = []float64{10, 20, 30}
	h.gridVoltagePerPhase = []float64{229, 230, 231}
	h.gridFrequency = 50

	ticker := time.NewTicker(3000 * time.Millisecond)
	go func() {
		for range ticker.C {
			_ = h.ucgcpmgcp.SetPower(math.Abs(h.gridPower))
			_ = h.ucmumpc.SetPower(math.Abs(h.gridPower))
			switch h.gridPower {
			case 3000:
				h.gridPower = 3100
			case 3100:
				h.gridPower = -3000
			case -3000:
				h.gridPower = 2900
			case 2900:
				h.gridPower = 3000
			}

			_ = h.ucmumpc.SetPowerPerPhase(h.gridPowerPerPhase)
			tmp := h.gridPowerPerPhase[0]
			h.gridPowerPerPhase[0] = h.gridPowerPerPhase[1]
			h.gridPowerPerPhase[1] = h.gridPowerPerPhase[2]
			h.gridPowerPerPhase[2] = tmp

			_ = h.ucgcpmgcp.SetEnergyConsumed(h.gridConsumedEnergy)
			_ = h.ucmumpc.SetEnergyConsumed(h.gridConsumedEnergy)
			h.gridConsumedEnergy++

			_ = h.ucgcpmgcp.SetEnergyFeedIn(h.gridFeedInEnergy)
			_ = h.ucmumpc.SetEnergyProduced(h.gridFeedInEnergy)
			h.gridFeedInEnergy--

			_ = h.ucgcpmgcp.SetCurrentPerPhase(h.gridCurrentPerPhase)
			_ = h.ucmumpc.SetCurrentPerPhase(h.gridCurrentPerPhase)
			tmp = h.gridCurrentPerPhase[0]
			h.gridCurrentPerPhase[0] = h.gridCurrentPerPhase[1]
			h.gridCurrentPerPhase[1] = h.gridCurrentPerPhase[2]
			h.gridCurrentPerPhase[2] = tmp

			_ = h.ucgcpmgcp.SetVoltagePerPhase(h.gridVoltagePerPhase)
			_ = h.ucmumpc.SetVoltagePerPhase(h.gridVoltagePerPhase)
			tmp = h.gridVoltagePerPhase[0]
			h.gridVoltagePerPhase[0] = h.gridVoltagePerPhase[1]
			h.gridVoltagePerPhase[1] = h.gridVoltagePerPhase[2]
			h.gridVoltagePerPhase[2] = tmp

			_ = h.ucgcpmgcp.SetFrequency(math.Abs(h.gridFrequency))
			_ = h.ucmumpc.SetFrequency(math.Abs(h.gridFrequency))
			switch h.gridFrequency {
			case 50:
				h.gridFrequency = 51
			case 51:
				h.gridFrequency = -50
			case -50:
				h.gridFrequency = 49
			case 49:
				h.gridFrequency = 50
			}
		}
	}()

	h.myService.RegisterRemoteSKI(remoteSki)

	h.myService.Start()
	// defer h.myService.Shutdown()
}

// Monitoring Appliance MGCP Event Handler

func (h *smartmeter) OnMGCPEvent(ski string, device spineapi.DeviceRemoteInterface, entity spineapi.EntityRemoteInterface, event api.EventType) {
}

// Monitored Unit MPC Event Handler

func (h *smartmeter) OnMPCEvent(ski string, device spineapi.DeviceRemoteInterface, entity spineapi.EntityRemoteInterface, event api.EventType) {
}

// EEBUSServiceHandler

func (h *smartmeter) RemoteSKIConnected(service api.ServiceInterface, ski string) {}

func (h *smartmeter) RemoteSKIDisconnected(service api.ServiceInterface, ski string) {}

func (h *smartmeter) VisibleRemoteServicesUpdated(service api.ServiceInterface, entries []shipapi.RemoteService) {
}

func (h *smartmeter) ServiceShipIDUpdate(ski string, shipdID string) {}

func (h *smartmeter) ServicePairingDetailUpdate(ski string, detail *shipapi.ConnectionStateDetail) {
	if ski == remoteSki && detail.State() == shipapi.ConnectionStateRemoteDeniedTrust {
		fmt.Println("The remote service denied trust. Exiting.")
		h.myService.CancelPairingWithSKI(ski)
		h.myService.UnregisterRemoteSKI(ski)
		h.myService.Shutdown()
		os.Exit(0)
	}
}

func (h *smartmeter) AllowWaitingForTrust(ski string) bool {
	return ski == remoteSki
}

// main app
func usage() {
	fmt.Println("First Run:")
	fmt.Println("  go run /examples/smartmeter/main.go <serverport>")
	fmt.Println()
	fmt.Println("General Usage:")
	fmt.Println("  go run /examples/smartmeter/main.go <serverport> <remoteski> <crtfile> <keyfile>")
}

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}

	h := smartmeter{}
	h.run()

	// Clean exit to make sure mdns shutdown is invoked
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	// User exit
}

// Logging interface

func (h *smartmeter) Trace(args ...interface{}) {
	h.print("TRACE", args...)
}

func (h *smartmeter) Tracef(format string, args ...interface{}) {
	h.printFormat("TRACE", format, args...)
}

func (h *smartmeter) Debug(args ...interface{}) {
	h.print("DEBUG", args...)
}

func (h *smartmeter) Debugf(format string, args ...interface{}) {
	h.printFormat("DEBUG", format, args...)
}

func (h *smartmeter) Info(args ...interface{}) {
	h.print("INFO ", args...)
}

func (h *smartmeter) Infof(format string, args ...interface{}) {
	h.printFormat("INFO ", format, args...)
}

func (h *smartmeter) Error(args ...interface{}) {
	h.print("ERROR", args...)
}

func (h *smartmeter) Errorf(format string, args ...interface{}) {
	h.printFormat("ERROR", format, args...)
}

func (h *smartmeter) currentTimestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func (h *smartmeter) print(msgType string, args ...interface{}) {
	value := fmt.Sprintln(args...)
	fmt.Printf("%s %s %s", h.currentTimestamp(), msgType, value)
}

func (h *smartmeter) printFormat(msgType, format string, args ...interface{}) {
	value := fmt.Sprintf(format, args...)
	fmt.Println(h.currentTimestamp(), msgType, value)
}
