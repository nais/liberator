package nais_io_v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ExampleAlertForDocumentation() *Alert {
	boolp := func(i bool) *bool {
		return &i
	}
	return &Alert{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Alert",
			APIVersion: "nais.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "myalert",
			Namespace: "myteam",
			Labels: map[string]string{
				"team": "myteam",
			},
		},
		Spec: AlertSpec{
			Route: Route{
				GroupWait:      "30s",
				GroupInterval:  "5m",
				RepeatInterval: "3h",
				GroupBy:        []string{"<label_name>"},
			},
			Receivers: Receivers{
				Slack: Slack{
					Channel:      "#alert-channel",
					PrependText:  "Oh noes!",
					SendResolved: boolp(true),
					Username:     "Alertmanager",
					IconUrl:      "http://lorempixel.com/48/48",
					IconEmoji:    ":chart_with_upwards_trend:",
				},
				Email: Email{
					To:           "myteam@nav.no",
					SendResolved: false,
				},
				SMS: SMS{
					Recipients:   "12345678",
					SendResolved: boolp(false),
				},
			},
			Alerts: []Rule{
				{
					Alert:         "applikasjon nede",
					Description:   "App {{ $labels.app }} er nede i namespace {{ $labels.kubernetes_namespace }}",
					Expr:          "kube_deployment_status_replicas_available{deployment=\"<appname>\"} > 0",
					For:           "2m",
					Action:        "kubectl describe pod {{ $labels.kubernetes_pod_name }} -n {{ $labels.kubernetes_namespace }}` for events, og `kubectl logs {{ $labels.kubernetes_pod_name }} -n {{ $labels.kubernetes_namespace }}` for logger",
					Documentation: "https://doc.nais.io/observability/alerts/",
					SLA:           "Mellom 8 og 16",
					Severity:      "danger",
				},
			},
			InhibitRules: []InhibitRules{
				{
					Targets: map[string]string{
						"key": "value",
					},
					TargetsRegex: map[string]string{
						"key": "value(.)+",
					},
					Sources: map[string]string{
						"key": "value",
					},
					SourcesRegex: map[string]string{
						"key": "value(.)?",
					},
					Labels: []string{
						"label",
						"lebal",
					},
				},
			},
		},
	}
}
