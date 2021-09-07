/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ApplicationSpec defines the desired state of Application
type ApplicationSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Application. Edit application_types.go to remove/update
	// +optional
	ImagePolicy *ImagePolicyConfig `json:"imagePolicy,omitempty"`

	// +optional
	// +kubebuilder:default={min: 2, minAvailable: 1}
	Replicas Replicas `json:"replicas"`

	Pod PodConfig `json:"pod"`

	Azure *AzureConfig `json:"azure,omitempty"`

	// +optional
	// +kubebuilder:default={public: {default: { enabled: true }}}
	Ingress *IngressConfig `json:"ingress,omitempty"`

	// +optional
	Egress *EgressConfig `json:"egress,omitempty"`

	// +kubebuilder:default=false
	UnsecureDebugDisableAllAccessPolicies bool `json:"unsecuredebugdisableallaccesspolicies,omitempty"`
}

type IngressConfig struct {
	Public   map[string]PublicIngressConfig   `json:"public,omitempty"`
	Internal map[string]InternalIngressConfig `json:"internal,omitempty"`
}

type EgressConfig struct {
	External map[string]ExternalEgressConfig `json:"external,omitempty"`
	Internal map[string]InternalEgressConfig `json:"internal,omitempty"`
}

type PortConfig struct {
	// +kubebuilder:default="default"
	Name string `json:"name,omitempty"`
	// +kubebuilder:default=8080
	Port uint16 `json:"port"`
	// +kubebuilder:default="TCP"
	Protocol string `json:"protocol,omitempty"`
}

type PublicIngressConfig struct {
	// +optional
	// +kubebuilder:default=true
	Enabled bool `json:"enabled"`
	// +optional
	// +kubebuilder:default=8080
	Port uint16 `json:"port,omitempty"`
	// +optional
	// +kubebuilder:default=istio-ingressgateway
	Gateway          string `json:"gateway,omitempty"`
	HostPrefix       string `json:"hostPrefix,omitempty"`
	OverrideHostname string `json:"overrideHostname,omitempty"`
}

type InternalIngressConfig struct {
	// +optional
	// +kubebuilder:default=true
	Enabled bool `json:"enabled,omitempty"`
	// +optional
	// +kubebuilder:default="*"
	Application string `json:"application,omitempty"`
	// +optional
	Namespace string `json:"namespace,omitempty"`
	// +optional
	Ports []PortConfig `json:"ports,omitempty"`
	// +optional
	Methods []string `json:"methods,omitempty"`
	// +optional
	Paths []string `json:"paths,omitempty"`
}

type ExternalEgressConfig struct {
	Host string `json:"host"`
	// +optional
	Ports []PortConfig `json:"ports,omitempty"`
}

type InternalEgressConfig struct {
	// +optional
	// +kubebuilder:default=true
	Enabled bool `json:"enabled,omitempty"`
	// +optional
	// +kubebuilder:default="*"
	Application string `json:"application,omitempty"`
	// +optional
	Namespace string `json:"namespace,omitempty"`
	// +optional
	Ports []PortConfig `json:"ports,omitempty"`
}

type AzureConfig struct {
	PostgreDatabases []*PostgreDatabaseConfig `json:"postgresDatabase,omitempty"`
}

type PostgreDatabaseConfig struct {
	Name   string                 `json:"name"`
	Server string                 `json:"server"`
	Users  []*PostgreDatabaseUser `json:"users"`
}

func (p PostgreDatabaseConfig) AzureName(application Application) string {
	return fmt.Sprintf("pgd-%s-%s-%s", application.Namespace, application.Name, p.Name)
}
func (p PostgreDatabaseConfig) AzureServerName(application Application) string {
	return fmt.Sprintf("pgs-%s-%s", application.Namespace, p.Server)
}

type PostgreDatabaseUser struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

func (p PostgreDatabaseUser) AzureName(application Application) string {
	return fmt.Sprintf("pgu-%s-%s", application.Name, p.Name)
}

func (p PostgreDatabaseUser) SecretName(application Application) string {
	return fmt.Sprintf("postgresqluser-%s", p.AzureName(application))
}

type PodConfig struct {
	Image string `json:"image"`

	//TODO: defaults
	// +optional
	// +kubebuilder:default={limits: { cpu: "500m", memory: "512Mi"}, requests: { cpu: "200m", memory: "256Mi"}}
	Resource v1.ResourceRequirements `json:"resources"`
}

type Replicas struct {

	// +kubebuilder:default:=2
	Min int `json:"min"`
	Max int `json:"max,omitempty"`

	HpaTargetCPUUtilizationPercentage int `json:"hpaTargetCPUUtilizationPercentage,omitempty"`

	// +kubebuilder:default:=1
	MinAvailable int32 `json:"minAvailable"`
}

type ImagePolicyConfig struct {
	// +optional
	// +kubebuilder:default:=true
	Enabled bool `json:"enabled"`

	// +optional
	Branch string `json:"branch,omitempty"`

	// +optional
	Semver string `json:"semver,omitempty"`
}

// ApplicationStatus contains different NAIS status properties
type ApplicationStatus struct {
	SynchronizationTime     int64  `json:"synchronizationTime,omitempty"`
	RolloutCompleteTime     int64  `json:"rolloutCompleteTime,omitempty"`
	CorrelationID           string `json:"correlationID,omitempty"`
	DeploymentRolloutStatus string `json:"deploymentRolloutStatus,omitempty"`
	SynchronizationState    string `json:"synchronizationState,omitempty"`
	SynchronizationHash     string `json:"synchronizationHash,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:shortName=app
// Application is the Schema for the applications API
type Application struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApplicationSpec   `json:"spec,omitempty"`
	Status ApplicationStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ApplicationList contains a list of Application
type ApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Application `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Application{}, &ApplicationList{})
}


func (a Application) StandardObjectMeta() metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      a.Name,
		Namespace: a.Namespace,
		Labels:    a.StandardLabels(),
	}
}

func (a Application) StandardLabelSelector() map[string]string {
	return map[string]string{
		"app": a.Name,
	}
}
func (a Application) StandardLabels() map[string]string {
	return map[string]string{
		"app": a.Name,
	}
}
