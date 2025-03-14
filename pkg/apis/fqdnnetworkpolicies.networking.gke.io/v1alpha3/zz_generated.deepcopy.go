//go:build !ignore_autogenerated

// Code generated by controller-gen. DO NOT EDIT.

package fqdnnetworkpolicies_networking_gke_io_v1alpha3

import (
	"k8s.io/api/networking/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FQDNNetworkPolicy) DeepCopyInto(out *FQDNNetworkPolicy) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FQDNNetworkPolicy.
func (in *FQDNNetworkPolicy) DeepCopy() *FQDNNetworkPolicy {
	if in == nil {
		return nil
	}
	out := new(FQDNNetworkPolicy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FQDNNetworkPolicy) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FQDNNetworkPolicyEgressRule) DeepCopyInto(out *FQDNNetworkPolicyEgressRule) {
	*out = *in
	if in.Ports != nil {
		in, out := &in.Ports, &out.Ports
		*out = make([]v1.NetworkPolicyPort, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.To != nil {
		in, out := &in.To, &out.To
		*out = make([]FQDNNetworkPolicyPeer, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FQDNNetworkPolicyEgressRule.
func (in *FQDNNetworkPolicyEgressRule) DeepCopy() *FQDNNetworkPolicyEgressRule {
	if in == nil {
		return nil
	}
	out := new(FQDNNetworkPolicyEgressRule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FQDNNetworkPolicyIngressRule) DeepCopyInto(out *FQDNNetworkPolicyIngressRule) {
	*out = *in
	if in.Ports != nil {
		in, out := &in.Ports, &out.Ports
		*out = make([]v1.NetworkPolicyPort, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.From != nil {
		in, out := &in.From, &out.From
		*out = make([]FQDNNetworkPolicyPeer, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FQDNNetworkPolicyIngressRule.
func (in *FQDNNetworkPolicyIngressRule) DeepCopy() *FQDNNetworkPolicyIngressRule {
	if in == nil {
		return nil
	}
	out := new(FQDNNetworkPolicyIngressRule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FQDNNetworkPolicyList) DeepCopyInto(out *FQDNNetworkPolicyList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]FQDNNetworkPolicy, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FQDNNetworkPolicyList.
func (in *FQDNNetworkPolicyList) DeepCopy() *FQDNNetworkPolicyList {
	if in == nil {
		return nil
	}
	out := new(FQDNNetworkPolicyList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FQDNNetworkPolicyList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FQDNNetworkPolicyPeer) DeepCopyInto(out *FQDNNetworkPolicyPeer) {
	*out = *in
	if in.FQDNs != nil {
		in, out := &in.FQDNs, &out.FQDNs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FQDNNetworkPolicyPeer.
func (in *FQDNNetworkPolicyPeer) DeepCopy() *FQDNNetworkPolicyPeer {
	if in == nil {
		return nil
	}
	out := new(FQDNNetworkPolicyPeer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FQDNNetworkPolicySpec) DeepCopyInto(out *FQDNNetworkPolicySpec) {
	*out = *in
	in.PodSelector.DeepCopyInto(&out.PodSelector)
	if in.Ingress != nil {
		in, out := &in.Ingress, &out.Ingress
		*out = make([]FQDNNetworkPolicyIngressRule, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Egress != nil {
		in, out := &in.Egress, &out.Egress
		*out = make([]FQDNNetworkPolicyEgressRule, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.PolicyTypes != nil {
		in, out := &in.PolicyTypes, &out.PolicyTypes
		*out = make([]v1.PolicyType, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FQDNNetworkPolicySpec.
func (in *FQDNNetworkPolicySpec) DeepCopy() *FQDNNetworkPolicySpec {
	if in == nil {
		return nil
	}
	out := new(FQDNNetworkPolicySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FQDNNetworkPolicyStatus) DeepCopyInto(out *FQDNNetworkPolicyStatus) {
	*out = *in
	if in.LastSyncTime != nil {
		in, out := &in.LastSyncTime, &out.LastSyncTime
		*out = (*in).DeepCopy()
	}
	if in.NextSyncTime != nil {
		in, out := &in.NextSyncTime, &out.NextSyncTime
		*out = (*in).DeepCopy()
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FQDNNetworkPolicyStatus.
func (in *FQDNNetworkPolicyStatus) DeepCopy() *FQDNNetworkPolicyStatus {
	if in == nil {
		return nil
	}
	out := new(FQDNNetworkPolicyStatus)
	in.DeepCopyInto(out)
	return out
}
