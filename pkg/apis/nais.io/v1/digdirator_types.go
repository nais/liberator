package nais_io_v1

import (
	"github.com/nais/liberator/pkg/hash"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DigdiratorStatus defines the observed state of Current Client
type DigdiratorStatus struct {
	// SynchronizationState denotes the last known state of the Instance during synchronization
	SynchronizationState string `json:"synchronizationState,omitempty"`
	// SynchronizationTime is the last time the Status subresource was updated
	SynchronizationTime *metav1.Time `json:"synchronizationTime,omitempty"`
	// SynchronizationHash is the hash of the Instance object
	SynchronizationHash string `json:"synchronizationHash,omitempty"`
	// SynchronizationSecretName is the SecretName set in the last successful synchronization
	SynchronizationSecretName string `json:"synchronizationSecretName,omitempty"`
	// ClientID is the corresponding client ID for this client at Digdir
	ClientID string `json:"clientID,omitempty"`
	// CorrelationID is the ID referencing the processing transaction last performed on this resource
	CorrelationID string `json:"correlationID,omitempty"`
	// KeyIDs is the list of key IDs for valid JWKs registered for the client at Digdir
	KeyIDs []string `json:"keyIDs,omitempty"`
}

func (in *DigdiratorStatus) GetSynchronizationHash() string {
	return in.SynchronizationHash
}

func (in *DigdiratorStatus) SetHash(hash string) {
	in.SynchronizationHash = hash
}

func (in *DigdiratorStatus) SetStateSynchronized() {
	now := metav1.Now()
	in.SynchronizationTime = &now
	in.SynchronizationState = EventSynchronized
}

func (in *DigdiratorStatus) GetClientID() string {
	return in.ClientID
}

func (in *DigdiratorStatus) SetClientID(clientID string) {
	in.ClientID = clientID
}

func (in *DigdiratorStatus) SetCorrelationID(correlationID string) {
	in.CorrelationID = correlationID
}

func (in *DigdiratorStatus) GetKeyIDs() []string {
	return in.KeyIDs
}

func (in *DigdiratorStatus) SetKeyIDs(keyIDs []string) {
	in.KeyIDs = keyIDs
}

func (in *DigdiratorStatus) SetSynchronizationState(state string) {
	in.SynchronizationState = state
}

func (in *DigdiratorStatus) GetSynchronizationSecretName() string {
	return in.SynchronizationSecretName
}

func (in *DigdiratorStatus) SetSynchronizationSecretName(name string) {
	in.SynchronizationSecretName = name
}

func init() {
	SchemeBuilder.Register(
		&MaskinportenClient{},
		&MaskinportenClientList{},
	)
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName=maskinportenclient

// +kubebuilder:printcolumn:name="Secret Ref",type=string,JSONPath=`.spec.secretName`
// +kubebuilder:printcolumn:name="ClientID",type=string,JSONPath=`.status.clientID`
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="Created",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="Synchronized",type="date",JSONPath=".status.synchronizationTime"

// MaskinportenClient is the Schema for the MaskinportenClient API
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MaskinportenClient struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MaskinportenClientSpec `json:"spec,omitempty"`
	Status DigdiratorStatus       `json:"status,omitempty"`
}

// MaskinportenClientSpec defines the desired state of MaskinportenClient
type MaskinportenClientSpec struct {
	// Scopes is a object of used end exposed scopes by application
	Scopes MaskinportenScope `json:"scopes,omitempty"`
	// SecretName is the name of the resulting Secret resource to be created
	SecretName string `json:"secretName"`
}

// MaskinportenClientList contains a list of MaskinportenClient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +nais:doc:Availability="team namespaces"
type MaskinportenClientList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MaskinportenClient `json:"items"`
}

type MaskinportenScope struct {
	// This is the Schema for the consumes and exposes API.
	// `consumes` is a list of scopes that your client can request access to.
	ConsumedScopes []ConsumedScope `json:"consumes,omitempty"`
	// `exposes` is a list of scopes your application want to expose to other organization where access to the scope is based on organization number.
	ExposedScopes []ExposedScope `json:"exposes,omitempty"`
}

type ConsumedScope struct {
	// The scope consumed by the application to gain access to an external organization API.
	// Ensure that the NAV organization has been granted access to the scope prior to requesting access.
	// +nais:doc:Link="https://doc.nais.io/security/auth/maskinporten/#consume-scopes"
	// +kubebuilder:validation:Required
	Name string `json:"name"`
}

type ExposedScope struct {
	// If Enabled the configured scope is available to be used and consumed by organizations granted access.
	// +nais:doc:Link="https://doc.nais.io/naisjob/reference/#maskinportenscopesexposesconsumers"
	// +kubebuilder:validation:Required
	Enabled bool `json:"enabled"`
	// The actual subscope combined with `Product`.
	// Ensure that `<Product><Name>` matches `Pattern`.
	// +nais:doc:Default="false"
	// +kubebuilder:validation:Pattern=`^([a-zæøå0-9]+\/?)+(\:[a-zæøå0-9]+)*[a-zæøå0-9]+(\.[a-zæøå0-9]+)*$`
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// The product-area your application belongs to e.g. arbeid, helse ...
	// This will be included in the final scope `nav:<Product><Name>`.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^[a-z0-9]+$`
	Product string `json:"product"`
	// Max time in seconds for a issued access_token.
	// Default is `30` sec.
	// +nais:doc:Default="30"
	// +kubebuilder:validation:Minimum=30
	// +kubebuilder:validation:Maximum=680
	AtMaxAge *int `json:"atMaxAge,omitempty"`
	// Whitelisting of integration's allowed.
	// Default is `maskinporten`
	// +nais:doc:Default="maskinporten"
	// +nais:doc:Link="https://docs.digdir.no/maskinporten_guide_apitilbyder.html#scope-begrensninger"
	// +kubebuilder:validation:MinItems=1
	AllowedIntegrations []string `json:"allowedIntegrations,omitempty"`
	// External consumers granted access to this scope and able to request access_token.
	Consumers []ExposedScopeConsumer `json:"consumers,omitempty"`
}

