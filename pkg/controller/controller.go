package controller

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func CreateWithManager[R NaisResource](mgr ctrl.Manager, inner NaisReconciler[R], logger log.FieldLogger) (reconcile.Reconciler, error) {
	r := &reconciler[R]{
		Client:        mgr.GetClient(),
		eventRecorder: mgr.GetEventRecorderFor(inner.Name()),
		inner:         inner,
		logger:        logger.WithField("controller", inner.Name()),
		finalizer:     fmt.Sprintf("%s.nais.io", inner.Name()),
	}
	return r, ctrl.NewControllerManagedBy(mgr).
		For(r.inner.New()).
		Complete(r)
}

type reconciler[R NaisResource] struct {
	client.Client
	inner           NaisReconciler[R]
	logger          log.FieldLogger
	requeueInterval time.Duration
	eventRecorder   record.EventRecorder
	finalizer       string
}

type NaisReconciler[R NaisResource] interface {
	// New returns a pointer to a new instance of the resource type that this reconciler handles.
	New() R

	// Name returns the name of this reconciler in all lowercase characters.
	Name() string

	// LogDetail returns a set of fields that identify the resource in log messages.
	LogDetail(resource R) log.Fields

	// Process should make any necessary changes to bring the real world in line with the resource specification.
	//
	// When Process encounters an error, it is responsible for informing the user/operator about the error.
	// This should be done via the logger, and by emitting an event attached to the resource if possible.
	//
	// If the reconciler needs to set custom fields on the resource status, these fields must be set on the status
	// before returning. The status will be saved to the cluster after Process is complete.
	Process(resource R, logger log.FieldLogger, eventRecorder record.EventRecorder) ReconcileResult

	// Delete is called to clean up any external dependencies when the resource is scheduled for deletion.
	//
	// When Delete encounters an error, it is responsible for informing the user/operator about the error.
	// This should be done via the logger, and by emitting an event attached to the resource if possible.
	//
	// Returns true if deletion was successful, false otherwise.
	Delete(resource R, logger log.FieldLogger, eventRecorder record.EventRecorder) error

	// NeedsProcessing is called to determine if the resource needs to be processed.
	// When NeedsProcessing is called, the resource has already been checked for changes in hash, so most reconcilers should just return false.
	NeedsProcessing(resource R, logger log.FieldLogger, eventRecorder record.EventRecorder) bool
}

func (r *reconciler[R]) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := r.logger.WithFields(log.Fields{
		"namespace": req.Namespace,
		"name":      req.Name,
	})

	logger.Infof("Reconciling")
	defer func() {
		logger.Infof("Finished reconciling")
	}()

	fail := func(err error, requeue bool) (ctrl.Result, error) {
		logger.Error(err)
		cr := &ctrl.Result{}
		if requeue {
			cr.RequeueAfter = r.requeueInterval
		}
		return *cr, nil
	}

	resource := r.inner.New()

	err := r.Get(ctx, req.NamespacedName, resource)
	switch {
	case errors.IsNotFound(err):
		return fail(fmt.Errorf("resource deleted from cluster; noop"), false)
	case err != nil:
		return fail(fmt.Errorf("unable to retrieve resource from cluster: %s", err), true)
	}

	logger = logger.WithFields(r.inner.LogDetail(resource))

	status := resource.GetStatus()
	status.UpdateObservationData(resource.GetGeneration())
	if status.SynchronizationTime != nil {
		logger.Infof("Last synchronization time: %s", status.SynchronizationTime)
	} else {
		logger.Info("Resource not synchronized before")
	}

	if resource.GetDeletionTimestamp() != nil {
		logger.Info("Resource is marked for deletion")
		err = r.inner.Delete(resource, logger, r.eventRecorder)
		if err != nil {
			return fail(fmt.Errorf("failed to delete resource: %w", err), true)
		}
		if controllerutil.RemoveFinalizer(resource, r.finalizer) {
			err = r.Update(ctx, resource)
			if err != nil {
				return fail(fmt.Errorf("failed to remove finalizer: %w", err), true)
			}
		}
		return ctrl.Result{}, nil
	}

	err = resource.ApplyDefaults()
	if err != nil {
		return fail(err, false)
	}

	hash, err := resource.Hash()
	if err != nil {
		return fail(err, false)
	}

	if !r.needsProcessing(resource, hash, logger) {
		return ctrl.Result{}, nil
	}

	result := r.inner.Process(resource, logger, r.eventRecorder)

	status.UpdateSynchronizationData(hash, result.State, resource.GetGeneration())
	conditions := append(result.Conditions, r.calculateSucess(result.Conditions, result.State))
	status.UpdateConditions(conditions, resource.GetGeneration())

	if controllerutil.AddFinalizer(resource, r.finalizer) {
		err = r.Update(ctx, resource)
		if err != nil {
			return fail(err, true)
		}
	}

	err = r.Status().Update(ctx, resource)
	if err != nil {
		return fail(err, true)
	}

	if result.State == SynchronizationStateFailed {
		return fail(fmt.Errorf("synchronisation failed"), result.Requeue)
	}

	return ctrl.Result{}, nil
}

func (r *reconciler[R]) needsProcessing(resource R, hash string, logger log.FieldLogger) bool {
	status := resource.GetStatus()
	if status.SynchronizationHash != hash {
		logger.Infof("Resource has changed since last synchronization: %s != %s", status.SynchronizationHash, hash)
		return true
	}

	return r.inner.NeedsProcessing(resource, logger, r.eventRecorder)
}

// calculateSucess uses the other conditions to calculate the overall success of the reconciliation
//
// If the reconciliation fails, reason and message are copied from the first failing condition
// Conditions should be a list of possible failure states
func (r *reconciler[R]) calculateSucess(conditions []metav1.Condition, state SynchronizationState) metav1.Condition {
	var succeededState metav1.ConditionStatus
	if state == SynchronizationStateFailed {
		succeededState = metav1.ConditionFalse
	} else {
		succeededState = metav1.ConditionTrue
	}

	succeeded := metav1.Condition{
		Type:    "Succeeded",
		Status:  succeededState,
		Reason:  "",
		Message: "",
	}

	for _, c := range conditions {
		if c.Status == metav1.ConditionTrue {
			succeeded.Reason = c.Reason
			succeeded.Message = c.Message
			break
		}
	}
	return succeeded
}
