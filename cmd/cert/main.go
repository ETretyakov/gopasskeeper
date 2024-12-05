package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

func prepareTLS(cert, key string) error {
	eg := &errgroup.Group{}
	eg.Go(func() error {
		if _, err := os.Stat(cert); err != nil {
			return err
		}
		return nil
	})
	eg.Go(func() error {
		if _, err := os.Stat(key); err != nil {
			return err
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Info().Msg("TLS certificate and (or) key not exists")
			if tlsErr := generateTLS(cert, key); tlsErr != nil {
				return fmt.Errorf("failed to generate TLS: %w", tlsErr)
			}
			log.Info().Msg("Successfully generated new TLS certificate and key")
			return nil
		}
		return fmt.Errorf("failed to check if cert and (or) key exists: %w", err)
	}

	log.Info().Msg("TLS certificate and key already exists, continue..")

	return nil
}

func generateTLS(cert, key string) error {
	certificate := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization: []string{"GophKeeper"},
			Country:      []string{"RU"},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		DNSNames:     []string{"localhost"},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(1, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return fmt.Errorf("failed to generate key: %w", err)
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, certificate, certificate, &privateKey.PublicKey, privateKey)
	if err != nil {
		return fmt.Errorf("failed to create cert: %w", err)
	}

	var certPEM bytes.Buffer
	err = pem.Encode(&certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	if err != nil {
		return fmt.Errorf("failed to encode certificate: %w", err)
	}

	var privateKeyPEM bytes.Buffer
	err = pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err != nil {
		return fmt.Errorf("failed to encode private key: %w", err)
	}

	certFile, err := os.OpenFile(cert, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open file %w", err)
	}
	defer func() {
		err = certFile.Close()
		if err != nil {
			log.Error().Err(err).Msg("failed to close file")
		}
	}()
	if _, err = certFile.Write(certPEM.Bytes()); err != nil {
		return fmt.Errorf("failed to write cert file: %w", err)
	}

	keyFile, err := os.OpenFile(key, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open file %w", err)
	}
	defer func() {
		err = keyFile.Close()
		if err != nil {
			log.Error().Err(err).Msg("failed to close file")
		}
	}()
	if _, err = keyFile.Write(privateKeyPEM.Bytes()); err != nil {
		return fmt.Errorf("failed to write key file: %w", err)
	}

	return nil
}

func main() {
	prepareTLS("certs/cert.pem", "certs/key")
}
