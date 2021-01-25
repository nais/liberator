package nais_io_v1

import (
	"regexp"
	"strconv"
	"time"

	hash "github.com/mitchellh/hashstructure"
	"gopkg.in/go-playground/validator.v9"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const LastSyncedHashAnnotation = "nais.io/lastSyncedHash"
const PrometheusRangeRegexPattern = "\\d+[smhdwy]"

var (
	Validate             *validator.Validate
	PrometheusRangeRegex = regexp.MustCompile(PrometheusRangeRegexPattern)
)

// +kubebuilder:object:root=true
type Alert struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec AlertSpec `json:"spec"`
}

// +kubebuilder:object:root=true
type AlertList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Alert `json:"items"`
}

type Slack struct {
	Channel     string `json:"channel"`
	PrependText string `json:"prependText"`
}

type Email struct {
	To           string `json:"to"`
	SendResolved bool   `json:"send_resolved"`
}

type SMS struct {
	Recipients   string `json:"recipients"`
	SendResolved bool   `json:"send_resolved"`
}

type Pushover struct {
	UserKey      string `json:"user_key"`
	SendResolved bool   `json:"send_resolved"`
}

type Receivers struct {
	Slack    Slack    `json:"slack"`
	Email    Email    `json:"email"`
	SMS      SMS      `json:"sms"`
	Pushover Pushover `json:"pushover"`
}

type Rule struct {
	Alert         string `json:"alert"`
	Description   string `json:"description"`
	Expr          string `json:"expr"`
	For           string `json:"for" validate:"required,prometheus_range"`
	Action        string `json:"action"`
	Documentation string `json:"documentation"`
	SLA           string `json:"sla"`
	Severity      string `json:"severity"`
	Priority      string `json:"priority"`
}

type InhibitRules struct {
	Targets      map[string]string `json:"targets"`
	TargetsRegex map[string]string `json:"targetsRegex"`
	Sources      map[string]string `json:"sources"`
	SourcesRegex map[string]string `json:"sourcesRegex"`
	Labels       []string          `json:"labels"`
}

type Route struct {
	GroupWait      string `json:"groupWait"`
	GroupInterval  string `json:"groupInterval"`
	RepeatInterval string `json:"repeatInterval"`
}

type AlertSpec struct {
	Route        Route          `json:"route"`
	Receivers    Receivers      `json:"receivers"`
	Alerts       []Rule         `json:"alerts" validate:"dive"`
	InhibitRules []InhibitRules `json:"inhibitRules"`
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

func (in *Alert) ValidateAlertFields() error {
	Validate = validator.New()
	_ = Validate.RegisterValidation("prometheus_range", IsValidPrometheusRange)
	return Validate.Struct(in)
}

func IsValidPrometheusRange(fl validator.FieldLevel) bool {
	return PrometheusRangeRegex.MatchString(fl.Field().String())
}
