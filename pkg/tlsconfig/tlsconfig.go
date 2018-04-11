package tlsconfig

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
)

// ClientDefault returns an hardened TLS configuration for mTLS connection
func ClientDefault(certPath, keyPath, caPath string) *tls.Config {
	// Load our TLS key pair to use for authentication
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		log.Fatalln("Unable to load cert", err)
	}

	// Load our CA certificate
	clientCACert, err := ioutil.ReadFile(caPath)
	if err != nil {
		log.Fatal("Unable to open cert", err)
	}

	// Append client certificate to cert pool
	clientCertPool := x509.NewCertPool()
	clientCertPool.AppendCertsFromPEM(clientCACert)

	// Build default TLS configuration
	config := &tls.Config{
		// Perfect Forward Secrecy + ECDSA only
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		},
		// Force server cipher suites
		PreferServerCipherSuites: true,
		// TLS 1.2 only
		MinVersion: tls.VersionTLS12,
		// Client certificate to use
		Certificates: []tls.Certificate{cert},
		// Root CA of the client certificate
		RootCAs: clientCertPool,
	}

	// Parse CommonName and SubjectAlternateName
	config.BuildNameToCertificate()

	// Return configuration
	return config
}
