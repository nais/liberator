package nais_io_v1

type AccessPolicyPortRule struct {
	// Human-readable identifier for this rule.
	Name string `json:"name"`
	// The port used for communication.
	Port uint32 `json:"port"`
	// The protocol used for communication.
	// +kubebuilder:validation:Enum=HTTP;HTTPS;GRPC;HTTP2;MONGO;TCP;TLS
	Protocol string `json:"protocol"`
}

type AccessPolicyExternalRule struct {
	// The _host_ that your application should be able to reach, i.e. without the protocol (e.g. `https://`).
	Host  string                 `json:"host"`
	// List of port rules for external communication. Must be specified if using protocols other than HTTPS.
	Ports []AccessPolicyPortRule `json:"ports,omitempty"`
}

type AccessPolicyRule struct {
	// The application's name.
	Application string `json:"application"`
	// The application's namespace. May be omitted if it should be in the same namespace as your application.F
	Namespace   string `json:"namespace,omitempty"`
	// The application's cluster. May be omitted if it should be in the same cluster as your application.
	Cluster     string `json:"cluster,omitempty"`
	// Permissions contains a set of permissions that are granted to the given application.
	// Currently only applicable for Azure AD clients.
	// +nais:doc:Link="https://doc.nais.io/security/auth/azure-ad#access-policy"
	Permissions *AccessPolicyPermissions `json:"permissions,omitempty"`
}

type AccessPolicyPermissions struct {
	// Scopes is a list of Azure AD scopes that are granted to the given application.
	// +nais:doc:Link="https://doc.nais.io/security/auth/azure-ad#scopes"
	Scopes []AccessPolicyPermission `json:"scopes,omitempty"`
	// Roles is a list of Azure AD roles that are granted to the given application.
	// +nais:doc:Link="https://doc.nais.io/security/auth/azure-ad#roles"
	Roles []AccessPolicyPermission `json:"roles,omitempty"`
}

// +kubebuilder:validation:Pattern=`^[a-z0-9-_./]+$`
type AccessPolicyPermission string

type AccessPolicyInbound struct {
	// List of NAIS applications that may access your application.
	// These settings apply both to Zero Trust network connectivity and token validity for Azure AD and TokenX tokens.
	Rules []AccessPolicyRule `json:"rules"`
}

type AccessPolicyOutbound struct {
	// List of NAIS applications that your application needs to access.
	// These settings apply to Zero Trust network connectivity.
	Rules    []AccessPolicyRule         `json:"rules,omitempty"`
	// List of external resources that your applications should be able to reach.
	// +nais:doc:Availability=GCP
	External []AccessPolicyExternalRule `json:"external,omitempty"`
}

type AccessPolicy struct {
	// Configures inbound access for your application.
	Inbound  *AccessPolicyInbound  `json:"inbound,omitempty"`
	// Configures outbound access for your application.
	Outbound *AccessPolicyOutbound `json:"outbound,omitempty"`
}

func (in AccessPolicyRule) MatchesCluster(clusterName string) bool {
	if len(in.Cluster) > 0 && in.Cluster != clusterName {
		return false
	}
	return true
}
