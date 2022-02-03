package nais_io_v1alpha1

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	hash "github.com/mitchellh/hashstructure"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	nais_io_v1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
)

const (
	DefaultSecretMountPath                 = "/var/run/secrets"
	DefaultJwkerMountPath                  = "/var/run/secrets/nais.io/jwker"
	DefaultAzureratorMountPath             = "/var/run/secrets/nais.io/azure"
	DefaultKafkaratorMountPath             = "/var/run/secrets/nais.io/kafka"
	DefaultDigdiratorIDPortenMountPath     = "/var/run/secrets/nais.io/idporten"
	DefaultDigdiratorMaskinportenMountPath = "/var/run/secrets/nais.io/maskinporten"
)

func init() {
	SchemeBuilder.Register(
		&Application{},
		&ApplicationList{},
	)
}

func GetDefaultMountPath(name string) string {
	return fmt.Sprintf("/var/run/configmaps/%s", name)
}

func GetDefaultPVCMountPath(name string) string {
	return fmt.Sprintf("/var/run/pvc/%s", name)
}

// Application defines a NAIS application.
//
// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="Team",type="string",JSONPath=".metadata.labels.team"
// +kubebuilder:printcolumn:name="State",type="string",JSONPath=".status.synchronizationState"
// +kubebuilder:resource:path="applications",shortName="app",singular="application"
type Application struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApplicationSpec   `json:"spec"`
	Status nais_io_v1.Status `json:"status,omitempty"`
}

