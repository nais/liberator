package storage_cnrm_cloud_google_com_v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	SchemeBuilder.Register(
		&StorageBucket{},
		&StorageBucketList{},
		&StorageBucketAccessControl{},
		&StorageBucketAccessControlList{},
	)
}

// +kubebuilder:object:root=true
type StorageBucket struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              StorageBucketSpec `json:"spec"`
}

type PublicAccessPrevention string

const (
	PublicAccessPreventionEnforced  PublicAccessPrevention = "enforced"
	PublicAccessPreventionInherited PublicAccessPrevention = "inherited"
)

type StorageBucketSpec struct {
	ResourceID               string           `json:"resourceID,omitempty"`
	Location                 string           `json:"location"`
	UniformBucketLevelAccess bool             `json:"uniformBucketLevelAccess,omitempty"`
	RetentionPolicy          *RetentionPolicy `json:"retentionPolicy,omitempty"`
	LifecycleRules           []LifecycleRules `json:"lifecycleRule,omitempty"`
	// +kubebuilder:validation:Enum=inherited;enforced
	PublicAccessPrevention PublicAccessPrevention `json:"publicAccessPrevention,omitempty"`
}

// +kubebuilder:object:root=true
type StorageBucketList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []StorageBucket `json:"items"`
}

// +kubebuilder:object:root=true
type StorageBucketAccessControl struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              StorageBucketAccessControlSpec `json:"spec"`
}

type StorageBucketAccessControlSpec struct {
	BucketRef BucketRef `json:"bucketRef"`
	Entity    string    `json:"entity"`
	Role      string    `json:"role"`
}

type BucketRef struct {
	Name string `json:"name"`
}

// +kubebuilder:object:root=true
type StorageBucketAccessControlList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []StorageBucketAccessControl `json:"items"`
}

type RetentionPolicy struct {
	RetentionPeriod int `json:"retentionPeriod,omitempty"`
}

type LifecycleRules struct {
	Action    Action    `json:"action"`
	Condition Condition `json:"condition"`
}

type Action struct {
	Type string `json:"type,omitempty"`
}

type Condition struct {
	Age                 int    `json:"age,omitempty"`
	CreatedBefore       string `json:"createdBefore,omitempty"`
	DaysSinceCustomTime int    `json:"daysSinceCustomTime,omitempty"`
	NumNewerVersions    int    `json:"numNewerVersions,omitempty"`
	WithState           string `json:"withState,omitempty"`
}
