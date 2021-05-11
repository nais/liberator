package nais_io_v1alpha1

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	hash "github.com/mitchellh/hashstructure"
	"github.com/nais/liberator/pkg/apis/nais.io/v1"
	log "github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	DeploymentCorrelationIDAnnotation      = "nais.io/deploymentCorrelationID"
	SkipDeploymentMessageAnnotation        = "nais.io/skipDeploymentMessage"
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

// Application defines a NAIS application.
//
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="Team",type="string",JSONPath=".metadata.labels.team"
// +kubebuilder:printcolumn:name="State",type="string",JSONPath=".status.synchronizationState"
// +kubebuilder:resource:path="applications",shortName="app",singular="application"
type Application struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApplicationSpec   `json:"spec"`
	Status ApplicationStatus `json:"status,omitempty"`
}

// +kubebuilder:validation:Pattern=`^https:\/\/.+$`
type Ingress string

// ApplicationSpec contains the NAIS manifest.
// Please keep this list sorted for clarity.
type ApplicationSpec struct {
	// By default, all traffic is disallowed between applications inside the cluster.
	// Configure access policies to allow specific applications.
	// +nais:doc:Availability=GCP
	AccessPolicy *nais_io_v1.AccessPolicy `json:"accessPolicy,omitempty"`

	Azure   *Azure   `json:"azure,omitempty"`
	Elastic *Elastic `json:"elastic,omitempty"`

	// Custom environment variables injected into your container.
	Env []EnvVar `json:"env,omitempty"`

	// Will expose all variables in ConfigMap or Secret resource as environment variables.
	// One of `configmap` or `secret` is required.
	// +nais:doc:Availability="team namespaces"
	EnvFrom []EnvFrom `json:"envFrom,omitempty"`

	// List of ConfigMap or Secret resources that will have their contents mounted into the containers as files.
	// Either `configmap` or `secret` is required.
	// +nais:doc:Availability="team namespaces"
	FilesFrom []FilesFrom `json:"filesFrom,omitempty"`

	// +nais:doc:Availability="GCP"
	GCP *GCP `json:"gcp,omitempty"`

	// Configures an ID-porten client for this application.
	// See [ID-porten](https://doc.nais.io/security/auth/idporten/) for more details.
	IDPorten *IDPorten `json:"idporten,omitempty"`

	// Your application's Docker image location and tag.
	// +nais:doc:Sample="ghcr.io/navikt/myapp:1.6.9"
	Image string `json:"image"`

	// List of URLs that will route HTTPS traffic to the application.
	// All URLs must start with `https://`. Domain availability differs according to which environment your application is running in,
	// see https://doc.nais.io/clusters/gcp/ and https://doc.nais.io/clusters/on-premises/.
	Ingresses []Ingress `json:"ingresses,omitempty"`

	// Enable Aiven Kafka for your application.
	Kafka *Kafka `json:"kafka,omitempty"`

	// If true, a HTTP endpoint will be available at `$ELECTOR_PATH` that returns the current leader.
	// See https://doc.nais.io/addons/leader-election/.
	LeaderElection bool `json:"leaderElection,omitempty"`

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

	// The port number which is exposed by the container and should receive traffic.
	Port int `json:"port,omitempty"`

	// A HTTP GET will be issued to this endpoint at least once before the pod is terminated.
	// See https://doc.nais.io/nais-application/#handles-termination-gracefully.
	PreStopHookPath string `json:"preStopHookPath,omitempty"`

	// Prometheus is used to [scrape metrics from the pod](https://doc.nais.io/observability/metrics/).
	Prometheus *PrometheusConfig `json:"prometheus,omitempty"`

	// Sometimes, applications are temporarily unable to serve traffic. For example, an application might need
	// to load large data or configuration files during startup, or depend on external services after startup.
	// In such cases, you don't want to kill the application, but you don’t want to send it requests either.
	// Kubernetes provides readiness probes to detect and mitigate these situations. A pod with containers
	// reporting that they are not ready does not receive traffic through Kubernetes Services.
	// Read more about this over at the [Kubernetes readiness documentation](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/).
	Readiness *Probe `json:"readiness,omitempty"`

	// The numbers of pods to run in parallel.
	Replicas *Replicas `json:"replicas,omitempty"`

	// When Containers have [resource requests](http://kubernetes.io/docs/user-guide/compute-resources/) specified,
	// the Kubernetes scheduler can make better decisions about which nodes to place pods on.
	Resources *ResourceRequirements `json:"resources,omitempty"`

	// Whether or not to enable a sidecar container for secure logging.
	SecureLogs *SecureLogs `json:"secureLogs,omitempty"`
	// How to connect to the default service in your application's container.
	Service *Service `json:"service,omitempty"`
	// Whether to skip injection of certificate authority bundle or not. Defaults to false.
	SkipCaBundle bool      `json:"skipCaBundle,omitempty"`
	// Startup probes will be available with Kubernetes 1.18 (in GCP, and 1.17 on-prem). Do not use this feature yet as it will not work.
	//
	// Sometimes, you have to deal with legacy applications that might require an additional startup time on their first
	// initialization. In such cases, it can be tricky to set up liveness probe parameters without compromising the fast
	// response to deadlocks that motivated such a probe. The trick is to set up a startup probe with the same command,
	// HTTP or TCP check, with a failureThreshold * periodSeconds long enough to cover the worst case startup time.
	Startup      *Probe    `json:"startup,omitempty"`
	Strategy     *Strategy `json:"strategy,omitempty"`
	TokenX       *TokenX   `json:"tokenx,omitempty"`
	Tracing      *Tracing  `json:"tracing,omitempty"`
	Vault        *Vault    `json:"vault,omitempty"`
	WebProxy     bool      `json:"webproxy,omitempty"`
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

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Application `json:"items"`
}

