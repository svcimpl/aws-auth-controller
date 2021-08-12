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
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	authv1alpha1 "svcimpl.com/aws-auth-controller/api/v1alpha1"
)

const (
	AwsAuthNamespace = "kube-system"
	AwsAuthName      = "aws-auth"
	//name of our custom finalizer
	FinalizerName = "auth.svcimpl.com/finalizer"
)

// AWSAuthReconciler reconciles a AWSAuth object
type AWSAuthReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=auth.svcimpl.com,resources=awsauths,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=auth.svcimpl.com,resources=awsauths/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=auth.svcimpl.com,resources=awsauths/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the AWSAuth object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *AWSAuthReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	//start := time.Now()
	log := ctrl.LoggerFrom(ctx)

	cmNsName := types.NamespacedName{
		Namespace: AwsAuthNamespace,
		Name:      AwsAuthName,
	}

	var customawsauth authv1alpha1.AWSAuth
	//Get CRD
	if err := r.Get(ctx, req.NamespacedName, &customawsauth); err != nil {
		log.Info("Delete reconcile received after finalization")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	log.Info("Successfully retrieved CustomAWSAuth")

	//Register the finalizer on the CRD if not already done
	if err := r.registerFinalizer(ctx, &customawsauth); err != nil {
		return ctrl.Result{}, err
	}

	//Our reconcile target
	var awsAuthConfigmap corev1.ConfigMap
	//Check if this reconcile call was a delete
	if !customawsauth.ObjectMeta.DeletionTimestamp.IsZero() {
		// Resource is marked for deletion
		// Clean-up the AWSAuth and remove the finalizer so that the source CustomAWSAuth can be deleted
		if err := r.Get(ctx, cmNsName, &awsAuthConfigmap); err != nil {
			if client.IgnoreNotFound(err) != nil { //we have a real error
				return ctrl.Result{}, err
			}
			// There is no config map for aws-auth, so nothing for us to clean-up
		} else {
			//Clean up the aws-auth configmap
			err := r.deleteRoleMapping(ctx, &awsAuthConfigmap, &customawsauth)
			if err != nil {
				return ctrl.Result{}, err
			}
			log.Info("RoleMapping deleted in aws-auth configmap")
		}
		// remove our finalizer from the list and update it.
		controllerutil.RemoveFinalizer(&customawsauth, FinalizerName)
		if err := r.Update(ctx, &customawsauth); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	//If we came here, it is not a delete reconcile

	//Get AWSAuth ConfigMap
	if err := r.Get(ctx, cmNsName, &awsAuthConfigmap); err != nil {
		if client.IgnoreNotFound(err) == nil { //no existing aws-auth configmap
			err := r.createOrUpdateRoleMapping(ctx, nil, &customawsauth)
			if err != nil {
				return ctrl.Result{}, err
			}
			log.Info("RoleMapping entry created in aws-auth configmap")
		} else { //some other error
			return ctrl.Result{}, err
		}
	}
	err := r.createOrUpdateRoleMapping(ctx, &awsAuthConfigmap, &customawsauth)
	if err != nil {
		return ctrl.Result{}, err
	}
	log.Info("RoleMapping entry inserted/updated in aws-auth configmap")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AWSAuthReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&authv1alpha1.AWSAuth{}).
		Complete(r)
}
