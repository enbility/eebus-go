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
	"github.com/enbility/eebus-go/service"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/eebus-go/usecases/cem/vabd"
	"github.com/enbility/eebus-go/usecases/cem/vapd"
	cslpc "github.com/enbility/eebus-go/usecases/cs/lpc"
	cslpp "github.com/enbility/eebus-go/usecases/cs/lpp"
	eglpc "github.com/enbility/eebus-go/usecases/eg/lpc"
	eglpp "github.com/enbility/eebus-go/usecases/eg/lpp"
	"github.com/enbility/eebus-go/usecases/ma/mgcp"
	shipapi "github.com/enbility/ship-go/api"
	"github.com/enbility/ship-go/cert"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

var remoteSki string

type hems struct {
	myService *service.Service

	uccslpc   ucapi.CsLPCInterface
	uccslpp   ucapi.CsLPPInterface
	uceglpc   ucapi.EgLPCInterface
	uceglpp   ucapi.EgLPPInterface
	ucmamgcp  ucapi.MaMGCPInterface
	uccemvabd ucapi.CemVABDInterface
	uccemvapd ucapi.CemVAPDInterface
}

func (h *hems) run() {
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
		"Demo", "Demo", "HEMS", "123456789",
		[]shipapi.DeviceCategoryType{shipapi.DeviceCategoryTypeEnergyManagementSystem},
		model.DeviceTypeTypeEnergyManagementSystem,
		[]model.EntityTypeType{model.EntityTypeTypeCEM},
		port, certificate, time.Second*4)
	if err != nil {
		log.Fatal(err)
	}
	configuration.SetAlternateIdentifier("Demo-HEMS-123456789")

	h.myService = service.NewService(configuration, h)
	h.myService.SetLogging(h)

	if err = h.myService.Setup(); err != nil {
		fmt.Println(err)
		return
	}

	localEntity := h.myService.LocalDevice().EntityForType(model.EntityTypeTypeCEM)
	h.uccslpc = cslpc.NewLPC(localEntity, h.OnLPCEvent)
	err = h.myService.AddUseCase(h.uccslpc)
	if err != nil {
		log.Fatal(err)
	}

	h.uccslpp = cslpp.NewLPP(localEntity, h.OnLPPEvent)
	err = h.myService.AddUseCase(h.uccslpp)
	if err != nil {
		log.Fatal(err)
	}

	h.uceglpc = eglpc.NewLPC(localEntity, nil)
	err = h.myService.AddUseCase(h.uceglpc)
	if err != nil {
		log.Fatal(err)
	}

	h.uceglpp = eglpp.NewLPP(localEntity, nil)
	err = h.myService.AddUseCase(h.uceglpp)
	if err != nil {
		log.Fatal(err)
	}

	h.ucmamgcp = mgcp.NewMGCP(localEntity, h.OnMGCPEvent)
	err = h.myService.AddUseCase(h.ucmamgcp)
	if err != nil {
		log.Fatal(err)
	}

	h.uccemvabd = vabd.NewVABD(localEntity, h.OnVABDEvent)
	err = h.myService.AddUseCase(h.uccemvabd)
	if err != nil {
		log.Fatal(err)
	}

	h.uccemvapd = vapd.NewVAPD(localEntity, h.OnVAPDEvent)
	err = h.myService.AddUseCase(h.uccemvapd)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize local server data
	_ = h.uccslpc.SetConsumptionNominalMax(32000)
	_ = h.uccslpc.SetConsumptionLimit(ucapi.LoadLimit{
		Value:        4200,
		IsChangeable: true,
		IsActive:     false,
	})
	_ = h.uccslpc.SetFailsafeConsumptionActivePowerLimit(4200, true)
	_ = h.uccslpc.SetFailsafeDurationMinimum(2*time.Hour, true)

	_ = h.uccslpp.SetProductionNominalMax(10000)
	_ = h.uccslpp.SetProductionLimit(ucapi.LoadLimit{
		Value:        10000,
		IsChangeable: true,
		IsActive:     false,
	})
	_ = h.uccslpp.SetFailsafeProductionActivePowerLimit(4200, true)
	_ = h.uccslpp.SetFailsafeDurationMinimum(2*time.Hour, true)

	if len(remoteSki) == 0 {
		os.Exit(0)
	}

	h.myService.RegisterRemoteSKI(remoteSki, "")

	h.myService.Start()
	// defer h.myService.Shutdown()
}

// Controllable System LPC Event Handler

func (h *hems) OnLPCEvent(ski string, device spineapi.DeviceRemoteInterface, entity spineapi.EntityRemoteInterface, event api.EventType) {
	switch event {
	case cslpc.WriteApprovalRequired:
		// get pending writes
		pendingWrites := h.uccslpc.PendingConsumptionLimits()

		// approve any write
		for msgCounter, write := range pendingWrites {
			fmt.Println("Approving LPC write with msgCounter", msgCounter, "and limit", write.Value, "W")
			h.uccslpc.ApproveOrDenyConsumptionLimit(msgCounter, true, "")
		}
	case cslpc.DataUpdateLimit:
		if currentLimit, err := h.uccslpc.ConsumptionLimit(); err == nil {
			fmt.Println("New LPC Limit set to", currentLimit.Value, "W")
		}
	}
}

// Controllable System LPP Event Handler

