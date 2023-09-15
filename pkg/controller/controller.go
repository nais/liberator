package controller

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func NewReconciler[R client.Object](c client.Client, inner NaisReconciler[R], scheme *runtime.Scheme, logger log.FieldLogger) *Reconciler[R] {
	return &Reconciler[R]{
		Client: c,
		inner:  inner,
		Scheme: scheme,
		logger: logger,
	}
}

type Reconciler[R client.Object] struct {
	client.Client
	inner           NaisReconciler[R]
	Scheme          *runtime.Scheme
	logger          log.FieldLogger
	RequeueInterval time.Duration
}

type NaisReconciler[R client.Object] interface {
	New() R
	LogDetail(resource R) log.Fields
	GetStatus(resource R) *NaisStatus
	SetStatus(resource R, status *NaisStatus)
	Process(resource R, logger log.FieldLogger) ReconcileResult
	Finalizer() string
}

func (r *Reconciler[R]) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
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
			cr.RequeueAfter = r.RequeueInterval
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

	status := r.inner.GetStatus(resource)
	if status != nil {
		logger.Infof("Last synchronization time: %s", status.SynchronizationTime)
	} else {
		logger.Info("Resource not synchronized before")
	}

	result := r.inner.Process(resource, logger)

	if result.Skipped {
		return ctrl.Result{}, nil
	}

	// TODO: Metrics
	// defer func() {
	// 	metrics.TopicsProcessed.With(prometheus.Labels{
	// 		metrics.LabelSyncState: result.Status.SynchronizationState,
	// 		metrics.LabelPool:      topic.Spec.Pool,
	// 	}).Inc()
	// }()

	if result.Error != nil {
		r.inner.SetStatus(resource, result.Status)
		err = r.Update(ctx, resource)
		if err != nil {
			logger.Errorf("Write resource status: %s", err)
		}
		return fail(result.Error, result.Requeue)
	}

	// If delete was finalized, mark resource as finally deleted by removing finalizer.
	// Otherwise, append finalizer to finalizers to ensure proper cleanup when resource is deleted
	if result.DeleteFinalized {
		controllerutil.RemoveFinalizer(resource, r.inner.Finalizer())
	} else {
		controllerutil.AddFinalizer(resource, r.inner.Finalizer())
	}

	// Write resource status; retry always
	r.inner.SetStatus(resource, result.Status)
	err = r.Update(ctx, resource)
	if err != nil {
		return fail(err, true)
	}

	logger.Infof("Resource object updated: %s", result.Status.SynchronizationState)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler[R]) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(r.inner.New()).
		Complete(r)
}
