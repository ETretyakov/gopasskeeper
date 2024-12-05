package secretstore

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	"github.com/pkg/errors"
	"google.golang.org/grpc/credentials"
)

func applyTLS(certPath string) (credentials.TransportCredentials, error) {
	caCert, err := os.ReadFile(certPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read certificate file")
	}

	caPool := x509.NewCertPool()
	if ok := caPool.AppendCertsFromPEM(caCert); !ok {
		return nil, errors.New("failed to append CA certificate to CA pool")
	}

	tlsConfig := &tls.Config{
		RootCAs:    caPool,
		MinVersion: tls.VersionTLS13,
	}

	creds := credentials.NewTLS(tlsConfig)

	return creds, nil
}
