package controllers

import (
	"context"
	"fmt"

	"github.com/lilshah/sandbox-operator/api/v1alpha1"
	finalizerUtil "github.com/stakater/operator-utils/util/finalizer"
	reconcilerUtil "github.com/stakater/operator-utils/util/reconciler"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *UserReconciler) handleDelete(ctx context.Context, req ctrl.Request, user *v1alpha1.User, sandbox *v1alpha1.SandBox) (ctrl.Result, error) {
	log := r.Log.WithValues("namespace", req.NamespacedName)
	log.Info(fmt.Sprintf("Deleting namespace: %s", req.Name))

	// Delete sandbox
	r.Client.Get(ctx, types.NamespacedName{Name: sandbox.ObjectMeta.Name, Namespace: req.Namespace}, sandbox)
	r.Client.Delete(ctx, sandbox)

	// Delete finalizer
	finalizerUtil.DeleteFinalizer(user, UserFinalizer)
	log.Info("Finalizer removed for namespace: " + user.ObjectMeta.Name)

	// Update the User instance
	err := r.Client.Update(ctx, user)
	if err != nil {
		log.Error(err, "Failed to update SandBox instance")
		return reconcilerUtil.ManageError(r.Client, user, err, false)
	}

	return ctrl.Result{}, nil
}
