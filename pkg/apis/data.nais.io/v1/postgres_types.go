package v1

import (
	"github.com/nais/liberator/pkg/apis/nais.io/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PostgresResources struct {
	// Disk size for the Postgres cluster.
	// +kubebuilder:validation:required
	DiskSize resource.Quantity `json:"diskSize"`

	// CPU resources for the Postgres cluster.
	// +kubebuilder:validation:required
	Cpu resource.Quantity `json:"cpu"`

	// Memory resources for the Postgres cluster.
	// +kubebuilder:validation:required
	Memory resource.Quantity `json:"memory"`
}

// +kubebuilder:validation:Enum=read;write;function;role;ddl;misc;misc_set;all
type PostgresAuditStatementClass string

type PostgresAudit struct {
	// Enable audit logging for the Postgres cluster.
	Enabled bool `json:"enabled,omitempty"`

	// Statement classes to log.
	// +nais:doc:Default="ddl,write"
	StatementClasses []PostgresAuditStatementClass `json:"statementClasses,omitempty"`
}

type PostgresCluster struct {
	Resources PostgresResources `json:"resources"`

	// Major version of Postgres to use.
	// +kubebuilder:validation:required
	// +kubebuilder:validation:Enum="17";"16"
	MajorVersion string `json:"majorVersion"`

	// High availability cluster.
	HighAvailability bool `json:"highAvailability,omitempty"`

	// Allow deletion of the Postgres cluster when the application is deleted.
	AllowDeletion bool `json:"allowDeletion,omitempty"`

	// Configure audit logging for the Postgres cluster.
	Audit *PostgresAudit `json:"audit,omitempty"`
}

type PostgresExtension struct {
	// Name of the Postgres extension to enable.
	// +kubebuilder:validation:required
	Name string `json:"name"`
}

type PostgresDatabase struct {
	// Collation for the Postgres database.
	// +kubebuilder:validation:Enum=nb_NO;en_US
	Collation string `json:"collation,omitempty"`

	// Extensions to enable in the Postgres database.
	Extensions []PostgresExtension `json:"extensions,omitempty"`
}

// PostgresSpec defines the desired state of Postgres
type PostgresSpec struct {
	// Cluster configures the Postgres cluster
	Cluster PostgresCluster `json:"cluster"`

	// Database configures the Postgres database.
	Database *PostgresDatabase `json:"database,omitempty"`

	// MaintenanceWindow configures the maintenance window for the Postgres cluster.
	MaintenanceWindow *nais_io_v1.Maintenance `json:"maintenanceWindow,omitempty"`
}

// PostgresStatus defines the observed state of Postgres.
type PostgresStatus struct {
	SynchronizationTime     int64  `json:"synchronizationTime,omitempty"`
	RolloutCompleteTime     int64  `json:"rolloutCompleteTime,omitempty"`
	CorrelationID           string `json:"correlationID,omitempty"`
	DeploymentRolloutStatus string `json:"deploymentRolloutStatus,omitempty"`
	SynchronizationState    string `json:"synchronizationState,omitempty"`
	SynchronizationHash     string `json:"synchronizationHash,omitempty"`
	ObservedGeneration      int64  `json:"observedGeneration,omitempty"`

	// conditions represent the current state of the Postgres resource.
	// Each condition has a unique type and reflects the status of a specific aspect of the resource.
	//
	// Standard condition types include:
	// - "Available": the resource is fully functional
	// - "Progressing": the resource is being created or updated
	// - "Degraded": the resource failed to reach or maintain its desired state
	//
	// The status of each condition is one of True, False, or Unknown.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Postgres is the Schema for the postgres API
type Postgres struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of Postgres
	// +required
	Spec PostgresSpec `json:"spec"`

	// status defines the observed state of Postgres
	// +optional
	Status PostgresStatus `json:"status,omitempty,omitzero"`
}

func (p *Postgres) GetCorrelationId() string {
	return p.Annotations[nais_io_v1.DeploymentCorrelationIDAnnotation]
}

// +kubebuilder:object:root=true

// PostgresList contains a list of Postgres
type PostgresList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Postgres `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Postgres{}, &PostgresList{})
}
