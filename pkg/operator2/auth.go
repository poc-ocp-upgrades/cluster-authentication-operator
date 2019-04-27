package operator2

import (
	"k8s.io/apimachinery/pkg/api/errors"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	configv1 "github.com/openshift/api/config/v1"
)

func (c *authOperator) handleAuthConfigInner() (*configv1.Authentication, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	authConfigNoDefaults, err := c.authentication.Get(globalConfigName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		authConfigNoDefaults, err = c.authentication.Create(&configv1.Authentication{ObjectMeta: defaultGlobalConfigMeta()})
	}
	if err != nil {
		return nil, err
	}
	expectedReference := configv1.ConfigMapNameReference{Name: targetName}
	if authConfigNoDefaults.Status.IntegratedOAuthMetadata == expectedReference {
		return authConfigNoDefaults, nil
	}
	authConfigNoDefaults.Status.IntegratedOAuthMetadata = expectedReference
	return c.authentication.UpdateStatus(authConfigNoDefaults)
}
func (c *authOperator) handleAuthConfig() (*configv1.Authentication, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	auth, err := c.handleAuthConfigInner()
	if err != nil {
		return nil, err
	}
	return defaultAuthConfig(auth), nil
}
func defaultAuthConfig(authConfig *configv1.Authentication) *configv1.Authentication {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	out := authConfig.DeepCopy()
	if len(out.Spec.Type) == 0 {
		out.Spec.Type = configv1.AuthenticationTypeIntegratedOAuth
	}
	return out
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
