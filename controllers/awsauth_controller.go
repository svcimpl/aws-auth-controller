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
	_ "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	authv1alpha1 "svcimpl.com/aws-auth-controller/api/v1alpha1"
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

	cm_ns_name := types.NamespacedName{
		Namespace: "kube-system",
		Name:      "aws-auth",
	}

	var customawsauth authv1alpha1.AWSAuth
	// your logic here
	if err := r.Get(ctx, req.NamespacedName, &customawsauth); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	marshal, err := json.Marshal(customawsauth)
	if err != nil {
		return ctrl.Result{}, err
	}
	log.Info(fmt.Sprintf("CRD AWSAuth is :\n %s", string(marshal)))
	//Now do the reconcile based on what we get
	//We will get the AWSRole
	//We need to get the aws-auth config map
	var awsAuthConfigmap corev1.ConfigMap
	if err := r.Get(ctx, cm_ns_name, &awsAuthConfigmap); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	awsCm, err := json.Marshal(awsAuthConfigmap)
	if err != nil {
		return ctrl.Result{}, err
	}
	log.Info(fmt.Sprintf("AWS-Auth Configmap is is :\n %s", string(awsCm)))
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AWSAuthReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&authv1alpha1.AWSAuth{}).
		Complete(r)
}
