//go:build !ignore_autogenerated

// Code generated by controller-gen. DO NOT EDIT.

package storage_cnrm_cloud_google_com_v1beta1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Action) DeepCopyInto(out *Action) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Action.
func (in *Action) DeepCopy() *Action {
	if in == nil {
		return nil
	}
	out := new(Action)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BucketRef) DeepCopyInto(out *BucketRef) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BucketRef.
func (in *BucketRef) DeepCopy() *BucketRef {
	if in == nil {
		return nil
	}
	out := new(BucketRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Condition) DeepCopyInto(out *Condition) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Condition.
func (in *Condition) DeepCopy() *Condition {
	if in == nil {
		return nil
	}
	out := new(Condition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LifecycleRules) DeepCopyInto(out *LifecycleRules) {
	*out = *in
	out.Action = in.Action
	out.Condition = in.Condition
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LifecycleRules.
func (in *LifecycleRules) DeepCopy() *LifecycleRules {
	if in == nil {
		return nil
	}
	out := new(LifecycleRules)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RetentionPolicy) DeepCopyInto(out *RetentionPolicy) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RetentionPolicy.
func (in *RetentionPolicy) DeepCopy() *RetentionPolicy {
	if in == nil {
		return nil
	}
	out := new(RetentionPolicy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StorageBucket) DeepCopyInto(out *StorageBucket) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StorageBucket.
func (in *StorageBucket) DeepCopy() *StorageBucket {
	if in == nil {
		return nil
	}
	out := new(StorageBucket)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *StorageBucket) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StorageBucketAccessControl) DeepCopyInto(out *StorageBucketAccessControl) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StorageBucketAccessControl.
func (in *StorageBucketAccessControl) DeepCopy() *StorageBucketAccessControl {
	if in == nil {
		return nil
	}
	out := new(StorageBucketAccessControl)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *StorageBucketAccessControl) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StorageBucketAccessControlList) DeepCopyInto(out *StorageBucketAccessControlList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]StorageBucketAccessControl, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StorageBucketAccessControlList.
func (in *StorageBucketAccessControlList) DeepCopy() *StorageBucketAccessControlList {
	if in == nil {
		return nil
	}
	out := new(StorageBucketAccessControlList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *StorageBucketAccessControlList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StorageBucketAccessControlSpec) DeepCopyInto(out *StorageBucketAccessControlSpec) {
	*out = *in
	out.BucketRef = in.BucketRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StorageBucketAccessControlSpec.
func (in *StorageBucketAccessControlSpec) DeepCopy() *StorageBucketAccessControlSpec {
	if in == nil {
		return nil
	}
	out := new(StorageBucketAccessControlSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StorageBucketList) DeepCopyInto(out *StorageBucketList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]StorageBucket, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StorageBucketList.
func (in *StorageBucketList) DeepCopy() *StorageBucketList {
	if in == nil {
		return nil
	}
	out := new(StorageBucketList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *StorageBucketList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StorageBucketSpec) DeepCopyInto(out *StorageBucketSpec) {
	*out = *in
	if in.RetentionPolicy != nil {
		in, out := &in.RetentionPolicy, &out.RetentionPolicy
		*out = new(RetentionPolicy)
		**out = **in
	}
	if in.LifecycleRules != nil {
		in, out := &in.LifecycleRules, &out.LifecycleRules
		*out = make([]LifecycleRules, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StorageBucketSpec.
func (in *StorageBucketSpec) DeepCopy() *StorageBucketSpec {
	if in == nil {
		return nil
	}
	out := new(StorageBucketSpec)
	in.DeepCopyInto(out)
	return out
}
