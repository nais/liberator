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
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	hash "github.com/mitchellh/hashstructure"
	nais_io_v1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

func init() {
	SchemeBuilder.Register(&Application{}, &ApplicationList{})
}


// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="State",type="string",JSONPath=".status.synchronizationState"
// +kubebuilder:resource:path="applications",shortName="nap",singular="application"
// Application is the Schema for the applications API
type Application struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApplicationSpec   `json:"spec,omitempty"`
	Status ApplicationStatus `json:"status,omitempty"`
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

// ApplicationList contains a list of Application
type ApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Application `json:"items"`
}

// ApplicationSpec defines the desired state of Application
type ApplicationSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Application. Edit application_types.go to remove/update
	// +optional
	ImagePolicy *ImagePolicyConfig `json:"imagePolicy,omitempty"`

	// +optional
	Replicas nais_io_v1.Replicas `json:"replicas"`

	Pod PodConfig `json:"pod"`

	Azure *AzureConfig `json:"azure,omitempty"`

	// +optional
	Ingress *IngressConfig `json:"ingress,omitempty"`

	// +optional
	Egress *EgressConfig `json:"egress,omitempty"`

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
	Name string `json:"name,omitempty"`
	Port uint16 `json:"port"`
	Protocol string `json:"protocol,omitempty"`
}

type PublicIngressConfig struct {
	// +optional
	Enabled bool `json:"enabled"`
	// +optional
	Port uint16 `json:"port,omitempty"`
	// +optional
	Gateway          string `json:"gateway,omitempty"`
	HostPrefix       string `json:"hostPrefix,omitempty"`
	OverrideHostname string `json:"overrideHostname,omitempty"`
}

type InternalIngressConfig struct {
	// +optional
	Enabled bool `json:"enabled,omitempty"`
	// +optional
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
	Enabled bool `json:"enabled,omitempty"`
	// +optional
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
	Resource v1.ResourceRequirements `json:"resources"`

	MinAvailable int32 `json:"minAvailable"`

}



type ImagePolicyConfig struct {
	// +optional
	Enabled bool `json:"enabled"`

	// +optional
	Branch string `json:"branch,omitempty"`

	// +optional
	Semver string `json:"semver,omitempty"`
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


func (in *Application) GetObjectKind() schema.ObjectKind {
	return in
}

func (in Application) GetObjectReference() v1.ObjectReference {
	return v1.ObjectReference{
		APIVersion:      "v1alpha1",
		UID:             in.UID,
		Name:            in.Name,
		Kind:            "Application",
		ResourceVersion: in.ResourceVersion,
		Namespace:       in.Namespace,
	}
}

func (in Application) GetOwnerReference() metav1.OwnerReference {
	return metav1.OwnerReference{
		APIVersion: "application.nebula.skatteetaten.no/v1alpha1",
		Kind:       "Application",
		Name:       in.Name,
		UID:        in.UID,
	}
}

func (in Application) Hash() (string, error) {
	// struct including the relevant fields for
	// creating a hash of an Application object
	var changeCause string
	if in.Annotations != nil {
		changeCause = in.Annotations["kubernetes.io/change-cause"]
	}
	relevantValues := struct {
		AppSpec     ApplicationSpec
		Labels      map[string]string
		ChangeCause string
	}{
		in.Spec,
		nil,
		changeCause,
	}

	// Exempt labels starting with 'nais.io/' from hash generation.
	// This is neccessary to avoid app re-sync because of automated NAIS processes.
	for k, v := range in.Labels {
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

func (in Application) LogFields() log.Fields {
	return log.Fields{
		"namespace":       in.GetNamespace(),
		"resourceversion": in.GetResourceVersion(),
		"application":     in.GetName(),
		"correlation_id":  in.Status.CorrelationID,
	}
}

// If the application was not deployed with a correlation ID annotation,
// generate a random UUID and add it to annotations.
func (in Application) EnsureCorrelationID() error {
	if in.Annotations == nil {
		in.SetAnnotations(map[string]string{})
	}

	if len(in.Annotations[nais_io_v1.DeploymentCorrelationIDAnnotation]) != 0 {
		return nil
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("generate deployment correlation ID: %s", err)
	}

	in.Annotations[nais_io_v1.DeploymentCorrelationIDAnnotation] = id.String()

	return nil
}

func (in Application) CorrelationID() string {
	return in.Annotations[nais_io_v1.DeploymentCorrelationIDAnnotation]
}

func (in Application) SetDeploymentRolloutStatus(rolloutStatus string) {
	in.Status.DeploymentRolloutStatus = rolloutStatus
}

func (in Application) DefaultSecretPath(base string) nais_io_v1.SecretPath {
	return nais_io_v1.SecretPath{
		MountPath: nais_io_v1.DefaultVaultMountPath,
		KvPath:    fmt.Sprintf("%s/%s/%s", base, in.Name, in.Namespace),
	}
}

func (in Application) SkipDeploymentMessage() bool {
	if in.Annotations == nil {
		return false
	}
	skip, _ := strconv.ParseBool(in.Annotations[nais_io_v1.SkipDeploymentMessageAnnotation])
	return skip
}

func (in Application) ClientID(cluster string) string {
	return fmt.Sprintf("%s:%s:%s", cluster, in.ObjectMeta.Namespace, in.ObjectMeta.Name)
}