func (h *hems) OnLPPEvent(ski string, device spineapi.DeviceRemoteInterface, entity spineapi.EntityRemoteInterface, event api.EventType) {
	switch event {
	case cslpp.WriteApprovalRequired:
		// get pending writes
		pendingWrites := h.uccslpp.PendingProductionLimits()

		// approve any write
		for msgCounter, write := range pendingWrites {
			fmt.Println("Approving LPP write with msgCounter", msgCounter, "and limit", write.Value, "W")
			h.uccslpp.ApproveOrDenyProductionLimit(msgCounter, true, "")
		}
	case cslpp.DataUpdateLimit:
		if currentLimit, err := h.uccslpp.ProductionLimit(); err == nil {
			fmt.Println("New LPP Limit set to", currentLimit.Value, "W")
		}
	}
}

// Cem VABD Event Handler

func (h *hems) OnVABDEvent(ski string, device spineapi.DeviceRemoteInterface, entity spineapi.EntityRemoteInterface, event api.EventType) {
	switch event {
	case vabd.DataUpdateEnergyCharged:
		if energy, err := h.uccemvabd.EnergyCharged(entity); err == nil {
			fmt.Println("New VABD Energy Charged set to", energy, "Wh")
		}
	case vabd.DataUpdateEnergyDischarged:
		if energy, err := h.uccemvabd.EnergyDischarged(entity); err == nil {
			fmt.Println("New VABD Energy Discharged set to", energy, "Wh")
		}
	case vabd.DataUpdatePower:
		if power, err := h.uccemvabd.Power(entity); err == nil {
			fmt.Println("New VABD Power set to", power, "W")
		}
	case vabd.DataUpdateStateOfCharge:
		if soc, err := h.uccemvabd.StateOfCharge(entity); err == nil {
			fmt.Println("New VABD State of Charge set to", soc, "%")
		}
	}
}

// Cem VAPD Event Handler

func (h *hems) OnVAPDEvent(ski string, device spineapi.DeviceRemoteInterface, entity spineapi.EntityRemoteInterface, event api.EventType) {
	switch event {
	case vapd.DataUpdatePVYieldTotal:
		if yield, err := h.uccemvapd.PVYieldTotal(entity); err == nil {
			fmt.Println("New VAPD PV Yield Total set to", yield, "Wh")
		}
	case vapd.DataUpdatePowerNominalPeak:
		if peak, err := h.uccemvapd.PowerNominalPeak(entity); err == nil {
			fmt.Println("New VAPD Power Nominal Peak set to", peak, "W")
		}
	case vapd.DataUpdatePower:
		if power, err := h.uccemvapd.Power(entity); err == nil {
			fmt.Println("New VAPD Power set to", power, "W")
		}
	}
}

// Monitoring Appliance MGCP Event Handler

func (h *hems) OnMGCPEvent(ski string, device spineapi.DeviceRemoteInterface, entity spineapi.EntityRemoteInterface, event api.EventType) {
	switch event {
	case mgcp.DataUpdatePowerLimitationFactor:
		if factor, err := h.ucmamgcp.PowerLimitationFactor(entity); err == nil {
			fmt.Println("New MGCP Power Limitation Factor set to", factor)
		}
	case mgcp.DataUpdatePower:
		if power, err := h.ucmamgcp.Power(entity); err == nil {
			fmt.Println("New MGCP Power set to", power, "W")
		}
	case mgcp.DataUpdateEnergyFeedIn:
		if energy, err := h.ucmamgcp.EnergyFeedIn(entity); err == nil {
			fmt.Println("New MGCP Energy Feed-In set to", energy, "Wh")
		}
	case mgcp.DataUpdateEnergyConsumed:
		if energy, err := h.ucmamgcp.EnergyConsumed(entity); err == nil {
			fmt.Println("New MGCP Energy Consumed set to", energy, "Wh")
		}
	case mgcp.DataUpdateCurrentPerPhase:
		if current, err := h.ucmamgcp.CurrentPerPhase(entity); err == nil {
			fmt.Println("New MGCP Current per Phase set to", current, "A")
		}
	case mgcp.DataUpdateVoltagePerPhase:
		if voltage, err := h.ucmamgcp.VoltagePerPhase(entity); err == nil {
			fmt.Println("New MGCP Voltage per Phase set to", voltage, "V")
		}
	case mgcp.DataUpdateFrequency:
		if frequency, err := h.ucmamgcp.Frequency(entity); err == nil {
			fmt.Println("New MGCP Frequency set to", frequency, "Hz")
		}
	}
}

// EEBUSServiceHandler

func (h *hems) RemoteSKIConnected(service api.ServiceInterface, ski string) {}

func (h *hems) RemoteSKIDisconnected(service api.ServiceInterface, ski string) {}

func (h *hems) VisibleRemoteServicesUpdated(service api.ServiceInterface, entries []shipapi.RemoteService) {
}

func (h *hems) ServiceShipIDUpdate(ski string, shipdID string) {}

func (h *hems) ServicePairingDetailUpdate(ski string, detail *shipapi.ConnectionStateDetail) {
	if ski == remoteSki && detail.State() == shipapi.ConnectionStateRemoteDeniedTrust {
		fmt.Println("The remote service denied trust. Exiting.")
		h.myService.CancelPairingWithSKI(ski)
		h.myService.UnregisterRemoteSKI(ski)
		h.myService.Shutdown()
		os.Exit(0)
	}
}

func (h *hems) AllowWaitingForTrust(ski string) bool {
	return ski == remoteSki
}

// UCEvseCommisioningConfigurationCemDelegate

// handle device state updates from the remote EVSE device
func (h *hems) HandleEVSEDeviceState(ski string, failure bool, errorCode string) {
	fmt.Println("EVSE Error State:", failure, errorCode)
}

// main app
func usage() {
	fmt.Println("First Run:")
	fmt.Println("  go run /examples/hems/main.go <serverport>")
	fmt.Println()
	fmt.Println("General Usage:")
	fmt.Println("  go run /examples/hems/main.go <serverport> <remoteski> <crtfile> <keyfile>")
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
