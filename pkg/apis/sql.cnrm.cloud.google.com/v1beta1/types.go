package sql_cnrm_cloud_google_com_v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	SchemeBuilder.Register(
		&SQLInstance{},
		&SQLInstanceList{},
		&SQLDatabase{},
		&SQLDatabaseList{},
		&SQLUser{},
		&SQLUserList{},
	)
}

type SQLInstanceSpec struct {
	ResourceID      string              `json:"resourceID,omitempty"`
	DatabaseVersion string              `json:"databaseVersion"`
	Region          string              `json:"region"`
	Settings        SQLInstanceSettings `json:"settings"`
}

type MaintenanceWindow struct {
	Day  int `json:"day"`
	Hour int `json:"hour"`
}

type SQLDatabaseFlag struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type SQLInstanceSettings struct {
	AvailabilityType    string                           `json:"availabilityType"`
	BackupConfiguration SQLInstanceBackupConfiguration   `json:"backupConfiguration"`
	InsightsConfig      SQLInstanceInsightsConfiguration `json:"insightsConfig"`
	IpConfiguration     SQLInstanceIpConfiguration       `json:"ipConfiguration"`
	DiskAutoresize      bool                             `json:"diskAutoresize"`
	DiskSize            int                              `json:"diskSize"`
	DiskType            string                           `json:"diskType"`
	Tier                string                           `json:"tier"`
	MaintenanceWindow   *MaintenanceWindow               `json:"maintenanceWindow,omitempty"`
	DatabaseFlags       []SQLDatabaseFlag                `json:"databaseFlags"`
}

type SQLInstanceInsightsConfiguration struct {
	QueryInsightsEnabled  bool `json:"queryInsightsEnabled,omitempty"`
	QueryStringLength     int  `json:"queryStringLength,omitempty"`
	RecordApplicationTags bool `json:"recordApplicationTags"`
	RecordClientAddress   bool `json:"recordClientAddress"`
}

type SQLInstanceBackupConfiguration struct {
	Enabled                    bool                               `json:"enabled"`
	StartTime                  string                             `json:"startTime"`
	PointInTimeRecoveryEnabled bool                               `json:"pointInTimeRecoveryEnabled"`
	BackupRetentionSettings    *SQLInstanceBackupRetentionSetting `json:"backupRetentionSettings,omitempty"`
}

type SQLInstanceBackupRetentionSetting struct {
	RetainedBackups int `json:"retainedBackups"`
}

type PrivateNetworkRef struct {
	External string `json:"external"`
}

type SQLInstanceIpConfiguration struct {
	RequireSsl        bool               `json:"requireSsl"`
	PrivateNetworkRef *PrivateNetworkRef `json:"privateNetworkRef, omitempty"`
}

// +kubebuilder:object:root=true
type SQLInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              SQLInstanceSpec `json:"spec"`
}

// +kubebuilder:object:root=true
type SQLInstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SQLInstance `json:"items"`
}

type InstanceRef struct {
	Name string `json:"name"`
}

type SQLDatabaseSpec struct {
	ResourceID  string      `json:"resourceID,omitempty"`
	InstanceRef InstanceRef `json:"instanceRef"`
}

// +kubebuilder:object:root=true
type SQLDatabase struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              SQLDatabaseSpec `json:"spec"`
}

// +kubebuilder:object:root=true
type SQLDatabaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SQLDatabase `json:"items"`
}

type SecretRef struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type SqlUserPasswordSecretKeyRef struct {
	SecretKeyRef SecretRef `json:"secretKeyRef"`
}

type SqlUserPasswordValue struct {
	ValueFrom SqlUserPasswordSecretKeyRef `json:"valueFrom"`
}

type SQLUserSpec struct {
	ResourceID  string               `json:"resourceID,omitempty"`
	InstanceRef InstanceRef          `json:"instanceRef"`
	Host        string               `json:"host"`
	Password    SqlUserPasswordValue `json:"password"`
}

// +kubebuilder:object:root=true
type SQLUser struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              SQLUserSpec `json:"spec"`
}

// +kubebuilder:object:root=true
type SQLUserList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SQLUser `json:"items"`
}
