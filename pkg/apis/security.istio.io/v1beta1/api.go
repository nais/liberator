package security_istio_io_v1beta1

// Source specifies the source identities of a request. Fields in the source are
// ANDed together.
//
// For example, the following source matches if the principal is "admin" or "dev"
// and the namespace is "prod" or "test" and the ip is not "1.2.3.4".
//
// ```yaml
// principals: ["admin", "dev"]
// namespaces: ["prod", "test"]
// not_ipblocks: ["1.2.3.4"]
// ```
type Source struct {
	// Optional. A list of source peer identities (i.e. service account), which
	// matches to the "source.principal" attribute. This field requires mTLS enabled.
	//
	// If not set, any principal is allowed.
	Principals []string `protobuf:"bytes,1,rep,name=principals,proto3" json:"principals,omitempty"`
	// Optional. A list of negative match of source peer identities.
	NotPrincipals []string `protobuf:"bytes,5,rep,name=not_principals,json=notPrincipals,proto3" json:"not_principals,omitempty"`
	// Optional. A list of request identities (i.e. "iss/sub" claims), which
	// matches to the "request.auth.principal" attribute.
	//
	// If not set, any request principal is allowed.
	RequestPrincipals []string `protobuf:"bytes,2,rep,name=request_principals,json=requestPrincipals,proto3" json:"request_principals,omitempty"`
	// Optional. A list of negative match of request identities.
	NotRequestPrincipals []string `protobuf:"bytes,6,rep,name=not_request_principals,json=notRequestPrincipals,proto3" json:"not_request_principals,omitempty"`
	// Optional. A list of namespaces, which matches to the "source.namespace"
	// attribute. This field requires mTLS enabled.
	//
	// If not set, any namespace is allowed.
	Namespaces []string `protobuf:"bytes,3,rep,name=namespaces,proto3" json:"namespaces,omitempty"`
	// Optional. A list of negative match of namespaces.
	NotNamespaces []string `protobuf:"bytes,7,rep,name=not_namespaces,json=notNamespaces,proto3" json:"not_namespaces,omitempty"`
	// Optional. A list of IP blocks, which matches to the "source.ip" attribute.
	// Single IP (e.g. "1.2.3.4") and CIDR (e.g. "1.2.3.0/24") are supported.
	//
	// If not set, any IP is allowed.
	IpBlocks []string `protobuf:"bytes,4,rep,name=ip_blocks,json=ipBlocks,proto3" json:"ip_blocks,omitempty"`
	// Optional. A list of negative match of IP blocks.
	NotIpBlocks []string `protobuf:"bytes,8,rep,name=not_ip_blocks,json=notIpBlocks,proto3" json:"not_ip_blocks,omitempty"`
}

// Operation specifies the operations of a request. Fields in the operation are
// ANDed together.
//
// For example, the following operation matches if the host has suffix ".example.com"
// and the method is "GET" or "HEAD" and the path doesn't have prefix "/admin".
//
// ```yaml
// hosts: ["*.example.com"]
// methods: ["GET", "HEAD"]
// not_paths: ["/admin*"]
// ```
type Operation struct {
	// Optional. A list of hosts, which matches to the "request.host" attribute.
	//
	// If not set, any host is allowed. Must be used only with HTTP.
	Hosts []string `protobuf:"bytes,1,rep,name=hosts,proto3" json:"hosts,omitempty"`
	// Optional. A list of negative match of hosts.
	NotHosts []string `protobuf:"bytes,5,rep,name=not_hosts,json=notHosts,proto3" json:"not_hosts,omitempty"`
	// Optional. A list of ports, which matches to the "destination.port" attribute.
	//
	// If not set, any port is allowed.
	Ports []string `protobuf:"bytes,2,rep,name=ports,proto3" json:"ports,omitempty"`
	// Optional. A list of negative match of ports.
	NotPorts []string `protobuf:"bytes,6,rep,name=not_ports,json=notPorts,proto3" json:"not_ports,omitempty"`
	// Optional. A list of methods, which matches to the "request.method" attribute.
	// For gRPC service, this should be the fully-qualified name in the form of
	// "/package.service/method"
	//
	// If not set, any method is allowed. Must be used only with HTTP or gRPC.
	Methods []string `protobuf:"bytes,3,rep,name=methods,proto3" json:"methods,omitempty"`
	// Optional. A list of negative match of methods.
	NotMethods []string `protobuf:"bytes,7,rep,name=not_methods,json=notMethods,proto3" json:"not_methods,omitempty"`
	// Optional. A list of paths, which matches to the "request.url_path" attribute.
	//
	// If not set, any path is allowed. Must be used only with HTTP.
	Paths []string `protobuf:"bytes,4,rep,name=paths,proto3" json:"paths,omitempty"`
	// Optional. A list of negative match of paths.
	NotPaths []string `protobuf:"bytes,8,rep,name=not_paths,json=notPaths,proto3" json:"not_paths,omitempty"`
}

