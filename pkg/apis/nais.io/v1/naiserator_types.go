package nais_io_v1

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

const (
	DeploymentCorrelationIDAnnotation = "nais.io/deploymentCorrelationID"
	SkipDeploymentMessageAnnotation   = "nais.io/skipDeploymentMessage"
	DefaultVaultMountPath             = "/var/run/secrets/nais.io/vault"
)

type Azure struct {
	// Configures an Azure AD client for this application.
	// See [Azure AD](https://doc.nais.io/security/auth/azure-ad/) for more details.
	Application *AzureApplication `json:"application"`
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
	Claims *AzureAdClaims `json:"claims,omitempty"`
}

type Elastic struct {
	// Provisions an Elasticsearch instance and configures your application so it can access it.
	// Use the `instance_name` that you specified in the [navikt/aiven-iac](https://github.com/navikt/aiven-iac) repository.
	Instance string `json:"instance"`
}

type Influx struct {
	// Provisions an Influx instance and configures your application so it can access it.
	// Use the prefix: `influx-` + `team` that you specified in the [navikt/aiven-iac](https://github.com/navikt/aiven-iac) repository.
	// +nais:doc:Availability=GCP
	Instance string `json:"instance"`
}

// +kubebuilder:validation:Pattern=`^https:\/\/.+$`
type Ingress string

type IDPorten struct {
	// Whether to enable provisioning of an ID-porten client.
	// If enabled, an ID-porten client be provisioned.
	// +nais:doc:Availability="team namespaces"
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

type GCP struct {
	// Provision BigQuery datasets and give your application's pod mountable secrets for connecting to each dataset.
	// Datasets are immutable and cannot be changed.
	// +nais:doc:Link="https://cloud.google.com/bigquery/docs"
	// +nais:doc:Availability=GCP
	BigQueryDatasets []CloudBigQueryDataset `json:"bigQueryDatasets,omitempty"`
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

type EnvVars []EnvVar

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

type EnvVarSource struct {
	FieldRef ObjectFieldSelector `json:"fieldRef"`
}

type EnvFrom struct {
	// Name of the `ConfigMap` where environment variables are specified.
	// Required unless `secret` is set.
	ConfigMap string `json:"configmap,omitempty"`
	// Name of the `Secret` where environment variables are specified.
	// Required unless `configMap` is set.
	Secret string `json:"secret,omitempty"`
}

type ObjectFieldSelector struct {
	// Field value from the `Pod` spec that should be copied into the environment variable.
	// +kubebuilder:validation:Enum="";metadata.name;metadata.namespace;metadata.labels;metadata.annotations;spec.nodeName;spec.serviceAccountName;status.hostIP;status.podIP
	FieldPath string `json:"fieldPath"`
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

type Tracing struct {
	Enabled bool `json:"enabled"`
}

type TokenX struct {
	// If enabled, will provision and configure a TokenX client and inject an accompanying secret.
	Enabled bool `json:"enabled"`
	// If enabled, secrets for TokenX will be mounted as files only, i.e. not as environment variables.
	MountSecretsAsFilesOnly bool `json:"mountSecretsAsFilesOnly,omitempty"`
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
	Kind string `json:"kind"`
	// Kubernetes _Name_.
	Name string `json:"name,omitempty"`
}

type CloudIAMPermission struct {
	// Name of the GCP role to bind the resource to.
	Role string `json:"role"`
	// IAM resource to bind the role to.
	Resource CloudIAMResource `json:"resource"`
}

type Maskinporten struct {
	// If enabled, provisions and configures a Maskinporten client with consumed scopes and/or Exposed scopes with DigDir.
	// +nais:doc:Availability="team namespaces"
	Enabled bool `json:"enabled"`
	// Schema to configure Maskinporten clients with consumed scopes and/or exposed scopes.
	Scopes MaskinportenScope `json:"scopes,omitempty"`
}

type SecureLogs struct {
	// Whether to enable a sidecar container for secure logging.
	// If enabled, a volume is mounted in the pods where secure logs can be saved.
	Enabled bool `json:"enabled"`
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

// BigQueryPermission defines access level
type BigQueryPermission string

const (
	BigQueryPermissionRead      BigQueryPermission = "READ"
	BigQueryPermissionReadWrite BigQueryPermission = "READWRITE"
)

func (b BigQueryPermission) String() string {
	return string(b)
}

func (b BigQueryPermission) GoogleType() string {
	switch b {
	case BigQueryPermissionRead:
		return "READER"
	case BigQueryPermissionReadWrite:
		return "WRITER"
	}
	return ""
}

type CloudBigQueryDataset struct {
	// Name of the BigQuery Dataset.
	// The canonical name of the dataset will be `<TEAM_PROJECT_ID>:<NAME>`.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^[a-z0-9][a-z0-9_]+$`
	Name string `json:"name"`
	// Permission level given to application.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=READ;READWRITE
	Permission BigQueryPermission `json:"permission"`
	// When set to true will delete the dataset, when the application resource is deleted.
	// NB: If no tables exist in the bigquery dataset, it _will_ delete the dataset even if this value is set/defaulted to `false`.
	// Default value is `false`.
	CascadingDelete bool `json:"cascadingDelete,omitempty"`
	// Human-readable description of what this BigQuery dataset contains, or is used for.
	// Will be visible in the GCP Console.
	Description string `json:"description,omitempty"`
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

func (envVar EnvVar) ToKubernetes() corev1.EnvVar {
	if envVar.ValueFrom != nil && envVar.ValueFrom.FieldRef.FieldPath != "" {
		return corev1.EnvVar{
			Name: envVar.Name,
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{FieldPath: envVar.ValueFrom.FieldRef.FieldPath},
			},
		}
	} else {
		return corev1.EnvVar{Name: envVar.Name, Value: envVar.Value}
	}
}

// Maps environment variables from ApplicationSpec to the ones we use in CreateSpec
func (envVars EnvVars) ToKubernetes() []corev1.EnvVar {
	var newEnvVars []corev1.EnvVar
	for _, envVar := range envVars {
		newEnvVars = append(newEnvVars, envVar.ToKubernetes())
	}

	return newEnvVars
}
