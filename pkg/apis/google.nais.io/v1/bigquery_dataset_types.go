package google_nais_io_v1

import (
	"encoding/json"
	"fmt"
	"strings"

	hash "github.com/mitchellh/hashstructure"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DatasetAccess struct {
	// +kubebuilder:validation:Enum=READER;WRITER;OWNER
	Role string `json:"role"`

	/* An email address of a user to grant access to. For example:
	fred@example.com. */
	UserByEmail string `json:"userByEmail"`
}

// BigQueryDatasetSpec defines the desired state of BigQueryDataset
type BigQueryDatasetSpec struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	// +kubebuilder:validation:Enum=europe-north1
	Location        string          `json:"location"`
	Access          []DatasetAccess `json:"access,omitempty"`
	Project         string          `json:"project"`
	CascadingDelete bool            `json:"cascadingDelete,omitempty"`
}

// BigQueryDatasetStatus defines the observed state of BigQueryDataset
type BigQueryDatasetStatus struct {
	SynchronizationHash string             `json:"synchronizationHash,omitempty"`
	CreationTime        int                `json:"creationTime,omitempty"`
	LastModifiedTime    int                `json:"lastModifiedTime,omitempty"`
	Conditions          []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// BigQueryDataset is the Schema for the bigquerydatasets API
type BigQueryDataset struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BigQueryDatasetSpec   `json:"spec"`
	Status BigQueryDatasetStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// BigQueryDatasetList contains a list of BigQueryDataset
type BigQueryDatasetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BigQueryDataset `json:"items"`
}

func (b BigQueryDataset) Hash() (string, error) {
	// struct including the relevant fields for
	// creating a hash of an Application object
	var changeCause string
	if b.Annotations != nil {
		changeCause = b.Annotations["kubernetes.io/change-cause"]
	}
	relevantValues := struct {
		Spec        BigQueryDatasetSpec
		Labels      map[string]string
		ChangeCause string
	}{
		b.Spec,
		nil,
		changeCause,
	}

	// Exempt labels starting with 'nais.io/' from hash generation.
	// This is neccessary to avoid app re-sync because of automated NAIS processes.
	for k, v := range b.Labels {
		if !strings.HasPrefix(k, "nais.io/") {
			if relevantValues.Labels == nil {
				// cannot be done in initializer, as this would change existing hashes
				// fixme: do this in initializer when breaking backwards compatibility in hash
				relevantValues.Labels = make(map[string]string)
			}
			relevantValues.Labels[k] = v
		}
	}

	marshalled, err := json.Marshal(relevantValues)
	if err != nil {
		return "", err
	}
	h, err := hash.Hash(marshalled, nil)
	return fmt.Sprintf("%x", h), err
}

func init() {
	SchemeBuilder.Register(&BigQueryDataset{}, &BigQueryDatasetList{})
}
