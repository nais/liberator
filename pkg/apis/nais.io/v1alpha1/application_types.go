package nais_io_v1alpha1

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	hash "github.com/mitchellh/hashstructure"
	log "github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/nais/liberator/pkg/apis/nais.io/v1"
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
	// By default, no traffic is allowed between applications inside the cluster.
	// Configure access policies to explicitly allow communication between applications.
	// This is also used for granting inbound access in the context of Azure AD and TokenX clients.
	// +nais:doc:Link="https://doc.nais.io/appendix/zero-trust/";"https://doc.nais.io/security/auth/azure-ad/#pre-authorization";"https://doc.nais.io/security/auth/tokenx/#access-policies"
	AccessPolicy *nais_io_v1.AccessPolicy `json:"accessPolicy,omitempty"`

	// Provisions and configures Azure resources.
	Azure   *Azure   `json:"azure,omitempty"`

	// Override command when starting Docker image.
	Command []string `json:"command,omitempty"`

	Elastic *Elastic `json:"elastic,omitempty"`

	// Custom environment variables injected into your container.
	// Specify either `value` or `valueFrom`, but not both.
	Env []EnvVar `json:"env,omitempty"`

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

	// Configures an ID-porten client for this application.
	// See [ID-porten](https://doc.nais.io/security/auth/idporten/) for more details.
	IDPorten *IDPorten `json:"idporten,omitempty"`

	// Your application's Docker image location and tag.
	Image string `json:"image"`

	// List of URLs that will route HTTPS traffic to the application.
	// All URLs must start with `https://`. Domain availability differs according to which environment your application is running in.
	// +nais:doc:Link="https://doc.nais.io/clusters/gcp/";"https://doc.nais.io/clusters/on-premises/"
	Ingresses []Ingress `json:"ingresses,omitempty"`

	// Enable Aiven Kafka for your application.
	Kafka *Kafka `json:"kafka,omitempty"`

	// If true, an HTTP endpoint will be available at `$ELECTOR_PATH` that returns the current leader.
	// +nais:doc:Link="https://doc.nais.io/addons/leader-election/"
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
	// +nais:doc:Link="https://doc.nais.io/nais-application/#handles-termination-gracefully"
	PreStopHookPath string `json:"preStopHookPath,omitempty"`

	// Prometheus is used to [scrape metrics from the pod](https://doc.nais.io/observability/metrics/).
	// Use this configuration to override the default values.
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

	// Specify which port and protocol is used to connect to the application in the container.
	// Defaults to HTTP on port 80.
	Service *Service `json:"service,omitempty"`

	// Whether to skip injection of NAV certificate authority bundle or not. Defaults to false.
	SkipCaBundle bool `json:"skipCaBundle,omitempty"`

	// Startup probes will be available with Kubernetes 1.20. Do not use this feature yet as it will not work.
	//
	// Sometimes, you have to deal with legacy applications that might require an additional startup time on their first
	// initialization. In such cases, it can be tricky to set up liveness probe parameters without compromising the fast
	// response to deadlocks that motivated such a probe. The trick is to set up a startup probe with the same command,
	// HTTP or TCP check, with a `failureThreshold * periodSeconds` long enough to cover the worst case startup time.
	Startup *Probe `json:"startup,omitempty"`

	// Specifies the strategy used to replace old Pods by new ones.
	Strategy *Strategy `json:"strategy,omitempty"`

	// Provisions and configures a TokenX client for your application.
	// +nais:doc:Link="https://doc.nais.io/security/auth/tokenx/"
	TokenX *TokenX `json:"tokenx,omitempty"`

	// *DEPRECATED*. Do not use. Will be removed in a future version.
	Tracing *Tracing `json:"tracing,omitempty"`

	// Provides secrets management, identity-based access, and encrypting application data for auditing of secrets
	// for applications, systems, and users.
	// +nais:doc:Link="https://github.com/navikt/vault-iac/tree/master/doc"
	// +nais:doc:Availability="on-premises"
	Vault *Vault `json:"vault,omitempty"`

	// Inject web proxy configuration to the application using the `$HTTP_PROXY`, `$HTTPS_PROXY` and `$NO_PROXY` environment variables.
	// +nais:doc:Availability="on-premises"
	WebProxy bool `json:"webproxy,omitempty"`
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
	Enabled bool `json:"enabled"`
}

