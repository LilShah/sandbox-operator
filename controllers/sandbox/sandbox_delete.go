package controllers

import (
	"context"
	"fmt"

	"github.com/lilshah/sandbox-operator/api/v1alpha1"
	finalizerUtil "github.com/stakater/operator-utils/util/finalizer"
	reconcilerUtil "github.com/stakater/operator-utils/util/reconciler"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *SandBoxReconciler) handleDelete(ctx context.Context, req ctrl.Request, sandbox *v1alpha1.SandBox) (ctrl.Result, error) {
	log := r.Log.WithValues("namespace", req.NamespacedName)
	log.Info(fmt.Sprintf("Deleting namespace: %s", req.Name))

	// Delete namespace
	namespace := &corev1.Namespace{}
	r.Client.Get(ctx, types.NamespacedName{Name: sandbox.ObjectMeta.Name, Namespace: req.Namespace}, namespace)
	r.Client.Delete(ctx, namespace)

	// Delete finalizer
	finalizerUtil.DeleteFinalizer(sandbox, SandBoxFinalizer)
	log.Info("Finalizer removed for namespace: " + sandbox.ObjectMeta.Name)

	// Update the SandBox instance
	err := r.Client.Update(ctx, sandbox)
	if err != nil {
		log.Error(err, "Failed to update SandBox instance")
		return reconcilerUtil.ManageError(r.Client, sandbox, err, false)
	}

	return ctrl.Result{}, nil
}
