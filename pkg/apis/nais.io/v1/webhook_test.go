package nais_io_v1

import (
	"testing"
	"time"

	aiven_io_v1alpha1 "github.com/nais/liberator/pkg/apis/aiven.io/v1alpha1"
	aiven_nais_io_v1 "github.com/nais/liberator/pkg/apis/aiven.nais.io/v1"
	data_nais_io_v1 "github.com/nais/pgrator/pkg/api/datav1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func fakeKubeClient(objs ...client.Object) client.Client {
	scheme := runtime.NewScheme()

	// Add necessary schemes for the test
	_ = clientgoscheme.AddToScheme(scheme)
	_ = aiven_io_v1alpha1.AddToScheme(scheme)
	_ = data_nais_io_v1.AddToScheme(scheme)
	_ = AddToScheme(scheme)

	return fake.NewClientBuilder().WithScheme(scheme).WithObjects(objs...).Build()
}

func TestJobValidator_ValidateCreate(t *testing.T) {
	t.Run("valid naisjob without aiven references", func(t *testing.T) {
		validator := &JobValidator{Client: fakeKubeClient()}
		nj := &Naisjob{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-job",
				Namespace: "test-ns",
			},
			Spec: NaisjobSpec{
				Image:    "nginx:latest",
				Schedule: "0 * * * *",
			},
		}

		warnings, err := validator.ValidateCreate(t.Context(), nj)
		assert.NoError(t, err)
		assert.Empty(t, warnings)
	})

	t.Run("naisjob name too long", func(t *testing.T) {
		validator := &JobValidator{Client: fakeKubeClient()}
		nj := &Naisjob{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "this-is-a-very-long-name-that-exceeds-the-maximum-allowed-length-for-kubernetes-label-values-which-is-63-characters",
				Namespace: "test-ns",
			},
			Spec: NaisjobSpec{
				Image:    "nginx:latest",
				Schedule: "0 * * * *",
			},
		}

		warnings, err := validator.ValidateCreate(t.Context(), nj)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "name length must be no more than 63 characters")
		assert.Empty(t, warnings)
	})

	t.Run("valid TTL duration", func(t *testing.T) {
		validator := &JobValidator{Client: fakeKubeClient()}
		ttl := "12h"
		nj := &Naisjob{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-job",
				Namespace: "test-ns",
			},
			Spec: NaisjobSpec{
				Image:    "nginx:latest",
				Schedule: "0 * * * *",
				TTL:      ttl,
			},
		}

		warnings, err := validator.ValidateCreate(t.Context(), nj)
		assert.NoError(t, err)
		assert.Empty(t, warnings)
	})

	t.Run("invalid TTL duration", func(t *testing.T) {
		validator := &JobValidator{Client: fakeKubeClient()}
		ttl := "invalid-duration"
		nj := &Naisjob{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-job",
				Namespace: "test-ns",
			},
			Spec: NaisjobSpec{
				Image:    "nginx:latest",
				Schedule: "0 * * * *",
				TTL:      ttl,
			},
		}

		warnings, err := validator.ValidateCreate(t.Context(), nj)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "TTL is not a valid duration")
		assert.Empty(t, warnings)
	})

	t.Run("opensearch reference exists", func(t *testing.T) {
		namespace := "test-ns"
		instance := "my-opensearch"
		fullyQualifiedName := aiven_nais_io_v1.OpenSearchFullyQualifiedName(instance, namespace)

		opensearch := &aiven_io_v1alpha1.OpenSearch{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fullyQualifiedName,
				Namespace: namespace,
			},
		}

		validator := &JobValidator{Client: fakeKubeClient(opensearch)}
		nj := &Naisjob{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-job",
				Namespace: namespace,
			},
			Spec: NaisjobSpec{
				Image:    "nginx:latest",
				Schedule: "0 * * * *",
				OpenSearch: &OpenSearch{
					Instance: instance,
				},
			},
		}

		warnings, err := validator.ValidateCreate(t.Context(), nj)
		assert.NoError(t, err)
		assert.Empty(t, warnings)
	})

	t.Run("opensearch reference not found", func(t *testing.T) {
		validator := &JobValidator{Client: fakeKubeClient()}
		nj := &Naisjob{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-job",
				Namespace: "test-ns",
			},
			Spec: NaisjobSpec{
				Image:    "nginx:latest",
				Schedule: "0 * * * *",
				OpenSearch: &OpenSearch{
					Instance: "nonexistent-opensearch",
				},
			},
		}

		warnings, err := validator.ValidateCreate(t.Context(), nj)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "OpenSearch 'nonexistent-opensearch' does not exist")
		assert.Empty(t, warnings)
	})

	t.Run("valkey reference exists", func(t *testing.T) {
		namespace := "test-ns"
		instance := "my-valkey"
		fullyQualifiedName := aiven_nais_io_v1.ValkeyFullyQualifiedName(instance, namespace)

		valkey := &aiven_io_v1alpha1.Valkey{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fullyQualifiedName,
				Namespace: namespace,
			},
		}

		validator := &JobValidator{Client: fakeKubeClient(valkey)}
		nj := &Naisjob{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-job",
				Namespace: namespace,
			},
			Spec: NaisjobSpec{
				Image:    "nginx:latest",
				Schedule: "0 * * * *",
				Valkey: []Valkey{
					{
						Instance: instance,
					},
				},
			},
		}

		warnings, err := validator.ValidateCreate(t.Context(), nj)
		assert.NoError(t, err)
		assert.Empty(t, warnings)
	})

	t.Run("valkey reference not found", func(t *testing.T) {
		validator := &JobValidator{Client: fakeKubeClient()}
		nj := &Naisjob{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-job",
				Namespace: "test-ns",
			},
			Spec: NaisjobSpec{
				Image:    "nginx:latest",
				Schedule: "0 * * * *",
				Valkey: []Valkey{
					{
						Instance: "nonexistent-valkey",
					},
				},
			},
		}

		warnings, err := validator.ValidateCreate(t.Context(), nj)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Valkey 'nonexistent-valkey' does not exist")
		assert.Empty(t, warnings)
	})

	t.Run("multiple valkey references all exist", func(t *testing.T) {
		namespace := "test-ns"
		instance1 := "valkey-1"
		instance2 := "valkey-2"

		valkey1 := &aiven_io_v1alpha1.Valkey{
			ObjectMeta: metav1.ObjectMeta{
				Name:      aiven_nais_io_v1.ValkeyFullyQualifiedName(instance1, namespace),
				Namespace: namespace,
			},
		}
		valkey2 := &aiven_io_v1alpha1.Valkey{
			ObjectMeta: metav1.ObjectMeta{
				Name:      aiven_nais_io_v1.ValkeyFullyQualifiedName(instance2, namespace),
				Namespace: namespace,
			},
		}

		validator := &JobValidator{Client: fakeKubeClient(valkey1, valkey2)}
		nj := &Naisjob{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-job",
				Namespace: namespace,
			},
			Spec: NaisjobSpec{
				Image:    "nginx:latest",
				Schedule: "0 * * * *",
				Valkey: []Valkey{
					{Instance: instance1},
					{Instance: instance2},
				},
			},
		}

		warnings, err := validator.ValidateCreate(t.Context(), nj)
		assert.NoError(t, err)
		assert.Empty(t, warnings)
	})

	t.Run("combined opensearch and valkey references", func(t *testing.T) {
		namespace := "test-ns"
		osInstance := "my-opensearch"
		vkInstance := "my-valkey"

		opensearch := &aiven_io_v1alpha1.OpenSearch{
			ObjectMeta: metav1.ObjectMeta{
				Name:      aiven_nais_io_v1.OpenSearchFullyQualifiedName(osInstance, namespace),
				Namespace: namespace,
			},
		}
		valkey := &aiven_io_v1alpha1.Valkey{
			ObjectMeta: metav1.ObjectMeta{
				Name:      aiven_nais_io_v1.ValkeyFullyQualifiedName(vkInstance, namespace),
				Namespace: namespace,
			},
		}

		validator := &JobValidator{Client: fakeKubeClient(opensearch, valkey)}
		nj := &Naisjob{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-job",
				Namespace: namespace,
			},
			Spec: NaisjobSpec{
				Image:    "nginx:latest",
				Schedule: "0 * * * *",
				OpenSearch: &OpenSearch{
					Instance: osInstance,
				},
				Valkey: []Valkey{
					{Instance: vkInstance},
				},
			},
		}

		warnings, err := validator.ValidateCreate(t.Context(), nj)
		assert.NoError(t, err)
		assert.Empty(t, warnings)
	})

	t.Run("postgres reference exists", func(t *testing.T) {
		namespace := "test-ns"
		clusterName := "my-postgres-cluster"

		postgres := &data_nais_io_v1.Postgres{
			ObjectMeta: metav1.ObjectMeta{
				Name:      clusterName,
				Namespace: namespace,
			},
		}

		validator := &JobValidator{Client: fakeKubeClient(postgres)}
		nj := &Naisjob{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-job",
				Namespace: namespace,
			},
			Spec: NaisjobSpec{
				Image:    "nginx:latest",
				Schedule: "0 * * * *",
				Postgres: &Postgres{
					ClusterName: clusterName,
				},
			},
		}

		warnings, err := validator.ValidateCreate(t.Context(), nj)
		assert.NoError(t, err)
		assert.Empty(t, warnings)
	})

	t.Run("postgres reference not found", func(t *testing.T) {
		validator := &JobValidator{Client: fakeKubeClient()}
		nj := &Naisjob{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-job",
				Namespace: "test-ns",
			},
			Spec: NaisjobSpec{
				Image:    "nginx:latest",
				Schedule: "0 * * * *",
				Postgres: &Postgres{
					ClusterName: "no-such-postgres",
				},
			},
		}

		warnings, err := validator.ValidateCreate(t.Context(), nj)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Postgres 'no-such-postgres' does not exist")
		assert.Empty(t, warnings)
	})
}

