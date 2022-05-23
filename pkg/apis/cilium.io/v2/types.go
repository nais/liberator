package cilium_io_v2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	SchemeBuilder.Register(
		&NetworkPolicy{},
		&NetworkPolicyList{},
	)
}

type NetworkPolicySpec struct {
	EndpointSelector *metav1.LabelSelector `json:"endpointSelector,omitempty"`
	Ingress          []Ingress             `json:"ingress,omitempty"`
	Egress           *Egress               `json:"egress,omitempty"`
}

type Ingress struct {
	FromEntities  []string              `json:"fromEntities,omitempty"`
	FromEndpoints *metav1.LabelSelector `json:"fromEndpoints,omitempty"`
	ToPorts       []Ports               `json:"toPorts,omitempty"`
}

type Ports struct {
	Ports []Port `json:"ports,omitempty"`
}

type Port struct {
	Port int32 `json:"port,omitempty"`
}

type CIDRSet struct {
	CIDR   string   `json:"cidr,omitempty"`
	Except []string `json:"except,omitempty"`
}

type FQDN struct {
	MatchName string `json:"matchName,omitempty"`
}

type Egress struct {
	ToCIDRSet   []CIDRSet               `json:"toCIDRSet,omitempty"`
	ToPorts     []Ports                 `json:"toPorts,omitempty"`
	ToFQDNs     []FQDN                  `json:"toFQDNs,omitempty"`
	ToEndpoints []*metav1.LabelSelector `json:"toEndpoints,omitempty"`
}

// +kubebuilder:object:root=true
type NetworkPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              NetworkPolicySpec `json:"spec"`
}

// +kubebuilder:object:root=true
type NetworkPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NetworkPolicy `json:"items"`
}
