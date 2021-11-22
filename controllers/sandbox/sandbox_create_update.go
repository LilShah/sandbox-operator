package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/lilshah/sandbox-operator/api/v1alpha1"
	reconcilerUtil "github.com/stakater/operator-utils/util/reconciler"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *SandBoxReconciler) handleCreateUpdate(ctx context.Context, req ctrl.Request, sandbox *v1alpha1.SandBox) (ctrl.Result, error) {
	log := r.Log.WithValues("sandbox", req.NamespacedName)
	log.Info("Creating/updating SandBox: " + sandbox.ObjectMeta.Name)
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: sandbox.ObjectMeta.Name,
		},
	}
	err := r.Client.Create(ctx, namespace)
	if errors.IsAlreadyExists(err) {
		log.Info(fmt.Sprintf("Namespace \"%s\" already exists for sandbox \"%s\"", namespace.ObjectMeta.Name, sandbox.ObjectMeta.Name))
		return reconcilerUtil.RequeueAfter(60 * time.Second)
	}

	if err != nil {
		log.Error(err, "Failed to create namespace")
		return reconcilerUtil.RequeueWithError(err)
	}
	return reconcilerUtil.RequeueAfter(60 * time.Second)
}
