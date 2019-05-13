package operator2

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

const (
	injectCABundleAnnotationName  = "service.alpha.openshift.io/inject-cabundle"
	injectCABundleAnnotationValue = "true"
)

func (c *authOperator) handleServiceCA() (*corev1.ConfigMap, *corev1.Secret, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cm := c.configMaps.ConfigMaps(targetName)
	secret := c.secrets.Secrets(targetName)
	serviceCA, err := cm.Get(serviceCAName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		serviceCA, err = cm.Create(defaultServiceCA())
	}
	if err != nil {
		return nil, nil, err
	}
	if len(serviceCA.Data[serviceCAKey]) == 0 {
		return nil, nil, fmt.Errorf("config map has no service ca data: %#v", serviceCA)
	}
	if err := isValidServiceCA(serviceCA); err != nil {
		klog.Infof("deleting invalid service CA config map: %#v", serviceCA)
		opts := &metav1.DeleteOptions{Preconditions: &metav1.Preconditions{UID: &serviceCA.UID}}
		if err := cm.Delete(serviceCA.Name, opts); err != nil && !errors.IsNotFound(err) {
			klog.Infof("failed to delete invalid service CA config map: %v", err)
		}
		return nil, nil, err
	}
	servingCert, err := secret.Get(servingCertName, metav1.GetOptions{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get serving cert: %v", err)
	}
	return serviceCA, servingCert, nil
}
func isValidServiceCA(ca *corev1.ConfigMap) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if ca.Annotations[injectCABundleAnnotationName] != injectCABundleAnnotationValue {
		return fmt.Errorf("config map missing injection annotation: %#v", ca)
	}
	return nil
}
func defaultServiceCA() *corev1.ConfigMap {
	_logClusterCodePath()
	defer _logClusterCodePath()
	meta := defaultMeta()
	meta.Name = serviceCAName
	meta.Annotations[injectCABundleAnnotationName] = injectCABundleAnnotationValue
	return &corev1.ConfigMap{ObjectMeta: meta}
}
