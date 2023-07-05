package nais_io_v1alpha1

import (
	"testing"

	"k8s.io/apimachinery/pkg/api/errors"

	nais_io_v1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
)

var expectedErrors = []string{
	"spec.gcp.bigQueryDatasets.0.permission",
}

func inputApp() *Application {
	app := getAppDefaults()

	app.Spec.GCP = &nais_io_v1.GCP{
		BigQueryDatasets: []nais_io_v1.CloudBigQueryDataset{
			{
				Name:       "name",
				Permission: nais_io_v1.BigQueryPermissionRead,
			},
		},
	}

	return app
}

func oldApp() *Application {
	app := getAppDefaults()

	app.Spec.GCP = &nais_io_v1.GCP{
		BigQueryDatasets: []nais_io_v1.CloudBigQueryDataset{
			{
				Name:       "name",
				Permission: nais_io_v1.BigQueryPermissionReadWrite,
			},
		},
	}

	return app
}

func TestWebhookValidateUpdateError(t *testing.T) {
	input := inputApp()
	old := oldApp()
	_, err := input.ValidateUpdate(old)
	if err == nil {
		t.Fatal("no error returned")
	}

	causes := err.(*errors.StatusError).Status().Details.Causes

	if len(expectedErrors) != len(causes) {
		t.Errorf("expected %v errors, got %v", len(expectedErrors), len(causes))
	}

	found := map[string]bool{}
	for _, terr := range causes {
		found[terr.Field] = true
	}
	for _, expected := range expectedErrors {
		if !found[expected] {
			t.Errorf("expected error: %q", expected)
		}
		delete(found, expected)
	}
	for val := range found {
		t.Errorf("got %q, but did not expect it", val)
	}
}

func TestWebhookValidateUpdate(t *testing.T) {
	input := inputApp()
	old := inputApp()
	_, err := input.ValidateUpdate(old)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
}

func TestWebhookValidateCreate(t *testing.T) {
	input := inputApp()
	_, err := input.ValidateCreate()
	if err != nil {
		t.Fatal("unexpected error", err)
	}
}

func TestWebhookValidateCreateTooLongName(t *testing.T) {
	input := inputApp()
	input.SetName("this-is-a-very-long-name-that-is-longer-than-63-characters-which-is-the-maximum-length-of-a-kubernetes-resource-name")
	_, err := input.ValidateCreate()
	if err.Error() != "Application name length must be no more than 63 characters" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestWebhookValidateTTL(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		input := inputApp()
		input.Spec.TTL = "12h"
		_, err := input.ValidateCreate()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		input := inputApp()
		input.Spec.TTL = "invalid"
		_, err := input.ValidateCreate()
		want := "TTL is not a valid duration: \"invalid\". Example of valid duration is '12h'"
		if err.Error() != want {
			t.Errorf("got: %s, want: %s", err, want)
		}
	})
}
