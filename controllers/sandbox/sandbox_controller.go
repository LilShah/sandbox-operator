/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	finalizerUtil "github.com/stakater/operator-utils/util/finalizer"
	reconcilerUtil "github.com/stakater/operator-utils/util/reconciler"

	"github.com/go-logr/logr"
	cachev1alpha1 "github.com/lilshah/sandbox-operator/api/v1alpha1"
)

const (
	SandBoxFinalizer string = "tenantoperator.stakater.com/namespace"
)

// SandBoxReconciler reconciles a SandBox object
type SandBoxReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=cache.my.domain,resources=sandboxes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cache.my.domain,resources=sandboxes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cache.my.domain,resources=sandboxes/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the SandBox object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *SandBoxReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// your logic here
	log := r.Log.WithValues("sandbox", req.NamespacedName)
	log.Info("Reconciling SandBox: " + req.Name)

	// Fetch the SandBox instance
	instance := &cachev1alpha1.SandBox{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("SandBox resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		log.Error(err, "unable to fetch SandBox")
		return ctrl.Result{}, err
	}

	// Check if the SandBox instance is marked to be deleted, which is
	// indicated by the deletion timestamp being set.
	isSandBoxMarkedToBeDeleted := instance.GetDeletionTimestamp() != nil
	if isSandBoxMarkedToBeDeleted {
		log.Info("SandBox %b is marked to be deleted, waiting for finalizers", req.Name)
		if finalizerUtil.HasFinalizer(instance, SandBoxFinalizer) {
			return r.handleDelete(ctx, req, instance)
		}
	}
	// Add finalizer if not present
	if !finalizerUtil.HasFinalizer(instance, SandBoxFinalizer) {
		log.Info("Adding finalizer for SandBox: " + req.Name)
		finalizerUtil.AddFinalizer(instance, SandBoxFinalizer)
		err := r.Client.Update(ctx, instance)
		if err != nil {
			return reconcilerUtil.RequeueWithError(err)
		}
		return ctrl.Result{}, nil
	}
	// If the SandBox instance is not marked to be deleted, then it must be
	// created or updated, so enqueue a reconcile request.

	return r.handleCreateUpdate(ctx, req, instance)
}

// SetupWithManager sets up the controller with the Manager.
func (r *SandBoxReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cachev1alpha1.SandBox{}).
		Complete(r)
}
