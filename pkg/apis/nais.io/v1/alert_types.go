package nais_io_v1

import (
	"strconv"
	"time"

	hash "github.com/mitchellh/hashstructure"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const LastSyncedHashAnnotation = "nais.io/lastSyncedHash"

func init() {
	SchemeBuilder.Register(
		&Alert{},
		&AlertList{},
	)
}

type Slack struct {
	Channel     string `json:"channel"`
	PrependText string `json:"prependText,omitempty"`
	SendResolved bool  `json:"send_resolved,omitempty"`
}

type Email struct {
	To           string `json:"to"`
	SendResolved bool   `json:"send_resolved,omitempty"`
}

type SMS struct {
	Recipients   string `json:"recipients"`
	SendResolved bool   `json:"send_resolved,omitempty"`
}

type Pushover struct {
	UserKey      string `json:"user_key"`
	SendResolved bool   `json:"send_resolved,omitempty"`
}

type Receivers struct {
	Slack    Slack    `json:"slack,omitempty"`
	Email    Email    `json:"email,omitempty"`
	SMS      SMS      `json:"sms,omitempty"`
	Pushover Pushover `json:"pushover,omitempty"`
}

type Rule struct {
	// +kubebuilder:validation:Required
	Alert         string `json:"alert"`
	Description   string `json:"description,omitempty"`
	// +kubebuilder:validation:Required
	Expr          string `json:"expr"`
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern="^\\d+[smhdwy]$"
	For           string `json:"for"`
	// +kubebuilder:validation:Required
	Action        string `json:"action"`
	Documentation string `json:"documentation,omitempty"`
	SLA           string `json:"sla,omitempty"`
	// +kubebuilder:validation:Pattern="^$|good|warning|danger|#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})"
	Severity      string `json:"severity,omitempty"`
	Priority      string `json:"priority,omitempty"`
}

type InhibitRules struct {
	Targets      map[string]string `json:"targets,omitempty"`
	TargetsRegex map[string]string `json:"targetsRegex,omitempty"`
	Sources      map[string]string `json:"sources,omitempty"`
	SourcesRegex map[string]string `json:"sourcesRegex,omitempty"`
	Labels       []string          `json:"labels,omitempty"`
}

type Route struct {
	// +kubebuilder:validation:Pattern="([0-9]+(ms|[smhdwy]))?"
	GroupWait      string `json:"groupWait,omitempty"`
	// +kubebuilder:validation:Pattern="([0-9]+(ms|[smhdwy]))?"
	GroupInterval  string `json:"groupInterval,omitempty"`
	// +kubebuilder:validation:Pattern="([0-9]+(ms|[smhdwy]))?"
	RepeatInterval string `json:"repeatInterval,omitempty"`
}

type AlertSpec struct {
	Route        Route          `json:"route,omitempty"`
	// +kubebuilder:validation:Required
	Receivers    Receivers      `json:"receivers,omitempty"`
	// +kubebuilder:validation:Required
	Alerts       []Rule         `json:"alerts,omitempty"`
	InhibitRules []InhibitRules `json:"inhibitRules,omitempty"`
}

// AlertStatus defines the observed state of Alerterator
type AlertStatus struct {
	SynchronizationTime  int64  `json:"synchronizationTime,omitempty"`
	SynchronizationState string `json:"synchronizationState,omitempty"`
	SynchronizationHash  string `json:"synchronizationHash,omitempty"`
}

// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Slack channel",type="string",JSONPath=".spec.receivers.slack.channel"
// +kubebuilder:object:root=true
type Alert struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AlertSpec   `json:"spec"`
	Status AlertStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type AlertList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Alert `json:"items"`
}

func (in *Alert) CreateEvent(reason, message, typeStr string) *corev1.Event {
	return &corev1.Event{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "alerterator-event",
			Namespace:    in.Namespace,
		},
		ReportingController: "alerterator",
		ReportingInstance:   "alerterator",
		Action:              reason,
		Reason:              reason,
		InvolvedObject:      in.GetObjectReference(),
		Source:              corev1.EventSource{Component: "alerterator"},
		Message:             message,
		EventTime:           metav1.MicroTime{Time: time.Now()},
		FirstTimestamp:      metav1.Time{Time: time.Now()},
		LastTimestamp:       metav1.Time{Time: time.Now()},
		Type:                typeStr,
	}
}

func (in *Alert) GetObjectKind() schema.ObjectKind {
	return in
}

func (in *Alert) GetObjectReference() corev1.ObjectReference {
	return corev1.ObjectReference{
		APIVersion:      "v1",
		UID:             in.UID,
		Name:            in.Name,
		Kind:            "Alert",
		ResourceVersion: in.ResourceVersion,
		Namespace:       in.Namespace,
	}
}

func (in *Alert) GetOwnerReference() metav1.OwnerReference {
	return metav1.OwnerReference{
		APIVersion: "v1",
		Kind:       "Alert",
		Name:       in.Name,
		UID:        in.UID,
	}
}

// NilFix initializes all slices from their nil defaults.
//
// This is done in order to workaround the k8s client serializer
// which crashes when these fields are uninitialized.
func (in *Alert) NilFix() {
	if in.Spec.Alerts == nil {
		in.Spec.Alerts = make([]Rule, 0)
	}
	if in.Spec.InhibitRules == nil {
		in.Spec.InhibitRules = make([]InhibitRules, 0)
	}
}

func (in Alert) Hash() (string, error) {
	// struct including the relevant fields for
	// creating a hash of an Application object
	relevantValues := struct {
		Spec   AlertSpec
		Labels map[string]string
	}{
		in.Spec,
		in.Labels,
	}

	h, err := hash.Hash(relevantValues, nil)
	return strconv.FormatUint(h, 10), err
}

func (in *Alert) LastSyncedHash() string {
	a := in.GetAnnotations()
	if a == nil {
		a = make(map[string]string)
	}
	return a[LastSyncedHashAnnotation]
}

func (in *Alert) SetLastSyncedHash(hash string) {
	a := in.GetAnnotations()
	if a == nil {
		a = make(map[string]string)
	}
	a[LastSyncedHashAnnotation] = hash
	in.SetAnnotations(a)
}
