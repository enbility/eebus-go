package main

import (
	"crypto/ecdsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
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

type configuration struct {
	port      uint
	certPath  string
	keyPath   string
	remoteSKI string
	secret    string
}

func (c configuration) String() string {
	return fmt.Sprintf("port: %v\ncertPath: %v\nkeyPath: %v\nremoteSKI: %v\nsecret: %v",
		c.port, c.certPath, c.keyPath, c.remoteSKI, c.secret)
}

func (c configuration) Valid() bool {
	return (c.port > 0 && c.port < 65536) //&&
	// c.certPath != "" && // they can be empty to generate new certificate
	// c.keyPath != "" &&
	// (len(c.remoteSKI) > 0)
}

var config configuration

type hems struct {
	myService *service.Service

	uccslpc   ucapi.CsLPCInterface
	uccslpp   ucapi.CsLPPInterface
	uceglpc   ucapi.EgLPCInterface
	uceglpp   ucapi.EgLPPInterface
	ucmamgcp  ucapi.MaMGCPInterface
	uccemvabd ucapi.CemVABDInterface
	uccemvapd ucapi.CemVAPDInterface

	ringBuffer *ExampleRingBufferPersistence
}

func (h *hems) run() {
	var err error
	var certificate tls.Certificate

	if config.certPath != "" && config.keyPath != "" {
		certificate, err = tls.LoadX509KeyPair(config.certPath, config.keyPath)
		if err != nil {
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

	var pairingConfig *shipapi.PairingConfig
	if len(config.secret) == 0 {
		fmt.Println("Secret not specified.. PairingConfig not set")
	} else {
		secretBytes, err := hex.DecodeString(config.secret)
		if err != nil {
			log.Fatal(err)
		}
		pairingConfig = shipapi.NewPairingConfig(shipapi.PairingModeListener, shipapi.PairingSecret(secretBytes))
	}

	h.ringBuffer = &ExampleRingBufferPersistence{}

	configuration, err := api.NewConfiguration(
		"Demo", "Demo", "HEMS", "123456789",
		[]shipapi.DeviceCategoryType{shipapi.DeviceCategoryTypeEnergyManagementSystem},
		model.DeviceTypeTypeEnergyManagementSystem,
		[]model.EntityTypeType{model.EntityTypeTypeCEM},
		int(config.port), certificate, time.Second*4, pairingConfig, h.ringBuffer)
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

	// NOTE: Per the LPP spec, APPL (Active Power Limit) values for production
	// must be <= 0. The eebus-go stack does not transform positive values to
	// negative values, so the caller must provide the correct sign.
	_ = h.uccslpp.SetProductionNominalMax(10000)
	_ = h.uccslpp.SetProductionLimit(ucapi.LoadLimit{
		Value:        -10000,
		IsChangeable: true,
		IsActive:     false,
	})
	_ = h.uccslpp.SetFailsafeProductionActivePowerLimit(4200, true)
	_ = h.uccslpp.SetFailsafeDurationMinimum(2*time.Hour, true)

	if len(config.remoteSKI) == 0 && pairingConfig == nil {
		os.Exit(0)
	}

	if len(config.remoteSKI) == 40 {
		h.myService.RegisterRemoteService(shipapi.NewServiceIdentity(config.remoteSKI, "", ""))
	}

	_ = h.myService.Start()
	if qrCodeText, err := h.myService.QRCodeText(); err == nil {
		fmt.Println("Service QR Code:", qrCodeText)
	}
	// defer h.myService.Shutdown()
}

// Controllable System LPC Event Handler

func (h *hems) OnLPCEvent(ski string, device spineapi.DeviceRemoteInterface, entity spineapi.EntityRemoteInterface, event api.EventType) {
	switch event {
	case cslpc.LimitWriteApprovalRequired:
		// get pending writes
		pendingWrites := h.uccslpc.PendingConsumptionLimits()

		// approve any write
		for msgCounter, write := range pendingWrites {
			fmt.Println("Approving LPC limit write with msgCounter", msgCounter, "and limit", write.Value, "W")
			h.uccslpc.ApproveOrDenyConsumptionLimit(msgCounter, true, "")
		}
	case cslpc.ConfigurationWriteApprovalRequired:
		pendingDeviceConfigWrites := h.uccslpc.PendingDeviceConfigurations()

		for msgCounter, configs := range pendingDeviceConfigWrites {
			fmt.Printf("Approving LPC device config write with msgCounter %d for features: ", msgCounter)
			for _, config := range configs {
				fmt.Printf("%s ", string(config.KeyName))
			}
			fmt.Print("\n")
			h.uccslpc.ApproveOrDenyDeviceConfiguration(msgCounter, true, "")
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
	case cslpp.LimitWriteApprovalRequired:
		// get pending limit writes
		pendingWrites := h.uccslpp.PendingProductionLimits()

		// approve any write
		for msgCounter, write := range pendingWrites {
			fmt.Println("Approving LPP limit write with msgCounter", msgCounter, "and limit", write.Value, "W")
			h.uccslpp.ApproveOrDenyProductionLimit(msgCounter, true, "")
		}
	case cslpp.ConfigurationWriteApprovalRequired:
		// get pending device config writes
		pendingDeviceConfigWrites := h.uccslpp.PendingDeviceConfigurations()

		// approve any write
		for msgCounter, configs := range pendingDeviceConfigWrites {
			fmt.Printf("Approving LPP device config write with msgCounter %d for features: ", msgCounter)
			for _, config := range configs {
				if config.Description.KeyName != nil {
					fmt.Printf("%s ", string(*config.Description.KeyName))
				}
			}
			fmt.Print("\n")
			h.uccslpp.ApproveOrDenyDeviceConfiguration(msgCounter, true, "")
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

func (h *hems) RemoteServiceConnected(service api.ServiceInterface, identity shipapi.ServiceIdentity) {
}

func (h *hems) RemoteServiceDisconnected(service api.ServiceInterface, identity shipapi.ServiceIdentity) {
}

func (h *hems) VisibleRemoteMdnsServicesUpdated(service api.ServiceInterface, entries []shipapi.RemoteMdnsService) {
}

func (h *hems) ServiceUpdated(identity shipapi.ServiceIdentity) {}

func (h *hems) ServicePairingDetailUpdate(identity shipapi.ServiceIdentity, detail *shipapi.ConnectionStateDetail) {
	if identity.SKI == config.remoteSKI && detail.State() == shipapi.ConnectionStateRemoteDeniedTrust {
		fmt.Println("The remote service denied trust. Exiting.")
		h.myService.CancelPairing(identity)
		h.myService.UnregisterRemoteService(identity)
		h.myService.Shutdown()
		os.Exit(0)
	}
}

func (h *hems) AllowWaitingForTrust(identity shipapi.ServiceIdentity) bool {
	if len(config.remoteSKI) != 40 {
		// no remoteSki given, give a chance for automatic pairing
		return true
	}
	return identity.SKI == config.remoteSKI
}

// ServiceAutoTrustFailed implements api.ServiceReaderInterface.
func (h *hems) ServiceAutoTrustFailed(service api.ServiceInterface, identity shipapi.ServiceIdentity, reason error) {
	log.Printf("Auto Trust failed for identity %v: %v\n", identity.String(), reason)
}

// ServiceAutoTrustRemoved implements api.ServiceReaderInterface.
func (h *hems) ServiceAutoTrustRemoved(service api.ServiceInterface, identity shipapi.ServiceIdentity, reason string) {
	log.Printf("Auto Trust removed for identity %v: %v\n", identity.String(), reason)
}

// ServiceAutoTrusted implements api.ServiceReaderInterface.
func (h *hems) ServiceAutoTrusted(service api.ServiceInterface, identity shipapi.ServiceIdentity) {
	log.Printf("Auto Trust successful for identity %v\n", identity.String())
}

// UCEvseCommisioningConfigurationCemDelegate

// handle device state updates from the remote EVSE device
func (h *hems) HandleEVSEDeviceState(ski string, failure bool, errorCode string) {
	fmt.Println("EVSE Error State:", failure, errorCode)
}

// main app

func main() {
	flag.StringVar(&config.certPath, "certpath", "", "./path/to/cert.pem")
	flag.StringVar(&config.keyPath, "keypath", "", "./path/to/key.pem")
	flag.StringVar(&config.remoteSKI, "remoteski", "", "remote SKI")
	flag.UintVar(&config.port, "port", 0, "server port")
	flag.StringVar(&config.secret, "secret", "", "secret hexadecimal")
	flag.Parse()

	if !config.Valid() {
		flag.Usage()
		log.Fatal("invalid configuration")
	}

	fmt.Println("--------------")
	fmt.Println("Configuration")
	fmt.Println(config)
	fmt.Println("--------------")

	h := hems{}
	h.run()

	// Clean exit to make sure mdns shutdown is invoked
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	if h.myService != nil {
		h.myService.Shutdown()
	}
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

// Ring buffer persistence

type ExampleRingBufferPersistence struct {
}

// LoadRingBuffer implements api.RingBufferPersistence.
func (r *ExampleRingBufferPersistence) LoadRingBuffer() ([]shipapi.DigestEntry, int, error) {
	fmt.Printf("   📂 Loading ring buffer from...\n")

	// For demonstration: return empty state (no previous storage)
	// Real applications should load from persistent storage:
	//
	//   data, err := os.ReadFile(r.filename)
	//   if os.IsNotExist(err) {
	//       // No previous data - return empty buffer (library will initialize)
	//       return make([]api.DigestEntry, 100), 0, nil
	//   }
	//   if err != nil {
	//       return nil, 0, fmt.Errorf("failed to read ring buffer: %w", err)
	//   }
	//
	//   var state struct {
	//       Entries   []api.DigestEntry `json:"entries"`
	//       NextIndex int              `json:"nextIndex"`
	//   }
	//   if err := json.Unmarshal(data, &state); err != nil {
	//       return nil, 0, fmt.Errorf("failed to parse ring buffer: %w", err)
	//   }
	//
	//   return state.Entries, state.NextIndex, nil

	// Demo: return empty buffer (library will manage the ring buffer logic)
	fmt.Printf("   📂 No previous data, library will create fresh ring buffer\n")
	return make([]shipapi.DigestEntry, 100), 0, nil
}

// SaveRingBuffer implements api.RingBufferPersistence.
func (r *ExampleRingBufferPersistence) SaveRingBuffer(entries []shipapi.DigestEntry, nextIndex int) error {
	fmt.Printf("   💾 Library requests save: %d entries, nextIndex=%d to\n", len(entries), nextIndex)

	// For demonstration: just log what the library is providing
	// Real applications should save the library's data to persistent storage:
	//
	//   state := struct {
	//       Entries   []api.DigestEntry `json:"entries"`
	//       NextIndex int              `json:"nextIndex"`
	//   }{
	//       Entries:   entries,    // Complete ring buffer from library
	//       NextIndex: nextIndex,  // Current position from library
	//   }
	//
	//   data, err := json.Marshal(state)
	//   if err != nil {
	//       return fmt.Errorf("failed to marshal ring buffer: %w", err)
	//   }
	//
	//   // Atomic write pattern for safety
	//   tempFile := r.filename + ".tmp"
	//   if err := os.WriteFile(tempFile, data, 0600); err != nil {
	//       return fmt.Errorf("failed to write ring buffer: %w", err)
	//   }
	//
	//   return os.Rename(tempFile, r.filename)

	fmt.Printf("   💾 Demo: Would save library's ring buffer state to\n")
	return nil
}
