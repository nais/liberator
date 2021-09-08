// +build !ignore_autogenerated

// Code generated by controller-gen. DO NOT EDIT.

package aiven_nais_io_v1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AivenApplication) DeepCopyInto(out *AivenApplication) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AivenApplication.
func (in *AivenApplication) DeepCopy() *AivenApplication {
	if in == nil {
		return nil
	}
	out := new(AivenApplication)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AivenApplication) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AivenApplicationBuilder) DeepCopyInto(out *AivenApplicationBuilder) {
	*out = *in
	in.application.DeepCopyInto(&out.application)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AivenApplicationBuilder.
func (in *AivenApplicationBuilder) DeepCopy() *AivenApplicationBuilder {
	if in == nil {
		return nil
	}
	out := new(AivenApplicationBuilder)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AivenApplicationCondition) DeepCopyInto(out *AivenApplicationCondition) {
	*out = *in
	in.LastUpdateTime.DeepCopyInto(&out.LastUpdateTime)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AivenApplicationCondition.
func (in *AivenApplicationCondition) DeepCopy() *AivenApplicationCondition {
	if in == nil {
		return nil
	}
	out := new(AivenApplicationCondition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AivenApplicationList) DeepCopyInto(out *AivenApplicationList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]AivenApplication, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AivenApplicationList.
func (in *AivenApplicationList) DeepCopy() *AivenApplicationList {
	if in == nil {
		return nil
	}
	out := new(AivenApplicationList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AivenApplicationList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AivenApplicationSpec) DeepCopyInto(out *AivenApplicationSpec) {
	*out = *in
	if in.ExpiresAt != nil {
		in, out := &in.ExpiresAt, &out.ExpiresAt
		*out = (*in).DeepCopy()
	}
	if in.Kafka != nil {
		in, out := &in.Kafka, &out.Kafka
		*out = new(KafkaSpec)
		**out = **in
	}
	if in.Elastic != nil {
		in, out := &in.Elastic, &out.Elastic
		*out = new(ElasticSpec)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AivenApplicationSpec.
func (in *AivenApplicationSpec) DeepCopy() *AivenApplicationSpec {
	if in == nil {
		return nil
	}
	out := new(AivenApplicationSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AivenApplicationStatus) DeepCopyInto(out *AivenApplicationStatus) {
	*out = *in
	if in.SynchronizationTime != nil {
		in, out := &in.SynchronizationTime, &out.SynchronizationTime
		*out = (*in).DeepCopy()
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]AivenApplicationCondition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AivenApplicationStatus.
func (in *AivenApplicationStatus) DeepCopy() *AivenApplicationStatus {
	if in == nil {
		return nil
	}
	out := new(AivenApplicationStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ElasticSpec) DeepCopyInto(out *ElasticSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ElasticSpec.
func (in *ElasticSpec) DeepCopy() *ElasticSpec {
	if in == nil {
		return nil
	}
	out := new(ElasticSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KafkaSpec) DeepCopyInto(out *KafkaSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KafkaSpec.
func (in *KafkaSpec) DeepCopy() *KafkaSpec {
	if in == nil {
		return nil
	}
	out := new(KafkaSpec)
	in.DeepCopyInto(out)
	return out
}
