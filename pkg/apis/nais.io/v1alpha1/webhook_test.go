package nais_io_v1alpha1

import (
	"testing"

	aiven_io_v1alpha1 "github.com/nais/liberator/pkg/apis/aiven.io/v1alpha1"
	nais_io_v1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
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
	_ = AddToScheme(scheme)

	return fake.NewClientBuilder().WithScheme(scheme).WithObjects(objs...).Build()
}

func TestApplicationValidator_ValidateCreate(t *testing.T) {
	t.Run("valid application without aiven references", func(t *testing.T) {
		validator := &ApplicationValidator{Client: fakeKubeClient()}
		app := &Application{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-app",
				Namespace: "test-ns",
			},
			Spec: ApplicationSpec{
				Image: "nginx:latest",
			},
		}

		warnings, err := validator.ValidateCreate(t.Context(), app)
		assert.NoError(t, err)
		assert.Empty(t, warnings)
	})

	t.Run("application name too long", func(t *testing.T) {
		validator := &ApplicationValidator{Client: fakeKubeClient()}
		app := &Application{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "this-is-a-very-long-name-that-exceeds-the-maximum-allowed-length-for-kubernetes-label-values-which-is-63-characters",
				Namespace: "test-ns",
			},
			Spec: ApplicationSpec{
				Image: "nginx:latest",
			},
		}

		warnings, err := validator.ValidateCreate(t.Context(), app)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "name length must be no more than 63 characters")
		assert.Empty(t, warnings)
	})

	t.Run("valid TTL duration", func(t *testing.T) {
		validator := &ApplicationValidator{Client: fakeKubeClient()}
		ttl := "12h"
		app := &Application{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-app",
				Namespace: "test-ns",
			},
			Spec: ApplicationSpec{
				Image: "nginx:latest",
				TTL:   ttl,
			},
		}

		warnings, err := validator.ValidateCreate(t.Context(), app)
		assert.NoError(t, err)
		assert.Empty(t, warnings)
	})

	t.Run("invalid TTL duration", func(t *testing.T) {
		validator := &ApplicationValidator{Client: fakeKubeClient()}
		ttl := "invalid-duration"
		app := &Application{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-app",
				Namespace: "test-ns",
			},
			Spec: ApplicationSpec{
				Image: "nginx:latest",
				TTL:   ttl,
			},
		}

		warnings, err := validator.ValidateCreate(t.Context(), app)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "TTL is not a valid duration")
		assert.Empty(t, warnings)
	})

	t.Run("opensearch reference exists", func(t *testing.T) {
		namespace := "test-ns"
		instance := "my-opensearch"
		fullyQualifiedName := aiven_io_v1alpha1.OpenSearchFullyQualifiedName(instance, namespace)

		opensearch := &aiven_io_v1alpha1.OpenSearch{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fullyQualifiedName,
				Namespace: namespace,
			},
		}

		validator := &ApplicationValidator{Client: fakeKubeClient(opensearch)}
		app := &Application{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-app",
				Namespace: namespace,
			},
			Spec: ApplicationSpec{
				Image: "nginx:latest",
				OpenSearch: &nais_io_v1.OpenSearch{
					Instance: instance,
				},
			},
		}

		warnings, err := validator.ValidateCreate(t.Context(), app)
		assert.NoError(t, err)
		assert.Empty(t, warnings)
	})

	t.Run("opensearch reference not found", func(t *testing.T) {
		validator := &ApplicationValidator{Client: fakeKubeClient()}
		app := &Application{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-app",
				Namespace: "test-ns",
			},
			Spec: ApplicationSpec{
				Image: "nginx:latest",
				OpenSearch: &nais_io_v1.OpenSearch{
					Instance: "nonexistent-opensearch",
				},
			},
		}

		warnings, err := validator.ValidateCreate(t.Context(), app)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "referenced OpenSearch instance 'nonexistent-opensearch' not found")
		assert.Empty(t, warnings)
	})

	t.Run("valkey reference exists", func(t *testing.T) {
		namespace := "test-ns"
		instance := "my-valkey"
		fullyQualifiedName := aiven_io_v1alpha1.ValkeyFullyQualifiedName(instance, namespace)

		valkey := &aiven_io_v1alpha1.Valkey{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fullyQualifiedName,
				Namespace: namespace,
			},
		}

		validator := &ApplicationValidator{Client: fakeKubeClient(valkey)}
		app := &Application{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-app",
				Namespace: namespace,
			},
			Spec: ApplicationSpec{
				Image: "nginx:latest",
				Valkey: []nais_io_v1.Valkey{
					{
						Instance: instance,
					},
				},
			},
		}

		warnings, err := validator.ValidateCreate(t.Context(), app)
		assert.NoError(t, err)
		assert.Empty(t, warnings)
	})

	t.Run("valkey reference not found", func(t *testing.T) {
		validator := &ApplicationValidator{Client: fakeKubeClient()}
		app := &Application{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-app",
				Namespace: "test-ns",
			},
			Spec: ApplicationSpec{
				Image: "nginx:latest",
				Valkey: []nais_io_v1.Valkey{
					{
						Instance: "nonexistent-valkey",
					},
				},
			},
		}

		warnings, err := validator.ValidateCreate(t.Context(), app)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "referenced Valkey instance 'nonexistent-valkey' not found")
		assert.Empty(t, warnings)
	})

	t.Run("multiple valkey references all exist", func(t *testing.T) {
		namespace := "test-ns"
		instance1 := "valkey-1"
		instance2 := "valkey-2"

		valkey1 := &aiven_io_v1alpha1.Valkey{
			ObjectMeta: metav1.ObjectMeta{
				Name:      aiven_io_v1alpha1.ValkeyFullyQualifiedName(instance1, namespace),
				Namespace: namespace,
			},
		}
		valkey2 := &aiven_io_v1alpha1.Valkey{
			ObjectMeta: metav1.ObjectMeta{
				Name:      aiven_io_v1alpha1.ValkeyFullyQualifiedName(instance2, namespace),
				Namespace: namespace,
			},
		}

		validator := &ApplicationValidator{Client: fakeKubeClient(valkey1, valkey2)}
		app := &Application{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-app",
				Namespace: namespace,
			},
			Spec: ApplicationSpec{
				Image: "nginx:latest",
				Valkey: []nais_io_v1.Valkey{
					{Instance: instance1},
					{Instance: instance2},
				},
			},
		}

		warnings, err := validator.ValidateCreate(t.Context(), app)
		assert.NoError(t, err)
		assert.Empty(t, warnings)
	})

	t.Run("combined opensearch and valkey references", func(t *testing.T) {
		namespace := "test-ns"
		osInstance := "my-opensearch"
		vkInstance := "my-valkey"

		opensearch := &aiven_io_v1alpha1.OpenSearch{
			ObjectMeta: metav1.ObjectMeta{
				Name:      aiven_io_v1alpha1.OpenSearchFullyQualifiedName(osInstance, namespace),
				Namespace: namespace,
			},
		}
		valkey := &aiven_io_v1alpha1.Valkey{
			ObjectMeta: metav1.ObjectMeta{
				Name:      aiven_io_v1alpha1.ValkeyFullyQualifiedName(vkInstance, namespace),
				Namespace: namespace,
			},
		}

		validator := &ApplicationValidator{Client: fakeKubeClient(opensearch, valkey)}
		app := &Application{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-app",
				Namespace: namespace,
			},
			Spec: ApplicationSpec{
				Image: "nginx:latest",
				OpenSearch: &nais_io_v1.OpenSearch{
					Instance: osInstance,
				},
				Valkey: []nais_io_v1.Valkey{
					{Instance: vkInstance},
				},
			},
		}

		warnings, err := validator.ValidateCreate(t.Context(), app)
		assert.NoError(t, err)
		assert.Empty(t, warnings)
	})
}