type TokenX struct {
	// If enabled, will provision and configure a TokenX client and inject an accompanying secret.
	Enabled bool `json:"enabled"`
	// If enabled, secrets for TokenX will be mounted as files only, i.e. not as environment variables.
	MountSecretsAsFilesOnly bool `json:"mountSecretsAsFilesOnly,omitempty"`
}

type IDPorten struct {
	// Whether to enable provisioning of an ID-porten client.
	// If enabled, an ID-porten client be provisioned.
	Enabled bool `json:"enabled"`
	// AccessTokenLifetime is the lifetime in seconds for any issued access token from ID-porten.
	//
	// If unspecified, defaults to `3600` seconds (1 hour).
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=3600
	AccessTokenLifetime *int `json:"accessTokenLifetime,omitempty"`
	// ClientURI is the URL shown to the user at ID-porten when displaying a 'back' button or on errors.
	ClientURI string `json:"clientURI,omitempty"`
	// FrontchannelLogoutPath is a valid path for your application where ID-porten sends a request to whenever the user has
	// initiated a logout elsewhere as part of a single logout (front channel logout) process.
	//
	// If unspecified, defaults to `/oauth2/logout`.
	// +nais:doc:Link="https://doc.nais.io/security/auth/idporten/#front-channel-logout";"https://docs.digdir.no/oidc_func_sso.html#2-h%C3%A5ndtere-utlogging-fra-id-porten"
	// +kubebuilder:validation:Pattern=`^\/.*$`
	FrontchannelLogoutPath string `json:"frontchannelLogoutPath,omitempty"`
	// *DEPRECATED*. Prefer using `frontchannelLogoutPath`.
	// +nais:doc:Link="https://doc.nais.io/security/auth/idporten/#front-channel-logout";"https://docs.digdir.no/oidc_func_sso.html#2-h%C3%A5ndtere-utlogging-fra-id-porten"
	FrontchannelLogoutURI string `json:"frontchannelLogoutURI,omitempty"`
	// PostLogoutRedirectURIs are valid URIs that ID-porten will allow redirecting the end-user to after a single logout
	// has been initiated and performed by the application.
	//
	// If unspecified, will default to `[ "https://www.nav.no" ]`
	// +nais:doc:Link="https://doc.nais.io/security/auth/idporten/#self-initiated-logout";"https://docs.digdir.no/oidc_func_sso.html#1-utlogging-fra-egen-tjeneste"
	PostLogoutRedirectURIs []string `json:"postLogoutRedirectURIs,omitempty"`
	// RedirectPath is a valid path that ID-porten redirects back to after a successful authorization request.
	//
	// If unspecified, will default to `/oauth2/callback`.
	// +kubebuilder:validation:Pattern=`^\/.*$`
	RedirectPath string `json:"redirectPath,omitempty"`
	// *DEPRECATED*. Prefer using `redirectPath`.
	// +kubebuilder:validation:Pattern=`^https:\/\/.+$`
	RedirectURI string `json:"redirectURI,omitempty"`
	// SessionLifetime is the maximum lifetime in seconds for any given user's session in your application.
	// The timeout starts whenever the user is redirected from the `authorization_endpoint` at ID-porten.
	//
	// If unspecified, defaults to `7200` seconds (2 hours).
	// Note: Attempting to refresh the user's `access_token` beyond this timeout will yield an error.
	// +kubebuilder:validation:Minimum=3600
	// +kubebuilder:validation:Maximum=7200
	SessionLifetime *int `json:"sessionLifetime,omitempty"`
}

