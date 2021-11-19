package controllers

import (
	"context"
	"time"

	"github.com/lilshah/sandbox-operator/api/v1alpha1"
	reconcilerUtil "github.com/stakater/operator-utils/util/reconciler"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *SandBoxReconciler) handleCreateUpdate(ctx context.Context, req ctrl.Request, sandbox *v1alpha1.SandBox) (ctrl.Result, error) {
	log := r.Log.WithValues("sandbox", req.NamespacedName)
	log.Info("Creating/updating SandBox: " + sandbox.ObjectMeta.Name)
	namespace := &corev1.Namespace{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: sandbox.ObjectMeta.Name}, namespace)
	if err != nil {
		log.Error(err, "Failed to get namespace")
		return reconcilerUtil.RequeueAfter(60 * time.Second)
	}
	r.Client.Create(ctx, namespace)
	return reconcilerUtil.RequeueAfter(60 * time.Second)
}
