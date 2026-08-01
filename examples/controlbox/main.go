package main

import (
	"crypto/ecdsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"
	"time"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/service"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/eebus-go/usecases/eg/lpc"
	"github.com/enbility/eebus-go/usecases/eg/lpp"
	shipapi "github.com/enbility/ship-go/api"
	"github.com/enbility/ship-go/cert"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type controlbox struct {
	myService *service.Service

	uclpc ucapi.EgLPCInterface
	uclpp ucapi.EgLPPInterface

	isConnected bool
}

type configuration struct {
	port           uint
	certPath       string
	keyPath        string
	remoteSKIs     remoteSKIList
	pairingTargets targetList
}

func (c configuration) String() string {
	return fmt.Sprintf("port: %v\ncertPath: %v\nkeyPath: %v\nremoteSKIs: %v\npairingTargets: %v",
		c.port, c.certPath, c.keyPath, c.remoteSKIs, c.pairingTargets)
}

func (c configuration) Valid() bool {
	return (c.port > 0 && c.port < 65536) //&&
	// c.certPath != "" && // they can be empty to generate new certificate
	// c.keyPath != "" &&
	// (len(c.remoteSKIs) > 0 || len(c.pairingTargets) > 0)
}

var config configuration

type remoteSKIList []string

func (r *remoteSKIList) String() string {
	return fmt.Sprint(*r)
}

func (r *remoteSKIList) Set(value string) error {
	if len(value) != 40 {
		return fmt.Errorf("invalid SKI")
	}
	_, err := hex.DecodeString(value)
	if err != nil {
		return err
	}
	*r = append(*r, value)
	return nil
}

type targetList []shipapi.PairingTarget

func (t *targetList) String() string {
	return fmt.Sprint(*t)
}

func (t *targetList) Set(value string) error {
	value = strings.TrimSpace(value)

	// Case 1: QR input
	if strings.HasPrefix(value, "SHIP;") {
		pt, err := NewPairingTargetFromQrCode(value)
		if err != nil {
			return fmt.Errorf("failed to decode QR input: %v", err)
		}

		*t = append(*t, pt)
		return nil
	}

	// Case 2: key=value parsing (existing logic)
	pt := shipapi.PairingTarget{}

	parts := strings.Split(value, ",")
	for _, p := range parts {
		kv := strings.SplitN(p, "=", 2)
		if len(kv) != 2 {
			return fmt.Errorf("invalid key=value pair: %q", p)
		}

		key := strings.ToLower(strings.TrimSpace(kv[0]))
		val := strings.TrimSpace(kv[1])

		switch key {
		case "ski":
			pt.SKI = val
		case "fingerprint":
			pt.Fingerprint = val
		case "shipid":
			pt.ShipID = val
		case "secret":
			decoded, err := base64.StdEncoding.DecodeString(val)
			if err != nil {
				return fmt.Errorf("invalid base64 for secret: %v", err)
			}
			pt.Secret = decoded
		default:
			return fmt.Errorf("unknown field: %s", key)
		}
	}

	// Validation
	if len(pt.Secret) != 16 || pt.Fingerprint == "" || pt.ShipID == "" {
		return fmt.Errorf("missing required fields (Secret, Fingerprint, ShipID)")
	}

	*t = append(*t, pt)
	return nil
}

func (h *controlbox) run() {
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

	pairingConfig := shipapi.NewPairingConfig(shipapi.PairingModeAnnouncer, nil)

	configuration, err := api.NewConfiguration(
		"Demo", "Demo", "ControlBox", "123456789",
		[]shipapi.DeviceCategoryType{shipapi.DeviceCategoryTypeGridConnectionHub},
		model.DeviceTypeTypeElectricitySupplySystem,
		[]model.EntityTypeType{model.EntityTypeTypeGridGuard},
		int(config.port), certificate, time.Second*60, pairingConfig, nil)
	if err != nil {
		log.Fatal(err)
	}
	configuration.SetAlternateIdentifier("Demo-ControlBox-123456789")

	h.myService = service.NewService(configuration, h)
	h.myService.SetLogging(h)

	if err = h.myService.Setup(); err != nil {
		fmt.Println(err)
		return
	}

	localEntity := h.myService.LocalDevice().EntityForType(model.EntityTypeTypeGridGuard)
	h.uclpc = lpc.NewLPC(localEntity, h.OnLPCEvent)
	err = h.myService.AddUseCase(h.uclpc)
	if err != nil {
		log.Fatal(err)
	}

	h.uclpp = lpp.NewLPP(localEntity, h.OnLPPEvent)
	err = h.myService.AddUseCase(h.uclpp)
	if err != nil {
		log.Fatal(err)
	}

	for _, remoteSki := range config.remoteSKIs {
		h.myService.RegisterRemoteService(shipapi.NewServiceIdentity(remoteSki, "", ""))
	}

	_ = h.myService.Start()
	for _, pairingTarget := range config.pairingTargets {
		identity := shipapi.NewServiceIdentity("", pairingTarget.Fingerprint, pairingTarget.ShipID)
		h.myService.RegisterRemoteService(identity)
		h.myService.StartAnnouncementTo(pairingTarget)
	}
	// defer h.myService.Shutdown()
}

// EEBUSServiceHandler

func (h *controlbox) RemoteServiceConnected(service api.ServiceInterface, identity shipapi.ServiceIdentity) {
	h.isConnected = true
}

func (h *controlbox) RemoteServiceDisconnected(service api.ServiceInterface, identity shipapi.ServiceIdentity) {
	h.isConnected = false
}

func (h *controlbox) VisibleRemoteMdnsServicesUpdated(service api.ServiceInterface, entries []shipapi.RemoteMdnsService) {
}

func (h *controlbox) ServiceUpdated(identity shipapi.ServiceIdentity) {
}

func (h *controlbox) ServicePairingDetailUpdate(identity shipapi.ServiceIdentity, detail *shipapi.ConnectionStateDetail) {
	if slices.Contains(config.remoteSKIs, identity.SKI) && detail.State() == shipapi.ConnectionStateRemoteDeniedTrust {
		fmt.Println("The remote service denied trust. Exiting.")
		h.myService.CancelPairing(identity)
		h.myService.UnregisterRemoteService(identity)
		h.myService.Shutdown()
		os.Exit(0)
	}
}

func (h *controlbox) AllowWaitingForTrust(identity shipapi.ServiceIdentity) bool {
	return slices.Contains(config.remoteSKIs, identity.SKI)
}

// ServiceAutoTrustFailed implements api.ServiceReaderInterface.
func (h *controlbox) ServiceAutoTrustFailed(service api.ServiceInterface, identity shipapi.ServiceIdentity, reason error) {
}

// ServiceAutoTrustRemoved implements api.ServiceReaderInterface.
func (h *controlbox) ServiceAutoTrustRemoved(service api.ServiceInterface, identity shipapi.ServiceIdentity, reason string) {
}

// ServiceAutoTrusted implements api.ServiceReaderInterface.
func (h *controlbox) ServiceAutoTrusted(service api.ServiceInterface, identity shipapi.ServiceIdentity) {
}

// LPC Event Handler

func (h *controlbox) sendConsumptionLimit(entity spineapi.EntityRemoteInterface) {
	scenarios := h.uclpc.AvailableScenariosForEntity(entity)
	if !slices.Contains(scenarios, 1) {
		return
	}

	fmt.Println("Sending a consumption limit in 5s...")
	time.AfterFunc(time.Second*5, func() {
		limit := ucapi.LoadLimit{
			Duration: time.Hour*1 + time.Minute*2 + time.Second*3,
			IsActive: true,
			Value:    100,
		}

		resultCB := func(msg model.ResultDataType, msgCounter model.MsgCounterType) {
			if *msg.ErrorNumber == model.ErrorNumberTypeNoError {
				fmt.Println("Consumption limit accepted.")
			} else {
				fmt.Println("Consumption limit rejected. Code", *msg.ErrorNumber, "Description", *msg.Description)
			}
		}
		msgCounter, err := h.uclpc.WriteConsumptionLimit(entity, limit, resultCB)
		if err != nil {
			fmt.Println("Failed to send consumption limit", err)
			return
		}
		fmt.Println("Sent consumption limit to", entity.Device().Ski(), "with msgCounter", msgCounter)
	})
}
func (h *controlbox) OnLPCEvent(ski string, device spineapi.DeviceRemoteInterface, entity spineapi.EntityRemoteInterface, event api.EventType) {
	if !h.isConnected {
		return
	}

	switch event {
	case lpc.UseCaseSupportUpdate:
		h.sendConsumptionLimit(entity)
	case lpc.DataUpdateLimit:
		if currentLimit, err := h.uclpc.ConsumptionLimit(entity); err == nil {
			fmt.Println("New consumption limit received", currentLimit.Value, "W")
		}
	default:
		return
	}
}

// LPP Event Handler

func (h *controlbox) sendProductionLimit(entity spineapi.EntityRemoteInterface) {
	scenarios := h.uclpc.AvailableScenariosForEntity(entity)
	if !slices.Contains(scenarios, 1) {
		return
	}

	fmt.Println("Sending a production limit in 5s...")
	time.AfterFunc(time.Second*5, func() {
		// NOTE: Per the LPP spec, APPL (Active Power Limit) values for production
		// must be <= 0. The eebus-go stack does not transform positive values to
		// negative values, so the caller must provide the correct sign.
		limit := ucapi.LoadLimit{
			Duration: time.Hour*1 + time.Minute*2 + time.Second*3,
			IsActive: true,
			Value:    -100,
		}

		resultCB := func(msg model.ResultDataType, msgCounter model.MsgCounterType) {
			if *msg.ErrorNumber == model.ErrorNumberTypeNoError {
				fmt.Println("Production limit accepted.")
			} else {
				fmt.Println("Production limit rejected. Code", *msg.ErrorNumber, "Description", *msg.Description)
			}
		}
		msgCounter, err := h.uclpp.WriteProductionLimit(entity, limit, resultCB)
		if err != nil {
			fmt.Println("Failed to send production limit", err)
			return
		}
		fmt.Println("Sent production limit to", entity.Device().Ski(), "with msgCounter", msgCounter)
	})
}
func (h *controlbox) OnLPPEvent(ski string, device spineapi.DeviceRemoteInterface, entity spineapi.EntityRemoteInterface, event api.EventType) {
	if !h.isConnected {
		return
	}

	switch event {
	case lpp.UseCaseSupportUpdate:
		h.sendProductionLimit(entity)
	case lpp.DataUpdateLimit:
		if currentLimit, err := h.uclpp.ProductionLimit(entity); err == nil {
			fmt.Println("New production limit received", currentLimit.Value, "W")
		}
	default:
		return
	}
}

// main app

func main() {
	flag.Var(&config.pairingTargets, "target", "target in format SKI=...,Fingerprint=...,ShipID=...,Secret=hex")
	flag.StringVar(&config.certPath, "certpath", "", "./path/to/cert.pem")
	flag.StringVar(&config.keyPath, "keypath", "", "./path/to/key.pem")
	flag.Var(&config.remoteSKIs, "remoteski", "remote SKI")
	flag.UintVar(&config.port, "port", 0, "server port")
	flag.Parse()

	if !config.Valid() {
		flag.Usage()
		log.Fatal("invalid configuration")
	}

	fmt.Println("--------------")
	fmt.Println("Configuration")
	fmt.Println(config)
	fmt.Println("--------------")

	h := controlbox{}
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

func (h *controlbox) Trace(args ...interface{}) {
	// h.print("TRACE", args...)
}

func (h *controlbox) Tracef(format string, args ...interface{}) {
	// h.printFormat("TRACE", format, args...)
}

func (h *controlbox) Debug(args ...interface{}) {
	// h.print("DEBUG", args...)
}

func (h *controlbox) Debugf(format string, args ...interface{}) {
	// h.printFormat("DEBUG", format, args...)
}

func (h *controlbox) Info(args ...interface{}) {
	h.print("INFO ", args...)
}

func (h *controlbox) Infof(format string, args ...interface{}) {
	h.printFormat("INFO ", format, args...)
}

func (h *controlbox) Error(args ...interface{}) {
	h.print("ERROR", args...)
}

func (h *controlbox) Errorf(format string, args ...interface{}) {
	h.printFormat("ERROR", format, args...)
}

func (h *controlbox) currentTimestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func (h *controlbox) print(msgType string, args ...interface{}) {
	value := fmt.Sprintln(args...)
	fmt.Printf("%s %s %s", h.currentTimestamp(), msgType, value)
}

func (h *controlbox) printFormat(msgType, format string, args ...interface{}) {
	value := fmt.Sprintf(format, args...)
	fmt.Println(h.currentTimestamp(), msgType, value)
}

func NewPairingTargetFromQrCode(qrString string) (shipapi.PairingTarget, error) {
	var ret shipapi.PairingTarget

	// Clean up the input
	qrString = strings.TrimSpace(qrString)

	// Validate SHIP format
	if !strings.HasPrefix(qrString, "SHIP;") {
		return ret, fmt.Errorf("QR string must start with 'SHIP;', got: %s", qrString)
	}

	if !strings.HasSuffix(qrString, "ENDSHIP;") {
		return ret, fmt.Errorf("QR string must end with 'ENDSHIP;', got: %s", qrString)
	}

	// Remove SHIP; prefix and ENDSHIP; suffix
	content := qrString[5 : len(qrString)-8] // Remove "SHIP;" and "ENDSHIP;"

	// Split into key-value pairs
	fields := strings.Split(content, ";")

	for _, field := range fields {
		if field == "" {
			continue
		}

		parts := strings.SplitN(field, ":", 2)
		if len(parts) != 2 {
			continue // Skip malformed fields
		}

		key := strings.ToUpper(strings.TrimSpace(parts[0]))
		value := strings.TrimSpace(parts[1])

		switch key {
		case "SKI":
			ret.SKI = value
		case "ID":
			ret.ShipID = value
		case "FPH256":
			ret.Fingerprint = value
		case "SPSEC":
			// Validate and decode secret
			secret, err := hex.DecodeString(value)
			if err != nil {
				return ret, fmt.Errorf("invalid secret: %w", err)
			}
			ret.Secret = secret
		default:
			// ignore unknown keys
			continue
		}
	}

	// Validate required fields
	if ret.SKI == "" {
		return ret, fmt.Errorf("missing required 'SKI' field")
	}
	if ret.ShipID == "" {
		return ret, fmt.Errorf("missing required 'ID' field (SHIP ID)")
	}
	if ret.Fingerprint == "" {
		return ret, fmt.Errorf("missing required 'FPH256' field (certificate fingerprint)")
	}
	if len(ret.Secret) == 0 {
		return ret, fmt.Errorf("missing required 'SPSEC' field (pairing secret)")
	}

	return ret, nil
}
