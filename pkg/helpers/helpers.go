package helpers

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

type CertificateInfo struct {
	Algorithm string
	BitLength int
	Mode      string
}

func ValidateCertificate(tlsCrt, tlsKey []byte) error {
	// Decode the certificate
	certBlock, _ := pem.Decode(tlsCrt)
	if certBlock == nil || certBlock.Type != "CERTIFICATE" {
		return fmt.Errorf("failed to decode certificate PEM block")
	}

	// Parse the certificate
	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse certificate: %v", err)
	}

	// Decode the private key
	keyBlock, _ := pem.Decode(tlsKey)
	if keyBlock == nil {
		return fmt.Errorf("failed to decode key PEM block")
	}

	// Parse the private key
	_, err = x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
	if err != nil {
		_, err = x509.ParseECPrivateKey(keyBlock.Bytes)
		if err != nil {
			_, err = x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
			if err != nil {
				return fmt.Errorf("failed to parse private key: %v", err)
			}
		}
	}

	// Create a certificate pool and add the certificate
	certPool := x509.NewCertPool()
	certPool.AddCert(cert)

	// Verify the certificate
	opts := x509.VerifyOptions{
		Roots: certPool,
	}
	if _, err := cert.Verify(opts); err != nil {
		return fmt.Errorf("failed to verify certificate: %v", err)
	}

	return nil
}

func GetCertificateInfo(tlsCrt, tlsKey []byte) (*CertificateInfo, error) {
	// Decode the certificate
	certBlock, _ := pem.Decode(tlsCrt)
	if certBlock == nil || certBlock.Type != "CERTIFICATE" {
		return &CertificateInfo{}, fmt.Errorf("failed to decode certificate PEM block")
	}

	// Parse the certificate
	_, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return &CertificateInfo{}, fmt.Errorf("failed to parse certificate: %v", err)
	}

	// Decode the private key
	keyBlock, _ := pem.Decode(tlsKey)
	if keyBlock == nil {
		return &CertificateInfo{}, fmt.Errorf("failed to decode key PEM block")
	}

	// Parse the private key
	var keyType string
	var bitLength int

	key, err := x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
	if err != nil {
		key, err = x509.ParseECPrivateKey(keyBlock.Bytes)
		if err != nil {
			key, err = x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
			if err != nil {
				return &CertificateInfo{}, fmt.Errorf("failed to parse private key: %v", err)
			}
		}
	}

	switch k := key.(type) {
	case *rsa.PrivateKey:
		keyType = "RSA"
		bitLength = k.N.BitLen()
	case *ecdsa.PrivateKey:
		keyType = "ECDSA"
		bitLength = k.Params().BitSize
	case ed25519.PrivateKey:
		keyType = "Ed25519"
		bitLength = len(k) * 8
	default:
		return &CertificateInfo{}, fmt.Errorf("unknown key type")
	}

	// For the mode, you typically get this from the encryption context, not directly from the certificate or key.
	// Here, we'll assume a common mode used with the given key type.
	var mode string
	if keyType == "RSA" {
		mode = "CBC" // Common mode for RSA
	} else if keyType == "ECDSA" {
		mode = "GCM" // Common mode for ECDSA
	} else if keyType == "Ed25519" {
		mode = "Stream" // Common mode for Ed25519
	}

	return &CertificateInfo{
		Algorithm: keyType,
		BitLength: bitLength,
		Mode:      mode,
	}, err
}
