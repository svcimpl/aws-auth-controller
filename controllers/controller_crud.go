package controllers

import (
	"context"
	yaml "gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	authv1alpha1 "svcimpl.com/aws-auth-controller/api/v1alpha1"
	"svcimpl.com/aws-auth-controller/mapper"
)

func (r *AWSAuthReconciler) createOrUpdateRoleMapping(ctx context.Context, cm *corev1.ConfigMap, mr *authv1alpha1.AWSAuth) error {
	var create bool

	if cm == nil {
		//We are creating the object map
		cm = createEmptyAWSAuthConfigMap()
		create = true
	}
	authData, err := parseAWSAuthConfigMap(cm)
	if err != nil {
		return err
	}
	ram := &mapper.RoleAuthMap{
		RoleARN:  mr.Spec.AWSRole,
		Username: mr.Spec.KubernetesUserName,
		Groups:   []string{},
	}
	if create {
		roles := append(authData.MapRoles, ram)
		mapRoles, err := yaml.Marshal(roles)
		if err != nil {
			return err
		}
		cm.Data["mapRoles"] = string(mapRoles)
		err = r.Create(ctx, cm)
		if err != nil {
			return err
		}
	} else { //Update
		var i = 0
		var found bool
		roles := authData.MapRoles
		for _, mrval := range roles {
			if mr.Spec.AWSRole == mrval.RoleARN {
				found = true
				break
			}
			i++
		}
		if found {
			roles = append(roles[:i], roles[i+1:]...)
		}
		roles = append(roles, ram)
		mapRoles, err := yaml.Marshal(roles)
		if err != nil {
			return err
		}
		cm.Data["mapRoles"] = string(mapRoles)
		err = r.Update(ctx, cm)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *AWSAuthReconciler) deleteRoleMapping(ctx context.Context, cm *corev1.ConfigMap, mr *authv1alpha1.AWSAuth) error {
	//Get the RoleMapping from the configmap
	authData, err := parseAWSAuthConfigMap(cm)
	if err != nil {
		return err
	}
	roles := authData.MapRoles
	var found bool
	var i = 0
	//Identify the role to delete, delete it and persist the cm
	for _, mrval := range roles {
		if mr.Spec.AWSRole == mrval.RoleARN {
			found = true
			break
		}
		i++
	}
	if found {
		roles = append(roles[:i], roles[i+1:]...)
	}
	mapRoles, err := yaml.Marshal(roles)
	if err != nil {
		return err
	}
	cm.Data["mapRoles"] = string(mapRoles)
	//persist cm
	err = r.Update(ctx, cm)
	if err != nil {
		return err
	}
	return nil
}

func (r *AWSAuthReconciler) registerFinalizer(ctx context.Context, customAuth *authv1alpha1.AWSAuth) error {
	if customAuth.ObjectMeta.DeletionTimestamp.IsZero() {
		//Object is not under deletion. Just make sure that there is a finalizer
		//registered on the resource and if not, register one
		if !containsString(customAuth.GetFinalizers(), FinalizerName) {
			controllerutil.AddFinalizer(customAuth, FinalizerName)
			if err := r.Update(ctx, customAuth); err != nil {
				return err
			}
		}
	}
	return nil
}
