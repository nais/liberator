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
	"github.com/imdario/mergo"
	hash "github.com/mitchellh/hashstructure"
	nais_io_v1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
	nais_io_v1alpha1 "github.com/nais/liberator/pkg/apis/nais.io/v1alpha1"
	"github.com/nais/liberator/pkg/intutil"
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
	Status nais_io_v1.Status `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

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
	Replicas *nais_io_v1.Replicas `json:"replicas"`

	Pod PodConfig `json:"pod"`

	Azure *AzureConfig `json:"azure,omitempty"`

	// +optional
	Ingress *IngressConfig `json:"ingress,omitempty"`

	// +optional
	Egress *EgressConfig `json:"egress,omitempty"`

	UnsecureDebugDisableAllAccessPolicies bool `json:"unsecuredebugdisableallaccesspolicies,omitempty"`

	//Set this flag if the application is build onPrem, this will add the default volume mounts an AuroraApplication requires
	// +optional
	AuroraApplication bool `json:"auroraApplication,omitempty"`
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
	Name     string `json:"name,omitempty"`
	Port     uint16 `json:"port"`
	Protocol string `json:"protocol,omitempty"`
}

type PublicIngressConfig struct {
	// +optional
	Disabled bool `json:"disabled,omitempty"`
	// +optional
	Port uint16 `json:"port,omitempty"`
	// +optional
	ServicePort uint16 `json:"serviceport,omitempty"`
	// +optional
	Gateway          string `json:"gateway,omitempty"`
	HostPrefix       string `json:"hostPrefix,omitempty"`
	OverrideHostname string `json:"overrideHostname,omitempty"`
}

type InternalIngressConfig struct {
	// +optional
	Disabled bool `json:"disabled,omitempty"`
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
	// +optional
	Disabled bool `json:"disabled,omitempty"`

	Host string `json:"host"`
	// +optional
	Ports []PortConfig `json:"ports,omitempty"`
}

type InternalEgressConfig struct {
	// +optional
	Disabled bool `json:"disabled,omitempty"`
	// +optional
	Application string `json:"application,omitempty"`
	// +optional
	Namespace string `json:"namespace,omitempty"`
	// +optional
	Ports []PortConfig `json:"ports,omitempty"`
}

type AzureConfig struct {
	ResourceGroup string `json:"resourceGroup"`

	PostgreDatabases map[string]*PostgreDatabaseConfig `json:"postgresDatabase,omitempty"`
	StorageAccount   map[string]*StorageAccountConfig  `json:"storageAccount,omitempty"`
	CosmosDB map[string]*CosmosDBConfig `json:"cosmosDb,omitempty"`
}

type CosmosDBConfig struct {
	// +optional
	Disabled bool `json:"disabled,omitempty"`

	Name   string                 `json:"name"`

	//For now if you set monoDbVersion you will get mongo, else you will get GlocalDb with sql
	MongoDBVersion string `json:"mongoDBVersion,omitempty"`

}

type StorageAccountConfig struct {
	// +optional
	Disabled bool `json:"disabled,omitempty"`

	Name string `json:"name"`
}

type PostgreDatabaseConfig struct {
	// +optional
	Disabled bool `json:"disabled,omitempty"`

	Name   string                          `json:"name"`
	Server string                          `json:"server"`
	Users  map[string]*PostgreDatabaseUser `json:"users"`
}

type PostgreDatabaseUser struct {
	// +optional
	Disabled bool `json:"disabled,omitempty"`

	Name string `json:"name"`
	Role string `json:"role"`
}

type PodConfig struct {
	Image string `json:"image"`

	Command []string `json:"command,omitempty"`

	// Custom environment variables injected into your container.
	// Specify either `value` or `valueFrom`, but not both.
	Env nais_io_v1.EnvVars `json:"env,omitempty"`

	// EnvFrom exposes all variables in the ConfigMap or Secret resources as environment variables.
	// One of `configMap` or `secret` is required.
	//
	// Environment variables will take the form `KEY=VALUE`, where `key` is the ConfigMap or Secret key.
	// You can specify as many keys as you like in a single ConfigMap or Secret.
	//
	// The ConfigMap and Secret resources must live in the same Kubernetes namespace as the Application resource.
	// +nais:doc:Availability="team namespaces"
	EnvFrom []nais_io_v1.EnvFrom `json:"envFrom,omitempty"`

	// List of ConfigMap or Secret resources that will have their contents mounted into the containers as files.
	// Either `configMap` or `secret` is required.
	//
	// Files will take the path `<mountPath>/<key>`, where `key` is the ConfigMap or Secret key.
	// You can specify as many keys as you like in a single ConfigMap or Secret, and they will all
	// be mounted to the same directory.
	//
	// The ConfigMap and Secret resources must live in the same Kubernetes namespace as the Application resource.
	// +nais:doc:Availability="team namespaces"
	FilesFrom []nais_io_v1.FilesFrom `json:"filesFrom,omitempty"`

	// Prometheus is used to [scrape metrics from the pod](https://doc.nais.io/observability/metrics/).
	// Use this configuration to override the default values.
	Prometheus *PrometheusConfig `json:"prometheus,omitempty"`

	// Sometimes, applications are temporarily unable to serve traffic. For example, an application might need
	// to load large data or configuration files during startup, or depend on external services after startup.
	// In such cases, you don't want to kill the application, but you don’t want to send it requests either.
	// Kubernetes provides readiness probes to detect and mitigate these situations. A pod with containers
	// reporting that they are not ready does not receive traffic through Kubernetes Services.
	// Read more about this over at the [Kubernetes readiness documentation](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/).
	Readiness *nais_io_v1.Probe `json:"readiness,omitempty"`

	// Many applications running for long periods of time eventually transition to broken states,
	// and cannot recover except by being restarted. Kubernetes provides liveness probes to detect
	// and remedy such situations. Read more about this over at the
	// [Kubernetes probes documentation](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/).
	Liveness *nais_io_v1.Probe `json:"liveness,omitempty"`

	// Kubernetes uses startup probes to know when a container application has started. If such a probe is configured,
	// it disables liveness and readiness checks until it succeeds, making sure those probes don't interfere with the
	// application startup. This can be used to adopt liveness checks on slow starting containers, avoiding them getting
	// killed by Kubernetes before they are up and running.
	Startup *nais_io_v1.Probe `json:"startup,omitempty"`

	// When Containers have [resource requests](http://kubernetes.io/docs/user-guide/compute-resources/) specified,
	// the Kubernetes scheduler can make better decisions about which nodes to place pods on.
	Resources *nais_io_v1.ResourceRequirements `json:"resources,omitempty"`

	MinAvailable int32 `json:"minAvailable"`
}

type PrometheusConfig struct {

	//Set this to true to disable
	// +optional
	Disabled bool `json:"disabled,omitempty"`

	//Set the port you want to expose metrics on. Default is '8080'
	// +kubebuilder:default:=8080
	Port    string `json:"port,omitempty"`

	//Set the path metrics are exposed on. Default is '/metrics'
	// +kubebuilder:default:=/metricss
	Path    string `json:"path,omitempty"`

}

type ImagePolicyConfig struct {
	// +optional
	Disabled bool `json:"disabled,omitempty"`

	// +optional
	Branch string `json:"branch,omitempty"`

	// +optional
	Semver string `json:"semver,omitempty"`
}

func (in *Application) GetObjectKind() schema.ObjectKind {
	return in
}

func (in Application) GetObjectReference() v1.ObjectReference {
	return v1.ObjectReference{
		APIVersion:      "nebula.skatteetaten.no/v1alpha1",
		UID:             in.UID,
		Name:            in.Name,
		Kind:            "Application",
		ResourceVersion: in.ResourceVersion,
		Namespace:       in.Namespace,
	}
}

func (in Application) GetOwnerReference() metav1.OwnerReference {
	return metav1.OwnerReference{
		APIVersion: "nebula.skatteetaten.no/v1alpha1",
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
		if !strings.HasPrefix(k, "skatteetaten.no") {
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

// ApplyDefaults sets default values where they are missing from an Application spec.
func (app *Application) ApplyDefaults() error {
	replicasIsZero := app.replicasDefined() && app.replicasIsZero()

	err := mergo.Merge(app, getAppDefaults())
	if err != nil {
		return err
	}

	if replicasIsZero {
		app.Spec.Replicas.Min = intutil.Intp(0)
		app.Spec.Replicas.Max = intutil.Intp(0)
	}

	return nil
}

func (app *Application) replicasDefined() bool {
	if app.Spec.Replicas != nil && app.Spec.Replicas.Min != nil && app.Spec.Replicas.Max != nil {
		return true
	}
	return false
}

func (app *Application) replicasIsZero() bool {
	return *app.Spec.Replicas.Min == 0 && *app.Spec.Replicas.Max == 0
}

func getAppDefaults() *Application {
	return &Application{
		Spec: ApplicationSpec{
			Ingress:&IngressConfig {
				Public: map[string]PublicIngressConfig{
					"default": {
						Port:    8080,
						ServicePort: 80,
						Gateway: "istio-ingressgateway",
					},
				},
			},
			Pod: PodConfig{
				Prometheus: &PrometheusConfig{
					Path: "/metrics",
					Port: "8080",
				},
				Liveness: &nais_io_v1.Probe{
					PeriodSeconds:    nais_io_v1alpha1.DefaultProbePeriodSeconds,
					Timeout:          nais_io_v1alpha1.DefaultProbeTimeoutSeconds,
					FailureThreshold: nais_io_v1alpha1.DefaultProbeFailureThreshold,
				},
				Resources: &nais_io_v1.ResourceRequirements{
					Limits: &nais_io_v1.ResourceSpec{
						Cpu:    "500m",
						Memory: "512Mi",
					},
					Requests: &nais_io_v1.ResourceSpec{
						Cpu:    "200m",
						Memory: "256Mi",
					},
				},
			},
			Replicas: &nais_io_v1.Replicas{
				Min:                    intutil.Intp(2),
				Max:                    intutil.Intp(4),
				CpuThresholdPercentage: 50,
			},
		},
	}
}