type AzureApplication struct {
	// Whether to enable provisioning of an Azure AD application.
	// If enabled, an Azure AD application will be provisioned.
	Enabled bool `json:"enabled"`
	// ReplyURLs is a list of allowed redirect URLs used when performing OpenID Connect flows for authenticating end-users.
	// +nais:doc:Link="https://doc.nais.io/security/auth/azure-ad/#reply-urls"
	ReplyURLs []string `json:"replyURLs,omitempty"`
	// A Tenant represents an organization in Azure AD.
	//
	// If unspecified, will default to `trygdeetaten.no` for development clusters and `nav.no` for production clusters.
	// +nais:doc:Link="https://doc.nais.io/security/auth/azure-ad/#tenants"
	// +kubebuilder:validation:Enum=nav.no;trygdeetaten.no
	Tenant string `json:"tenant,omitempty"`
	// Claims defines additional configuration of the emitted claims in tokens returned to the Azure AD application.
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
	Path string `json:"path"`
	// Port for the startup probe.
	Port int `json:"port,omitempty"`
	// Number of seconds after the container has started before startup probes are initiated.
	InitialDelay int `json:"initialDelay,omitempty"`
	// How often (in seconds) to perform the probe.
	PeriodSeconds int `json:"periodSeconds,omitempty"`
	// When a Pod starts, and the probe fails, Kubernetes will try _failureThreshold_ times before giving up.
	// Giving up in case of a startup probe means restarting the Pod.
	FailureThreshold int `json:"failureThreshold,omitempty"`
	// Number of seconds after which the probe times out.
	Timeout int `json:"timeout,omitempty"`
}

type PrometheusConfig struct {
	Enabled bool   `json:"enabled,omitempty"`
	Port    string `json:"port,omitempty"`
	Path    string `json:"path,omitempty"`
}

type Replicas struct {
	// The minimum amount of running replicas for a deployment.
	Min int `json:"min,omitempty"`
	// The pod autoscaler will increase replicas when required up to the maximum.
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
	// Limit defines the maximum amount of resources a container can use before getting evicted.
	Limits *ResourceSpec `json:"limits,omitempty"`
	// Request defines the amount of resources a container is allocated on startup.
	Requests *ResourceSpec `json:"requests,omitempty"`
}

type ObjectFieldSelector struct {
	// Field value from the `Pod` spec that should be copied into the environment variable.
	// +kubebuilder:validation:Enum="";metadata.name;metadata.namespace;metadata.labels;metadata.annotations;spec.nodeName;spec.serviceAccountName;status.hostIP;status.podIP
	FieldPath string `json:"fieldPath"`
}

type EnvVarSource struct {
	FieldRef ObjectFieldSelector `json:"fieldRef"`
}

type CloudStorageBucket struct {
	// The name of the bucket
	Name string `json:"name"`
	// Allows deletion of bucket. Set to true if you want to delete the bucket.
	CascadingDelete bool `json:"cascadingDelete,omitempty"`
	// The number of days to hold objects in the bucket before it is allowed to delete them.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=36500
	RetentionPeriodDays *int `json:"retentionPeriodDays,omitempty"`
	// Conditions for the bucket to use when selecting objects to delete in cleanup.
	// +nais:doc:Link="https://cloud.google.com/storage/docs/lifecycle"
	LifecycleCondition *LifecycleCondition `json:"lifecycleCondition,omitempty"`
}

type LifecycleCondition struct {
	// Condition is satisfied when the object reaches the specified age in days. These will be deleted.
	Age int `json:"age,omitempty"`
	// Condition is satisfied when the object is created before midnight on the specified date. These will be deleted.
	CreatedBefore string `json:"createdBefore,omitempty"`
	// Condition is satisfied when the object has the specified number of newer versions.
	// The older versions will be deleted.
	NumNewerVersions int `json:"numNewerVersions,omitempty"`
	// Condition is satisfied when the object has the specified state.
	// +kubebuilder:validation:Enum="";LIVE;ARCHIVED;ANY
	WithState string `json:"withState,omitempty"`
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
	// Database name.
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// Prefix to add to environment variables made available for database connection.
	EnvVarPrefix string `json:"envVarPrefix,omitempty"`
	// The users created to allow database access.
	Users []CloudSqlDatabaseUser `json:"users,omitempty"`
}

