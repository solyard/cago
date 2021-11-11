package config

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func ReadConfig() {
	var envPresented bool

	_, envPresented = os.LookupEnv("WATCHER_NAMESPACE")
	if !envPresented {
		log.Fatal("Cannot find WATCHER_NAMESPACE variable. Must be missed in configuration!")
		panic("")
	}

	_, envPresented = os.LookupEnv("CA_CERTS_PATH")
	if !envPresented {
		log.Fatal("Cannot find CA_CERTS_PATH variable. Must be missed in configuration!")
	}

	_, envPresented = os.LookupEnv("POLLING_INTERVAL_SECONDS")
	if !envPresented {
		log.Fatal("Cannot find POLLING_INTERVAL_SECONDS variable. Must be missed in configuration!")
	}
}
