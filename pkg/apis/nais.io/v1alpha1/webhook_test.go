package nais_io_v1alpha1

import (
	"testing"

	nais_io_v1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
	"k8s.io/apimachinery/pkg/api/errors"
)

var expectedErrors = []string{
	"spec.gcp.bigQueryDatasets",
}

func inputApp() *Application {
	app := getAppDefaults()

	app.Spec.GCP = &nais_io_v1.GCP{
		BigQueryDatasets: []nais_io_v1.CloudBigQueryDataset{
			{
				Name: "newname",
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
				Name: "name",
			},
		},
	}

	return app
}

func TestWebhookValidateUpdateError(t *testing.T) {
	input := inputApp()
	old := oldApp()
	err := input.ValidateUpdate(old)
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
	err := input.ValidateUpdate(old)
	if err != nil {
		t.Fatal("unexpected error", err)
	}
}
