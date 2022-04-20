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

	// If relelvent specify an ImagePolicy that will track a given branch og semver range into a stream of releases you can subscribe to
	// +optional
	ImagePolicy *ImagePolicyConfig `json:"imagePolicy,omitempty"`

	// The number of replicas to start with and maxReplicas to scale to using HPA
	// +optional
	Replicas *nais_io_v1.Replicas `json:"replicas"`

	// How the application pod of this application should run
	Pod PodConfig `json:"pod"`

	// Configure logging
	// +optional
	Logging *LogConfig `json:"logging,omitempty"`

	// Integrations this application has with Azure
	Azure *AzureConfig `json:"azure,omitempty"`

	// Configure zero-trust for incoming traficc
	// +optional
	Ingress *IngressConfig `json:"ingress,omitempty"`

	// Configure zero-trust for outgoing trafix
	// +optional
	Egress *EgressConfig `json:"egress,omitempty"`

	// Flag to disable all zero-trust configuration to debug
	UnsecureDebugDisableAllAccessPolicies bool `json:"unsecuredebugdisableallaccesspolicies,omitempty"`

	// Set this flag if the application is build onPrem, this will add the default volume mounts an AuroraApplication requires
	// +optional
	AuroraApplication bool `json:"auroraApplication,omitempty"`
}

type SplunkLoggingConfig struct {
	// Which splunk index are you logging to
	SplunkIndex string `json:"splunkIndex,omitempty"`
	// Is this configuration active? Defaults to true.
	Enabled *bool `json:"enabled,omitempty"`
	// File pattern which identifies file to watch
	FilePattern string `json:"filePattern,omitempty"`
	// We have number of predefined source types, see fluentbit-splunk-injector doc for details. Taxerator
	// will blindly make the config, but it might not be valid in fluentbit-sidecar-injector
	SourceType string `json:"sourceType,omitempty"`
}

type LogConfig struct {
	// Define logging to splunk
	Splunk map[string]*SplunkLoggingConfig `json:"splunk,omitempty"`
	// Shall we add spring-boot logging config or not? Defaults to true.
	ActivateDefaultSpringBootLogbackConfig *bool `json:"activateDefaultSpringBootLogbackConfig,omitempty"`
	// LogDirectory indicates where the log volume shall exist and should be mounted. Defaults to /workspace/logs
	LogDirectory string `json:"logDirectory,omitempty"`
}

type IngressConfig struct {
	// Zero-trust incoming configuration for traffic coming from outside the mesh
	Public map[string]PublicIngressConfig `json:"public,omitempty"`

	// Zero-trust incoming configuration for traffic coming from inside the mesh
	Internal map[string]InternalIngressConfig `json:"internal,omitempty"`
}

type EgressConfig struct {

	// Zero-trust outgoing configuration for traffic going out of the mesh
	External map[string]ExternalEgressConfig `json:"external,omitempty"`

	// Zero-trust outgoing configuration for traffic inside the mesh
	Internal map[string]InternalEgressConfig `json:"internal,omitempty"`
}

type PortConfig struct {
	// Name of the port
	Name string `json:"name,omitempty"`
	// The port number
	Port uint16 `json:"port"`
	// The protocol of the port
	Protocol string `json:"protocol,omitempty"`
}

type PublicIngressConfig struct {

	// Set this to true to disable this ingress rule
	// +optional
	Disabled bool `json:"disabled,omitempty"`

	// The port to use, the default value if not specified is '8080'
	// +optional
	Port uint16 `json:"port,omitempty"`

	// The port in the service to point to, the default value if not specified is '80'
	// +optional
	ServicePort uint16 `json:"serviceport,omitempty"`

	// The ingress gateway to use. the default value is 'istio-ingressgateway'
	// +optional
	Gateway string `json:"gateway,omitempty"`

	// A prefix to use for the host part
	HostPrefix string `json:"hostPrefix,omitempty"`

	// Fully qualified host that can be set to override host generation
	OverrideHostname string `json:"overrideHostname,omitempty"`
}

