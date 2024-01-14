package cert

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"fmt"
	"math/big"
	"time"
)

var CiperSuites = []uint16{
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256, // SHIP 9.1: required cipher suite
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256, // SHIP 9.1: optional cipher suite
}

// Create a ship compatible self signed certificate
// organizationalUnit is the OU of the certificate
// organization is the O of the certificate
// country is the C of the certificate
// commonName is the CN of the certificate
// Example for commonName: "deviceModel-deviceSerialNumber"
func CreateCertificate(organizationalUnit, organization, country, commonName string) (tls.Certificate, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return tls.Certificate{}, err
	}

	// Create the EEBUS service SKI using the private key
	asn1, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return tls.Certificate{}, err
	}
	// SHIP 12.2: Required to be created according to RFC 3280 4.2.1.2
	ski := sha1.Sum(asn1)

	subject := pkix.Name{
		OrganizationalUnit: []string{organizationalUnit},
		Organization:       []string{organization},
		Country:            []string{country},
		CommonName:         commonName,
	}

	// Create a random serial big int value
	maxValue := new(big.Int)
	maxValue.Exp(big.NewInt(2), big.NewInt(130), nil).Sub(maxValue, big.NewInt(1))
	serialNumber, err := rand.Int(rand.Reader, maxValue)
	if err != nil {
		return tls.Certificate{}, err
	}

	template := x509.Certificate{
		SignatureAlgorithm:    x509.ECDSAWithSHA256,
		SerialNumber:          serialNumber,
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
		Certificate:                  [][]byte{certBytes},
		PrivateKey:                   privateKey,
		SupportedSignatureAlgorithms: []tls.SignatureScheme{tls.ECDSAWithP256AndSHA256},
	}

	return tlsCertificate, nil
}

func SkiFromCertificate(cert *x509.Certificate) (string, error) {
	// check if the clients certificate provides a SKI
	subjectKeyId := cert.SubjectKeyId
	if len(subjectKeyId) != 20 {
		return "", errors.New("Client certificate does not provide a SKI")
	}

	return fmt.Sprintf("%0x", subjectKeyId), nil
}