// ApplicationSpec contains the NAIS manifest.
// Please keep this list sorted for clarity.
type ApplicationSpec struct {
	// By default, no traffic is allowed between applications inside the cluster.
	// Configure access policies to explicitly allow communication between applications.
	// This is also used for granting inbound access in the context of Azure AD and TokenX clients.
	// +nais:doc:Link="https://doc.nais.io/appendix/zero-trust/";"https://doc.nais.io/security/auth/azure-ad/access-policy";"https://doc.nais.io/security/auth/tokenx/#access-policies"
	AccessPolicy *nais_io_v1.AccessPolicy `json:"accessPolicy,omitempty"`

	// Provisions and configures Azure resources.
	Azure *nais_io_v1.Azure `json:"azure,omitempty"`

	// Configuration for automatic cleanup of failing pods
	Cleanup *nais_io_v1.Cleanup `json:"cleanup,omitempty"`

	// Override command when starting Docker image.
	Command []string `json:"command,omitempty"`

	// To get your own Elastic Search instance head over to the IaC-repo to provision each instance.
	// See [navikt/aiven-iac](https://github.com/navikt/aiven-iac) repository.
	Elastic *nais_io_v1.Elastic `json:"elastic,omitempty"`

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

	// +nais:doc:Availability="GCP"
	GCP *nais_io_v1.GCP `json:"gcp,omitempty"`

	// Configures an ID-porten client for this application.
	// See [ID-porten](https://doc.nais.io/security/auth/idporten/) for more details.
	IDPorten *nais_io_v1.IDPorten `json:"idporten,omitempty"`

	// Your application's Docker image location and tag.
	Image string `json:"image"`

	// List of URLs that will route HTTPS traffic to the application.
	// All URLs must start with `https://`. Domain availability differs according to which environment your application is running in.
	// +nais:doc:Link="https://doc.nais.io/clusters/gcp/";"https://doc.nais.io/clusters/on-premises/"
	Ingresses []nais_io_v1.Ingress `json:"ingresses,omitempty"`

	// An InfluxDB via Aiven. A typical use case for influxdb is to store metrics from your application and visualize them in Grafana.
	// +nais:doc:Availability="GCP"
	Influx *nais_io_v1.Influx `json:"influx,omitempty"`

	// Set up Aiven Kafka for your application.
	// +nais:doc:Link="https://doc.nais.io/persistence/kafka/"
	Kafka *nais_io_v1.Kafka `json:"kafka,omitempty"`

	// If true, an HTTP endpoint will be available at `$ELECTOR_PATH` that returns the current leader.
	// +nais:doc:Link="https://doc.nais.io/addons/leader-election/"
	LeaderElection bool `json:"leaderElection,omitempty"`

	// Many applications running for long periods of time eventually transition to broken states,
	// and cannot recover except by being restarted. Kubernetes provides liveness probes to detect
	// and remedy such situations. Read more about this over at the
	// [Kubernetes probes documentation](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/).
	Liveness *nais_io_v1.Probe `json:"liveness,omitempty"`

	// Format of the logs from the container. Use this if the container doesn't support
	// JSON logging and the log is in a special format that need to be parsed.
	// +kubebuilder:validation:Enum="";accesslog;accesslog_with_processing_time;accesslog_with_referer_useragent;capnslog;logrus;gokit;redis;glog;simple;influxdb;log15
	Logformat string `json:"logformat,omitempty"`

	// Extra filters for modifying log content. This can e.g. be used for setting loglevel based on http status code.
	// +kubebuilder:validation:Enum=http_loglevel;dns_loglevel
	Logtransform string `json:"logtransform,omitempty"`

	// Configures a Maskinporten client for this application.
	// See [Maskinporten](https://doc.nais.io/security/auth/maskinporten/) for more details.
	Maskinporten *nais_io_v1.Maskinporten `json:"maskinporten,omitempty"`

	// To get your own OpenSearch instance head over to the IaC-repo to provision each instance.
	// See [navikt/aiven-iac](https://github.com/navikt/aiven-iac) repository.
	OpenSearch *nais_io_v1.OpenSearch `json:"openSearch,omitempty"`

	// The port number which is exposed by the container and should receive traffic.
	// Note that ports under 1024 are unavailable.
	Port int `json:"port,omitempty"`

	// PreStopHook is called immediately before a container is terminated due to an API request or management event such as liveness/startup probe failure, preemption, resource contention, etc.
	// The handler is not called if the container crashes or exits by itself.
	// The reason for termination is passed to the handler.
	// +nais:doc:Link="https://doc.nais.io/naisjob/#handles-termination-gracefully";"https://kubernetes.io/docs/concepts/containers/container-lifecycle-hooks/#container-hooks"
	PreStopHook *nais_io_v1.PreStopHook `json:"preStopHook,omitempty"`

	// A HTTP GET will be issued to this endpoint at least once before the pod is terminated.
	// This feature is deprecated and will be removed in the next major version (nais.io/v1).
	// +nais:doc:Link="https://doc.nais.io/nais-application/#handles-termination-gracefully"
	PreStopHookPath string `json:"preStopHookPath,omitempty"`

	// Prometheus is used to [scrape metrics from the pod](https://doc.nais.io/observability/metrics/).
	// Use this configuration to override the default values.
	Prometheus *nais_io_v1.PrometheusConfig `json:"prometheus,omitempty"`

	// Sometimes, applications are temporarily unable to serve traffic. For example, an application might need
	// to load large data or configuration files during startup, or depend on external services after startup.
	// In such cases, you don't want to kill the application, but you don’t want to send it requests either.
	// Kubernetes provides readiness probes to detect and mitigate these situations. A pod with containers
	// reporting that they are not ready does not receive traffic through Kubernetes Services.
	// Read more about this over at the [Kubernetes readiness documentation](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/).
	Readiness *nais_io_v1.Probe `json:"readiness,omitempty"`

	// The numbers of pods to run in parallel.
	Replicas *nais_io_v1.Replicas `json:"replicas,omitempty"`

	// When Containers have [resource requests](http://kubernetes.io/docs/user-guide/compute-resources/) specified,
	// the Kubernetes scheduler can make better decisions about which nodes to place pods on.
	Resources *nais_io_v1.ResourceRequirements `json:"resources,omitempty"`

	// Whether or not to enable a sidecar container for secure logging.
	SecureLogs *nais_io_v1.SecureLogs `json:"secureLogs,omitempty"`

	// Specify which port and protocol is used to connect to the application in the container.
	// Defaults to HTTP on port 80.
	Service *nais_io_v1.Service `json:"service,omitempty"`

	// Whether to skip injection of NAV certificate authority bundle or not. Defaults to false.
	SkipCaBundle bool `json:"skipCaBundle,omitempty"`

	// Kubernetes uses startup probes to know when a container application has started. If such a probe is configured,
	// it disables liveness and readiness checks until it succeeds, making sure those probes don't interfere with the
	// application startup. This can be used to adopt liveness checks on slow starting containers, avoiding them getting
	// killed by Kubernetes before they are up and running.
	Startup *nais_io_v1.Probe `json:"startup,omitempty"`

	// Specifies the strategy used to replace old Pods by new ones.
	Strategy *nais_io_v1.Strategy `json:"strategy,omitempty"`

	// Provisions and configures a TokenX client for your application.
	// +nais:doc:Link="https://doc.nais.io/security/auth/tokenx/"
	TokenX *nais_io_v1.TokenX `json:"tokenx,omitempty"`

	// Provides secrets management, identity-based access, and encrypting application data for auditing of secrets
	// for applications, systems, and users.
	// +nais:doc:Link="https://github.com/navikt/vault-iac/tree/master/doc"
	// +nais:doc:Availability="on-premises"
	Vault *nais_io_v1.Vault `json:"vault,omitempty"`

	// Inject on-premises web proxy configuration into the application pod.
	// Most Linux applications should auto-detect these settings from the `$HTTP_PROXY`, `$HTTPS_PROXY` and `$NO_PROXY` environment variables (and their lowercase counterparts).
	// Java applications can start the JVM using parameters from the `$JAVA_PROXY_OPTIONS` environment variable.
	// +nais:doc:Availability="on-premises"
	WebProxy bool `json:"webproxy,omitempty"`
}

