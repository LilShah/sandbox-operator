package controllers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lilshah/sandbox-operator/api/v1alpha1"
	reconcilerUtil "github.com/stakater/operator-utils/util/reconciler"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *UserReconciler) handleCreateUpdate(ctx context.Context, req ctrl.Request, user *v1alpha1.User) (ctrl.Result, error) {
	log := r.Log.WithValues("user", req.NamespacedName)
	log.Info("Creating/updating user: " + user.ObjectMeta.Name)

	// Create sandbox for user
	for i := 0; i < user.Spec.SandBoxCount; i++ {
		sandboxInstance := &v1alpha1.SandBox{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "sb-" + user.ObjectMeta.Name + "-" + strconv.Itoa(i),
				Namespace: req.Namespace,
			},
			Spec: v1alpha1.SandBoxSpec{
				Name: "SB-" + user.ObjectMeta.Name + "-" + strconv.Itoa(i),
				Type: "T1",
			},
		}
		err := r.Client.Create(ctx, sandboxInstance)
		if err != nil {
			if errors.IsAlreadyExists(err) {
				log.Info(fmt.Sprintf("Sandbox \"%s\" already exists for user \"%s\"", sandboxInstance.ObjectMeta.Name, user.ObjectMeta.Name))
				continue
				// return reconcilerUtil.RequeueAfter(60 * time.Second)
			}
			log.Error(err, "Failed to create sandbox for user: "+user.ObjectMeta.Name)
			return reconcilerUtil.RequeueWithError(err)
		}
	}

	return reconcilerUtil.ManageSuccess(r.Client, user)
}