type CloudSqlDatabaseUser struct {
	// User name.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern="^[_a-zA-Z][_a-zA-Z0-9]+$"
	Name string `json:"name"`
}

type CloudSqlInstance struct {
	// PostgreSQL version.
	// +kubebuilder:validation:Enum=POSTGRES_11;POSTGRES_12
	// +kubebuilder:validation:Required
	Type CloudSqlInstanceType `json:"type"`
	// The name of the instance, if omitted the database name will be used.
	Name string `json:"name,omitempty"`
	// Server tier, i.e. how much CPU and memory allocated.
	// Available tiers can be retrieved on the command line
	// by running `gcloud sql tiers list`.
	// +kubebuilder:validation:Pattern="db-.+"
	Tier string `json:"tier,omitempty"`
	// Disk type to use for storage in the database.
	// +kubebuilder:validation:Enum=SSD;HDD
	DiskType CloudSqlInstanceDiskType `json:"diskType,omitempty"`
	// When set to true this will set up standby database for failover.
	HighAvailability bool `json:"highAvailability,omitempty"`
	// How much hard drive space to allocate for the SQL server, in gigabytes.
	// +kubebuilder:validation:Minimum=10
	DiskSize int `json:"diskSize,omitempty"`
	// When set to true, GCP will automatically increase storage by XXX for the database when
	// disk usage is above the high water mark.
	// +nais:doc:Link="https://cloud.google.com/sql/docs/postgres/instance-settings#threshold"
	DiskAutoresize bool `json:"diskAutoresize,omitempty"`
	// If specified, run automatic backups of the SQL database at the given hour.
	// Note that this will backup the whole SQL instance, and not separate databases.
	// Restores are done using the Google Cloud Console.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=23
	AutoBackupHour *int `json:"autoBackupHour,omitempty"`
	// Desired maintenance window for database updates.
	Maintenance *Maintenance `json:"maintenance,omitempty"`
	// List of databases that should be created on this Postgres server.
	// +kubebuilder:validation:Required
	Databases []CloudSqlDatabase `json:"databases,omitempty"`
	// Remove the entire Postgres server including all data when the Kubernetes resource is deleted.
	// *THIS IS A DESTRUCTIVE OPERATION*! Set cascading delete only when you want to remove data forever.
	CascadingDelete bool `json:"cascadingDelete,omitempty"`
	// Sort order for `ORDER BY ...` clauses.
	Collation string `json:"collation,omitempty"`
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
	// +nais:doc:Link="https://doc.nais.io/persistence/buckets/"
	// +nais:doc:Availability=GCP
	Buckets []CloudStorageBucket `json:"buckets,omitempty"`
	// Provision database instances and connect them to your application.
	// +nais:doc:Link="https://doc.nais.io/persistence/postgres/";"https://cloud.google.com/sql/docs/postgres/instance-settings#impact"
	// +nais:doc:Availability=GCP
	SqlInstances []CloudSqlInstance `json:"sqlInstances,omitempty"`
	// List of _additional_ permissions that should be granted to your application for accessing external GCP resources that have not been provisioned through NAIS.
	// +nais:doc:Link="https://cloud.google.com/config-connector/docs/reference/resource-docs/iam/iampolicymember#external_organization_level_policy_member"
	// +nais:doc:Availability=GCP
	Permissions []CloudIAMPermission `json:"permissions,omitempty"`
}

type EnvVar struct {
	// Environment variable name. May only contain letters, digits, and the underscore `_` character.
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// Environment variable value. Numbers and boolean values must be quoted.
	// Required unless `valueFrom` is specified.
	Value string `json:"value,omitempty"`
	// Dynamically set environment variables based on fields found in the Pod spec.
	// +nais:doc:Link="https://kubernetes.io/docs/tasks/inject-data-application/environment-variable-expose-pod-information/"
	ValueFrom *EnvVarSource `json:"valueFrom,omitempty"`
}