func TestJobValidator_ValidateUpdate(t *testing.T) {
	t.Run("valid update without spec changes", func(t *testing.T) {
		validator := &JobValidator{Client: fakeKubeClient()}
		oldNj := &Naisjob{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-job",
				Namespace: "test-ns",
			},
			Spec: NaisjobSpec{
				Image:    "nginx:latest",
				Schedule: "0 * * * *",
			},
		}
		newNj := oldNj.DeepCopy()

		warnings, err := validator.ValidateUpdate(t.Context(), oldNj, newNj)
		assert.NoError(t, err)
		assert.Empty(t, warnings)
	})

	t.Run("update with aiven references", func(t *testing.T) {
		namespace := "test-ns"
		instance := "my-opensearch"
		fullyQualifiedName := aiven_nais_io_v1.OpenSearchFullyQualifiedName(instance, namespace)

		opensearch := &aiven_io_v1alpha1.OpenSearch{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fullyQualifiedName,
				Namespace: namespace,
			},
		}

		validator := &JobValidator{Client: fakeKubeClient(opensearch)}
		oldNj := &Naisjob{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-job",
				Namespace: namespace,
			},
			Spec: NaisjobSpec{
				Image:    "nginx:latest",
				Schedule: "0 * * * *",
			},
		}
		newNj := oldNj.DeepCopy()
		newNj.Spec.OpenSearch = &OpenSearch{
			Instance: instance,
		}

		warnings, err := validator.ValidateUpdate(t.Context(), oldNj, newNj)
		assert.NoError(t, err)
		assert.Empty(t, warnings)
	})

	t.Run("update with invalid aiven reference", func(t *testing.T) {
		validator := &JobValidator{Client: fakeKubeClient()}
		oldNj := &Naisjob{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-job",
				Namespace: "test-ns",
			},
			Spec: NaisjobSpec{
				Image:    "nginx:latest",
				Schedule: "0 * * * *",
			},
		}
		newNj := oldNj.DeepCopy()
		newNj.Spec.Valkey = []Valkey{
			{Instance: "nonexistent-valkey"},
		}

		warnings, err := validator.ValidateUpdate(t.Context(), oldNj, newNj)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Valkey 'nonexistent-valkey' does not exist")
		assert.Empty(t, warnings)
	})

	t.Run("update with invalid opensearch reference", func(t *testing.T) {
		validator := &JobValidator{Client: fakeKubeClient()}
		oldNj := &Naisjob{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-job",
				Namespace: "test-ns",
			},
			Spec: NaisjobSpec{
				Image:    "nginx:latest",
				Schedule: "0 * * * *",
			},
		}
		newNj := oldNj.DeepCopy()
		newNj.Spec.OpenSearch = &OpenSearch{
			Instance: "nonexistent-opensearch",
		}

		warnings, err := validator.ValidateUpdate(t.Context(), oldNj, newNj)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "OpenSearch 'nonexistent-opensearch' does not exist")
		assert.Empty(t, warnings)
	})

	t.Run("invalid old object type", func(t *testing.T) {
		validator := &JobValidator{Client: fakeKubeClient()}
		oldNj := &Image{} // Wrong type
		newNj := &Naisjob{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-job",
				Namespace: "test-ns",
			},
		}

		warnings, err := validator.ValidateUpdate(t.Context(), oldNj, newNj)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expected a Naisjob")
		assert.Empty(t, warnings)
	})

	t.Run("invalid new object type", func(t *testing.T) {
		validator := &JobValidator{Client: fakeKubeClient()}
		oldNj := &Naisjob{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-job",
				Namespace: "test-ns",
			},
		}
		newNj := &Image{} // Wrong type

		warnings, err := validator.ValidateUpdate(t.Context(), oldNj, newNj)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expected a Naisjob")
		assert.Empty(t, warnings)
	})
}

