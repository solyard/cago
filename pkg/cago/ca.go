package cago

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"os"
	"time"

	log "github.com/gookit/slog"
)

var revokedCerts []pkix.RevokedCertificate

func signCACertsCheck() {
	_, privateKey, _ := loadCACerts()
	log.Info("CA certs loaded successfull. Now make test sign...")
	msg := []byte("verifiable message")
	msgHash := sha256.New()
	_, err := msgHash.Write(msg)
	if err != nil {
		panic(err)
	}
	msgHashSum := msgHash.Sum(nil)
	signed, err := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, msgHashSum, nil)
	if err != nil {
		log.Errorf("Error while sing test data with CA Private Key. Error ", err.Error())
		log.Exit(1)
	}

	log.Debugf("Test Copleted. Private Key - OK! Signature: ", string(signed))
}

func loadCACerts() (*x509.Certificate, *rsa.PrivateKey, []byte) {
	caCertificateFile, err := os.ReadFile(fmt.Sprintf("%v/tls.crt", os.Getenv("CA_CERTS_PATH")))
	if err != nil {
		log.Errorf("Error while read root CA. Error: %v", err)
		log.Exit(1)
	}
	caCertificatePem, _ := pem.Decode(caCertificateFile)

	caCertificate, err := x509.ParseCertificate(caCertificatePem.Bytes)
	if err != nil {
		log.Errorf("Error while parse CA certificate! Error: %v", err)
		log.Exit(1)
	}

	if (caCertificate.KeyUsage & x509.KeyUsageCRLSign) == 0 {
		log.Error("You are try to used the CA cert without signKeyUsage bit!")
		caCertificate.KeyUsage |= x509.KeyUsageCRLSign
	}

	privateKeyFile, err := os.ReadFile(fmt.Sprintf("%v/tls.key", os.Getenv("CA_CERTS_PATH")))
	if err != nil {
		log.Errorf("Error while read CA private key. Error: %v", err)
		log.Exit(1)
	}
	privateKeyPem, _ := pem.Decode(privateKeyFile)
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyPem.Bytes)
	if err != nil {
		log.Errorf("Error whle parse CA private key. Error: %v", err)
		log.Exit(1)
	}

	return caCertificate, privateKey, caCertificateFile
}

// GenerateCRL recieve array of certificates that need to be revoked by CA.
func GenerateCRL(revokedCert *CertificateMetadata) ([]byte, []byte) {

	caCertificate, privateKey, cacertfile := loadCACerts()

	if revokedCert.Revoked {
		certificate, _ := pem.Decode([]byte(revokedCert.Certificate))
		crt, err := x509.ParseCertificate(certificate.Bytes)
		if err != nil {
			log.Printf("Error while parse certificate. Error: %v", err)
			log.Exit(1)
		}

		clientRevocation := pkix.RevokedCertificate{
			SerialNumber:   crt.SerialNumber,
			RevocationTime: time.Now(),
		}

		revokedCerts = append(revokedCerts, clientRevocation)
	}

	tbsCertList := &x509.RevocationList{
		SignatureAlgorithm:  0,
		Number:              big.NewInt(1),
		ThisUpdate:          time.Now(),
		NextUpdate:          time.Now().Add(time.Hour * time.Duration(86400)),
		RevokedCertificates: revokedCerts,
	}

	var reader io.Reader
	crl, err := x509.CreateRevocationList(reader, tbsCertList, caCertificate, crypto.Signer(privateKey))
	if err != nil {
		log.Errorf("Error while create CRL. Error: %v", err)
		log.Exit(1)
		return nil, nil
	}
	crlPEM := pem.EncodeToMemory(&pem.Block{Type: "X509 CRL", Bytes: crl})

	return crlPEM, cacertfile
}

func getCRLFile(location string) []byte {
	if location == "local" {
		crl, err := os.ReadFile("crl.pem")
		if err != nil {
			log.Error("Cannot find crl.pem. Error: ", err.Error())
			log.Exit(1)
		}

		return crl
	} else {
		crldata, err := GetCRLFromKubernetes()
		if err != nil {
			log.Error("Error while get crl file from Kubernetes Cluster. Error: ", err)
			log.Exit(1)
		}
		crl, err := base64.StdEncoding.DecodeString(string(crldata))
		if err != nil {
			log.Error("Error while decode base64 content from Kubernetes Cluster. Error: ", err)
			log.Exit(1)
		}

		return crl
	}

}

func checkCRL(location string) {
	crl := getCRLFile(location)
	crlparsed, err := x509.ParseCRL(crl)
	if err != nil {
		log.Errorf("Error while parse CRL file ", err)
		log.Exit(1)
	}

	for _, v := range crlparsed.TBSCertList.RevokedCertificates {
		log.Infof("\nRevoked certs:\n  Serial: %v\n  Revocation Time: %v", v.SerialNumber, v.RevocationTime)
	}
}
