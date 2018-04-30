package core

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"time"
)

//Modified https://github.com/Shyp/generate-tls-cert
func GenP256(hosts []string) (key, cert []byte, err error) {
	notBefore := time.Now()
	notAfter := notBefore.Add(time.Hour * 24 * 366)
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return
	}
	serverKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return
	}

	serverTemplate := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"WebSocks"},
			CommonName:   "WebSocks Server CA",
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA: true,
	}

	for _, host := range hosts {
		if ip := net.ParseIP(host); ip != nil {
			serverTemplate.IPAddresses = append(serverTemplate.IPAddresses, ip)
		} else {
			serverTemplate.DNSNames = append(serverTemplate.DNSNames, host)
		}
	}

	serverCert, err := x509.CreateCertificate(rand.Reader, &serverTemplate, &serverTemplate, &serverKey.PublicKey, serverKey)
	if err != nil {
		return
	}

	x509Key, err := x509.MarshalECPrivateKey(serverKey)
	if err != nil {
		return
	}

	key = pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: x509Key})
	cert = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: serverCert})
	return
}
