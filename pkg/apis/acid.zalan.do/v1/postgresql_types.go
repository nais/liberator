package v1

import (
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient

// Postgresql defines PostgreSQL Custom Resource Definition Object.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Postgresql struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PostgresSpec   `json:"spec"`
	Status PostgresStatus `json:"status"`
	Error  string         `json:"-"`
}

// PostgresSpec defines the specification for the PostgreSQL CRD.
// This is a reduced version of the spec the operator uses.
// See https://github.com/zalando/postgres-operator/blob/master/pkg/apis/acid.zalan.do/v1/postgresql_type.go for original
type PostgresSpec struct {
	PostgresqlParam `json:"postgresql"`
	Volume          `json:"volume,omitempty"`
	Patroni         `json:"patroni,omitempty"`
	*Resources      `json:"resources,omitempty"`

	EnableConnectionPooler        *bool             `json:"enableConnectionPooler,omitempty"`
	EnableReplicaConnectionPooler *bool             `json:"enableReplicaConnectionPooler,omitempty"`
	ConnectionPooler              *ConnectionPooler `json:"connectionPooler,omitempty"`

	TeamID      string `json:"teamId"`
	DockerImage string `json:"dockerImage,omitempty"`

	Users                          map[string]UserFlags `json:"users,omitempty"`
	UsersIgnoringSecretRotation    []string             `json:"usersIgnoringSecretRotation,omitempty"`
	UsersWithSecretRotation        []string             `json:"usersWithSecretRotation,omitempty"`
	UsersWithInPlaceSecretRotation []string             `json:"usersWithInPlaceSecretRotation,omitempty"`

	NumberOfInstances      int32                       `json:"numberOfInstances"`
	MaintenanceWindows     []MaintenanceWindow         `json:"maintenanceWindows,omitempty"`
	Clone                  *CloneDescription           `json:"clone,omitempty"`
	PreparedDatabases      map[string]PreparedDatabase `json:"preparedDatabases,omitempty"`
	NodeAffinity           *v1.NodeAffinity            `json:"nodeAffinity,omitempty"`
	Tolerations            []v1.Toleration             `json:"tolerations,omitempty"`
	Sidecars               []Sidecar                   `json:"sidecars,omitempty"`
	InitContainers         []v1.Container              `json:"initContainers,omitempty"`
	PodPriorityClassName   string                      `json:"podPriorityClassName,omitempty"`
	ShmVolume              *bool                       `json:"enableShmVolume,omitempty"`
	EnableLogicalBackup    bool                        `json:"enableLogicalBackup,omitempty"`
	LogicalBackupRetention string                      `json:"logicalBackupRetention,omitempty"`
	LogicalBackupSchedule  string                      `json:"logicalBackupSchedule,omitempty"`
	TLS                    *TLSDescription             `json:"tls,omitempty"`
	AdditionalVolumes      []AdditionalVolume          `json:"additionalVolumes,omitempty"`
	Env                    []v1.EnvVar                 `json:"env,omitempty"`
}

// PostgresqlList defines a list of PostgreSQL clusters.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PostgresqlList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Postgresql `json:"items"`
}

// PreparedDatabase describes elements to be bootstrapped
type PreparedDatabase struct {
	PreparedSchemas map[string]PreparedSchema `json:"schemas,omitempty"`
	DefaultUsers    bool                      `json:"defaultUsers,omitempty" defaults:"false"`
	Extensions      map[string]string         `json:"extensions,omitempty"`
	SecretNamespace string                    `json:"secretNamespace,omitempty"`
}

// PreparedSchema describes elements to be bootstrapped per schema
type PreparedSchema struct {
	DefaultRoles *bool `json:"defaultRoles,omitempty" defaults:"true"`
	DefaultUsers bool  `json:"defaultUsers,omitempty" defaults:"false"`
}

// MaintenanceWindow describes the time window when the operator is allowed to do maintenance on a cluster.
type MaintenanceWindow struct {
	Everyday  bool         `json:"everyday,omitempty"`
	Weekday   time.Weekday `json:"weekday,omitempty"`
	StartTime metav1.Time  `json:"startTime,omitempty"`
	EndTime   metav1.Time  `json:"endTime,omitempty"`
}

// Volume describes a single volume in the manifest.
type Volume struct {
	Selector      *metav1.LabelSelector `json:"selector,omitempty"`
	Size          string                `json:"size"`
	StorageClass  string                `json:"storageClass,omitempty"`
	SubPath       string                `json:"subPath,omitempty"`
	IsSubPathExpr *bool                 `json:"isSubPathExpr,omitempty"`
	Iops          *int64                `json:"iops,omitempty"`
	Throughput    *int64                `json:"throughput,omitempty"`
	VolumeType    string                `json:"type,omitempty"`
}

