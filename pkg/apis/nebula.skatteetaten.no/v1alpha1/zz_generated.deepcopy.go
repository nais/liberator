//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	v1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Application) DeepCopyInto(out *Application) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Application.
func (in *Application) DeepCopy() *Application {
	if in == nil {
		return nil
	}
	out := new(Application)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Application) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ApplicationList) DeepCopyInto(out *ApplicationList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Application, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ApplicationList.
func (in *ApplicationList) DeepCopy() *ApplicationList {
	if in == nil {
		return nil
	}
	out := new(ApplicationList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ApplicationList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ApplicationSpec) DeepCopyInto(out *ApplicationSpec) {
	*out = *in
	if in.ImagePolicy != nil {
		in, out := &in.ImagePolicy, &out.ImagePolicy
		*out = new(ImagePolicyConfig)
		**out = **in
	}
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = new(v1.Replicas)
		(*in).DeepCopyInto(*out)
	}
	in.Pod.DeepCopyInto(&out.Pod)
	if in.Logging != nil {
		in, out := &in.Logging, &out.Logging
		*out = new(LogConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.Azure != nil {
		in, out := &in.Azure, &out.Azure
		*out = new(AzureConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.Ingress != nil {
		in, out := &in.Ingress, &out.Ingress
		*out = new(IngressConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.Egress != nil {
		in, out := &in.Egress, &out.Egress
		*out = new(EgressConfig)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ApplicationSpec.
func (in *ApplicationSpec) DeepCopy() *ApplicationSpec {
	if in == nil {
		return nil
	}
	out := new(ApplicationSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AzureConfig) DeepCopyInto(out *AzureConfig) {
	*out = *in
	if in.PostgreDatabases != nil {
		in, out := &in.PostgreDatabases, &out.PostgreDatabases
		*out = make(map[string]*PostgreDatabaseConfig, len(*in))
		for key, val := range *in {
			var outVal *PostgreDatabaseConfig
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = new(PostgreDatabaseConfig)
				(*in).DeepCopyInto(*out)
			}
			(*out)[key] = outVal
		}
	}
	if in.StorageAccount != nil {
		in, out := &in.StorageAccount, &out.StorageAccount
		*out = make(map[string]*StorageAccountConfig, len(*in))
		for key, val := range *in {
			var outVal *StorageAccountConfig
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = new(StorageAccountConfig)
				**out = **in
			}
			(*out)[key] = outVal
		}
	}
	if in.CosmosDB != nil {
		in, out := &in.CosmosDB, &out.CosmosDB
		*out = make(map[string]*CosmosDBConfig, len(*in))
		for key, val := range *in {
			var outVal *CosmosDBConfig
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = new(CosmosDBConfig)
				**out = **in
			}
			(*out)[key] = outVal
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AzureConfig.
func (in *AzureConfig) DeepCopy() *AzureConfig {
	if in == nil {
		return nil
	}
	out := new(AzureConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CosmosDBConfig) DeepCopyInto(out *CosmosDBConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CosmosDBConfig.
func (in *CosmosDBConfig) DeepCopy() *CosmosDBConfig {
	if in == nil {
		return nil
	}
	out := new(CosmosDBConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EgressConfig) DeepCopyInto(out *EgressConfig) {
	*out = *in
	if in.External != nil {
		in, out := &in.External, &out.External
		*out = make(map[string]ExternalEgressConfig, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
	if in.Internal != nil {
		in, out := &in.Internal, &out.Internal
		*out = make(map[string]InternalEgressConfig, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EgressConfig.
func (in *EgressConfig) DeepCopy() *EgressConfig {
	if in == nil {
		return nil
	}
	out := new(EgressConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExternalEgressConfig) DeepCopyInto(out *ExternalEgressConfig) {
	*out = *in
	if in.Ports != nil {
		in, out := &in.Ports, &out.Ports
		*out = make([]PortConfig, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExternalEgressConfig.
func (in *ExternalEgressConfig) DeepCopy() *ExternalEgressConfig {
	if in == nil {
		return nil
	}
	out := new(ExternalEgressConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ImagePolicyConfig) DeepCopyInto(out *ImagePolicyConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ImagePolicyConfig.
func (in *ImagePolicyConfig) DeepCopy() *ImagePolicyConfig {
	if in == nil {
		return nil
	}
	out := new(ImagePolicyConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IngressConfig) DeepCopyInto(out *IngressConfig) {
	*out = *in
	if in.Public != nil {
		in, out := &in.Public, &out.Public
		*out = make(map[string]PublicIngressConfig, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Internal != nil {
		in, out := &in.Internal, &out.Internal
		*out = make(map[string]InternalIngressConfig, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IngressConfig.
func (in *IngressConfig) DeepCopy() *IngressConfig {
	if in == nil {
		return nil
	}
	out := new(IngressConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InternalEgressConfig) DeepCopyInto(out *InternalEgressConfig) {
	*out = *in
	if in.Ports != nil {
		in, out := &in.Ports, &out.Ports
		*out = make([]PortConfig, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InternalEgressConfig.
func (in *InternalEgressConfig) DeepCopy() *InternalEgressConfig {
	if in == nil {
		return nil
	}
	out := new(InternalEgressConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InternalIngressConfig) DeepCopyInto(out *InternalIngressConfig) {
	*out = *in
	if in.Ports != nil {
		in, out := &in.Ports, &out.Ports
		*out = make([]PortConfig, len(*in))
		copy(*out, *in)
	}
	if in.Methods != nil {
		in, out := &in.Methods, &out.Methods
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Paths != nil {
		in, out := &in.Paths, &out.Paths
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InternalIngressConfig.
func (in *InternalIngressConfig) DeepCopy() *InternalIngressConfig {
	if in == nil {
		return nil
	}
	out := new(InternalIngressConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LogConfig) DeepCopyInto(out *LogConfig) {
	*out = *in
	if in.Splunk != nil {
		in, out := &in.Splunk, &out.Splunk
		*out = make(map[string]*SplunkLoggingConfig, len(*in))
		for key, val := range *in {
			var outVal *SplunkLoggingConfig
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = new(SplunkLoggingConfig)
				(*in).DeepCopyInto(*out)
			}
			(*out)[key] = outVal
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LogConfig.
func (in *LogConfig) DeepCopy() *LogConfig {
	if in == nil {
		return nil
	}
	out := new(LogConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PodConfig) DeepCopyInto(out *PodConfig) {
	*out = *in
	if in.Command != nil {
		in, out := &in.Command, &out.Command
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Env != nil {
		in, out := &in.Env, &out.Env
		*out = make(v1.EnvVars, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.EnvFrom != nil {
		in, out := &in.EnvFrom, &out.EnvFrom
		*out = make([]v1.EnvFrom, len(*in))
		copy(*out, *in)
	}
	if in.FilesFrom != nil {
		in, out := &in.FilesFrom, &out.FilesFrom
		*out = make([]v1.FilesFrom, len(*in))
		copy(*out, *in)
	}
	if in.Prometheus != nil {
		in, out := &in.Prometheus, &out.Prometheus
		*out = new(PrometheusConfig)
		**out = **in
	}
	if in.Readiness != nil {
		in, out := &in.Readiness, &out.Readiness
		*out = new(v1.Probe)
		**out = **in
	}
	if in.Liveness != nil {
		in, out := &in.Liveness, &out.Liveness
		*out = new(v1.Probe)
		**out = **in
	}
	if in.Startup != nil {
		in, out := &in.Startup, &out.Startup
		*out = new(v1.Probe)
		**out = **in
	}
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = new(v1.ResourceRequirements)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PodConfig.
func (in *PodConfig) DeepCopy() *PodConfig {
	if in == nil {
		return nil
	}
	out := new(PodConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PortConfig) DeepCopyInto(out *PortConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PortConfig.
func (in *PortConfig) DeepCopy() *PortConfig {
	if in == nil {
		return nil
	}
	out := new(PortConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PostgreDatabaseConfig) DeepCopyInto(out *PostgreDatabaseConfig) {
	*out = *in
	if in.Users != nil {
		in, out := &in.Users, &out.Users
		*out = make(map[string]*PostgreDatabaseUser, len(*in))
		for key, val := range *in {
			var outVal *PostgreDatabaseUser
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = new(PostgreDatabaseUser)
				**out = **in
			}
			(*out)[key] = outVal
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PostgreDatabaseConfig.
func (in *PostgreDatabaseConfig) DeepCopy() *PostgreDatabaseConfig {
	if in == nil {
		return nil
	}
	out := new(PostgreDatabaseConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PostgreDatabaseUser) DeepCopyInto(out *PostgreDatabaseUser) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PostgreDatabaseUser.
func (in *PostgreDatabaseUser) DeepCopy() *PostgreDatabaseUser {
	if in == nil {
		return nil
	}
	out := new(PostgreDatabaseUser)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PrometheusConfig) DeepCopyInto(out *PrometheusConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PrometheusConfig.
func (in *PrometheusConfig) DeepCopy() *PrometheusConfig {
	if in == nil {
		return nil
	}
	out := new(PrometheusConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PublicIngressConfig) DeepCopyInto(out *PublicIngressConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PublicIngressConfig.
func (in *PublicIngressConfig) DeepCopy() *PublicIngressConfig {
	if in == nil {
		return nil
	}
	out := new(PublicIngressConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SplunkLoggingConfig) DeepCopyInto(out *SplunkLoggingConfig) {
	*out = *in
	if in.Enabled != nil {
		in, out := &in.Enabled, &out.Enabled
		*out = new(bool)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SplunkLoggingConfig.
func (in *SplunkLoggingConfig) DeepCopy() *SplunkLoggingConfig {
	if in == nil {
		return nil
	}
	out := new(SplunkLoggingConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StorageAccountConfig) DeepCopyInto(out *StorageAccountConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StorageAccountConfig.
func (in *StorageAccountConfig) DeepCopy() *StorageAccountConfig {
	if in == nil {
		return nil
	}
	out := new(StorageAccountConfig)
	in.DeepCopyInto(out)
	return out
}
