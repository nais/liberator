package nais_io_v1alpha1_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/mitchellh/hashstructure"
	nais_io_v1alpha1 "github.com/nais/liberator/pkg/apis/nais.io/v1alpha1"
	"github.com/nais/liberator/pkg/hash"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// Change this value to accept re-synchronization of ALL application resources when deploying a new version.
	applicationHash = "3b136519de70ef72"
	aivenGeneration = 0
)

func TestApplication_Hash(t *testing.T) {
	apps := []*nais_io_v1alpha1.Application{
		{Spec: nais_io_v1alpha1.ApplicationSpec{}},
		{Spec: nais_io_v1alpha1.ApplicationSpec{}, ObjectMeta: metav1.ObjectMeta{Annotations: map[string]string{"a": "b", "team": "banan"}}},
		{Spec: nais_io_v1alpha1.ApplicationSpec{}, ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"a": "b", "team": "banan"}}},
	}
	hashes := make([]string, len(apps))
	for i := range apps {
		err := apps[i].ApplyDefaults()
		if err != nil {
			panic(err)
		}
		hashes[i], err = apps[i].Hash(aivenGeneration)
		if err != nil {
			panic(err)
		}
	}

	assert.Equal(t, hashes[0], hashes[1], "matches, as annotations are ignored")
	assert.NotEqual(t, hashes[1], hashes[2], "should not match")
}

// Test that updating the application spec with new, default-null values does not trigger a hash change.
func TestHashJSONMarshalling(t *testing.T) {
	type a struct {
		Foo string `json:"foo"`
	}
	type oldspec struct {
		A *a `json:"a,omitempty"`
	}
	type newspec struct {
		A *a     `json:"a,omitempty"`
		B *a     `json:"b,omitempty"` // new field added to crd spec
		C string `json:"c,omitempty"` // new field added to crd spec
	}
	old := &oldspec{}
	neu := &newspec{}
	oldMarshal, _ := json.Marshal(old)
	oldHashI, _ := hashstructure.Hash(oldMarshal, nil)
	oldHash := fmt.Sprintf("%x", oldHashI)
	newHash, _ := hash.Hash(neu)
	assert.Equal(t, newHash, oldHash)
}

func TestNewCRD(t *testing.T) {
	app := &nais_io_v1alpha1.Application{
		Spec: nais_io_v1alpha1.ApplicationSpec{
			Image: "example/image:tag",
		},
	}
	err := app.ApplyDefaults()
	if err != nil {
		panic(err)
	}
	hash, err := app.Hash(aivenGeneration)
	assert.NoError(t, err)
	assert.Equalf(t, applicationHash, hash, "Your Application default value changes will trigger a FULL REDEPLOY of ALL APPLICATIONS in ALL NAMESPACES across ALL CLUSTERS. If this is what you really want, change the `applicationHash` constant in this test file to `%s`.", hash)
}
