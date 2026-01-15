package mockserver

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestCertificates holds all certificate files for a test scenario.
type TestCertificates struct {
	CACertFile     string
	CAKeyFile      string
	ServerCertFile string
	ServerKeyFile  string
	ClientCertFile string
	ClientKeyFile  string
}

// SetupTestCertificates creates a full PKI for testing mTLS.
// Uses t.TempDir() for automatic cleanup.
func SetupTestCertificates(t testing.TB) *TestCertificates {
	t.Helper()
	tempDir := t.TempDir()

	cfg := &TestCertificates{
		CACertFile:     filepath.Join(tempDir, "ca.crt"),
		CAKeyFile:      filepath.Join(tempDir, "ca.key"),
		ServerCertFile: filepath.Join(tempDir, "server.crt"),
		ServerKeyFile:  filepath.Join(tempDir, "server.key"),
		ClientCertFile: filepath.Join(tempDir, "client.crt"),
		ClientKeyFile:  filepath.Join(tempDir, "client.key"),
	}

	// Generate CA
	if err := generateTestCA(cfg.CACertFile, cfg.CAKeyFile); err != nil {
		t.Fatalf("failed to generate CA: %v", err)
	}

	// Generate server certificate signed by CA
	if err := generateTestCertSignedByCA(cfg.ServerCertFile, cfg.ServerKeyFile, cfg.CACertFile, cfg.CAKeyFile, true); err != nil {
		t.Fatalf("failed to generate server certificate: %v", err)
	}

	// Generate client certificate signed by CA
	if err := generateTestCertSignedByCA(cfg.ClientCertFile, cfg.ClientKeyFile, cfg.CACertFile, cfg.CAKeyFile, false); err != nil {
		t.Fatalf("failed to generate client certificate: %v", err)
	}

	return cfg
}

// generateTestCA generates a CA certificate and key for testing.
func generateTestCA(certFile, keyFile string) error {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Test CA"},
			CommonName:   "Test CA",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		IsCA:                  true,
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return err
	}

	return writeCertAndKey(certFile, keyFile, derBytes, priv)
}

// generateTestCertSignedByCA generates a certificate signed by the given CA.
func generateTestCertSignedByCA(certFile, keyFile, caCertFile, caKeyFile string, isServer bool) error {
	// Load CA
	caCertPEM, err := os.ReadFile(caCertFile)
	if err != nil {
		return err
	}
	caKeyPEM, err := os.ReadFile(caKeyFile)
	if err != nil {
		return err
	}

	caCertBlock, _ := pem.Decode(caCertPEM)
	caCert, err := x509.ParseCertificate(caCertBlock.Bytes)
	if err != nil {
		return err
	}

	caKeyBlock, _ := pem.Decode(caKeyPEM)
	caKey, err := x509.ParsePKCS1PrivateKey(caKeyBlock.Bytes)
	if err != nil {
		return err
	}

	// Generate new key
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject: pkix.Name{
			Organization: []string{"Test"},
			CommonName:   "localhost",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost"},
		IPAddresses:           []net.IP{net.IPv4(127, 0, 0, 1)},
	}

	if isServer {
		template.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	} else {
		template.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, caCert, &priv.PublicKey, caKey)
	if err != nil {
		return err
	}

	return writeCertAndKey(certFile, keyFile, derBytes, priv)
}

func writeCertAndKey(certFile, keyFile string, certDER []byte, key *rsa.PrivateKey) error {
	certOut, err := os.Create(certFile)
	if err != nil {
		return err
	}
	defer func() { _ = certOut.Close() }()
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certDER}); err != nil {
		return err
	}

	keyOut, err := os.Create(keyFile)
	if err != nil {
		return err
	}
	defer func() { _ = keyOut.Close() }()
	if err := pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}); err != nil {
		return err
	}

	return nil
}
