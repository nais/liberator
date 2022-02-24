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
	// The channel or user to send notifications to.
	// Can be specified with and without `#`.
	Channel string `json:"channel"`
	// Text to prepend every Slack message with severity `danger`.
	PrependText string `json:"prependText,omitempty"`
	// Whether or not to notify about resolved alerts.
	// +nais:doc:Default="true"
	SendResolved *bool `json:"send_resolved,omitempty"`
	// Set your bot's user name.
	Username string `json:"username,omitempty"`
	// URL to an image to use as the icon for this message
	IconUrl string `json:"icon_url,omitempty"`
	// Emoji to use as the icon for this message
	IconEmoji string `json:"icon_emoji,omitempty"`
}

type Email struct {
	To string `json:"to"`
	// Whether or not to notify about resolved alerts.
	// +nais:doc:Default="false"
	SendResolved bool `json:"send_resolved,omitempty"`
}

type SMS struct {
	Recipients string `json:"recipients"`
	// Whether or not to notify about resolved alerts.
	// +nais:doc:Default="true"
	SendResolved *bool `json:"send_resolved,omitempty"`
}

type TLSConfig struct {
	// Disable validation of the server certificate.
	// +nais:doc:Default="false"
	InsecureSkipVerify bool `json:"insecure_skip_verify,omitempty"`
}

type HttpConfig struct {
	// Optional proxy URL.
	ProxyUrl string `json:"proxy_url,omitempty"`
	// Configures the TLS settings.
	TLSConfig TLSConfig `json:"tls_config,omitempty"`
}

type Webhook struct {
	// The endpoint to send HTTP POST requests to.
	URL string `json:"url"`
	// The maximum number of alerts to include in a single webhook message. Alerts
	// above this threshold are truncated. When leaving this at its default value of
	// 0, all alerts are included.
	// +nais:doc:Default="0"
	MaxAlerts int `json:"max_alerts"`
	// Whether or not to notify about resolved alerts.
	// +nais:doc:Default="true"
	SendResolved *bool `json:"send_resolved,omitempty"`
	// A http_config allows configuring the HTTP client that the receiver uses to communicate with HTTP-based API services.
	HttpConfig HttpConfig `json:"http_config,omitempty"`
}

type Receivers struct {
	// Slack notifications are sent via Slack webhooks.
	Slack Slack `json:"slack,omitempty"`
	// Alerts via e-mails
	Email Email `json:"email,omitempty"`
	// Alerts via SMS
	SMS SMS `json:"sms,omitempty"`
	// Alerts via custom web application
	Webhook Webhook `json:"webhook,omitempty"`
}

type Rule struct {
	// The name of the alert.
	// +kubebuilder:validation:Required
	Alert string `json:"alert"`
	// Simple description of the triggered alert.
	Description string `json:"description,omitempty"`
	// Prometheus expression that triggers an alert.
	// Explore expressions in the [Prometheus](https://docs.nais.io/observability/alerts/#writing-the-expr)-interface
	// +kubebuilder:validation:Required
	Expr string `json:"expr"`
	// Duration before the alert should trigger.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern="^\\d+[smhdwy]$"
	For string `json:"for"`
	// What human actions are needed to resolve or investigate this alert.
	// +kubebuilder:validation:Required
	Action string `json:"action"`
	// URL for documentation for this alert.
	Documentation string `json:"documentation,omitempty"`
	// Time before a human should resolve the alert.
	SLA string `json:"sla,omitempty"`
	// Alert level for Slack messages.
	// +nais:doc:Default="danger"
	// +kubebuilder:validation:Pattern="^$|good|warning|danger|#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})"
	Severity string `json:"severity,omitempty"`
}

type InhibitRules struct {
	// Matchers that have to be fulfilled in the alerts to be muted.
	// These are key/value pairs.
	Targets map[string]string `json:"targets,omitempty"`
	// Regex matchers that have to be fulfilled in the alerts to be muted.
	// These are key/value pairs, where the value can be a regex.
	TargetsRegex map[string]string `json:"targetsRegex,omitempty"`
	// Matchers for which one or more alerts have to exist for the inhibition to take effect.
	// These are key/value pairs.
	Sources map[string]string `json:"sources,omitempty"`
	// Regex matchers for which one or more alerts have to exist for the inhibition to take effect.
	// These are key/value pairs, where the value can be a regex.
	SourcesRegex map[string]string `json:"sourcesRegex,omitempty"`
	// Labels that must have an equal value in the source and target alert for the inhibition to take effect.
	Labels []LabelName `json:"labels,omitempty"`
}

type Route struct {
	// How long to initially wait to send a notification for a group of alerts.
	// Allows to wait for an inhibiting alert to arrive or collect more initial alerts for the same group.
	// +nais:doc:Default="10s"
	// +kubebuilder:validation:Pattern="((([0-9]+)y)?(([0-9]+)w)?(([0-9]+)d)?(([0-9]+)h)?(([0-9]+)m)?(([0-9]+)s)?(([0-9]+)ms)?|0)"
	GroupWait string `json:"groupWait,omitempty"`
	// How long to wait before sending a notification about new alerts that are added to a group of alerts for which an initial notification has already been sent.
	// +nais:doc:Default="5m"
	// +kubebuilder:validation:Pattern="((([0-9]+)y)?(([0-9]+)w)?(([0-9]+)d)?(([0-9]+)h)?(([0-9]+)m)?(([0-9]+)s)?(([0-9]+)ms)?|0)"
	GroupInterval string `json:"groupInterval,omitempty"`
	// How long to wait before sending a notification again if it has already been sent successfully for an alert.
	// +nais:doc:Default="1h"
	// +kubebuilder:validation:Pattern="((([0-9]+)y)?(([0-9]+)w)?(([0-9]+)d)?(([0-9]+)h)?(([0-9]+)m)?(([0-9]+)s)?(([0-9]+)ms)?|0)"
	RepeatInterval string `json:"repeatInterval,omitempty"`
	// The labels by which incoming alerts are grouped together.
	GroupBy []LabelName `json:"group_by,omitempty"`
}

// +kubebuilder:validation:Pattern="[a-zA-Z_][a-zA-Z0-9_]*"
type LabelName string

type AlertSpec struct {
	Route Route `json:"route,omitempty"`
	// A list of notification recievers. You can use one or more of: e-mail, slack, sms.
	// There needs to be at least one receiver.
	// +kubebuilder:validation:Required
	Receivers Receivers `json:"receivers,omitempty"`
	// +kubebuilder:validation:Required
	Alerts []Rule `json:"alerts,omitempty"`
	// A list of inhibit rules.
	// An inhibition rule mutes an alert (target) matching a set of matchers when an alert (source) exists that matches another set of matchers.
	// Both target and source alerts must have the same label values for the label names in the labels list.
	// +nais:doc:Link="https://prometheus.io/docs/alerting/latest/configuration/#inhibit_rule"
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
