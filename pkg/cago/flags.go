package cago

import (
	"flag"
	"os"

	log "github.com/gookit/slog"
)

func ReadInput() {
	pemFile := flag.String("pem", "", "Absolute path to the .pem file or .pem multilined content")
	checkCA := flag.Bool("checkca", false, "Detect CA files and check ability to revoke certificates and generate CRL")
	checkCRL := flag.Bool("checkcrl", false, "Get current CRL file and get Revoked certs list!")
	revokeCheck := flag.Bool("revoke", false, "Must be set true to revoke certificate, otherwise will do nothing")
	isFile := flag.Bool("file", false, "If you try to give .pem file into input set this flag")
	kubernetes := flag.Bool("kubernetes", false, "If set to true will try to connect to Kubernetes cluster and create secret with CRL inside os.GetEnv('CRL_NAMESPACE')")

	flag.Parse()

	processInput(pemFile, checkCA, revokeCheck, isFile, checkCRL, kubernetes)
}

func processInput(pem *string, cacheck *bool, revoke *bool, isfile *bool, checkcrl *bool, kubernetes *bool) {
	if *isfile {
		if *revoke {
			pemFile, err := os.ReadFile(*pem)
			if err != nil {
				log.Errorf("Error while reading file with provided path: %s. Error: %v", *pem, err.Error())
			}
			log.Info("User provided path to file", *pem)
			revokeCertificate(pemFile, *kubernetes)
		}
	}

	if !*isfile {
		if *revoke {
			pemFile := *pem
			revokeCertificate([]byte(pemFile), *kubernetes)
		}
	}

	if *cacheck {
		signCACertsCheck()
	}
	if *checkcrl {
		if !*kubernetes {
			checkCRL("local")
		} else {
			checkCRL("kubernetes")
		}
	}
}
