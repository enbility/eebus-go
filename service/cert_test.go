package service

import (
	"crypto/x509"
	"testing"

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

	ski, err := skiFromCertificate(leaf)
	assert.Nil(c.T(), err)
	assert.NotNil(c.T(), ski)
}
