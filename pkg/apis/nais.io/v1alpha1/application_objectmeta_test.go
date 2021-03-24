package nais_io_v1alpha1

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reflect"
	"testing"
)

func TestApplication_CreateObjectMeta(t *testing.T) {
	const app, namespace, team, key, value = "myapp", "mynamespace", "myteam", "key", "value"

	tests := []struct {
		name string
		in   *Application
		want map[string]string
	}{
		{
			"test object meta plain",
			&Application{
				ObjectMeta: v1.ObjectMeta{
					Name:      app,
					Namespace: namespace,
					Labels:    map[string]string{
						"team": team,
					},
				},
			},
			map[string]string{
				"app": app,
				"team": team,
			},
		},
		{
			"test object meta custom label",
			&Application{
				ObjectMeta: v1.ObjectMeta{
					Name:      app,
					Namespace: namespace,
					Labels:    map[string]string{
						"team": team,
						key: value,
					},
				},
			},
			map[string]string{
				"app": app,
				"team": team,
				key: value,
			},
		},
		{
			"test object meta app label not overrideable",
			&Application{
				ObjectMeta: v1.ObjectMeta{
					Name:      app,
					Namespace: namespace,
					Labels:    map[string]string{
						"team": team,
						"app": "ignored",
					},
				},
			},
			map[string]string{
				"app": app,
				"team": team,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.in.CreateObjectMeta()
			if !reflect.DeepEqual(got.Labels, tt.want) {
				t.Errorf("CreateObjectMeta().Labels = %v, want %v", got.Labels, tt.want)
			}
		})
	}
}
