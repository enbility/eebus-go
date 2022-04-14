package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/DerAndereAndi/eebus-go/ship"
	"github.com/DerAndereAndi/eebus-go/util"
	"github.com/gorilla/websocket"
)

func createCertificate() (tls.Certificate, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return tls.Certificate{}, err
	}

	// Create the EEBUS service SKI using the private key
	asn1, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return tls.Certificate{}, err
	}
	ski := sha1.Sum(asn1)

	skiString := fmt.Sprintf("%0x", ski)
	fmt.Println("Local SKI: ", skiString)

	subject := pkix.Name{
		OrganizationalUnit: []string{"Demo"},
		Organization:       []string{"Demo"},
		Country:            []string{"DE"},
	}

	template := x509.Certificate{
		SignatureAlgorithm:    x509.ECDSAWithSHA256,
		SerialNumber:          big.NewInt(1),
		Subject:               subject,
		NotBefore:             time.Now(),                                // Valid starting now
		NotAfter:              time.Now().Add(time.Hour * 24 * 365 * 10), // Valid for 10 years
		KeyUsage:              x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IsCA:                  true,
		SubjectKeyId:          ski[:],
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return tls.Certificate{}, err
	}

	tlsCertificate := tls.Certificate{
		Certificate: [][]byte{certBytes},
		PrivateKey:  privateKey,
	}

	return tlsCertificate, nil
}

// TLSConnection creates an encrypted websocket connection
func TLSConnection(cert tls.Certificate) func(uri string) (*websocket.Conn, error) {
	return func(uri string) (*websocket.Conn, error) {
		var ciperSuites = []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
		}

		dialer := &websocket.Dialer{
			Proxy:            http.ProxyFromEnvironment,
			HandshakeTimeout: 5 * time.Second,
			TLSClientConfig: &tls.Config{
				Certificates:       []tls.Certificate{cert},
				InsecureSkipVerify: true,
				CipherSuites:       ciperSuites,
			},
			Subprotocols: []string{"ship"},
		}

		conn, _, err := dialer.Dial(uri, nil)

		return conn, err
	}
}

func sendModel(ws *websocket.Conn, model interface{}) error {
	msg, err := json.Marshal(model)
	if err != nil {
		return err
	}

	eebusMsg, err := util.JsonIntoEEBUSJson(msg)
	if err != nil {
		return err
	}

	fmt.Println("Send: ", string(eebusMsg))

	// Wrap the message into a binary message with the ship header
	shipMsg := []byte{byte(ship.MsgTypeControl)}
	shipMsg = append(shipMsg, eebusMsg...)
	err = ws.WriteMessage(websocket.BinaryMessage, shipMsg)
	if err != nil {
		return err
	}

	return nil
}

func readMessage(ws *websocket.Conn) ([]byte, error) {
	_, b, err := ws.ReadMessage()
	if err != nil {
		return nil, err
	}

	if len(b) < 2 {
		return nil, fmt.Errorf("Invalid ship message length")
	}
	if uint(b[0]) < uint(ship.MsgTypeControl) {
		return nil, fmt.Errorf("Invalid ship message type")
	}

	// Remove header byte
	b = b[1:]

	fmt.Println("Recv: ", string(b))

	data := util.JsonFromEEBUSJson(b)

	return data, nil
}

func readModel[T any](ws *websocket.Conn, t T) (interface{}, error) {
	data, err := readMessage(ws)
	if err != nil {
		return nil, err
	}

	var model T

	if err := json.Unmarshal([]byte(data), &model); err != nil {
		return nil, err
	}

	return model, nil
}

