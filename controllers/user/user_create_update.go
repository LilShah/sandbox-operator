package controllers

import (
	"context"
	"time"

	"github.com/lilshah/sandbox-operator/api/v1alpha1"
	reconcilerUtil "github.com/stakater/operator-utils/util/reconciler"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *UserReconciler) handleCreateUpdate(ctx context.Context, req ctrl.Request, user *v1alpha1.User, sandbox *v1alpha1.SandBox) (ctrl.Result, error) {
	log := r.Log.WithValues("user", req.NamespacedName)
	log.Info("Creating/updating user: " + user.ObjectMeta.Name)

	// namespace := &corev1.Namespace{}
	// err := r.Client.Get(ctx, types.NamespacedName{Name: user.ObjectMeta.Name}, namespace)
	// if err != nil {
	// 	log.Error(err, "Failed to get namespace")
	// 	return reconcilerUtil.RequeueAfter(60 * time.Second)
	// }
	// r.Client.Create(ctx, namespace)
	return reconcilerUtil.RequeueAfter(60 * time.Second)
}
