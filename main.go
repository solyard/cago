package main

import (
	"os"
	"strconv"
	"time"

	"github.com/cert-manager-issuer/pkg/config"
	"github.com/cert-manager-issuer/pkg/crlmanager"
	"github.com/cert-manager-issuer/pkg/watcher"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

func checkRevokedCertificatesAndRegenerateCRL() {
	revokedCertificates := watcher.StartKubernetesCertificateWatcher()
	crl, ca := crlmanager.GenerateCRL(revokedCertificates)
	watcher.WriteCRL(crl, ca)
}
func main() {
	config.ReadConfig()
	pollingInterval := os.Getenv("POLLING_INTERVAL_SECONDS")
	pollingIntervalInt, err := strconv.Atoi(pollingInterval)
	if err != nil {
		log.Fatalf("Error while converting POLLING_INTERVAL_SECONDS to integer. Error: %v", err)
	}
	for {
		time.Sleep(time.Duration(pollingIntervalInt) * time.Second)
		log.Infof("Starting Kubernetes Certificate watcher with polling interval: %vs", pollingIntervalInt)
		checkRevokedCertificatesAndRegenerateCRL()
	}
}
