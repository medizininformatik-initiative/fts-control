package utils

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

// applyCertificateAuth applies certificate authentication to an existing client.
// It loads the client certificate and key, and optionally a custom CA certificate
// for server verification. Returns a new client with TLS configuration; does not
// modify the original client.
func applyCertificateAuth(client *Client, cert *CertificateAuthConfig) (*Client, error) {
	tlsCert, err := tls.LoadX509KeyPair(cert.CertFile, cert.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load certificate: %w", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
	}

	if cert.CAFile != "" {
		caCertPool, err := loadCACertPool(cert.CAFile)
		if err != nil {
			return nil, err
		}
		tlsConfig.RootCAs = caCertPool
	}

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	return client.WithHTTPClient(httpClient), nil
}

// createCertificateClient creates a new client configured for mutual TLS.
// This is a convenience function that creates a new base client and applies
// certificate authentication to it.
// Deprecated: Use applyCertificateAuth with an existing client for better
// compatibility with the client injection mechanism.
func createCertificateClient(cert *CertificateAuthConfig) (*Client, error) {
	return applyCertificateAuth(NewClient(), cert)
}

// loadCACertPool loads a CA certificate from a PEM file and returns a certificate pool.
func loadCACertPool(caFile string) (*x509.CertPool, error) {
	caCert, err := os.ReadFile(caFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA certificate: %w", err)
	}
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to parse CA certificate from %s", caFile)
	}
	return caCertPool, nil
}

// validateCertificateAuth validates the certificate authentication configuration.
// It checks that required files exist and that the certificate/key pair is valid.
// It also validates that the base URL uses HTTPS when certificate auth is configured.
func validateCertificateAuth(cert *CertificateAuthConfig) error {
	if cert.CertFile == "" {
		return fmt.Errorf("auth.certificate.cert_file is required")
	}
	if cert.KeyFile == "" {
		return fmt.Errorf("auth.certificate.key_file is required")
	}
	if err := validateFileExists(cert.CertFile, "certificate file"); err != nil {
		return err
	}
	if err := validateFileExists(cert.KeyFile, "certificate key file"); err != nil {
		return err
	}
	if cert.CAFile != "" {
		if err := validateFileExists(cert.CAFile, "CA certificate file"); err != nil {
			return err
		}
		// Validate CA certificate is valid PEM
		if _, err := loadCACertPool(cert.CAFile); err != nil {
			return err
		}
	}
	// Validate the certificate and key pair can be loaded
	if _, err := tls.LoadX509KeyPair(cert.CertFile, cert.KeyFile); err != nil {
		return fmt.Errorf("failed to load certificate: %w", err)
	}
	return nil
}

// validateCertificateAuthRequiresHTTPS validates that the base URL uses HTTPS
// when certificate authentication is configured.
func validateCertificateAuthRequiresHTTPS(baseURL string) error {
	if baseURL == "" {
		return fmt.Errorf("api.base_url is required when using certificate authentication")
	}
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return fmt.Errorf("invalid api.base_url: %w", err)
	}
	if parsedURL.Scheme != "https" {
		return fmt.Errorf("certificate authentication requires HTTPS, but api.base_url uses %s", parsedURL.Scheme)
	}
	return nil
}
