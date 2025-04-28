package nais_io_v1

type AccessPolicyPortRule struct {
	// The port used for communication.
	Port uint32 `json:"port"`
}

type AccessPolicyExternalRule struct {
	// The _host_ that your application should be able to reach, i.e. without the protocol (e.g. `https://`). "Host" and "IPv4" are mutually exclusive
	// +kubebuilder:validation:Pattern=`^([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])(\.([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]{0,61}[a-zA-Z0-9]))*$`
	Host string `json:"host,omitempty"`
	// The IPv4 address that your application should be able to reach. "IPv4" and "Host" are mutually exclusive
	// +kubebuilder:validation:Pattern=`^(([0-9])|([1-9][0-9])|(1([0-9]{2}))|(2[0-4][0-9])|(25[0-5]))((\.(([0-9])|([1-9][0-9])|(1([0-9]{2}))|(2[0-4][0-9])|(25[0-5]))){3})$`
	IPv4 string `json:"ipv4,omitempty"`
	// List of port rules for external communication. Must be specified if using protocols other than HTTPS.
	Ports []AccessPolicyPortRule `json:"ports,omitempty"`
}

type AccessPolicyRule struct {
	// The application's name.
	Application string `json:"application"`
	// The application's namespace. May be omitted if it should be in the same namespace as your application.
	Namespace string `json:"namespace,omitempty"`
	// The application's cluster. May be omitted if it should be in the same cluster as your application.
	Cluster string `json:"cluster,omitempty"`
}

// +k8s:deepcopy-gen=false
type AccessPolicyBaseRules interface {
	GetRules() []AccessPolicyRule
}

type AccessPolicyRules []AccessPolicyRule

func (in AccessPolicyRules) GetRules() []AccessPolicyRule {
	return in
}

type AccessPolicyInboundRules []AccessPolicyInboundRule

func (in AccessPolicyInboundRules) GetRules() []AccessPolicyRule {
	rules := make([]AccessPolicyRule, len(in))

	for i, rule := range in {
		rules[i] = rule.AccessPolicyRule
	}

	return rules
}

type AccessPolicyInboundRule struct {
	AccessPolicyRule `json:",inline"`
	// Permissions contains a set of permissions that are granted to the given application.
	// Currently only applicable for Azure AD clients.
	// +nais:doc:Tenants="nav"
	Permissions *AccessPolicyPermissions `json:"permissions,omitempty"`
}

type AccessPolicyPermissions struct {
	// Scopes is a set of custom permission scopes that are granted to a given application.
	Scopes []AccessPolicyPermission `json:"scopes,omitempty"`
	// Roles is a set of custom permission roles that are granted to a given application.
	Roles []AccessPolicyPermission `json:"roles,omitempty"`
}

// +kubebuilder:validation:Pattern=`^[a-z0-9-_./]+$`
type AccessPolicyPermission string

type AccessPolicyInbound struct {
	// List of Nais applications that may access your application.
	// These settings apply both to Zero Trust network connectivity and token validity for Azure AD and TokenX tokens.
	Rules AccessPolicyInboundRules `json:"rules"`
}

type AccessPolicyOutbound struct {
	// List of Nais applications that your application needs to access.
	// These settings apply to Zero Trust network connectivity.
	Rules AccessPolicyRules `json:"rules,omitempty"`
	// List of external resources that your applications should be able to reach.
	// +nais:doc:Availability=GCP
	External []AccessPolicyExternalRule `json:"external,omitempty"`
}

type AccessPolicy struct {
	// Configures inbound access for your application.
	Inbound *AccessPolicyInbound `json:"inbound,omitempty"`
	// Configures outbound access for your application.
	Outbound *AccessPolicyOutbound `json:"outbound,omitempty"`
}

func (in AccessPolicyRule) MatchesCluster(clusterName string) bool {
	if len(in.Cluster) > 0 && in.Cluster != clusterName {
		return false
	}
	return true
}