type InternalIngressConfig struct {

	// Set this to true to disable this ingress rule
	// +optional
	Disabled bool `json:"disabled,omitempty"`

	// The application to limit traffic from
	// +optional
	Application string `json:"application,omitempty"`

	// The namespace to limit traffic from
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// The ports that are allowed to be called
	// +optional
	Ports []PortConfig `json:"ports,omitempty"`

	// The HTTP verbs that are allowed
	// +optional
	Methods []string `json:"methods,omitempty"`

	// The paths that are affected by this rule
	// +optional
	Paths []string `json:"paths,omitempty"`
}

type ExternalEgressConfig struct {

	// Set this to true to disable
	// +optional
	Disabled bool `json:"disabled,omitempty"`

	// The hostname to allow traffic to
	Host string `json:"host"`

	// The ports to allow traffic to, the default is HTTPS on port 443
	// +optional
	Ports []PortConfig `json:"ports,omitempty"`
}

type InternalEgressConfig struct {

	// Set this to true to disable
	// +optional
	Disabled bool `json:"disabled,omitempty"`

	// The internal application to allow traffic to
	// +optional
	Application string `json:"application,omitempty"`

	// The internal namespace to allow traffic to
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// Ports to alloww traffic to
	// +optional
	Ports []PortConfig `json:"ports,omitempty"`
}

type AzureConfig struct {

	// The resource group in azure to provision resources to. This field is required.
	// This is the logical name of the resource group and will be prefixed with a namepace to make it global
	ResourceGroup string `json:"resourceGroup"`

	// PostgresDatabases in azure using
	// If a single database/user is configured it will be marked as primary by default
	// The primary user will get env vars populated with the SPRING_DATASORUCE prefix.
	// Alle users will get env vars populated with a prefix of SPRING_DATASOURCE_<USER>
	// The variables populated are
	// - <PREFIX>_URL
	// - <PREFIX>_USERNAME
	// - <PREFIX>_PASSWORD
	// The prefix can be overwritten if needed by setting overridePrefix
	PostgreDatabases map[string]*PostgreDatabaseConfig `json:"postgresDatabase,omitempty"`

	// Storageaccounts in azure
	// If a single storageaccount is configured it will be marked as primary by default
	// The primary storageaccount with have the prefix AZURE_STORAGE
	// All storageaccounts will have prefix AZURE_STORAGE_<NAME>
	// The env variables populated are
	// - <PREFIX>_CONNECTIONSTRING
	// The prefix can be overwritten if needed by setting overridePrefix
	StorageAccount map[string]*StorageAccountConfig `json:"storageAccount,omitempty"`

	// CosmodDB databases in azure
	// If a single cosmosdb is configured it will be marked as primary by default
	// The primary cosmosdb with have the prefix SPRING_DATA_MONGODB if it has mongodb or COSMOSDB if not
	// All cosmosDb will have prefix SPRING_DATA_MONGODB_<NAME> if mongodb or COSMODDB_<NAME> if not
	// The env variables populated are:
	// - <PREFIX>_URI
	// - <PREFIX>_DATABASE
	// The prefix can be overwritten if needed by setting overridePrefix
	CosmosDB map[string]*CosmosDBConfig `json:"cosmosDb,omitempty"`
}

type CosmosDBConfig struct {
	// Set this to true to disable
	// +optional
	Disabled bool `json:"disabled,omitempty"`

	// the name of the cosmosdb
	Name string `json:"name"`

	// If you specify multiple resources mark this as the primary one.
	// +optional
	Primary bool `json:"primary,omitempty"`

	// The version of MongoDB you want in your CosmodDB api.
	MongoDBVersion string `json:"mongoDBVersion,omitempty"`

	// Override the generated prefix
	Prefix string `json:"overridePrefix,omitempty"`
}