// AdditionalVolume specs additional optional volumes for statefulset
type AdditionalVolume struct {
	Name             string          `json:"name"`
	MountPath        string          `json:"mountPath"`
	SubPath          string          `json:"subPath,omitempty"`
	IsSubPathExpr    *bool           `json:"isSubPathExpr,omitempty"`
	TargetContainers []string        `json:"targetContainers"`
	VolumeSource     v1.VolumeSource `json:"volumeSource"`
}

// PostgresqlParam describes PostgreSQL version and pairs of configuration parameter name - values.
type PostgresqlParam struct {
	PgVersion  string            `json:"version"`
	Parameters map[string]string `json:"parameters,omitempty"`
}

// ResourceDescription describes CPU and memory resources defined for a cluster.
type ResourceDescription struct {
	CPU          *string `json:"cpu,omitempty"`
	Memory       *string `json:"memory,omitempty"`
	HugePages2Mi *string `json:"hugepages-2Mi,omitempty"`
	HugePages1Gi *string `json:"hugepages-1Gi,omitempty"`
}

// Resources describes requests and limits for the cluster resouces.
type Resources struct {
	ResourceRequests ResourceDescription `json:"requests,omitempty"`
	ResourceLimits   ResourceDescription `json:"limits,omitempty"`
}

// Patroni contains Patroni-specific configuration
type Patroni struct {
	InitDB                map[string]string            `json:"initdb,omitempty"`
	PgHba                 []string                     `json:"pg_hba,omitempty"`
	TTL                   uint32                       `json:"ttl,omitempty"`
	LoopWait              uint32                       `json:"loop_wait,omitempty"`
	RetryTimeout          uint32                       `json:"retry_timeout,omitempty"`
	Slots                 map[string]map[string]string `json:"slots,omitempty"`
	SynchronousMode       bool                         `json:"synchronous_mode,omitempty"`
	SynchronousModeStrict bool                         `json:"synchronous_mode_strict,omitempty"`
	SynchronousNodeCount  uint32                       `json:"synchronous_node_count,omitempty" defaults:"1"`
	FailsafeMode          *bool                        `json:"failsafe_mode,omitempty"`
}

// TLSDescription specs TLS properties
type TLSDescription struct {
	SecretName      string `json:"secretName,omitempty"`
	CertificateFile string `json:"certificateFile,omitempty"`
	PrivateKeyFile  string `json:"privateKeyFile,omitempty"`
	CAFile          string `json:"caFile,omitempty"`
	CASecretName    string `json:"caSecretName,omitempty"`
}

// CloneDescription describes which cluster the new should clone and up to which point in time
type CloneDescription struct {
	ClusterName  string `json:"cluster,omitempty"`
	UID          string `json:"uid,omitempty"`
	EndTimestamp string `json:"timestamp,omitempty"`
}

// Sidecar defines a container to be run in the same pod as the Postgres container.
type Sidecar struct {
	*Resources  `json:"resources,omitempty"`
	Name        string             `json:"name,omitempty"`
	DockerImage string             `json:"image,omitempty"`
	Ports       []v1.ContainerPort `json:"ports,omitempty"`
	Env         []v1.EnvVar        `json:"env,omitempty"`
	Command     []string           `json:"command,omitempty"`
}

// UserFlags defines flags (such as superuser, nologin) that could be assigned to individual users
type UserFlags []string

// PostgresStatus contains status of the PostgreSQL cluster (running, creation failed etc.)
type PostgresStatus struct {
	PostgresClusterStatus string `json:"PostgresClusterStatus"`
}

// ConnectionPooler Options for connection pooler
//
// TODO: prepared snippets of configuration, one can choose via type, e.g.
// pgbouncer-large (with higher resources) or odyssey-small (with smaller
// resources)
// Type              string `json:"type,omitempty"`
//
// TODO: figure out what other important parameters of the connection pooler it
// makes sense to expose. E.g. pool size (min/max boundaries), max client
// connections etc.
type ConnectionPooler struct {
	NumberOfInstances *int32 `json:"numberOfInstances,omitempty"`
	Schema            string `json:"schema,omitempty"`
	User              string `json:"user,omitempty"`
	Mode              string `json:"mode,omitempty"`
	DockerImage       string `json:"dockerImage,omitempty"`
	MaxDBConnections  *int32 `json:"maxDBConnections,omitempty"`

	*Resources `json:"resources,omitempty"`
}