type EnvFrom struct {
	// Name of the `ConfigMap` where environment variables are specified.
	// Required unless `secret` is set.
	ConfigMap string `json:"configmap,omitempty"`
	// Name of the `Secret` where environment variables are specified.
	// Required unless `configMap` is set.
	Secret string `json:"secret,omitempty"`
}

type FilesFrom struct {
	// Name of the `ConfigMap` that contains files that should be mounted into the container.
	// Required unless `secret` is set.
	ConfigMap string `json:"configmap,omitempty"`
	// Name of the `Secret` that contains files that should be mounted into the container.
	// Required unless `configMap` is set.
	// If mounting multiple secrets, `mountPath` *MUST* be set to avoid collisions.
	Secret string `json:"secret,omitempty"`
	// Filesystem path inside the pod where files are mounted.
	// The directory will be created if it does not exist. If the directory exists,
	// any files in the directory will be made unaccessible.
	//
	// Defaults to `/var/run/configmaps/<NAME>` or `/var/run/secrets`, depending on which of them is specified.
	MountPath string `json:"mountPath,omitempty"`
}

type SecretPath struct {
	// File system path that the secret will be mounted into.
	MountPath string `json:"mountPath"`
	// Path to Vault key/value store that should be mounted into the file system.
	KvPath string `json:"kvPath"`
	// Format of the secret that should be processed.
	// +kubebuilder:validation:Enum=flatten;json;yaml;env;properties;""
	Format string `json:"format,omitempty"`
}

type Vault struct {
	// If set to true, fetch secrets from Vault and inject into the pods.
	Enabled bool `json:"enabled,omitempty"`
	// If enabled, the sidecar will automatically refresh the token's Time-To-Live before it expires.
	Sidecar bool `json:"sidecar,omitempty"`
	// List of secret paths to be read from Vault and injected into the pod's filesystem.
	// Overriding the `paths` array is optional, and will give you fine-grained control over which Vault paths that will be mounted on the file system.
	//
	// By default, the list will contain an entry with
	//
	// `kvPath: /kv/<environment>/<zone>/<application>/<namespace>`
	// `mountPath: /var/run/secrets/nais.io/vault`
	//
	// that will always be attempted to be mounted.
	Paths []SecretPath `json:"paths,omitempty"`
}

type Strategy struct {
	// Specifies the strategy used to replace old Pods by new ones.
	// `RollingUpdate` is the default value.
	// +kubebuilder:validation:Enum=Recreate;RollingUpdate
	Type string `json:"type"`
}

type Service struct {
	// +kubebuilder:validation:Enum=http;redis;tcp;grpc
	// Which protocol the backend service runs on. Default is `http`.
	Protocol string `json:"protocol,omitempty"`
	// Port for the default service. Default port is 80.
	Port int32 `json:"port"`
}

type Kafka struct {
	// Configures your application to access an Aiven Kafka cluster.
	// +nais:doc:Link="https://doc.nais.io/addons/kafka/"
	// +kubebuilder:validation:Enum=nav-dev;nav-prod;nav-infrastructure
	Pool string `json:"pool"`
}

type CloudIAMResource struct {
	// Kubernetes _APIVersion_.
	APIVersion string `json:"apiVersion"`
	// Kubernetes _Kind_.
	Kind       string `json:"kind"`
	// Kubernetes _Name_.
	Name       string `json:"name,omitempty"`
}

type CloudIAMPermission struct {
	// Name of the GCP role to bind the resource to.
	Role string `json:"role"`
	// IAM resource to bind the role to.
	Resource CloudIAMResource `json:"resource"`
}

type Maskinporten struct {
	// If enabled, provisions and configures a Maskinporten client at DigDir.
	Enabled bool                           `json:"enabled"`
	// Ensure that the NAV organization has been granted access to the scope prior to requesting access.
	Scopes  nais_io_v1.MaskinportenScope `json:"scopes,omitempty"`
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