func shipHandshake(ws *websocket.Conn) error {
	// CMI_STATE_CLIENT_SEND
	shipInit := []byte{byte(ship.MsgTypeInit), 0x00}

	err := ws.WriteMessage(websocket.BinaryMessage, shipInit)
	if err != nil {
		log.Fatal(err)
	}

	// CMI_STATE_CLIENT_WAIT
	_, msg, err := ws.ReadMessage()
	if err != nil {
		log.Fatal(err)
	}

	// CMI_STATE_CLIENT_EVALUATE
	if msg[0] != byte(ship.MsgTypeInit) {
		log.Fatal("Expected init message")
	}

	waitingDuration := uint(60000)

	// SME_HELLO_STATE_READY_INIT
	helloMsg := ship.ConnectionHello{
		ConnectionHello: ship.ConnectionHelloType{
			Phase:   ship.ConnectionHelloPhaseTypeReady,
			Waiting: &waitingDuration,
		},
	}

	if err = sendModel(ws, helloMsg); err != nil {
		log.Fatal(err)
	}

	// SME_HELLO_STATE_READY_LISTEN
	data, err := readModel(ws, ship.ConnectionHello{})
	if err != nil {
		log.Fatal(err)
	}

	helloMsg = data.(ship.ConnectionHello)
	switch helloMsg.ConnectionHello.Phase {
	case ship.ConnectionHelloPhaseTypeReady:
		fmt.Println("Got ready message")
	case ship.ConnectionHelloPhaseTypeAborted:
		log.Fatal("Connection aborted")
	default:
		log.Fatal("Unexpected connection hello phase: ", helloMsg.ConnectionHello.Phase)
	}

	// HELLO_OK
	protocolHandshake := ship.MessageProtocolHandshake{
		MessageProtocolHandshake: ship.MessageProtocolHandshakeType{
			HandshakeType: ship.ProtocolHandshakeTypeTypeAnnounceMax,
			Version:       ship.Version{Major: 1, Minor: 0},
			Formats: ship.MessageProtocolFormatsType{
				Format: []ship.MessageProtocolFormatType{ship.MessageProtocolFormatTypeUTF8},
			},
		},
	}

	if err = sendModel(ws, protocolHandshake); err != nil {
		log.Fatal(err)
	}

	data, err = readModel(ws, ship.MessageProtocolHandshake{})
	if err != nil {
		log.Fatal(err)
	}

	protocolFormat := data.(ship.MessageProtocolHandshake)
	if protocolFormat.MessageProtocolHandshake.HandshakeType != ship.ProtocolHandshakeTypeTypeSelect {
		log.Fatal("Invalid protocol handshake response: ", protocolFormat)
	}

	protocolHandshake.MessageProtocolHandshake.HandshakeType = ship.ProtocolHandshakeTypeTypeSelect
	if err = sendModel(ws, protocolHandshake); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Got protocol handshake")

	// PIN State
	pinState := ship.ConnectionPinState{
		ConnectionPinState: ship.ConnectionPinStateType{
			PinState: ship.PinStateTypeNone,
		},
	}

	if err = sendModel(ws, pinState); err != nil {
		log.Fatal(err)
	}

	data, err = readModel(ws, ship.ConnectionPinState{})
	if err != nil {
		log.Fatal(err)
	}

	pinState = data.(ship.ConnectionPinState)
	switch pinState.ConnectionPinState.PinState {
	case ship.PinStateTypeNone:
		fmt.Println("Got pin state: none")
	case ship.PinStateTypeRequired:
		log.Fatal("Got pin state: required (unsupported)")
	case ship.PinStateTypeOptional:
		log.Fatal("Got pin state: optional (unsupported)")
	case ship.PinStateTypePinOk:
		fmt.Println("Got pin state: ok (unsupported)")
	default:
		log.Fatal("Got invalid pin state: ", pinState.ConnectionPinState.PinState)
	}

	// Access Methods
	accessMethodsRequest := ship.AccessMethodsRequest{
		AccessMethodsRequest: ship.AccessMethodsRequestType{},
	}

	if err = sendModel(ws, accessMethodsRequest); err != nil {
		log.Fatal(err)
	}

	for {
		data, err := readMessage(ws)
		if err != nil {
			log.Fatal(err)
		}

		dataString := string(data)

		if strings.Contains(dataString, "\"accessMethodsRequest\":{") {
			fmt.Println("Got access methods request")
			methodsId := "Test"

			accessMethods := ship.AccessMethods{
				AccessMethods: ship.AccessMethodsType{
					Id: &methodsId,
				},
			}
			if err = sendModel(ws, accessMethods); err != nil {
				log.Fatal(err)
			}
		} else if strings.Contains(dataString, "\"accessMethods\":{") {
			fmt.Println("Got access methods")
			break
		} else {
			log.Fatal("access methods: invalid response: ", dataString)
		}
	}

	// Spine
	if _, err = readMessage(ws); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Got SHIP message with SPINE payload")

	ws.Close()

	return nil
}

func setupWebsocketClient(certificate tls.Certificate, host, port string) {
	var ciperSuites = []uint16{
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
	}

	dialer := &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 5 * time.Second,
		TLSClientConfig: &tls.Config{
			Certificates:       []tls.Certificate{certificate},
			InsecureSkipVerify: true,
			CipherSuites:       ciperSuites,
		},
		Subprotocols: []string{"ship"},
	}

	address := fmt.Sprintf("wss://%s:%s", host, port)
	ws, _, err := dialer.Dial(address, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	if err = shipHandshake(ws); err != nil {
		log.Fatal(err)
	}
}

func connectToEEBUSService(host, port string) {
	certificate, err := createCertificate()
	if err != nil {
		log.Fatal(err)
	}
	setupWebsocketClient(certificate, host, port)
}

func usage() {
	fmt.Println("Usage: {} <command>", os.Args[0])
	fmt.Println("Commands:")
	fmt.Println("  connect <host> <port>")
}

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}

	command := os.Args[1]
	switch command {
	case "connect":
		if len(os.Args) < 4 {
			fmt.Println("Usage: {} connect <host> <port>", os.Args[0])
			return
		}
		host := os.Args[2]
		port := os.Args[3]
		connectToEEBUSService(host, port)
	default:
		usage()
	}
}
