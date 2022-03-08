package v1alpha1_test

import (
	"testing"

	"github.com/nais/liberator/pkg/apis/nebula.skatteetaten.no/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func TestNewCRD(t *testing.T) {
	app := &v1alpha1.Application{}
	err := app.ApplyDefaults()
	assert.NoError(t, err)
}

func TestNewCRDWithAzure(t *testing.T) {
	app := &v1alpha1.Application{
		Spec: v1alpha1.ApplicationSpec{
			Azure: &v1alpha1.AzureConfig{
				ResourceGroup: "test",
			},
		},
	}
	err := app.ApplyDefaults()
	assert.NoError(t, err)
}

func TestLoggingCRD(t *testing.T) {
	app := &v1alpha1.Application{
		Spec: v1alpha1.ApplicationSpec{
			Logging: &v1alpha1.LogConfig{
				Splunk: map[string]*v1alpha1.SplunkLoggingConfig{
					"somekey": {
						SplunkIndex: "index",
						FilePattern: "*.jadda",
						SourceType:  "unit-test",
					},
				},
			},
		},
	}

	err := app.ApplyDefaults()
	for _, splunk := range app.Spec.Logging.Splunk {
		assert.Equal(t, true, *splunk.Enabled, "Expecting unset enabled field to be set as true")
	}
	assert.NoError(t, err)
}
