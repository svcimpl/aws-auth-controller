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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AWSAuthSpec defines the desired state of AWSAuth
type AWSAuthSpec struct {
	//The AWSRole that is assumed by AWS Principal
	AWSRole string `json:"aws_role,omitempty"`
	//The namespace that this AWSRole has access to
	TargetNamespace string `json:"target_namespace,omitempty"`
	//Kubernetes role name for the RBAC
	KubernetesRole string `json:"kubernetes_role,omitempty"`
	//Kubernetes username for RBAC RoleBinding to AWS AUTH ConfigMap
	KubernetesUserName string `json:"kubernetes_user_name,omitempty"`
}

// AWSAuthStatus defines the observed state of AWSAuth
type AWSAuthStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// AWSAuth is the Schema for the awsauths API
type AWSAuth struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AWSAuthSpec   `json:"spec,omitempty"`
	Status AWSAuthStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AWSAuthList contains a list of AWSAuth
type AWSAuthList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AWSAuth `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AWSAuth{}, &AWSAuthList{})
}
