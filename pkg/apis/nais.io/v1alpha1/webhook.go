package nais_io_v1alpha1

import (
	"context"
	"fmt"
	"time"

	"github.com/nais/liberator/pkg/webhookvalidator"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var _ webhook.CustomValidator = &Application{}

func (a *Application) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		WithValidator(&Application{}).
		For(a).
		Complete()
}

// The generated manifest is invalid, so we use kubebuilder to make the initial manifest, and then update with annotations and correct name manually
// DISABLE: +kubebuilder:webhook:verbs=create;update,path=/validate-nais-io-v1alpha1-applications,mutating=false,failurePolicy=fail,groups=nais.io,resources=applications,versions=v1alpha1,name=validation.applications.nais.io

func (a *Application) ValidateCreate(ctx context.Context, obj runtime.Object) (warnings admission.Warnings, err error) {
	if len(a.GetName()) > validation.LabelValueMaxLength {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("Application name length must be no more than %d characters", validation.LabelValueMaxLength))
	}

	if a.Spec.TTL != "" {
		if _, err := time.ParseDuration(a.Spec.TTL); err != nil {
			return nil, apierrors.NewBadRequest(fmt.Sprintf("TTL is not a valid duration: %q. Example of valid duration is '12h'", a.Spec.TTL))
		}
	}

	return nil, nil
}

func (a *Application) ValidateUpdate(ctx context.Context, oldObj runtime.Object, newObj runtime.Object) (warnings admission.Warnings, err error) {
	// Type-cast from runtime.Object to Application
	oldA, ok := oldObj.(*Application)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected an Application but got a %T", oldObj))
	}
	// Perform actual comparison
	if err := webhookvalidator.NaisCompare(a.Spec, oldA.Spec, field.NewPath("spec")); err != nil {
		if allErrs, ok := err.(errors.Aggregate); ok {
			return nil, apierrors.NewInvalid(
				schema.GroupKind{Group: GroupVersion.Group, Kind: "Application"},
				a.Name,
				fromAggregate(allErrs),
			)
		}

		return nil, err
	}

	return nil, nil
}

func (a *Application) ValidateDelete(ctx context.Context, obj runtime.Object) (warnings admission.Warnings, err error) {
	return nil, nil
}

func fromAggregate(agg errors.Aggregate) field.ErrorList {
	errs := agg.Errors()
	list := make(field.ErrorList, len(errs))
	for i := range errs {
		list[i] = errs[i].(*field.Error)
	}
	return list
}