type Azure struct {
	// Configures an Azure AD client for this application.
	// See [Azure AD](https://doc.nais.io/security/auth/azure-ad/) for more details.
	Application *AzureApplication `json:"application"`
}

type Tracing struct {
	// if enabled, a rule allowing egress to app:jaeger will be appended to
	// NetworkPolicy
	Enabled bool `json:"enabled"`
}

type TokenX struct {
	// if enabled, the application will have a jwker secret injected
	Enabled bool `json:"enabled"`
	// if enabled, secrets for TokenX will be mounted as files only, i.e. not as env.
	MountSecretsAsFilesOnly bool `json:"mountSecretsAsFilesOnly,omitempty"`
}

type IDPorten struct {
	Enabled   bool   `json:"enabled"`
	ClientURI string `json:"clientURI,omitempty"`
	// +kubebuilder:validation:Pattern=`^https:\/\/.+$`
	RedirectURI string `json:"redirectURI,omitempty"`
	// +kubebuilder:validation:Pattern=`^\/.*$`
	RedirectPath string `json:"redirectPath,omitempty"`
	// +kubebuilder:validation:Pattern=`^\/.*$`
	FrontchannelLogoutPath string   `json:"frontchannelLogoutPath,omitempty"`
	FrontchannelLogoutURI  string   `json:"frontchannelLogoutURI,omitempty"`
	PostLogoutRedirectURIs []string `json:"postLogoutRedirectURIs,omitempty"`
	// +kubebuilder:validation:Minimum=3600
	// +kubebuilder:validation:Maximum=7200
	SessionLifetime *int `json:"sessionLifetime,omitempty"`
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=3600
	AccessTokenLifetime *int `json:"accessTokenLifetime,omitempty"`
}

type AzureApplication struct {
	Enabled   bool     `json:"enabled"`
	ReplyURLs []string `json:"replyURLs,omitempty"`
	// +kubebuilder:validation:Enum=nav.no;trygdeetaten.no
	Tenant string `json:"tenant,omitempty"`
	// Claims defines additional configuration of the emitted claims in tokens returned to the AzureAdApplication
	Claims *nais_io_v1.AzureAdClaims `json:"claims,omitempty"`
}

type SecureLogs struct {
	// Whether to enable a sidecar container for secure logging.
	// If enabled, a volume is mounted in the pods where secure logs can be saved.
	Enabled bool `json:"enabled"`
}

// Liveness probe and readiness probe definitions.
type Probe struct {
	// HTTP endpoint path that signals 200 OK if the application has started successfully.
	Path             string `json:"path"`
	// Port for the startup probe.
	// Default: .spec.port
	Port             int    `json:"port,omitempty"`
	// Number of seconds after the container has started before startup probes are initiated.
	// Default: 20
	InitialDelay     int    `json:"initialDelay,omitempty"`
	// How often (in seconds) to perform the probe.
	// Default: 10
	PeriodSeconds    int    `json:"periodSeconds,omitempty"`
	// When a Pod starts and the probe fails, Kubernetes will try _failureThreshold_ times before giving up. Giving up in
	// case of a startup probe means restarting the Pod.
	FailureThreshold int    `json:"failureThreshold,omitempty"`
	// Number of seconds after which the probe times out.
	// Default: 1
	Timeout          int    `json:"timeout,omitempty"`
}

