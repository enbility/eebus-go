package cert

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestCertSuite(t *testing.T) {
	suite.Run(t, new(CertSuite))
}

type CertSuite struct {
	suite.Suite
}

func (c *CertSuite) Test_CreateCertificate() {
	cert, err := CreateCertificate("", "Org", "DE", "CN")
	assert.Nil(c.T(), err)
	assert.NotNil(c.T(), cert)
}

func (c *CertSuite) Test_SkiFromCertificate() {
	cert, err := CreateCertificate("", "Org", "DE", "CN")
	assert.Nil(c.T(), err)

	leaf, err := x509.ParseCertificate(cert.Certificate[0])
	assert.Nil(c.T(), err)

	ski, err := SkiFromCertificate(leaf)
	assert.Nil(c.T(), err)
	assert.NotEqual(c.T(), "", ski)

	cert, err = createInvalidCertificate("unit", "org", "DE", "CN")
	assert.Nil(c.T(), err)

	leaf, err = x509.ParseCertificate(cert.Certificate[0])
	assert.Nil(c.T(), err)

	ski, err = SkiFromCertificate(leaf)
	assert.NotNil(c.T(), err)
	assert.Equal(c.T(), "", ski)
}

func createInvalidCertificate(organizationalUnit, organization, country, commonName string) (tls.Certificate, error) {
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
		SubjectKeyId:          ski[:19],
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
