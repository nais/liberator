package tlsutil

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"

	"k8s.io/client-go/util/keyutil"
)

func tlsBytesFromFiles(certPath, keyPath, caPath string) (cert, key, ca []byte, err error) {
	cert, err = ioutil.ReadFile(certPath)
	if err != nil {
		err = fmt.Errorf("read TLS certificate file %s: %s", certPath, err)
		return
	}

	key, err = ioutil.ReadFile(keyPath)
	if err != nil {
		err = fmt.Errorf("read TLS key file %s: %s", keyPath, err)
		return
	}

	ca, err = ioutil.ReadFile(caPath)
	if err != nil {
		err = fmt.Errorf("read TLS CA certificate file %s: %s", caPath, err)
		return
	}

	return
}

// Return a TLS config object based on files containing certificate, private key, and root CA.
func TLSConfigFromFiles(certPath, keyPath, caPath string) (*tls.Config, error) {
	cert, key, ca, err := tlsBytesFromFiles(certPath, keyPath, caPath)
	if err != nil {
		return nil, err
	}

	return TLSConfigFromBytes(cert, key, ca)
}

func TLSConfigFromBytes(certificate, key, ca []byte) (*tls.Config, error) {
	cert, _ := pem.Decode(certificate)
	if cert == nil {
		return nil, fmt.Errorf("unable to parse certificate: no PEM data found")
	}

	privateKey, err := keyutil.ParsePrivateKeyPEM(key)
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %s", err)
	}

	cablock, _ := pem.Decode(ca)
	if cablock == nil {
		return nil, fmt.Errorf("unable to parse CA certificate: no PEM data found")
	}

	cacert, err := x509.ParseCertificate(cablock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse CA certificate: %s", err)
	}

	certpool := x509.NewCertPool()
	certpool.AddCert(cacert)

	return &tls.Config{
		Certificates: []tls.Certificate{
			{
				Certificate: [][]byte{cert.Bytes},
				PrivateKey:  privateKey,
			},
		},
		RootCAs:            certpool,
		InsecureSkipVerify: false,
	}, nil
}
