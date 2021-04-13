package nais_io_v1alpha1_test

import (
	"encoding/json"
	"testing"

	"github.com/mitchellh/hashstructure"
	"github.com/nais/liberator/pkg/apis/nais.io/v1alpha1"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// Change this value to accept re-synchronization of ALL application resources when deploying a new version.
	applicationHash = "56c407b7c74b1ecc"
)

func TestApplication_Hash(t *testing.T) {
	apps := []*nais_io_v1alpha1.Application{
		{Spec: nais_io_v1alpha1.ApplicationSpec{}},
		{Spec: nais_io_v1alpha1.ApplicationSpec{}, ObjectMeta: v1.ObjectMeta{Annotations: map[string]string{"a": "b", "team": "banan"}}},
		{Spec: nais_io_v1alpha1.ApplicationSpec{}, ObjectMeta: v1.ObjectMeta{Labels: map[string]string{"a": "b", "team": "banan"}}},
	}
	hashes := make([]string, len(apps))
	for i := range apps {
		err := nais_io_v1alpha1.ApplyDefaults(apps[i])
		if err != nil {
			panic(err)
		}
		hashes[i], err = apps[i].Hash()
		if err != nil {
			panic(err)
		}
	}

	assert.Equal(t, hashes[0], hashes[1], "matches, as annotations is ignored")
	assert.NotEqual(t, hashes[1], hashes[2], "should not match")
}

func TestApplication_CreateAppNamespaceHash(t *testing.T) {
	application := &nais_io_v1alpha1.Application{
		ObjectMeta: v1.ObjectMeta{
			Name:      "reallylongapplicationname",
			Namespace: "evenlongernamespacename",
		},
	}
	application2 := &nais_io_v1alpha1.Application{
		ObjectMeta: v1.ObjectMeta{
			Name:      "reallylongafoo",
			Namespace: "evenlongerbar",
		},
	}
	application3 := &nais_io_v1alpha1.Application{
		ObjectMeta: v1.ObjectMeta{
			Name:      "short",
			Namespace: "name",
		},
	}
	appNameHash := application.CreateAppNamespaceHash()
	appNameHash2 := application2.CreateAppNamespaceHash()
	appNameHash3 := application3.CreateAppNamespaceHash()
	assert.Equal(t, "reallylonga-evenlonger-siqwsiq", appNameHash)
	assert.Equal(t, "reallylonga-evenlonger-piuqbfq", appNameHash2)
	assert.Equal(t, "short-name-hqb7npi", appNameHash3)
	assert.True(t, len(appNameHash2) <= 30)
	assert.True(t, len(appNameHash3) >= 6)
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
	newMarshal, _ := json.Marshal(neu)
	oldHash, _ := hashstructure.Hash(oldMarshal, nil)
	newHash, _ := hashstructure.Hash(newMarshal, nil)
	assert.Equal(t, newHash, oldHash)
}

func TestNewCRD(t *testing.T) {
	app := &nais_io_v1alpha1.Application{}
	err := nais_io_v1alpha1.ApplyDefaults(app)
	if err != nil {
		panic(err)
	}
	hash, err := app.Hash()
	assert.NoError(t, err)
	assert.Equalf(t, applicationHash, hash, "Your Application default value changes will trigger a FULL REDEPLOY of ALL APPLICATIONS in ALL NAMESPACES across ALL CLUSTERS. If this is what you really want, change the `applicationHash` constant in this test file to `%s`.", hash)
}