func TestJobValidator_ValidateDelete(t *testing.T) {
	t.Run("delete always succeeds", func(t *testing.T) {
		validator := &JobValidator{Client: fakeKubeClient()}
		nj := &Naisjob{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-job",
				Namespace: "test-ns",
			},
			Spec: NaisjobSpec{
				Image:    "nginx:latest",
				Schedule: "0 * * * *",
			},
		}

		warnings, err := validator.ValidateDelete(t.Context(), nj)
		assert.NoError(t, err)
		assert.Empty(t, warnings)
	})
}

func TestJobValidator_ValidateCreate_InvalidObjectType(t *testing.T) {
	validator := &JobValidator{Client: fakeKubeClient()}

	// Pass wrong object type
	wrongObj := &Image{}
	warnings, err := validator.ValidateCreate(t.Context(), wrongObj)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "expected a Naisjob")
	assert.Empty(t, warnings)
}

func TestJobMutator_Default(t *testing.T) {
	t.Run("sets kill-after label when TTL is specified", func(t *testing.T) {
		mutator := &JobMutator{}
		nj := &Naisjob{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-job",
				Namespace: "test-ns",
			},
			Spec: NaisjobSpec{
				Image:    "nginx:latest",
				Schedule: "0 * * * *",
				TTL:      "1h",
			},
		}

		err := mutator.Default(t.Context(), nj)
		require.NoError(t, err)
		require.NotNil(t, nj.Labels)

		killAfter, exists := nj.Labels[LabelKillAfter]
		require.True(t, exists, "kill-after label should be set")

		_, err = time.Parse(time.RFC3339, killAfter)
		require.NoError(t, err, "kill-after label should be valid RFC3339 timestamp")
	})

	t.Run("does not set label when TTL is empty", func(t *testing.T) {
		mutator := &JobMutator{}
		nj := &Naisjob{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-job",
				Namespace: "test-ns",
			},
			Spec: NaisjobSpec{
				Image:    "nginx:latest",
				Schedule: "0 * * * *",
			},
		}

		err := mutator.Default(t.Context(), nj)
		require.NoError(t, err)

		if nj.Labels != nil {
			_, exists := nj.Labels[LabelKillAfter]
			assert.False(t, exists, "kill-after label should not be set when TTL is empty")
		}
	})

	t.Run("overwrites existing kill-after label", func(t *testing.T) {
		mutator := &JobMutator{}
		existingTime := "2025-01-01T12:00:00Z"
		nj := &Naisjob{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-job",
				Namespace: "test-ns",
				Labels: map[string]string{
					LabelKillAfter: existingTime,
				},
			},
			Spec: NaisjobSpec{
				Image:    "nginx:latest",
				Schedule: "0 * * * *",
				TTL:      "1h",
			},
		}

		err := mutator.Default(t.Context(), nj)
		require.NoError(t, err)

		assert.NotEqual(t, existingTime, nj.Labels[LabelKillAfter], "kill-after label should be overwritten")
	})

	t.Run("handles invalid TTL gracefully", func(t *testing.T) {
		mutator := &JobMutator{}
		nj := &Naisjob{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-job",
				Namespace: "test-ns",
			},
			Spec: NaisjobSpec{
				Image:    "nginx:latest",
				Schedule: "0 * * * *",
				TTL:      "invalid",
			},
		}

		err := mutator.Default(t.Context(), nj)
		require.NoError(t, err)

		if nj.Labels != nil {
			_, exists := nj.Labels[LabelKillAfter]
			assert.False(t, exists, "kill-after label should not be set for invalid TTL")
		}
	})

	t.Run("returns error for wrong object type", func(t *testing.T) {
		mutator := &JobMutator{}
		wrongObj := &Image{}

		err := mutator.Default(t.Context(), wrongObj)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "expected a Naisjob")
	})
}
