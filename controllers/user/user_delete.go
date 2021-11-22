package controllers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lilshah/sandbox-operator/api/v1alpha1"
	finalizerUtil "github.com/stakater/operator-utils/util/finalizer"
	reconcilerUtil "github.com/stakater/operator-utils/util/reconciler"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *UserReconciler) handleDelete(ctx context.Context, req ctrl.Request, user *v1alpha1.User) (ctrl.Result, error) {
	log := r.Log.WithValues("namespace", req.NamespacedName)
	log.Info(fmt.Sprintf("Deleting user: %s", req.Name))

	// Delete sandboxes for user
	for i := 0; i < user.Spec.SandBoxCount; i++ {
		sandboxInstance := &v1alpha1.SandBox{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "sb-" + req.Name + "-" + strconv.Itoa(i),
				Namespace: req.Namespace,
			},
			Spec: v1alpha1.SandBoxSpec{
				Name: "SB" + user.ObjectMeta.Name + "-" + strconv.Itoa(i),
				Type: "T1",
			},
		}
		err := r.Client.Get(ctx, types.NamespacedName{Name: sandboxInstance.ObjectMeta.Name, Namespace: req.Namespace}, sandboxInstance)
		if err != nil {
			log.Error(err, fmt.Sprintf("Failed to get sandbox \"%s\" for user \"%s\"", sandboxInstance.ObjectMeta.Name, user.ObjectMeta.Name))
			return reconcilerUtil.RequeueWithError(err)
		}
		err = r.Client.Delete(ctx, sandboxInstance)
		if err != nil {
			log.Error(err, "Failed to delete sandbox %s", sandboxInstance.ObjectMeta.Name)
			return reconcilerUtil.RequeueWithError(err)
		}
	}

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
