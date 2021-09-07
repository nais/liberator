package scheme

import (
	"fmt"

	aiven_nais_io_v1 "github.com/nais/liberator/pkg/apis/aiven.nais.io/v1"
	azure_microsoft_com_v1alpha1 "github.com/nais/liberator/pkg/apis/azure.microsoft.com/v1alpha1"
	azure_microsoft_com_v1alpha2 "github.com/nais/liberator/pkg/apis/azure.microsoft.com/v1alpha2"
	azure_microsoft_com_v1beta1 "github.com/nais/liberator/pkg/apis/azure.microsoft.com/v1beta1"
	bigquery_cnrm_cloud_google_com_v1beta1 "github.com/nais/liberator/pkg/apis/bigquery.cnrm.cloud.google.com/v1beta1"
	iam_cnrm_cloud_google_com_v1beta1 "github.com/nais/liberator/pkg/apis/iam.cnrm.cloud.google.com/v1beta1"
	kafka_nais_io_v1 "github.com/nais/liberator/pkg/apis/kafka.nais.io/v1"
	nais_io_v1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
	nais_io_v1alpha1 "github.com/nais/liberator/pkg/apis/nais.io/v1alpha1"
	skatteetaten_no_v1alpha1 "github.com/nais/liberator/pkg/apis/nebula.skatteetaten.no/v1alpha1"
	sql_cnrm_cloud_google_com_v1beta1 "github.com/nais/liberator/pkg/apis/sql.cnrm.cloud.google.com/v1beta1"
	storage_cnrm_cloud_google_com_v1beta1 "github.com/nais/liberator/pkg/apis/storage.cnrm.cloud.google.com/v1beta1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"

	imagev1beta1 "github.com/fluxcd/image-reflector-controller/api/v1beta1"
	//networkingv1beta1 "istio.io/client-go/pkg/apis/networking/v1beta1"
	//securityv1beta1 "istio.io/client-go/pkg/apis/security/v1beta1"
)

// Creates a new runtime.Scheme object and adds a list of schemes to it.
// This function should be provided with a list of AddToScheme functions.
func Scheme(schemes ...func(*runtime.Scheme) error) (*runtime.Scheme, error) {
	scheme := runtime.NewScheme()
	for _, fn := range schemes {
		err := fn(scheme)
		if err != nil {
			return nil, err
		}
	}
	return scheme, nil
}

// Return a scheme with all native Kubernetes types and all CRDs supported by liberator.
func All() (*runtime.Scheme, error) {
	return Scheme(
		nais_io_v1alpha1.AddToScheme,
		nais_io_v1.AddToScheme,
		skatteetaten_no_v1alpha1.AddToScheme,
		iam_cnrm_cloud_google_com_v1beta1.AddToScheme,
		sql_cnrm_cloud_google_com_v1beta1.AddToScheme,
		bigquery_cnrm_cloud_google_com_v1beta1.AddToScheme,
		storage_cnrm_cloud_google_com_v1beta1.AddToScheme,
		clientgoscheme.AddToScheme,
		aiven_nais_io_v1.AddToScheme,
		kafka_nais_io_v1.AddToScheme,
		azure_microsoft_com_v1alpha1.AddToScheme,
		azure_microsoft_com_v1alpha2.AddToScheme,
		azure_microsoft_com_v1beta1.AddToScheme,
		imagev1beta1.AddToScheme,
	)
}

// Human-readable description of a Kubernetes object metadata.
func TypeName(resource runtime.Object) string {
	var kind, name, namespace string
	typ, err := meta.TypeAccessor(resource)
	if err == nil {
		kind = typ.GetKind()
	}
	obj, err := meta.Accessor(resource)
	if err == nil {
		name = obj.GetName()
		namespace = obj.GetNamespace()
	}
	return fmt.Sprintf("resource '%s' named '%s' in namespace '%s'", kind, name, namespace)
}

func Webhooks(mgr ctrl.Manager) error {
	type webhook interface {
		SetupWebhookWithManager(ctrl.Manager) error
	}

	webhooks := []webhook{
		// List of types that implements one or more webhook interfaces
		&nais_io_v1alpha1.Application{},
	}

	errors := []error{}
	for _, wh := range webhooks {
		if err := wh.SetupWebhookWithManager(mgr); err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("got %v errors: %+v", len(errors), errors)
	}

	return nil
}
