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
	"github.com/enbility/eebus-go/usecases/cs/lpc"
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

type evse struct {
	myService *service.Service

	uclpc *lpc.LPC

	isConnected bool

	ringBuffer *ExampleRingBufferPersistence
}

func (h *evse) run() {
	var err error
	var certificate tls.Certificate

	if config.certPath != "" && config.keyPath != "" {
		certificate, err = tls.LoadX509KeyPair(config.certPath, config.keyPath)
		if err != nil {
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
		"Demo", "Demo", "EVSE", "234567890",
		[]shipapi.DeviceCategoryType{shipapi.DeviceCategoryTypeEMobility},
		model.DeviceTypeTypeChargingStation,
		[]model.EntityTypeType{model.EntityTypeTypeEVSE},
		int(config.port), certificate, time.Second*4, pairingConfig, h.ringBuffer)
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

	localEntity := h.myService.LocalDevice().EntityForType(model.EntityTypeTypeEVSE)
	h.uclpc = lpc.NewLPC(localEntity, h.OnLPCEvent)
	err = h.myService.AddUseCase(h.uclpc)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize local server data
	_ = h.uclpc.SetConsumptionNominalMax(32000)
	_ = h.uclpc.SetConsumptionLimit(ucapi.LoadLimit{
		Value:        4200,
		IsChangeable: true,
		IsActive:     false,
	})
	_ = h.uclpc.SetFailsafeConsumptionActivePowerLimit(4200, true)
	_ = h.uclpc.SetFailsafeDurationMinimum(2*time.Hour, true)

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

// EEBUSServiceHandler

func (h *evse) RemoteServiceConnected(service api.ServiceInterface, identity shipapi.ServiceIdentity) {
	log.Println("Remote Service Connected:", identity.String())
	h.isConnected = true
}

func (h *evse) RemoteServiceDisconnected(service api.ServiceInterface, identity shipapi.ServiceIdentity) {
	log.Println("Remote Service Disconnected:", identity.String())
	h.isConnected = false
}

func (h *evse) VisibleRemoteMdnsServicesUpdated(service api.ServiceInterface, entries []shipapi.RemoteMdnsService) {
}

func (h *evse) ServiceUpdated(identity shipapi.ServiceIdentity) {
	log.Println("Service Updated:", identity.String())
}

func (h *evse) ServicePairingDetailUpdate(identity shipapi.ServiceIdentity, detail *shipapi.ConnectionStateDetail) {
	log.Println("Remote Service Pairing Detail Update:", identity.String())
	if identity.SKI == config.remoteSKI && detail.State() == shipapi.ConnectionStateRemoteDeniedTrust {
		fmt.Println("The remote service denied trust. Exiting.")
		h.myService.CancelPairing(identity)
		h.myService.UnregisterRemoteService(identity)
		h.myService.Shutdown()
		os.Exit(0)
	}
}

func (h *evse) AllowWaitingForTrust(identity shipapi.ServiceIdentity) bool {
	if len(config.remoteSKI) != 40 {
		// no remoteSki given, give a chance for automatic pairing
		return true
	}
	return identity.SKI == config.remoteSKI
}

// ServiceAutoTrustFailed implements api.ServiceReaderInterface.
func (h *evse) ServiceAutoTrustFailed(service api.ServiceInterface, identity shipapi.ServiceIdentity, reason error) {
	log.Printf("Auto Trust failed for identity %v: %v\n", identity.String(), reason)
}

// ServiceAutoTrustRemoved implements api.ServiceReaderInterface.
func (h *evse) ServiceAutoTrustRemoved(service api.ServiceInterface, identity shipapi.ServiceIdentity, reason string) {
	log.Printf("Auto Trust removed for identity %v: %v\n", identity.String(), reason)
}

// ServiceAutoTrusted implements api.ServiceReaderInterface.
func (h *evse) ServiceAutoTrusted(service api.ServiceInterface, identity shipapi.ServiceIdentity) {
	log.Printf("Auto Trust successful for identity %v\n", identity.String())
}

// LPC Event Handler

func (h *evse) OnLPCEvent(ski string, device spineapi.DeviceRemoteInterface, entity spineapi.EntityRemoteInterface, event api.EventType) {
	if !h.isConnected {
		return
	}

	switch event {
	case lpc.LimitWriteApprovalRequired:
		// get pending writes
		pendingWrites := h.uclpc.PendingConsumptionLimits()

		// approve any write
		for msgCounter, write := range pendingWrites {
			fmt.Println("Approving write with msgCounter", msgCounter, "and limit", write.Value, "W")
			h.uclpc.ApproveOrDenyConsumptionLimit(msgCounter, true, "")
		}
	case lpc.DataUpdateLimit:
		if currentLimit, err := h.uclpc.ConsumptionLimit(); err == nil {
			fmt.Println("New Limit set to", currentLimit.Value, "W")
		}
	}
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

	h := evse{}
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

func (h *evse) Trace(args ...interface{}) {
	// h.print("TRACE", args...)
}

func (h *evse) Tracef(format string, args ...interface{}) {
	// h.printFormat("TRACE", format, args...)
}

func (h *evse) Debug(args ...interface{}) {
	// h.print("DEBUG", args...)
}

func (h *evse) Debugf(format string, args ...interface{}) {
	// h.printFormat("DEBUG", format, args...)
}

func (h *evse) Info(args ...interface{}) {
	h.print("INFO ", args...)
}

func (h *evse) Infof(format string, args ...interface{}) {
	h.printFormat("INFO ", format, args...)
}

func (h *evse) Error(args ...interface{}) {
	h.print("ERROR", args...)
}

func (h *evse) Errorf(format string, args ...interface{}) {
	h.printFormat("ERROR", format, args...)
}

func (h *evse) currentTimestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func (h *evse) print(msgType string, args ...interface{}) {
	value := fmt.Sprintln(args...)
	fmt.Printf("%s %s %s", h.currentTimestamp(), msgType, value)
}

func (h *evse) printFormat(msgType, format string, args ...interface{}) {
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