type StorageAccountConfig struct {

	// Set this to true to disable
	// +optional
	Disabled bool `json:"disabled,omitempty"`

	// If you specify multiple resources mark this as the primary one.
	// +optional
	Primary bool `json:"primary,omitempty"`

	// Override the generated prefix
	Prefix string `json:"overridePrefix,omitempty"`

	// The name of the storageaccount. This is the logical name in the config, the actual name in azure is generated from this.
	Name string `json:"name"`
}

type PostgreDatabaseConfig struct {
	// +optional
	Disabled bool `json:"disabled,omitempty"`

	// The name of the database. This is the logical name in the config, the actual name in azure is generated from this.
	Name string `json:"name"`

	// The name of the server. For this to work there needs to be a postgresqlserver provisioned in the namespace
	// with the name pgs-<namespace>-<server>
	// Here namespace is the name of the namespace this resource is in and server is this name
	Server string `json:"server"`

	// All the users that are in this postgres database
	Users map[string]*PostgreDatabaseUser `json:"users"`
}

type PostgreDatabaseUser struct {

	// Set this to true to disable
	// +optional
	Disabled bool `json:"disabled,omitempty"`

	// The name of the user
	Name string `json:"name"`

	// The role this user has in the database. Right now only "azure_pg_admin" is supported
	Role string `json:"role"`

	// If you specify multiple resources mark this as the primary one.
	// +optional
	Primary bool `json:"primary,omitempty"`

	// Override the generated prefix
	Prefix string `json:"overridePrefix,omitempty"`
}

type PodConfig struct {

	// The image to run for this pod
	ImageName string `json:"imageName"`

	// The tag of the image to run
	// In order to follow a stream of releases use a imagePolicy comment like the following

	// ` tag: 0.8.3 # {"$imagepolicy": "example:supernova-0.x:tag"}`
	// where the `0.x` part is the spec.imagePolicy.policySuffix stream you want to follow
	Tag string `json:"tag"`

	// The commands to send if any
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

	// +kubebuilder:default:=1
	// Minimum available pods for PodDisruptionBudget, default is 1
	MinAvailable int32 `json:"minAvailable,omitempty"`
}

type PrometheusConfig struct {

	// Set this to true to disable
	// +optional
	Disabled bool `json:"disabled,omitempty"`

	// Set the port you want to expose metrics on. Default is '8080'
	// +kubebuilder:default:="8080"
	Port string `json:"port,omitempty"`

	// Set the path metrics are exposed on. Default is '/metrics'
	// +kubebuilder:default:="/metrics"
	Path string `json:"path,omitempty"`
}

type ImagePolicyConfig struct {

	// Set this to true to disable
	// +optional
	Disabled bool `json:"disabled,omitempty"`

	// Specify either Branch or Semver, not both
	// The branch you want to create this imagePolicy for. This will create a stream of images that can be listened to
	// +optional
	Branch string `json:"branch,omitempty"`

	// Specify either Branch or Semver, not both
	// The semver expression to create imagePolicy for. This will create a stream of images that can be listened to
	// see [flux documentation](https://fluxcd.io/docs/guides/image-update/) for examples
	// +optional
	Semver string `json:"semver,omitempty"`

	NameSuffix string `json:"nameSuffix"`
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
	// This is necessary to avoid app re-sync because of automated NAIS processes.
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

	if app.Spec.Logging != nil {
		for _, splunk := range app.Spec.Logging.Splunk {
			if splunk.Enabled == nil || *splunk.Enabled {
				splunk.Enabled = &[]bool{true}[0]
			}
		}
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
			Ingress: &IngressConfig{
				Public: map[string]PublicIngressConfig{
					"default": {
						Port:        8080,
						ServicePort: 80,
						Gateway:     "istio-ingressgateway",
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
				MinAvailable: 1,
			},
			Replicas: &nais_io_v1.Replicas{
				Min:                    intutil.Intp(2),
				Max:                    intutil.Intp(4),
				CpuThresholdPercentage: 50,
			},
		},
	}
}
