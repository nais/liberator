package nais_io_v1

import (
	"github.com/nais/liberator/pkg/hash"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Machine readable event "Reason" fields, used for determining synchronization state.
const (
	EventAddedFinalizer     = "AddedFinalizer"
	EventDeletedFinalizer   = "DeletedFinalizer"
	EventCreatedInAzure     = "CreatedInAzure"
	EventUpdatedInAzure     = "UpdatedInAzure"
	EventRotatedInAzure     = "RotatedInAzure"
	EventDeletedInAzure     = "DeletedInAzure"
	EventNotInTeamNamespace = "NotInTeamNamespace"
	EventSkipped            = "Skipped"
)

// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName=azureapp
// +kubebuilder:subresource:status

// AzureAdApplication is the Schema for the AzureAdApplications API
// +kubebuilder:printcolumn:name="Client ID",type=string,JSONPath=`.status.clientId`
// +kubebuilder:printcolumn:name="Tenant",type=string,JSONPath=`.status.synchronizationTenantName`
// +kubebuilder:printcolumn:name="Tenant ID",type=string,JSONPath=`.status.synchronizationTenant`,priority=1
// +kubebuilder:printcolumn:name="Secret Ref",type=string,JSONPath=`.spec.secretName`,priority=2
// +kubebuilder:printcolumn:name="Created",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="Synchronized",type="date",JSONPath=".status.synchronizationTime"
// +kubebuilder:printcolumn:name="Assigned",type=integer,description="Number of assigned pre-authorized apps",JSONPath=`.status.preAuthorizedApps.assignedCount`
// +kubebuilder:printcolumn:name="Unassigned",type=integer,description="Number of unassigned pre-authorized apps",JSONPath=`.status.preAuthorizedApps.unassignedCount`
type AzureAdApplication struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AzureAdApplicationSpec   `json:"spec,omitempty"`
	Status AzureAdApplicationStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AzureAdApplicationList contains a list of AzureAdApplication
type AzureAdApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AzureAdApplication `json:"items"`
}

// AzureAdApplicationSpec defines the desired state of AzureAdApplication
type AzureAdApplicationSpec struct {
	// AllowAllUsers denotes whether or not all users within the tenant should be allowed to access this AzureAdApplication.
	// If undefined will default to `true` when Spec.Claims.Groups is undefined, and `false` if Spec.Claims.Groups is defined.
	AllowAllUsers *bool          `json:"allowAllUsers,omitempty"`
	Claims        *AzureAdClaims `json:"claims,omitempty"`
	// LogoutUrl is the URL where Azure AD sends a request to have the application clear the user's session data.
	// This is required if single sign-out should work correctly. Must start with 'https'
	LogoutUrl                 string                    `json:"logoutUrl,omitempty"`
	PreAuthorizedApplications []AccessPolicyInboundRule `json:"preAuthorizedApplications,omitempty"`
	ReplyUrls                 []AzureAdReplyUrl         `json:"replyUrls,omitempty"`
	// SecretName is the name of the resulting Secret resource to be created
	SecretName string `json:"secretName"`
	// SecretKeyPrefix is an optional user-defined prefix applied to the keys in the secret output, replacing the default prefix.
	SecretKeyPrefix string `json:"secretKeyPrefix,omitempty"`
	// SecretProtected protects the secret's credentials from being revoked by the janitor even when not in use.
	SecretProtected bool `json:"secretProtected,omitempty"`
	// SinglePageApplication denotes whether or not this Azure AD application should be registered as a single-page-application for usage in client-side applications without access to secrets.
	SinglePageApplication *bool `json:"singlePageApplication,omitempty"`
	// Tenant is an optional alias for targeting a tenant matching an instance of Azurerator that targets said tenant.
	// Can be omitted if only running a single instance or targeting the default tenant.
	Tenant string `json:"tenant,omitempty"`
}

// AzureAdApplicationStatus defines the observed state of AzureAdApplication
type AzureAdApplicationStatus struct {
	// CertificateKeyIds is the list of key IDs for the latest valid certificate credentials in use
	CertificateKeyIds []string `json:"certificateKeyIds,omitempty"`
	// ClientId is the Azure application client ID
	ClientId string `json:"clientId,omitempty"`
	// CorrelationId is the ID referencing the processing transaction last performed on this resource
	CorrelationId string `json:"correlationId,omitempty"`
	// ObjectId is the Azure AD Application object ID
	ObjectId string `json:"objectId,omitempty"`
	// PasswordKeyIds is the list of key IDs for the latest valid password credentials in use
	PasswordKeyIds []string `json:"passwordKeyIds,omitempty"`
	// ServicePrincipalId is the Azure applications service principal object ID
	ServicePrincipalId string `json:"servicePrincipalId,omitempty"`
	// SynchronizationHash is the hash of the AzureAdApplication object
	SynchronizationHash string `json:"synchronizationHash,omitempty"`
	// SynchronizationSecretName is the SecretName set in the last successful synchronization
	SynchronizationSecretName string `json:"synchronizationSecretName,omitempty"`
	// SynchronizationSecretRotationTime is the last time the AzureAdApplication had its keys rotated.
	SynchronizationSecretRotationTime *metav1.Time `json:"synchronizationSecretRotationTime,omitempty"`
	// SynchronizationState denotes whether the provisioning of the AzureAdApplication has been successfully completed or not
	SynchronizationState string `json:"synchronizationState,omitempty"`
	// SynchronizationTenant is the ID of the tenant that the AzureAdApplication was synchronized to.
	SynchronizationTenant string `json:"synchronizationTenant,omitempty"`
	// SynchronizationTenantName is the an alias that identifies the tenant that the AzureAdApplication was synchronized to.
	SynchronizationTenantName string `json:"synchronizationTenantName,omitempty"`
	// SynchronizationTime is the last time the Status subresource was updated
	SynchronizationTime *metav1.Time `json:"synchronizationTime,omitempty"`
	// PreAuthorizedApps contains the list of desired pre-authorized apps defined in the spec, separated by their actual status in Azure AD.
	PreAuthorizedApps *AzureAdPreAuthorizedAppsStatus `json:"preAuthorizedApps,omitempty"`
}