// From includes a list or sources.
type Rule_From struct {
	// Source specifies the source of a request.
	Source *Source `protobuf:"bytes,1,opt,name=source,proto3" json:"source,omitempty"`
}

// To includes a list or operations.
type Rule_To struct {
	// Operation specifies the operation of a request.
	Operation *Operation `protobuf:"bytes,1,opt,name=operation,proto3" json:"operation,omitempty"`
}

// Rule matches requests from a list of sources that perform a list of operations subject to a
// list of conditions. A match occurs when at least one source, operation and condition
// matches the request. An empty rule is always matched.
//
// Any string field in the rule supports Exact, Prefix, Suffix and Presence match:
//
// - Exact match: "abc" will match on value "abc".
// - Prefix match: "abc*" will match on value "abc" and "abcd".
// - Suffix match: "*abc" will match on value "abc" and "xabc".
// - Presence match: "*" will match when value is not empty.
type Rule struct {
	// Optional. from specifies the source of a request.
	//
	// If not set, any source is allowed.
	From []*Rule_From `protobuf:"bytes,1,rep,name=from,proto3" json:"from,omitempty"`
	// Optional. to specifies the operation of a request.
	//
	// If not set, any operation is allowed.
	To []*Rule_To `protobuf:"bytes,2,rep,name=to,proto3" json:"to,omitempty"`
	// Optional. when specifies a list of additional conditions of a request.
	//
	// If not set, any condition is allowed.
	When []*Condition `protobuf:"bytes,3,rep,name=when,proto3" json:"when,omitempty"`
}

// Condition specifies additional required attributes.
type Condition struct {
	// The name of an Istio attribute.
	// See the [full list of supported attributes](https://istio.io/docs/reference/config/security/conditions/).
	Key string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	// Optional. A list of allowed values for the attribute.
	// Note: at least one of values or not_values must be set.
	Values []string `protobuf:"bytes,2,rep,name=values,proto3" json:"values,omitempty"`
	// Optional. A list of negative match of values for the attribute.
	// Note: at least one of values or not_values must be set.
	NotValues []string `protobuf:"bytes,3,rep,name=not_values,json=notValues,proto3" json:"not_values,omitempty"`
}

// from istio.io/api/security/v1beta1
type WorkloadSelector struct {
	// One or more labels that indicate a specific set of pods/VMs
	// on which a policy should be applied. The scope of label search is restricted to
	// the configuration namespace in which the resource is present.
	MatchLabels map[string]string `protobuf:"bytes,1,rep,name=match_labels,json=matchLabels,proto3" json:"matchLabels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

// from istio.io/api/security/v1beta1
type AuthorizationPolicySpec struct {
	// Optional. Workload selector decides where to apply the authorization policy.
	// If not set, the authorization policy will be applied to all workloads in the
	// same namespace as the authorization policy.
	Selector *WorkloadSelector `protobuf:"bytes,1,opt,name=selector,proto3" json:"selector,omitempty"`
	// Optional. A list of rules to match the request. A match occurs when at least
	// one rule matches the request.
	//
	// If not set, the match will never occur. This is equivalent to setting a
	// default of deny for the target workloads.
	Rules []*Rule `protobuf:"bytes,2,rep,name=rules,proto3" json:"rules,omitempty"`
	// Optional. The action to take if the request is matched with the rules.
	Action string `protobuf:"varint,3,opt,name=action,proto3,enum=istio.security.v1beta1.AuthorizationPolicy_Action" json:"action,omitempty"`
}
