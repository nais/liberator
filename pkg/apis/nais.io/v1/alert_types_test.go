package nais_io_v1_test

import (
	"testing"

	"github.com/nais/liberator/pkg/apis/nais.io/v1"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestApplication_Hash(t *testing.T) {
	a1, err := nais_io_v1.Alert{Spec: nais_io_v1.AlertSpec{}}.Hash()
	a2, _ := nais_io_v1.Alert{Spec: nais_io_v1.AlertSpec{}, ObjectMeta: corev1.ObjectMeta{Annotations: map[string]string{"a": "b", "team": "banana"}}}.Hash()
	a3, _ := nais_io_v1.Alert{Spec: nais_io_v1.AlertSpec{}, ObjectMeta: corev1.ObjectMeta{Labels: map[string]string{"a": "b", "team": "banana"}}}.Hash()

	assert.NoError(t, err)
	assert.Equal(t, a1, a2, "matches, as annotations is ignored")
	assert.NotEqual(t, a2, a3, "must not match ")
}

func TestNilFix(t *testing.T) {
	alert := nais_io_v1.Alert{}
	assert.Nil(t, alert.Spec.Alerts)
	alert.NilFix()
	assert.NotNil(t, alert.Spec.Receivers)
	assert.NotNil(t, alert.Spec.Alerts)
}

func TestAlertRulesValidationWithEmptyForValue(t *testing.T) {
	alert := GenerateAlertWithForValue("")
	err := alert.ValidateAlertFields()
	assert.Error(t, err)
}

func TestAlertRulesValidationWithValidForValue(t *testing.T) {
	alert := GenerateAlertWithForValue("1m")
	err := alert.ValidateAlertFields()
	assert.NoError(t, err)
}

func TestAlertRulersValidationWithInvalidForValue(t *testing.T) {
	alert := GenerateAlertWithForValue("foo")
	err := alert.ValidateAlertFields()
	assert.Error(t, err)
}

func GenerateAlertWithForValue(forValue string) nais_io_v1.Alert {
	return nais_io_v1.Alert{Spec: nais_io_v1.AlertSpec{Alerts: []nais_io_v1.Rule{
		{
			Alert:  "app is down",
			For:    forValue,
			Expr:   "kube_deployment_status_replicas_unavailable{deployment=\"my-app\"} > 0",
			Action: "kubectl describe pod -l app=my-app",
		},
	}}}
}
