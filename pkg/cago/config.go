package cago

import (
	"os"

	log "github.com/gookit/slog"
)

func ReadConfig() {
	var envPresented bool

	_, envPresented = os.LookupEnv("CRL_NAMESPACE")
	if !envPresented {
		log.Fatal("Cannot find CRL_NAMESPACE variable. Must be missed in configuration!")
		panic("")
	}

	_, envPresented = os.LookupEnv("CA_CERTS_PATH")
	if !envPresented {
		log.Fatal("Cannot find CA_CERTS_PATH variable. Must be missed in configuration!")
	}
}