// +kubebuilder:object:root=true
type ApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Application `json:"items"`
}

func (in *Application) GetObjectKind() schema.ObjectKind {
	return in
}

func (in *Application) GetObjectReference() corev1.ObjectReference {
	return corev1.ObjectReference{
		APIVersion:      "nais.io/v1alpha1",
		UID:             in.UID,
		Name:            in.Name,
		Kind:            "Application",
		ResourceVersion: in.ResourceVersion,
		Namespace:       in.Namespace,
	}
}

func (in *Application) GetOwnerReference() metav1.OwnerReference {
	return metav1.OwnerReference{
		APIVersion: "nais.io/v1alpha1",
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

func (in *Application) LogFields() log.Fields {
	return log.Fields{
		"namespace":       in.GetNamespace(),
		"resourceversion": in.GetResourceVersion(),
		"application":     in.GetName(),
		"correlation_id":  in.Status.CorrelationID,
	}
}

func (in *Application) CorrelationID() string {
	return in.Annotations[nais_io_v1.DeploymentCorrelationIDAnnotation]
}

func (in *Application) SetDeploymentRolloutStatus(rolloutStatus string) {
	in.Status.DeploymentRolloutStatus = rolloutStatus
}

func (in *Application) DefaultSecretPath(base string) nais_io_v1.SecretPath {
	return nais_io_v1.SecretPath{
		MountPath: nais_io_v1.DefaultVaultMountPath,
		KvPath:    fmt.Sprintf("%s/%s/%s", base, in.Name, in.Namespace),
	}
}

func (in *Application) SkipDeploymentMessage() bool {
	if in.Annotations == nil {
		return false
	}
	skip, _ := strconv.ParseBool(in.Annotations[nais_io_v1.SkipDeploymentMessageAnnotation])
	return skip
}

func (in *Application) ClientID(cluster string) string {
	return fmt.Sprintf("%s:%s:%s", cluster, in.ObjectMeta.Namespace, in.ObjectMeta.Name)
}