type PrometheusConfig struct {
	Enabled bool   `json:"enabled,omitempty"`
	Port    string `json:"port,omitempty"`
	Path    string `json:"path,omitempty"`
}

type Replicas struct {
	// The minimum amount of replicas acceptable for a successful deployment.
	Min int `json:"min,omitempty"`
	// The pod autoscaler will scale deployments on demand until this maximum has been reached.
	Max int `json:"max,omitempty"`
	// Amount of CPU usage before the autoscaler kicks in.
	CpuThresholdPercentage int `json:"cpuThresholdPercentage,omitempty"`
}

type ResourceSpec struct {
	// +kubebuilder:validation:Pattern=^\d+m?$
	Cpu string `json:"cpu,omitempty"`
	// +kubebuilder:validation:Pattern=^\d+[KMG]i$
	Memory string `json:"memory,omitempty"`
}

type ResourceRequirements struct {
	Limits   *ResourceSpec `json:"limits,omitempty"`
	Requests *ResourceSpec `json:"requests,omitempty"`
}

type ObjectFieldSelector struct {
	// +kubebuilder:validation:Enum="";metadata.name;metadata.namespace;metadata.labels;metadata.annotations;spec.nodeName;spec.serviceAccountName;status.hostIP;status.podIP
	FieldPath string `json:"fieldPath"`
}

type EnvVarSource struct {
	FieldRef ObjectFieldSelector `json:"fieldRef"`
}

type CloudStorageBucket struct {
	Name            string `json:"name"`
	CascadingDelete bool   `json:"cascadingDelete,omitempty"`
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=36500
	RetentionPeriodDays *int                `json:"retentionPeriodDays,omitempty"`
	LifecycleCondition  *LifecycleCondition `json:"lifecycleCondition,omitempty"`
}

type LifecycleCondition struct {
	Age              int    `json:"age,omitempty"`
	CreatedBefore    string `json:"createdBefore,omitempty"`
	NumNewerVersions int    `json:"numNewerVersions,omitempty"`
	WithState        string `json:"withState,omitempty"`
}

type CloudSqlInstanceType string

const (
	CloudSqlInstanceTypePostgres11 CloudSqlInstanceType = "POSTGRES_11"
	CloudSqlInstanceTypePostgres12 CloudSqlInstanceType = "POSTGRES_12"
)

type CloudSqlInstanceDiskType string

func (c CloudSqlInstanceDiskType) GoogleType() string {
	return fmt.Sprintf("PD_%s", c)
}

const (
	CloudSqlInstanceDiskTypeSSD CloudSqlInstanceDiskType = "SSD"
	CloudSqlInstanceDiskTypeHDD CloudSqlInstanceDiskType = "HDD"
)

type CloudSqlDatabase struct {
	// +kubebuilder:validation:Required
	Name         string                 `json:"name"`
	EnvVarPrefix string                 `json:"envVarPrefix,omitempty"`
	Users        []CloudSqlDatabaseUser `json:"users,omitempty"`
}

type CloudSqlDatabaseUser struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern="^[_a-zA-Z][_a-zA-Z0-9]+$"
	Name string `json:"name"`
}

type CloudSqlInstance struct {
	// +kubebuilder:validation:Enum=POSTGRES_11;POSTGRES_12
	// +kubebuilder:validation:Required
	Type CloudSqlInstanceType `json:"type"`
	Name string               `json:"name,omitempty"`
	// +kubebuilder:validation:Pattern="db-.+"
	Tier string `json:"tier,omitempty"`
	// +kubebuilder:validation:Enum=SSD;HDD
	DiskType         CloudSqlInstanceDiskType `json:"diskType,omitempty"`
	HighAvailability bool                     `json:"highAvailability,omitempty"`
	// +kubebuilder:validation:Minimum=10
	DiskSize       int  `json:"diskSize,omitempty"`
	DiskAutoresize bool `json:"diskAutoresize,omitempty"`
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=23
	AutoBackupHour *int         `json:"autoBackupHour,omitempty"`
	Maintenance    *Maintenance `json:"maintenance,omitempty"`
	// +kubebuilder:validation:Required
	Databases       []CloudSqlDatabase `json:"databases,omitempty"`
	CascadingDelete bool               `json:"cascadingDelete,omitempty"`
	Collation       string             `json:"collation,omitempty"`
}

type Maintenance struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=7
	Day int `json:"day,omitempty"`
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=23
	Hour *int `json:"hour,omitempty"` // must use pointer here to be able to distinguish between no value and value 0 from user.
}

type Elastic struct {
	// Provisions an Elasticsearch instance and configures your application so it can access it.
	// Use the `instance_name` that you specified in the [navikt/aiven-iac](https://github.com/navikt/aiven-iac) repository.
	// +nais:doc:Availability=GCP
	Instance string `json:"instance"`
}

