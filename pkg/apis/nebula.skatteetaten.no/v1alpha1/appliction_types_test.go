package v1alpha1_test

import (
	"testing"

	"github.com/nais/liberator/pkg/apis/nebula.skatteetaten.no/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func TestNewCRD(t *testing.T) {
	app := &v1alpha1.Application{}
	err := app.ApplyDefaults()
	assert.NoError(t,  err)
}

func TestNewCRDWithAzure(t *testing.T) {
	app := &v1alpha1.Application{
		Spec:       v1alpha1.ApplicationSpec{
			Azure:                                 &v1alpha1.AzureConfig{
				ResourceGroup:    "test",
			},
		},
	}
	err := app.ApplyDefaults()
	assert.NoError(t,  err)
}