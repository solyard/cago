package crlmanager

import (
	"crypto"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"os"
	"time"

	"github.com/cert-manager-issuer/internal/models"
	log "github.com/sirupsen/logrus"
)

var revokedCerts []pkix.RevokedCertificate

// GenerateCRL recieve array of certificates that need to be revoked by CA
func GenerateCRL(issuedCerts []*models.CertificateMetadata) ([]byte, []byte) {

	for _, cert := range issuedCerts {
		if cert.Revoked != "" {
			certDecoded, err := base64.StdEncoding.DecodeString(cert.Certificate)
			certificate, _ := pem.Decode([]byte(certDecoded))
			crt, err := x509.ParseCertificate(certificate.Bytes)
			if err != nil {
				log.Printf("Error while parse certificate. Error: %v", err)
			}

			clientRevocation := pkix.RevokedCertificate{
				SerialNumber:   crt.SerialNumber,
				RevocationTime: time.Now(),
			}

			revokedCerts = append(revokedCerts, clientRevocation)
		}
	}
	cacertfile, err := os.ReadFile(fmt.Sprintf("%v/tls.crt", os.Getenv("CA_CERTS_PATH")))
	if err != nil {
		log.Errorf("Error while read root CA. Error: %v", err)
	}

	cacertpem, _ := pem.Decode(cacertfile)
	crlSubject, err := x509.ParseCertificate(cacertpem.Bytes)
	if err != nil {
		log.Errorf("Error while parse CA certificate! Error: %v", err)
	}

	privkeyfile, err := os.ReadFile(fmt.Sprintf("%v/tls.key", os.Getenv("CA_CERTS_PATH")))
	if err != nil {
		log.Errorf("Error while read CA private key. Error: %v", err)
	}

	pk, _ := pem.Decode(privkeyfile)
	privkey, err := x509.ParsePKCS1PrivateKey(pk.Bytes)
	if err != nil {
		log.Errorf("Error whle parse CA private key. Error: %v", err)
	}

	tbsCertList := &x509.RevocationList{
		SignatureAlgorithm:  0,
		Number:              big.NewInt(1),
		ThisUpdate:          time.Now(),
		NextUpdate:          time.Now().Add(time.Hour * time.Duration(86400)),
		RevokedCertificates: revokedCerts,
	}
	var cryptorandom io.Reader
	crl, err := x509.CreateRevocationList(cryptorandom, tbsCertList, crlSubject, crypto.Signer(privkey))
	if err != nil {
		log.Errorf("Error while create CRL. Error: %v", err)
		return nil, nil
	}
	crlPEM := pem.EncodeToMemory(&pem.Block{Type: "X509 CRL", Bytes: crl})

	return crlPEM, cacertfile
}

func compareCRLWithPreviousVersion() {
	//TODO
}
