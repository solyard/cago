package cago

import (
	"context"
	"flag"
	"os"
	"path/filepath"

	log "github.com/gookit/slog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/applyconfigurations/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func getKubernetesClient() *kubernetes.Clientset {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Panicf("Error while loading kubecofnig file. Error: %v", err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panicf("Error while connecting to Kubernetes Cluster. Error: %v", err.Error())
	}

	return clientset
}

func GetCRLFromKubernetes() ([]byte, error) {
	namespaceForCRLSecret := os.Getenv("CRL_NAMESPACE")
	clientset := getKubernetesClient()

	secret, err := clientset.CoreV1().Secrets(namespaceForCRLSecret).Get(context.TODO(), "crl-ingress", metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return secret.Data["ca.crl"], nil
}

//writeCRLToKubernetes function make request to Kubernetes API to create Secret with CRL file
func writeCRLToKubernetes(crl []byte, ca []byte) {
	namespaceForCRLSecret := os.Getenv("CRL_NAMESPACE")
	clientset := getKubernetesClient()
	crlSecret := v1.Secret("crl-ingress", namespaceForCRLSecret)
	crlSecret.Data = make(map[string][]byte)
	crlSecret.Data["ca.crl"] = crl
	crlSecret.Data["ca.crt"] = ca
	_, err := clientset.CoreV1().Secrets(namespaceForCRLSecret).Apply(context.Background(), crlSecret, metav1.ApplyOptions{FieldManager: "secret"})
	if err != nil {
		log.Fatalf("Error creating Kubernetes Secret in %v namepsace. Error: %v", namespaceForCRLSecret, err)
	}
	log.Info("CRL file in Kubernetes Cluster was updated")
}