type AzureAdPreAuthorizedAppsStatus struct {
	// Assigned is the list of desired pre-authorized apps that have been pre-authorized to access this application.
	Assigned []AzureAdPreAuthorizedApp `json:"assigned,omitempty"`
	// AssignedCount is the size of the list in Assigned.
	AssignedCount *int `json:"assignedCount,omitempty"`
	// Unassigned is the list of desired pre-authorized apps that have _not_ been pre-authorized to access this application.
	Unassigned []AzureAdPreAuthorizedApp `json:"unassigned,omitempty"`
	// UnassignedCount is the size of the list in Unassigned.
	UnassignedCount *int `json:"unassignedCount,omitempty"`
}

type AzureAdPreAuthorizedApp struct {
	// AccessPolicyRule is the desired nais_io_v1.AccessPolicyRule matching the definition in AzureAdApplicationSpec.PreAuthorizedApplications.
	AccessPolicyRule *AccessPolicyRule `json:"accessPolicyRule,omitempty"`
	// Client ID is the actual client ID of the application found in Azure AD, if it exists.
	ClientID string `json:"clientId,omitempty"`
	// Object ID is the actual object ID of the service principal belonging to the application found in Azure AD, if it exists.
	ServicePrincipalObjectID string `json:"servicePrincipalObjectId,omitempty"`
	// Reason is a human-readable message that provides detailed information about the application and its status.
	Reason string `json:"reason,omitempty"`
}

// Claims defines additional configuration of the emitted claims in tokens returned to the Azure AD application.
type AzureAdClaims struct {
	// Extra is a list of additional claims to be mapped from an associated claim-mapping policy.
	// Currently, the only supported values are `NAVident` and `azp_name`.
	// +nais:doc:Link="https://doc.nais.io/security/auth/azure-ad/configuration#extra"
	Extra []AzureAdExtraClaim `json:"extra,omitempty"`
	// Groups is a list of Azure AD group IDs to be emitted in the 'Groups' claim.
	// This also restricts access to only contain users of the defined groups unless overridden by Spec.AllowAllUsers.
	// +nais:doc:Link="https://doc.nais.io/security/auth/azure-ad/access-policy#groups"
	// +nais:doc:Link="https://doc.nais.io/security/auth/azure-ad/access-policy#users"
	Groups []AzureAdGroup `json:"groups,omitempty"`
}

// +kubebuilder:validation:Enum=NAVident;azp_name
type AzureAdExtraClaim string

type AzureAdGroup struct {
	// ID is the actual `object ID` associated with the given group in Azure AD.
	ID string `json:"id,omitempty"`
}

// AzureAdReplyUrl defines the valid reply URLs for callbacks after OIDC flows for this application
type AzureAdReplyUrl struct {
	// +kubebuilder:validation:Pattern=`^https?:\/\/.+$`
	Url string `json:"url,omitempty"`
}

func (in *AzureAdApplication) GetObjectId() string {
	return in.Status.ObjectId
}

func (in *AzureAdApplication) GetServicePrincipalId() string {
	return in.Status.ServicePrincipalId
}

func (in *AzureAdApplication) GetClientId() string {
	return in.Status.ClientId
}

func (in *AzureAdApplication) Hash() (string, error) {
	relevantValues := struct {
		ReplyUrls                 []AzureAdReplyUrl
		PreAuthorizedApplications []AccessPolicyInboundRule
		LogoutUrl                 string
		Tenant                    string
		Claims                    *AzureAdClaims
		SecretKeyPrefix           string
		SinglePageApplication     *bool `json:"singlePageApplication,omitempty"`
		AllowAllUsers             *bool `json:"allowAllUsers,omitempty"`
	}{
		ReplyUrls:                 in.Spec.ReplyUrls,
		PreAuthorizedApplications: in.Spec.PreAuthorizedApplications,
		LogoutUrl:                 in.Spec.LogoutUrl,
		Tenant:                    in.Spec.Tenant,
		Claims:                    in.Spec.Claims,
		SecretKeyPrefix:           in.Spec.SecretKeyPrefix,
		SinglePageApplication:     in.Spec.SinglePageApplication,
		AllowAllUsers:             in.Spec.AllowAllUsers,
	}
	return hash.Hash(relevantValues)
}

func init() {
	SchemeBuilder.Register(&AzureAdApplication{}, &AzureAdApplicationList{})
}
