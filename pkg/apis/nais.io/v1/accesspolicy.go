package nais_io_v1

type AccessPolicyPortRule struct {
	Name string `json:"name"`
	Port uint32 `json:"port"`
	// +kubebuilder:validation:Enum=HTTP;HTTPS;GRPC;HTTP2;MONGO;TCP;TLS
	Protocol string `json:"protocol"`
}

type AccessPolicyExternalRule struct {
	Host  string                 `json:"host"`
	Ports []AccessPolicyPortRule `json:"ports,omitempty"`
}

type AccessPolicyRule struct {
	Application string `json:"application"`
	Namespace   string `json:"namespace,omitempty"`
	Cluster     string `json:"cluster,omitempty"`
}

type AccessPolicyInbound struct {
	// List of NAIS applications that may access your application.
	// These settings apply both to Zero Trust network connectivity and token validity for Azure AD and TokenX tokens.
	Rules []AccessPolicyRule `json:"rules"`
}

type AccessPolicyOutbound struct {
	// List of NAIS applications that your application needs to access.
	// These settings apply to Zero Trust network connectivity.
	Rules    []AccessPolicyRule         `json:"rules,omitempty"`
	// +nais:doc:Availability=GCP
	External []AccessPolicyExternalRule `json:"external,omitempty"`
}

type AccessPolicy struct {
	Inbound  *AccessPolicyInbound  `json:"inbound,omitempty"`
	Outbound *AccessPolicyOutbound `json:"outbound,omitempty"`
}

func (in AccessPolicyRule) MatchesCluster(clusterName string) bool {
	if len(in.Cluster) > 0 && in.Cluster != clusterName {
		return false
	}
	return true
}
