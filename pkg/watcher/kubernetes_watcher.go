package watcher

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/cert-manager-issuer/internal/models"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/applyconfigurations/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var namespaceForWatcher = os.Getenv("WATCHER_NAMESPACE")
var kuberentesAPIPath = fmt.Sprintf("/apis/cert-manager.io/v1/namespaces/%v/certificaterequests", namespaceForWatcher)
var kubernetesEndpoint = "https://kubernetes.default.svc:443"
var resourcesURL = fmt.Sprintf("%v%v", kubernetesEndpoint, kuberentesAPIPath)

func createHttpGETRequest(url string, bearer string) *http.Request {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Errorf("Error while creating httpNewRequest. Error: %v", err)
	}
	req.Header.Add("Authorization", bearer)
	return req
}

func createHttpClient() *http.Client {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	return &http.Client{}
}

func getToken() string {
	token, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		log.Errorf("Error while getting ServiceAccount token. Error: %v", err)
	}
	bearer := "Bearer " + string(token)
	return bearer
}

func getResources(url string) *models.CertificateRequestObjects {
	bearer := getToken()
	client, req := createHttpClient(), createHttpGETRequest(url, bearer)
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Error while making request in Kubernetes API:", err)
	}

	defer resp.Body.Close()

	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var certRequestObject *models.CertificateRequestObjects
	err = json.Unmarshal(body, &certRequestObject)
	if err != nil {
		log.Errorf("Error while unmarshal data from Cert Manager. Error: %v", err)
	}

	return certRequestObject
}

func certificateResourceWatcher() []*models.CertificateMetadata {

	// Ignore untrusted CA for Kubernetes API endpoint and make request to API
	certificateRequestObjects := getResources(resourcesURL)

	var certificates []*models.CertificateMetadata

	for _, value := range certificateRequestObjects.Items {
		certificates = append(certificates,
			&models.CertificateMetadata{Name: value.Metadata.Name,
				Certificate: value.Status.Certificate,
				Issued:      true,
				Revoked:     value.Metadata.Annotations.CertificateRevoked})
	}

	return certificates
}

// StartKubernetesCertificateWatcher starting making request to Kubernetes API for
// gathering certificates in CRD of Cert Manager
func StartKubernetesCertificateWatcher() []*models.CertificateMetadata {
	return certificateResourceWatcher()
}

func getKubernetesClent() *kubernetes.Clientset {
	var clientset *kubernetes.Clientset

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println("Connection to Kubernetes Cluster was successfull")
	}

	return clientset
}

// WriteCRL function make request to Kubernetes API to create Secret with CRL file
func WriteCRL(crl []byte, ca []byte) {
	clientset := getKubernetesClent()

	crlSecret := v1.Secret("crl-ingress", namespaceForWatcher)
	crlSecret.Data = make(map[string][]byte)
	crlSecret.Data["ca.crl"] = crl
	crlSecret.Data["ca.crt"] = ca
	_, err := clientset.CoreV1().Secrets(namespaceForWatcher).Apply(context.Background(), crlSecret, metav1.ApplyOptions{FieldManager: "secret"})
	if err != nil {
		log.Fatalf("Error creating Kubernetes Secret in %v namepsace. Error: %v", namespaceForWatcher, err)
	}
}
