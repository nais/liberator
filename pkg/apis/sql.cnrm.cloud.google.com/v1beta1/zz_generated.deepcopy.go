//go:build !ignore_autogenerated

// Code generated by controller-gen. DO NOT EDIT.

package sql_cnrm_cloud_google_com_v1beta1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstanceRef) DeepCopyInto(out *InstanceRef) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstanceRef.
func (in *InstanceRef) DeepCopy() *InstanceRef {
	if in == nil {
		return nil
	}
	out := new(InstanceRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MaintenanceWindow) DeepCopyInto(out *MaintenanceWindow) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MaintenanceWindow.
func (in *MaintenanceWindow) DeepCopy() *MaintenanceWindow {
	if in == nil {
		return nil
	}
	out := new(MaintenanceWindow)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PrivateNetworkRef) DeepCopyInto(out *PrivateNetworkRef) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PrivateNetworkRef.
func (in *PrivateNetworkRef) DeepCopy() *PrivateNetworkRef {
	if in == nil {
		return nil
	}
	out := new(PrivateNetworkRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SQLDatabase) DeepCopyInto(out *SQLDatabase) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SQLDatabase.
func (in *SQLDatabase) DeepCopy() *SQLDatabase {
	if in == nil {
		return nil
	}
	out := new(SQLDatabase)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SQLDatabase) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SQLDatabaseFlag) DeepCopyInto(out *SQLDatabaseFlag) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SQLDatabaseFlag.
func (in *SQLDatabaseFlag) DeepCopy() *SQLDatabaseFlag {
	if in == nil {
		return nil
	}
	out := new(SQLDatabaseFlag)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SQLDatabaseList) DeepCopyInto(out *SQLDatabaseList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]SQLDatabase, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SQLDatabaseList.
func (in *SQLDatabaseList) DeepCopy() *SQLDatabaseList {
	if in == nil {
		return nil
	}
	out := new(SQLDatabaseList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SQLDatabaseList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SQLDatabaseSpec) DeepCopyInto(out *SQLDatabaseSpec) {
	*out = *in
	out.InstanceRef = in.InstanceRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SQLDatabaseSpec.
func (in *SQLDatabaseSpec) DeepCopy() *SQLDatabaseSpec {
	if in == nil {
		return nil
	}
	out := new(SQLDatabaseSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SQLInstance) DeepCopyInto(out *SQLInstance) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SQLInstance.
func (in *SQLInstance) DeepCopy() *SQLInstance {
	if in == nil {
		return nil
	}
	out := new(SQLInstance)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SQLInstance) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SQLInstanceBackupConfiguration) DeepCopyInto(out *SQLInstanceBackupConfiguration) {
	*out = *in
	if in.BackupRetentionSettings != nil {
		in, out := &in.BackupRetentionSettings, &out.BackupRetentionSettings
		*out = new(SQLInstanceBackupRetentionSetting)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SQLInstanceBackupConfiguration.
func (in *SQLInstanceBackupConfiguration) DeepCopy() *SQLInstanceBackupConfiguration {
	if in == nil {
		return nil
	}
	out := new(SQLInstanceBackupConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SQLInstanceBackupRetentionSetting) DeepCopyInto(out *SQLInstanceBackupRetentionSetting) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SQLInstanceBackupRetentionSetting.
func (in *SQLInstanceBackupRetentionSetting) DeepCopy() *SQLInstanceBackupRetentionSetting {
	if in == nil {
		return nil
	}
	out := new(SQLInstanceBackupRetentionSetting)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SQLInstanceInsightsConfiguration) DeepCopyInto(out *SQLInstanceInsightsConfiguration) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SQLInstanceInsightsConfiguration.
func (in *SQLInstanceInsightsConfiguration) DeepCopy() *SQLInstanceInsightsConfiguration {
	if in == nil {
		return nil
	}
	out := new(SQLInstanceInsightsConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SQLInstanceIpConfiguration) DeepCopyInto(out *SQLInstanceIpConfiguration) {
	*out = *in
	if in.PrivateNetworkRef != nil {
		in, out := &in.PrivateNetworkRef, &out.PrivateNetworkRef
		*out = new(PrivateNetworkRef)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SQLInstanceIpConfiguration.
func (in *SQLInstanceIpConfiguration) DeepCopy() *SQLInstanceIpConfiguration {
	if in == nil {
		return nil
	}
	out := new(SQLInstanceIpConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SQLInstanceList) DeepCopyInto(out *SQLInstanceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]SQLInstance, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SQLInstanceList.
func (in *SQLInstanceList) DeepCopy() *SQLInstanceList {
	if in == nil {
		return nil
	}
	out := new(SQLInstanceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SQLInstanceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SQLInstanceSettings) DeepCopyInto(out *SQLInstanceSettings) {
	*out = *in
	in.BackupConfiguration.DeepCopyInto(&out.BackupConfiguration)
	out.InsightsConfig = in.InsightsConfig
	in.IpConfiguration.DeepCopyInto(&out.IpConfiguration)
	if in.MaintenanceWindow != nil {
		in, out := &in.MaintenanceWindow, &out.MaintenanceWindow
		*out = new(MaintenanceWindow)
		**out = **in
	}
	if in.DatabaseFlags != nil {
		in, out := &in.DatabaseFlags, &out.DatabaseFlags
		*out = make([]SQLDatabaseFlag, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SQLInstanceSettings.
func (in *SQLInstanceSettings) DeepCopy() *SQLInstanceSettings {
	if in == nil {
		return nil
	}
	out := new(SQLInstanceSettings)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SQLInstanceSpec) DeepCopyInto(out *SQLInstanceSpec) {
	*out = *in
	in.Settings.DeepCopyInto(&out.Settings)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SQLInstanceSpec.
func (in *SQLInstanceSpec) DeepCopy() *SQLInstanceSpec {
	if in == nil {
		return nil
	}
	out := new(SQLInstanceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SQLSSLCert) DeepCopyInto(out *SQLSSLCert) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	out.Spec = in.Spec
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SQLSSLCert.
func (in *SQLSSLCert) DeepCopy() *SQLSSLCert {
	if in == nil {
		return nil
	}
	out := new(SQLSSLCert)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SQLSSLCert) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SQLSSLCertList) DeepCopyInto(out *SQLSSLCertList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]SQLSSLCert, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SQLSSLCertList.
func (in *SQLSSLCertList) DeepCopy() *SQLSSLCertList {
	if in == nil {
		return nil
	}
	out := new(SQLSSLCertList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SQLSSLCertList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SQLSSLCertSpec) DeepCopyInto(out *SQLSSLCertSpec) {
	*out = *in
	out.InstanceRef = in.InstanceRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SQLSSLCertSpec.
func (in *SQLSSLCertSpec) DeepCopy() *SQLSSLCertSpec {
	if in == nil {
		return nil
	}
	out := new(SQLSSLCertSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SQLUser) DeepCopyInto(out *SQLUser) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SQLUser.
func (in *SQLUser) DeepCopy() *SQLUser {
	if in == nil {
		return nil
	}
	out := new(SQLUser)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SQLUser) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SQLUserList) DeepCopyInto(out *SQLUserList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]SQLUser, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SQLUserList.
func (in *SQLUserList) DeepCopy() *SQLUserList {
	if in == nil {
		return nil
	}
	out := new(SQLUserList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SQLUserList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SQLUserSpec) DeepCopyInto(out *SQLUserSpec) {
	*out = *in
	out.InstanceRef = in.InstanceRef
	out.Password = in.Password
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SQLUserSpec.
func (in *SQLUserSpec) DeepCopy() *SQLUserSpec {
	if in == nil {
		return nil
	}
	out := new(SQLUserSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretRef) DeepCopyInto(out *SecretRef) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretRef.
func (in *SecretRef) DeepCopy() *SecretRef {
	if in == nil {
		return nil
	}
	out := new(SecretRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SqlUserPasswordSecretKeyRef) DeepCopyInto(out *SqlUserPasswordSecretKeyRef) {
	*out = *in
	out.SecretKeyRef = in.SecretKeyRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SqlUserPasswordSecretKeyRef.
func (in *SqlUserPasswordSecretKeyRef) DeepCopy() *SqlUserPasswordSecretKeyRef {
	if in == nil {
		return nil
	}
	out := new(SqlUserPasswordSecretKeyRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SqlUserPasswordValue) DeepCopyInto(out *SqlUserPasswordValue) {
	*out = *in
	out.ValueFrom = in.ValueFrom
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SqlUserPasswordValue.
func (in *SqlUserPasswordValue) DeepCopy() *SqlUserPasswordValue {
	if in == nil {
		return nil
	}
	out := new(SqlUserPasswordValue)
	in.DeepCopyInto(out)
	return out
}