type GCP struct {
	// Provision cloud storage buckets and connect them to your application.
	Buckets []CloudStorageBucket `json:"buckets,omitempty"`
	// Provision database instances and connect them to your application.
	// See [PostgreSQL documentation](https://doc.nais.io/persistence/postgres/) for more details.
	SqlInstances []CloudSqlInstance `json:"sqlInstances,omitempty"`
	// List of _additional_ permissions that should be granted to your application for accessing external GCP resources that have not been provisioned through NAIS.
	// [Supported resources found here](https://cloud.google.com/config-connector/docs/reference/resource-docs/iam/iampolicymember#external_organization_level_policy_member).
	Permissions []CloudIAMPermission `json:"permissions,omitempty"`
}

type EnvVar struct {
	Name      string        `json:"name"`
	Value     string        `json:"value,omitempty"`
	ValueFrom *EnvVarSource `json:"valueFrom,omitempty"`
}

type EnvFrom struct {
	ConfigMap string `json:"configmap,omitempty"`
	Secret    string `json:"secret,omitempty"`
}

type FilesFrom struct {
	ConfigMap string `json:"configmap,omitempty"`
	Secret    string `json:"secret,omitempty"`
	MountPath string `json:"mountPath,omitempty"`
}

type SecretPath struct {
	MountPath string `json:"mountPath"`
	KvPath    string `json:"kvPath"`
	// +kubebuilder:validation:Enum=flatten;json;yaml;env;properties;""
	Format string `json:"format,omitempty"`
}

type Vault struct {
	Enabled bool         `json:"enabled,omitempty"`
	Sidecar bool         `json:"sidecar,omitempty"`
	Paths   []SecretPath `json:"paths,omitempty"`
}

type Strategy struct {
	// +kubebuilder:validation:Enum=Recreate;RollingUpdate
	Type string `json:"type"`
}

type Service struct {
	// +kubebuilder:validation:Enum=http;redis;tcp;grpc
	// Which protocol the backend service runs on. Default is http.
	Protocol string `json:"protocol,omitempty"`
	// Port for the default service. Default port is 80.
	Port int32 `json:"port"`
}

type Kafka struct {
	// Configures your application to access an Aiven Kafka cluster.
	// See https://doc.nais.io/addons/kafka/.
	// +kubebuilder:validation:Enum=nav-dev;nav-prod;nav-infrastructure
	Pool string `json:"pool"`
}

type CloudIAMResource struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Name       string `json:"name,omitempty"`
}

type CloudIAMPermission struct {
	Role     string           `json:"role"`
	Resource CloudIAMResource `json:"resource"`
}

type Maskinporten struct {
	Enabled bool                           `json:"enabled"`
	Scopes  []nais_io_v1.MaskinportenScope `json:"scopes,omitempty"`
}

func (in *Application) GetObjectKind() schema.ObjectKind {
	return in
}

func (in *Application) GetObjectReference() v1.ObjectReference {
	return v1.ObjectReference{
		APIVersion:      "v1alpha1",
		UID:             in.UID,
		Name:            in.Name,
		Kind:            "Application",
		ResourceVersion: in.ResourceVersion,
		Namespace:       in.Namespace,
	}
}

func (in *Application) GetOwnerReference() metav1.OwnerReference {
	return metav1.OwnerReference{
		APIVersion: "v1alpha1",
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
		in.Labels,
		changeCause,
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

// If the application was not deployed with a correlation ID annotation,
// generate a random UUID and add it to annotations.
func (in *Application) EnsureCorrelationID() error {
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

func (in *Application) CorrelationID() string {
	return in.Annotations[DeploymentCorrelationIDAnnotation]
}

func (in *Application) SetDeploymentRolloutStatus(rolloutStatus string) {
	in.Status.DeploymentRolloutStatus = rolloutStatus
}

func (in *Application) DefaultSecretPath(base string) SecretPath {
	return SecretPath{
		MountPath: DefaultVaultMountPath,
		KvPath:    fmt.Sprintf("%s/%s/%s", base, in.Name, in.Namespace),
	}
}

func (in *Application) SkipDeploymentMessage() bool {
	if in.Annotations == nil {
		return false
	}
	skip, _ := strconv.ParseBool(in.Annotations[SkipDeploymentMessageAnnotation])
	return skip
}

func (in *Application) ClientID(cluster string) string {
	return fmt.Sprintf("%s:%s:%s", cluster, in.ObjectMeta.Namespace, in.ObjectMeta.Name)
}
