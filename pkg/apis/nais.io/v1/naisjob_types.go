package nais_io_v1

import (
	"fmt"
	"strconv"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/nais/liberator/pkg/hash"
)

func init() {
	SchemeBuilder.Register(
		&Naisjob{},
		&NaisjobList{},
	)
}

// Naisjob defines a NAIS Naisjob.
//
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:printcolumn:name="Schedule",type="string",JSONPath=".spec.schedule"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="Team",type="string",JSONPath=".metadata.labels.team"
// +kubebuilder:printcolumn:name="State",type="string",JSONPath=".status.synchronizationState"
// +kubebuilder:resource:path="naisjobs",shortName="nj",singular="naisjob"
type Naisjob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NaisjobSpec `json:"spec"`
	Status Status      `json:"status,omitempty"`
}

// NaisjobSpec contains the NAIS manifest.
// Please keep this list sorted for clarity.
type NaisjobSpec struct {
	// By default, no traffic is allowed between naisjobs inside the cluster.
	// Configure access policies to explicitly allow communication between naisjobs.
	// This is also used for granting inbound access in the context of Azure AD and TokenX clients.
	// +nais:doc:Link="https://doc.nais.io/appendix/zero-trust/";"https://doc.nais.io/security/auth/azure-ad/access-policy";"https://doc.nais.io/security/auth/tokenx/#access-policies"
	AccessPolicy *AccessPolicy `json:"accessPolicy,omitempty"`

	// Once a Naisjob reaches activeDeadlineSeconds, all of its running Pods are terminated and the Naisjob status will become type: Failed with reason: DeadlineExceeded.
	// If set, this takes presedence over BackoffLimit.
	ActiveDeadlineSeconds *int64 `json:"activeDeadlineSeconds,omitempty"`

	// Provisions and configures Azure resources.
	Azure *AzureNaisJob `json:"azure,omitempty"`

	// Specify the number of retries before considering a Naisjob as failed
	BackoffLimit int32 `json:"backoffLimit,omitempty"`

	// Override command when starting Docker image.
	Command []string `json:"command,omitempty"`

	// A Job tracks the successful completions. When a specified number of successful completions is reached, the task (ie, Job) is complete.
	// +nais:doc:Default="1"
	// +nais:doc:Link="https://kubernetes.io/docs/concepts/workloads/controllers/job/#job-patterns"
	Completions *int32 `json:"completions,omitempty"`

	// Specifies how to treat concurrent executions of a job that is created by this Naisjob-cron.
	// +kubebuilder:validation:Enum=Forbid;Replace;Allow
	// +nais:doc:Default="Allow"
	// +nais:doc:Link="https://kubernetes.io/docs/tasks/job/automated-tasks-with-cron-jobs/#concurrency-policy"
	ConcurrencyPolicy string `json:"concurrencyPolicy,omitempty"`

	// Custom environment variables injected into your container.
	// Specify either `value` or `valueFrom`, but not both.
	Env EnvVars `json:"env,omitempty"`

	// EnvFrom exposes all variables in the ConfigMap or Secret resources as environment variables.
	// One of `configMap` or `secret` is required.
	//
	// Environment variables will take the form `KEY=VALUE`, where `key` is the ConfigMap or Secret key.
	// You can specify as many keys as you like in a single ConfigMap or Secret.
	//
	// The ConfigMap and Secret resources must live in the same Kubernetes namespace as the Naisjob resource.
	// +nais:doc:Availability="team namespaces"
	EnvFrom []EnvFrom `json:"envFrom,omitempty"`

	// Specify how many failed Jobs should be kept.
	FailedJobsHistoryLimit int32 `json:"failedJobsHistoryLimit,omitempty"`

	// List of ConfigMap or Secret resources that will have their contents mounted into the containers as files.
	// Either `configMap` or `secret` is required.
	//
	// Files will take the path `<mountPath>/<key>`, where `key` is the ConfigMap or Secret key.
	// You can specify as many keys as you like in a single ConfigMap or Secret, and they will all
	// be mounted to the same directory.
	//
	// The ConfigMap and Secret resources must live in the same Kubernetes namespace as the Naisjob resource.
	// +nais:doc:Availability="team namespaces"
	FilesFrom []FilesFrom `json:"filesFrom,omitempty"`

	// +nais:doc:Availability="GCP"
	GCP *GCP `json:"gcp,omitempty"`

	// Your Naisjob's Docker image location and tag.
	Image string `json:"image"`

	// An Influxdb via Aiven. A typical use case is to store metrics from your application and visualize them in Grafana.
	// See [navikt/aiven-iac](https://github.com/navikt/aiven-iac) repository
	// +nais:doc:Availability="GCP"
	Influx *Influx `json:"influx,omitempty"`

	// Enable Aiven Kafka for your Naisjob.
	Kafka *Kafka `json:"kafka,omitempty"`

	// Many Naisjobs running for long periods of time eventually transition to broken states,
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

	// Configures a Maskinporten client for this Naisjob.
	// See [Maskinporten](https://doc.nais.io/security/auth/maskinporten/) for more details.
	Maskinporten *Maskinporten `json:"maskinporten,omitempty"`

	// To get your own OpenSearch instance head over to the IaC-repo to provision each instance.
	// See [navikt/aiven-iac](https://github.com/navikt/aiven-iac) repository.
	OpenSearch *OpenSearch `json:"openSearch,omitempty"`

	// For running pods in parallel.
	// If it is specified as 0, then the Job is effectively paused until it is increased.
	// +nais:doc:Default="1"
	// +nais:doc:Link="https://kubernetes.io/docs/concepts/workloads/controllers/job/#controlling-parallelism"
	Parallelism *int32 `json:"parallelism,omitempty"`

	// PreStopHook is called immediately before a container is terminated due to an API request or management event such as liveness/startup probe failure, preemption, resource contention, etc.
	// The handler is not called if the container crashes or exits by itself.
	// The reason for termination is passed to the handler.
	// +nais:doc:Link="https://doc.nais.io/naisjob/#handles-termination-gracefully";"https://kubernetes.io/docs/concepts/containers/container-lifecycle-hooks/#container-hooks"
	PreStopHook *PreStopHook `json:"preStopHook,omitempty"`

	// Sometimes, Naisjobs are temporarily unable to serve traffic. For example, an Naisjob might need
	// to load large data or configuration files during startup, or depend on external services after startup.
	// In such cases, you don't want to kill the Naisjob, but you don’t want to send it requests either.
	// Kubernetes provides readiness probes to detect and mitigate these situations. A pod with containers
	// reporting that they are not ready does not receive traffic through Kubernetes Services.
	// Read more about this over at the [Kubernetes readiness documentation](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/).
	Readiness *Probe `json:"readiness,omitempty"`

	// RestartPolicy describes how the container should be restarted. Only one of the following restart policies may be specified.
	// If none of the following policies is specified, the default one is Never.
	// Read more about [Kubernetes handling pod and container failures](https://kubernetes.io/docs/concepts/workloads/controllers/job/#handling-pod-and-container-failures)
	// +kubebuilder:validation:Enum=OnFailure;Never
	RestartPolicy string `json:"restartPolicy,omitempty"`

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

	// Specify how many completed Jobs should be kept.
	SuccessfulJobsHistoryLimit int32 `json:"successfulJobsHistoryLimit,omitempty"`

	// The grace period is the duration in seconds after the processes running in the pod are sent a termination signal and the time when the processes are forcibly halted with a kill signal.
	// Set this value longer than the expected cleanup time for your process.
	// For most jobs, the default is more than enough. Defaults to 30 seconds.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=180
	TerminationGracePeriodSeconds *int64 `json:"terminationGracePeriodSeconds,omitempty"`

	// Specify the number of seconds to wait before removing the Job after it has finished (either Completed or Failed).
	// If the field is unset, this Job won't be cleaned up by the TTL controller after it finishes.
	// +nais:doc:Availability="on-premises"
	TTLSecondsAfterFinished *int32 `json:"ttlSecondsAfterFinished,omitempty"`

	// Provides secrets management, identity-based access, and encrypting application data for auditing of secrets
	// for applications, systems, and users.
	// +nais:doc:Link="https://github.com/navikt/vault-iac/tree/master/doc"
	// +nais:doc:Availability="on-premises"
	Vault *Vault `json:"vault,omitempty"`

	// Inject on-premises web proxy configuration into the job container.
	// Most Linux applications should auto-detect these settings from the `$HTTP_PROXY`, `$HTTPS_PROXY` and `$NO_PROXY` environment variables (and their lowercase counterparts).
	// Java applications can start the JVM using parameters from the `$JAVA_PROXY_OPTIONS` environment variable.
	// +nais:doc:Availability="on-premises"
	WebProxy bool `json:"webproxy,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type NaisjobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Naisjob `json:"items"`
}

func (in *Naisjob) GetObjectKind() schema.ObjectKind {
	return in
}

func (in *Naisjob) GetObjectReference() corev1.ObjectReference {
	return corev1.ObjectReference{
		APIVersion:      "nais.io/v1",
		UID:             in.UID,
		Name:            in.Name,
		Kind:            "Naisjob",
		ResourceVersion: in.ResourceVersion,
		Namespace:       in.Namespace,
	}
}

func (in *Naisjob) GetOwnerReference() metav1.OwnerReference {
	return metav1.OwnerReference{
		APIVersion: "nais.io/v1",
		Kind:       "Naisjob",
		Name:       in.Name,
		UID:        in.UID,
	}
}

func (in *Naisjob) Hash() (string, error) {
	return hash.Hash(in.Spec)
}

func (in *Naisjob) LogFields() log.Fields {
	return log.Fields{
		"namespace":       in.GetNamespace(),
		"resourceversion": in.GetResourceVersion(),
		"naisjob":         in.GetName(),
		"correlation_id":  in.Status.CorrelationID,
	}
}

// If the Naisjob was not deployed with a correlation ID annotation,
// generate a random UUID and add it to annotations.
func (in *Naisjob) EnsureCorrelationID() error {
	if in.Annotations == nil {
		in.SetAnnotations(map[string]string{})
	}

	if len(in.Annotations[DeploymentCorrelationIDAnnotation]) != 0 {
		return nil
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("generate deployment correlation ID: %s", err)
	}

	in.Annotations[DeploymentCorrelationIDAnnotation] = id.String()

	return nil
}

func (in *Naisjob) CorrelationID() string {
	return in.Annotations[DeploymentCorrelationIDAnnotation]
}

func (in *Naisjob) SetDeploymentRolloutStatus(rolloutStatus string) {
	in.Status.DeploymentRolloutStatus = rolloutStatus
}

func (in *Naisjob) DefaultSecretPath(base string) SecretPath {
	return SecretPath{
		MountPath: DefaultVaultMountPath,
		KvPath:    fmt.Sprintf("%s/%s/%s", base, in.Name, in.Namespace),
	}
}

func (in *Naisjob) SkipDeploymentMessage() bool {
	if in.Annotations == nil {
		return false
	}
	skip, _ := strconv.ParseBool(in.Annotations[SkipDeploymentMessageAnnotation])
	return skip
}

func (in *Naisjob) ClientID(cluster string) string {
	return fmt.Sprintf("%s:%s:%s", cluster, in.ObjectMeta.Namespace, in.ObjectMeta.Name)
}
