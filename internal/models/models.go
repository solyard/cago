package models

import "time"

type CertificateMetadata struct {
	Name        string
	Certificate string
	Issued      bool
	Revoked     string
}

type CertificateRequestObjects struct {
	APIVersion string `json:"apiVersion"`
	Items      []struct {
		APIVersion string `json:"apiVersion"`
		Kind       string `json:"kind"`
		Metadata   struct {
			Annotations struct {
				CertificateRevoked string `json:"certificate/revoked"`
			} `json:"annotations"`
			CreationTimestamp time.Time     `json:"creationTimestamp"`
			Generation        int           `json:"generation"`
			ManagedFields     []interface{} `json:"managedFields"`
			Name              string        `json:"name"`
			Namespace         string        `json:"namespace"`
			ResourceVersion   string        `json:"resourceVersion"`
			UID               string        `json:"uid"`
		} `json:"metadata"`
		Spec struct {
			Duration  string `json:"duration"`
			IssuerRef struct {
				Kind string `json:"kind"`
				Name string `json:"name"`
			} `json:"issuerRef"`
			Request string `json:"request"`
		} `json:"spec"`
		Status struct {
			Ca          string `json:"ca"`
			Certificate string `json:"certificate"`
			Conditions  []struct {
				LastTransitionTime time.Time `json:"lastTransitionTime"`
				Message            string    `json:"message"`
				Reason             string    `json:"reason"`
				Status             string    `json:"status"`
				Type               string    `json:"type"`
			} `json:"conditions"`
		} `json:"status"`
	} `json:"items"`
	Kind     string `json:"kind"`
	Metadata struct {
		Continue        string `json:"continue"`
		ResourceVersion string `json:"resourceVersion"`
	} `json:"metadata"`
}
