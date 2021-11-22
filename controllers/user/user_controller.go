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
	"fmt"

	"github.com/go-logr/logr"
	cachev1alpha1 "github.com/lilshah/sandbox-operator/api/v1alpha1"
	finalizerUtil "github.com/stakater/operator-utils/util/finalizer"
	reconcilerUtil "github.com/stakater/operator-utils/util/reconciler"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	UserFinalizer string = "tenantoperator.stakater.com/namespace"
)

// UserReconciler reconciles a User object
type UserReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=cache.my.domain,resources=users,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cache.my.domain,resources=users/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cache.my.domain,resources=users/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the User object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *UserReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// your logic here
	log := r.Log.WithValues("user", req.NamespacedName)
	log.Info("Reconciling User: " + req.Name)

	// Fetch the User userInstance
	userInstance := &cachev1alpha1.User{}
	err := r.Get(ctx, req.NamespacedName, userInstance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("User resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		log.Error(err, "unable to fetch User")
		return ctrl.Result{}, err
	}

	// Check if the SandBox instance is marked to be deleted, which is
	// indicated by the deletion timestamp being set.
	isUserMarkedToBeDeleted := userInstance.GetDeletionTimestamp() != nil
	if isUserMarkedToBeDeleted {
		log.Info(fmt.Sprintf("User %s is marked to be deleted, waiting for finalizers", req.Name))
		if finalizerUtil.HasFinalizer(userInstance, UserFinalizer) {
			return r.handleDelete(ctx, req, userInstance)
		}
	}

	// Add finalizer if not present
	if !finalizerUtil.HasFinalizer(userInstance, UserFinalizer) {
		log.Info("Adding finalizer for SandBox: " + req.Name)
		finalizerUtil.AddFinalizer(userInstance, UserFinalizer)
		err := r.Client.Update(ctx, userInstance)
		if err != nil {
			return reconcilerUtil.RequeueWithError(err)
		}
		return ctrl.Result{}, nil
	}
	// If the User instance is not marked to be deleted, then it must be
	// created or updated, so enqueue a reconcile request.

	return r.handleCreateUpdate(ctx, req, userInstance)
}

// SetupWithManager sets up the controller with the Manager.
func (r *UserReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cachev1alpha1.User{}).
		Complete(r)
}