type ExposedScopeConsumer struct {
	// The external business/organization number.
	// +kubebuilder:validation:Pattern=`^\d{9}$`
	Orgno string `json:"orgno"`
	// This is a describing field intended for clarity not used for any other purpose.
	// +optional
	Name string `json:"name,omitempty"`
}

func (in *MaskinportenClient) Hash() (string, error) {
	return hash.Hash(in.Spec)
}

func (in *MaskinportenClient) GetStatus() *DigdiratorStatus {
	return &in.Status
}

func (in *MaskinportenClient) SetStatus(new DigdiratorStatus) {
	in.Status = new
}

func (in *MaskinportenClient) GetConsumedScopes() []string {
	scopes := make([]string, 0)
	for _, scope := range in.Spec.Scopes.ConsumedScopes {
		scopes = append(scopes, scope.Name)
	}
	return scopes
}

func (in *MaskinportenClient) GetExposedScopes() map[string]ExposedScope {
	scopes := make(map[string]ExposedScope)
	for _, scope := range in.Spec.Scopes.ExposedScopes {
		scopes[scope.Name] = scope
	}
	return scopes
}

func init() {
	SchemeBuilder.Register(
		&IDPortenClient{},
		&IDPortenClientList{},
	)
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName=idportenclient

// +kubebuilder:printcolumn:name="Secret Ref",type=string,JSONPath=`.spec.secretName`
// +kubebuilder:printcolumn:name="ClientID",type=string,JSONPath=`.status.clientID`
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="Created",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="Synchronized",type="date",JSONPath=".status.synchronizationTime"

// IDPortenClient is the Schema for the IDPortenClients API
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type IDPortenClient struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IDPortenClientSpec `json:"spec,omitempty"`
	Status DigdiratorStatus   `json:"status,omitempty"`
}

// IDPortenClientList contains a list of IDPortenClient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type IDPortenClientList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IDPortenClient `json:"items"`
}

// IDPortenClientSpec defines the desired state of IDPortenClient
type IDPortenClientSpec struct {
	// AccessTokenLifetime is the maximum lifetime in seconds for the returned access_token from ID-porten.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=3600
	AccessTokenLifetime *int `json:"accessTokenLifetime,omitempty"`
	// ClientURI is the URL to the client to be used at DigDir when displaying a 'back' button or on errors
	ClientURI IDPortenURI `json:"clientURI,omitempty"`
	// IntegrationType is used to make sensible choices for your client.
	// Which type of integration you choose will provide guidance on which scopes you can use with the client.
	// A client can only have one integration type.
	//
	// NB! It is not possible to change the integration type after creation.
	//
	// +nais:doc:Immutable=true
	// +nais:doc:Default=idporten
	// +nais:doc:Link="https://docs.digdir.no/oidc_protocol_scope.html#scope-limitations"
	// +nais:doc:Link="https://docs.digdir.no/oidc_func_clientreg.html"
	// +kubebuilder:validation:Enum=krr;idporten;api_klient
	IntegrationType string `json:"IntegrationType,omitempty" nais:"immutable"`
	// FrontchannelLogoutURI is the URL that ID-porten sends a requests to whenever a logout is triggered by another application using the same session
	FrontchannelLogoutURI IDPortenURI `json:"frontchannelLogoutURI,omitempty"`
	// PostLogoutRedirectURI is a list of valid URIs that ID-porten may redirect to after logout
	PostLogoutRedirectURIs []IDPortenURI `json:"postLogoutRedirectURIs,omitempty"`
	// RedirectURI is the redirect URI to be registered at DigDir
	RedirectURI IDPortenURI `json:"redirectURI"`
	// SecretName is the name of the resulting Secret resource to be created
	SecretName string `json:"secretName"`
	// Register different oauth2 Scopes on your client.
	// You will not be able to add a scope to your client that conflicts with the client's IntegrationType.
	// For example, you can not add a scope that is limited to the IntegrationType `krr` of integrationType `idporten`, and vice versa.
	//
	// Default for IntegrationType `krr` = ("krr:global/kontaktinformasjon.read", "krr:global/digitalpost.read")
	// Default for IntegrationType `idporten` = ("openid", "profile")
	// IntegrationType `api_klient` have no Default, checkout Digdir documentation.
	//
	// +nais:doc:Link="https://docs.digdir.no/oidc_func_clientreg.html?h=api_klient#scopes"
	Scopes []string `json:"scopes,omitempty"`
	// SessionLifetime is the maximum session lifetime in seconds for a logged in end-user for this client.
	// +kubebuilder:validation:Minimum=3600
	// +kubebuilder:validation:Maximum=7200
	SessionLifetime *int `json:"sessionLifetime,omitempty"`
}

// +kubebuilder:validation:Pattern=`^https:\/\/.+$`
type IDPortenURI string

func (in *IDPortenClient) Hash() (string, error) {
	return hash.Hash(in.Spec)
}

func (in *IDPortenClient) GetStatus() *DigdiratorStatus {
	return &in.Status
}

func (in *IDPortenClient) SetStatus(new DigdiratorStatus) {
	in.Status = new
}
