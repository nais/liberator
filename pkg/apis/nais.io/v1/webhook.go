package nais_io_v1

import (
	"context"
	"fmt"
	"time"

	aiven_io_v1alpha1 "github.com/nais/liberator/pkg/apis/aiven.io/v1alpha1"
	"github.com/nais/liberator/pkg/webhookvalidator"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var _ webhook.CustomValidator = &JobValidator{}

type JobValidator struct {
	client.Client
}

func SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		WithValidator(&JobValidator{Client: mgr.GetClient()}).
		For(&Naisjob{}).
		Complete()
}

// The generated manifest is invalid, so we use kubebuilder to make the initial manifest, and then update with annotations and correct name manually
// DISABLE: +kubebuilder:webhook:verbs=create;update,path=/validate-nais-io-v1-naisjobs,mutating=false,failurePolicy=fail,groups=nais.io,resources=naisjobs,versions=v1,name=validation.naisjobs.nais.io

func (v *JobValidator) ValidateCreate(ctx context.Context, obj runtime.Object) (warnings admission.Warnings, err error) {
	nj, ok := obj.(*Naisjob)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a Naisjob but got a %T", obj))
	}

	if len(nj.GetName()) > validation.LabelValueMaxLength {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("Naisjob name length must be no more than %d characters", validation.LabelValueMaxLength))
	}

	if nj.Spec.TTL != "" {
		if _, err := time.ParseDuration(nj.Spec.TTL); err != nil {
			return nil, apierrors.NewBadRequest(fmt.Sprintf("TTL is not a valid duration: %q. Example of valid duration is '12h'", nj.Spec.TTL))
		}
	}

	if err := v.checkAivenReferences(ctx, nj); err != nil {
		return nil, err
	}

	return nil, nil
}

func (v *JobValidator) ValidateUpdate(ctx context.Context, oldObj runtime.Object, newObj runtime.Object) (warnings admission.Warnings, err error) {
	// Type-cast from runtime.Object to Naisjob
	oldA, ok := oldObj.(*Naisjob)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a Naisjob but got a %T", oldObj))
	}
	nj, ok := newObj.(*Naisjob)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a Naisjob but got a %T", newObj))
	}

	// Perform actual comparison
	if err := webhookvalidator.NaisCompare(nj.Spec, oldA.Spec, field.NewPath("spec")); err != nil {
		if allErrs, ok := err.(errors.Aggregate); ok {
			return nil, apierrors.NewInvalid(
				schema.GroupKind{Group: GroupVersion.Group, Kind: "Naisjob"},
				nj.Name,
				fromAggregate(allErrs),
			)
		}

		return nil, err
	}

	if err := v.checkAivenReferences(ctx, nj); err != nil {
		return nil, err
	}

	return nil, nil
}

func (v *JobValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (warnings admission.Warnings, err error) {
	return nil, nil
}

func (v *JobValidator) checkAivenReferences(ctx context.Context, nj *Naisjob) error {
	if nj.Spec.OpenSearch != nil && nj.Spec.OpenSearch.Instance != "" {
		fullyQualifiedName := aiven_io_v1alpha1.OpenSearchFullyQualifiedName(nj.Spec.OpenSearch.Instance, nj.Namespace)
		opensearch := &aiven_io_v1alpha1.OpenSearch{}
		if err := v.Get(ctx, client.ObjectKey{Name: fullyQualifiedName, Namespace: nj.Namespace}, opensearch); err != nil {
			if apierrors.IsNotFound(err) {
				return apierrors.NewBadRequest(fmt.Sprintf("referenced OpenSearch instance '%s' not found in namespace '%s'", nj.Spec.OpenSearch.Instance, nj.Namespace))
			}
			return apierrors.NewInternalError(fmt.Errorf("could not validate OpenSearch reference: %w", err))
		}
	}

	for _, valkey := range nj.Spec.Valkey {
		fullyQualifiedName := aiven_io_v1alpha1.ValkeyFullyQualifiedName(valkey.Instance, nj.Namespace)
		valkeyObj := &aiven_io_v1alpha1.Valkey{}
		if err := v.Get(ctx, client.ObjectKey{Name: fullyQualifiedName, Namespace: nj.Namespace}, valkeyObj); err != nil {
			if apierrors.IsNotFound(err) {
				return apierrors.NewBadRequest(fmt.Sprintf("referenced Valkey instance '%s' not found in namespace '%s'", valkey.Instance, nj.Namespace))
			}
			return apierrors.NewInternalError(fmt.Errorf("could not validate Valkey reference: %w", err))
		}
	}
	return nil
}

func fromAggregate(agg errors.Aggregate) field.ErrorList {
	errs := agg.Errors()
	list := make(field.ErrorList, len(errs))
	for i := range errs {
		list[i] = errs[i].(*field.Error)
	}
	return list
}
