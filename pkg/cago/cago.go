package cago

import (
	"os"
)

type CertificateMetadata struct {
	Name        string
	Certificate []byte
	Revoked     bool
}

func revokeCertificate(pem []byte, kubernetes bool) {
	revokedCertificate := CertificateMetadata{
		Name:        "Test",
		Certificate: pem,
		Revoked:     true,
	}

	crlPEM, cafile := GenerateCRL(&revokedCertificate)
	if !kubernetes {
		os.WriteFile("crl/crl.pem", crlPEM, 0744)
		os.WriteFile("crl/ca.crt", cafile, 0744)
	} else {
		writeCRLToKubernetes(crlPEM, cafile)
	}
}
