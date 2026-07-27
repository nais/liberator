package scheme

import (
	"testing"

	aiven_io_v1alpha1 "github.com/nais/liberator/pkg/apis/aiven.io/v1alpha1"
	aiven_nais_io_v1 "github.com/nais/liberator/pkg/apis/aiven.nais.io/v1"
	aiven_nais_io_v2 "github.com/nais/liberator/pkg/apis/aiven.nais.io/v2"
	bigquery_cnrm_cloud_google_com_v1beta1 "github.com/nais/liberator/pkg/apis/bigquery.cnrm.cloud.google.com/v1beta1"
	fdqnnetworkpolicies_networking_gke_io_v1alpha3 "github.com/nais/liberator/pkg/apis/fqdnnetworkpolicies.networking.gke.io/v1alpha3"
	google_nais_io_v1 "github.com/nais/liberator/pkg/apis/google.nais.io/v1"
	iam_cnrm_cloud_google_com_v1beta1 "github.com/nais/liberator/pkg/apis/iam.cnrm.cloud.google.com/v1beta1"
	kafka_nais_io_v1 "github.com/nais/liberator/pkg/apis/kafka.nais.io/v1"
	nais_io_v1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
	nais_io_v1alpha1 "github.com/nais/liberator/pkg/apis/nais.io/v1alpha1"
	sql_cnrm_cloud_google_com_v1beta1 "github.com/nais/liberator/pkg/apis/sql.cnrm.cloud.google.com/v1beta1"
	storage_cnrm_cloud_google_com_v1beta1 "github.com/nais/liberator/pkg/apis/storage.cnrm.cloud.google.com/v1beta1"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
)

func TestAddToSchemeRegistersMetaTypes(t *testing.T) {
	tests := []struct {
		name         string
		groupVersion schema.GroupVersion
		addToScheme  func(*runtime.Scheme) error
	}{
		{"aiven.io/v1alpha1", aiven_io_v1alpha1.GroupVersion, aiven_io_v1alpha1.AddToScheme},
		{"aiven.nais.io/v1", aiven_nais_io_v1.GroupVersion, aiven_nais_io_v1.AddToScheme},
		{"aiven.nais.io/v2", aiven_nais_io_v2.GroupVersion, aiven_nais_io_v2.AddToScheme},
		{"bigquery.cnrm.cloud.google.com/v1beta1", bigquery_cnrm_cloud_google_com_v1beta1.GroupVersion, bigquery_cnrm_cloud_google_com_v1beta1.AddToScheme},
		{"fqdnnetworkpolicies.networking.gke.io/v1alpha3", fdqnnetworkpolicies_networking_gke_io_v1alpha3.GroupVersion, fdqnnetworkpolicies_networking_gke_io_v1alpha3.AddToScheme},
		{"google.nais.io/v1", google_nais_io_v1.GroupVersion, google_nais_io_v1.AddToScheme},
		{"iam.cnrm.cloud.google.com/v1beta1", iam_cnrm_cloud_google_com_v1beta1.GroupVersion, iam_cnrm_cloud_google_com_v1beta1.AddToScheme},
		{"kafka.nais.io/v1", kafka_nais_io_v1.GroupVersion, kafka_nais_io_v1.AddToScheme},
		{"nais.io/v1", nais_io_v1.GroupVersion, nais_io_v1.AddToScheme},
		{"nais.io/v1alpha1", nais_io_v1alpha1.GroupVersion, nais_io_v1alpha1.AddToScheme},
		{"sql.cnrm.cloud.google.com/v1beta1", sql_cnrm_cloud_google_com_v1beta1.GroupVersion, sql_cnrm_cloud_google_com_v1beta1.AddToScheme},
		{"storage.cnrm.cloud.google.com/v1beta1", storage_cnrm_cloud_google_com_v1beta1.GroupVersion, storage_cnrm_cloud_google_com_v1beta1.AddToScheme},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scheme := runtime.NewScheme()
			require.NoError(t, clientgoscheme.AddToScheme(scheme))
			require.NoError(t, tt.addToScheme(scheme))
			_, err := scheme.ConvertToVersion(&metav1.CreateOptions{}, tt.groupVersion)
			require.NoError(t, err)
		})
	}
}
