package bigquery_cnrm_cloud_google_com_v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	SchemeBuilder.Register(
		&BigQueryDataset{},
		&BigQueryDatasetList{},
	)
}

type BigQueryDatasetAccess struct {
	// Hardcoded read-write/admin role for service account
	Role string `json:"role"`
	// Email of service account (from GCP Team project) used to access the dataset
	UserByEmail string `json:"userByEmail"`
}

// This map (and its members) are inspired by output of (the resource existed in GCP clusters already):
//   kubectl get crd bigquerydatasets.bigquery.cnrm.cloud.google.com -o yaml
type BigqueryDatasetSpec struct {
	// The datasetId of the resource. Used for creation and acquisition.
	ResourceID string `json:"resourceID"`
	// Physical location of GCP resource
	Location string `json:"location"`
	// Optional - Will also be shown in google cloud console (in browser)
	Description string `json:"description,omitempty"`
	// Email and role for service user given access to dataset
	Access []*BigQueryDatasetAccess `json:"access"`
}

// +kubebuilder:object:root=true
type BigQueryDataset struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              BigqueryDatasetSpec `json:"spec"`
}

// +kubebuilder:object:root=true
type BigQueryDatasetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BigqueryDatasetSpec `json:"items"`
}