func TestApplicationValidator_ValidateUpdate(t *testing.T) {
	t.Run("valid update without spec changes", func(t *testing.T) {
		validator := &ApplicationValidator{Client: fakeKubeClient()}
		oldApp := &Application{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-app",
				Namespace: "test-ns",
			},
			Spec: ApplicationSpec{
				Image: "nginx:latest",
			},
		}
		newApp := oldApp.DeepCopy()

		warnings, err := validator.ValidateUpdate(t.Context(), oldApp, newApp)
		assert.NoError(t, err)
		assert.Empty(t, warnings)
	})

	t.Run("update with spec changes that should fail validation", func(t *testing.T) {
		validator := &ApplicationValidator{Client: fakeKubeClient()}
		oldApp := &Application{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-app",
				Namespace: "test-ns",
			},
			Spec: ApplicationSpec{
				Image: "nginx:latest",
				GCP: &nais_io_v1.GCP{
					BigQueryDatasets: []nais_io_v1.CloudBigQueryDataset{
						{
							Name:       "dataset1",
							Permission: nais_io_v1.BigQueryPermissionReadWrite,
						},
					},
				},
			},
		}
		newApp := oldApp.DeepCopy()
		newApp.Spec.GCP.BigQueryDatasets[0].Permission = nais_io_v1.BigQueryPermissionRead

		warnings, err := validator.ValidateUpdate(t.Context(), oldApp, newApp)
		assert.Error(t, err)
		assert.Empty(t, warnings)
	})

	t.Run("update with aiven references", func(t *testing.T) {
		namespace := "test-ns"
		instance := "my-opensearch"
		fullyQualifiedName := aiven_io_v1alpha1.OpenSearchFullyQualifiedName(instance, namespace)

		opensearch := &aiven_io_v1alpha1.OpenSearch{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fullyQualifiedName,
				Namespace: namespace,
			},
		}

		validator := &ApplicationValidator{Client: fakeKubeClient(opensearch)}
		oldApp := &Application{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-app",
				Namespace: namespace,
			},
			Spec: ApplicationSpec{
				Image: "nginx:latest",
			},
		}
		newApp := oldApp.DeepCopy()
		newApp.Spec.OpenSearch = &nais_io_v1.OpenSearch{
			Instance: instance,
		}

		warnings, err := validator.ValidateUpdate(t.Context(), oldApp, newApp)
		assert.NoError(t, err)
		assert.Empty(t, warnings)
	})

	t.Run("update with invalid aiven reference", func(t *testing.T) {
		validator := &ApplicationValidator{Client: fakeKubeClient()}
		oldApp := &Application{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-app",
				Namespace: "test-ns",
			},
			Spec: ApplicationSpec{
				Image: "nginx:latest",
			},
		}
		newApp := oldApp.DeepCopy()
		newApp.Spec.Valkey = []nais_io_v1.Valkey{
			{Instance: "nonexistent-valkey"},
		}

		warnings, err := validator.ValidateUpdate(t.Context(), oldApp, newApp)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "referenced Valkey instance 'nonexistent-valkey' not found")
		assert.Empty(t, warnings)
	})

	t.Run("update with invalid opensearch reference", func(t *testing.T) {
		validator := &ApplicationValidator{Client: fakeKubeClient()}
		oldApp := &Application{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-app",
				Namespace: "test-ns",
			},
			Spec: ApplicationSpec{
				Image: "nginx:latest",
			},
		}
		newApp := oldApp.DeepCopy()
		newApp.Spec.OpenSearch = &nais_io_v1.OpenSearch{
			Instance: "nonexistent-opensearch",
		}

		warnings, err := validator.ValidateUpdate(t.Context(), oldApp, newApp)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "referenced OpenSearch instance 'nonexistent-opensearch' not found")
		assert.Empty(t, warnings)
	})

	t.Run("invalid old object type", func(t *testing.T) {
		validator := &ApplicationValidator{Client: fakeKubeClient()}
		oldApp := &nais_io_v1.Image{} // Wrong type
		newApp := &Application{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-app",
				Namespace: "test-ns",
			},
		}

		warnings, err := validator.ValidateUpdate(t.Context(), oldApp, newApp)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expected an Application")
		assert.Empty(t, warnings)
	})

	t.Run("invalid new object type", func(t *testing.T) {
		validator := &ApplicationValidator{Client: fakeKubeClient()}
		oldApp := &Application{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-app",
				Namespace: "test-ns",
			},
		}
		newApp := &nais_io_v1.Image{} // Wrong type

		warnings, err := validator.ValidateUpdate(t.Context(), oldApp, newApp)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expected an Application")
		assert.Empty(t, warnings)
	})
}

func TestApplicationValidator_ValidateDelete(t *testing.T) {
	t.Run("delete always succeeds", func(t *testing.T) {
		validator := &ApplicationValidator{Client: fakeKubeClient()}
		app := &Application{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-app",
				Namespace: "test-ns",
			},
			Spec: ApplicationSpec{
				Image: "nginx:latest",
			},
		}

		warnings, err := validator.ValidateDelete(t.Context(), app)
		assert.NoError(t, err)
		assert.Empty(t, warnings)
	})
}

func TestApplicationValidator_ValidateCreate_InvalidObjectType(t *testing.T) {
	validator := &ApplicationValidator{Client: fakeKubeClient()}

	// Pass wrong object type
	wrongObj := &nais_io_v1.Image{}
	warnings, err := validator.ValidateCreate(t.Context(), wrongObj)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "expected an Application")
	assert.Empty(t, warnings)
}
