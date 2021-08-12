package controllers

import (
	yaml "gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"svcimpl.com/aws-auth-controller/mapper"
)

// Helper functions to check and remove string from a slice of strings.
func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func createEmptyAWSAuthConfigMap() *corev1.ConfigMap {
	configMapObj := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      AwsAuthName,
			Namespace: AwsAuthNamespace,
		},
		Data: map[string]string{},
	}
	return configMapObj
}

func parseAWSAuthConfigMap(awsAuthCM *corev1.ConfigMap) (mapper.AWSAuthData, error) {
	var awsAuthData mapper.AWSAuthData
	err := yaml.Unmarshal([]byte(awsAuthCM.Data["mapRoles"]), &awsAuthData.MapRoles)
	if err != nil {
		return mapper.AWSAuthData{}, err
	}

	err = yaml.Unmarshal([]byte(awsAuthCM.Data["mapUsers"]), &awsAuthData.MapUsers)
	if err != nil {
		return mapper.AWSAuthData{}, err
	}
	return awsAuthData, nil
}
