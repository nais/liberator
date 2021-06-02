package nais_io_v1

import (
	"github.com/nais/liberator/pkg/hash"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Naisjob defines a NAIS application.
//
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="Team",type="string",JSONPath=".metadata.labels.team"
// +kubebuilder:printcolumn:name="State",type="string",JSONPath=".status.synchronizationState"
// +kubebuilder:resource:path="naisjobs",shortName="nj",singular="naisjob"
type Naisjob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NaisjobSpec   `json:"spec"`
	Status NaisjobStatus `json:"status,omitempty"`
}

// NaisjobSpec contains the NAIS manifest.
// Please keep this list sorted for clarity.
type NaisjobSpec struct {
	// By default, no traffic is allowed between naisjobs inside the cluster.
	// Configure access policies to explicitly allow communication between naisjobs.
	// This is also used for granting inbound access in the context of Azure AD and TokenX clients.
	// +nais:doc:Link="https://doc.nais.io/appendix/zero-trust/";"https://doc.nais.io/security/auth/azure-ad/#pre-authorization";"https://doc.nais.io/security/auth/tokenx/#access-policies"
	AccessPolicy *AccessPolicy `json:"accessPolicy,omitempty"`

	// Once a Naisjob reaches activeDeadlineSeconds, all of its running Pods are terminated and the Naisjob status will become type: Failed with reason: DeadlineExceeded.
	ActiveDeadlineSeconds int `json:"activeDeadlineSeconds,omitempty"`

	// Provisions and configures Azure resources.
	Azure *Azure `json:"azure,omitempty"`

	// Specify the number of retries before considering a Naisjob as failed
	BackoffLimit int `json:"backoffLimit,omitempty"`

	// Override command when starting Docker image.
	Command []string `json:"command,omitempty"`

	Elastic *Elastic `json:"elastic,omitempty"`

	// Custom environment variables injected into your container.
	// Specify either `value` or `valueFrom`, but not both.
	Env EnvVars `json:"env,omitempty"`

	// EnvFrom exposes all variables in the ConfigMap or Secret resources as environment variables.
	// One of `configMap` or `secret` is required.
	//
	// Environment variables will take the form `KEY=VALUE`, where `key` is the ConfigMap or Secret key.
	// You can specify as many keys as you like in a single ConfigMap or Secret.
	//
	// The ConfigMap and Secret resources must live in the same Kubernetes namespace as the Application resource.
	// +nais:doc:Availability="team namespaces"
	EnvFrom []EnvFrom `json:"envFrom,omitempty"`

	// List of ConfigMap or Secret resources that will have their contents mounted into the containers as files.
	// Either `configMap` or `secret` is required.
	//
	// Files will take the path `<mountPath>/<key>`, where `key` is the ConfigMap or Secret key.
	// You can specify as many keys as you like in a single ConfigMap or Secret, and they will all
	// be mounted to the same directory.
	//
	// The ConfigMap and Secret resources must live in the same Kubernetes namespace as the Application resource.
	// +nais:doc:Availability="team namespaces"
	FilesFrom []FilesFrom `json:"filesFrom,omitempty"`

	// +nais:doc:Availability="GCP"
	GCP *GCP `json:"gcp,omitempty"`

	// Your application's Docker image location and tag.
	Image string `json:"image"`

	// Enable Aiven Kafka for your application.
	Kafka *Kafka `json:"kafka,omitempty"`

	// Many applications running for long periods of time eventually transition to broken states,
	// and cannot recover except by being restarted. Kubernetes provides liveness probes to detect
	// and remedy such situations. Read more about this over at the
	// [Kubernetes probes documentation](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/).
	Liveness *Probe `json:"liveness,omitempty"`

	// Format of the logs from the container. Use this if the container doesn't support
	// JSON logging and the log is in a special format that need to be parsed.
	// +kubebuilder:validation:Enum="";accesslog;accesslog_with_processing_time;accesslog_with_referer_useragent;capnslog;logrus;gokit;redis;glog;simple;influxdb;log15
	Logformat string `json:"logformat,omitempty"`

	// Extra filters for modifying log content. This can e.g. be used for setting loglevel based on http status code.
	// +kubebuilder:validation:Enum=http_loglevel;dns_loglevel
	Logtransform string `json:"logtransform,omitempty"`

	// Configures a Maskinporten client for this application.
	// See [Maskinporten](https://doc.nais.io/security/auth/maskinporten/) for more details.
	Maskinporten *Maskinporten `json:"maskinporten,omitempty"`

	// A HTTP GET will be issued to this endpoint at least once before the pod is terminated.
	// +nais:doc:Link="https://doc.nais.io/nais-application/#handles-termination-gracefully"
	PreStopHookPath string `json:"preStopHookPath,omitempty"`

	// Sometimes, applications are temporarily unable to serve traffic. For example, an application might need
	// to load large data or configuration files during startup, or depend on external services after startup.
	// In such cases, you don't want to kill the application, but you don’t want to send it requests either.
	// Kubernetes provides readiness probes to detect and mitigate these situations. A pod with containers
	// reporting that they are not ready does not receive traffic through Kubernetes Services.
	// Read more about this over at the [Kubernetes readiness documentation](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/).
	Readiness *Probe `json:"readiness,omitempty"`

	// When Containers have [resource requests](http://kubernetes.io/docs/user-guide/compute-resources/) specified,
	// the Kubernetes scheduler can make better decisions about which nodes to place pods on.
	Resources *ResourceRequirements `json:"resources,omitempty"`

	// The [Cron](https://en.wikipedia.org/wiki/Cron) schedule for running the Naisjob.
	// If not specified, the Naisjob will be run as a one-shot Job.
	Schedule string `json:"schedule,omitempty"`

	// Whether or not to enable a sidecar container for secure logging.
	SecureLogs *SecureLogs `json:"secureLogs,omitempty"`

	// Whether to skip injection of NAV certificate authority bundle or not. Defaults to false.
	SkipCaBundle bool `json:"skipCaBundle,omitempty"`

	// Kubernetes uses startup probes to know when a container application has started. If such a probe is configured,
	// it disables liveness and readiness checks until it succeeds, making sure those probes don't interfere with the
	// application startup. This can be used to adopt liveness checks on slow starting containers, avoiding them getting
	// killed by Kubernetes before they are up and running.
	Startup *Probe `json:"startup,omitempty"`

	// Specify the number of seconds to wait before removing the Job after it has finished (either Completed or Failed).
	// +nais:doc:Availability="on-premises"
	TTLSecondsAfterFinished int `json:"ttlSecondsAfterFinished,omitempty"`

	// Provides secrets management, identity-based access, and encrypting application data for auditing of secrets
	// for applications, systems, and users.
	// +nais:doc:Link="https://github.com/navikt/vault-iac/tree/master/doc"
	// +nais:doc:Availability="on-premises"
	Vault *Vault `json:"vault,omitempty"`

	// Inject on-premises web proxy configuration into the application pod.
	// Most Linux applications should auto-detect these settings from the `$HTTP_PROXY`, `$HTTPS_PROXY` and `$NO_PROXY` environment variables (and their lowercase counterparts).
	// Java applications can start the JVM using parameters from the `$JAVA_PROXY_OPTIONS` environment variable.
	// +nais:doc:Availability="on-premises"
	WebProxy bool `json:"webproxy,omitempty"`
}

// NaisjobStatus contains different NAIS status properties
type NaisjobStatus struct {
	SynchronizationTime     int64  `json:"synchronizationTime,omitempty"`
	RolloutCompleteTime     int64  `json:"rolloutCompleteTime,omitempty"`
	CorrelationID           string `json:"correlationID,omitempty"`
	DeploymentRolloutStatus string `json:"deploymentRolloutStatus,omitempty"`
	SynchronizationState    string `json:"synchronizationState,omitempty"`
	SynchronizationHash     string `json:"synchronizationHash,omitempty"`
}

func (in *Naisjob) Hash() (string, error) {
	return hash.Hash(in.Spec)
}
