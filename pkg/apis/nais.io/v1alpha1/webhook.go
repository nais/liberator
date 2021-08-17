package nais_io_v1alpha1

import (
	"fmt"

	"github.com/nais/liberator/pkg/webhookvalidator"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

func (a *Application) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(a).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-nais-io-v1alpha1-application,mutating=false,failurePolicy=fail,groups=nais.io,resources=application,versions=v1alpha1,name=validation.application.nais.io

var _ webhook.Validator = &Application{}

func (a *Application) ValidateCreate() error {
	return nil
}

func (a *Application) ValidateUpdate(old runtime.Object) error {
	// Type-cast from runtime.Object to Application
	oldA, ok := old.(*Application)
	if !ok {
		return apierrors.NewBadRequest(fmt.Sprintf("expected an Application but got a %T", old))
	}
	// Perform actual comparison
	err := webhookvalidator.DeepComparison(a.Spec, oldA.Spec, field.NewPath("spec"))
	if err != nil {
		if allErrs, ok := err.(errors.Aggregate); ok {
			return apierrors.NewInvalid(
				schema.GroupKind{Group: GroupVersion.Group, Kind: "Application"},
				a.Name,
				fromAggregate(allErrs),
			)
		}
		return err
	}

	return nil
}

func (a *Application) ValidateDelete() error {
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
